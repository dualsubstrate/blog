package motifeditor

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/png"
	_ "image/png"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const upscaleTimeout = 10 * time.Minute

// UpscaleMode reports how UpscaleToWidth reached the target width.
const (
	UpscaleModeRealESRGAN  = "realesrgan+ffmpeg"
	UpscaleModeWidthResize = "width-resize"
	UpscaleModeCopy        = "copy"
)

// BinaryAvailable reports whether a CLI binary can be executed.
func BinaryAvailable(bin string) bool {
	if strings.TrimSpace(bin) == "" {
		return false
	}
	if strings.Contains(bin, string(os.PathSeparator)) {
		info, err := os.Stat(bin)
		if err != nil {
			return false
		}
		return !info.IsDir()
	}
	_, err := exec.LookPath(bin)
	return err == nil
}

// ImageSize returns width and height for a raster image file.
func ImageSize(path string) (width, height int, err error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()
	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return 0, 0, err
	}
	return cfg.Width, cfg.Height, nil
}

// UpscaleToWidth uses Real-ESRGAN then ffmpeg resize to the exact target width.
// When slideCount > 0, strips studio panorama gaps (600px slices + 5px separators) before processing.
// KeyColor (e.g. #013231): flatten transparency before Real-ESRGAN and restore via colorkey after resize.
func UpscaleToWidth(
	ctx context.Context,
	realesrganBin, model, ffmpegBin,
	inputPath, outputPath string,
	targetWidth int,
	slideCount, panoramaSlideWidth, panoramaGapPx int,
	keyColor string,
) (mode string, err error) {
	if targetWidth < 1 {
		return "", fmt.Errorf("target width must be positive")
	}
	keyColor = strings.TrimSpace(keyColor)

	if !BinaryAvailable(realesrganBin) {
		return "", fmt.Errorf(
			"Real-ESRGAN binary not found (%s). Install realesrgan-ncnn-vulkan and set REALESRGAN_BIN or pass -realesrgan. See tools/motif-editor/README.md",
			realesrganBin,
		)
	}

	var cleanup []string
	defer func() {
		for _, p := range cleanup {
			if p != inputPath && p != outputPath {
				os.Remove(p)
			}
		}
	}()

	work := inputPath
	if slideCount > 0 && panoramaSlideWidth > 0 {
		strippedPath, err := maybeStripPanoramaSeparators(inputPath, slideCount, panoramaSlideWidth, panoramaGapPx)
		if err != nil {
			return "", fmt.Errorf("strip panorama gaps: %w", err)
		}
		if strippedPath != inputPath {
			cleanup = append(cleanup, strippedPath)
			work = strippedPath
		}
	}
	if keyColor != "" {
		prep := filepath.Join(filepath.Dir(outputPath), fmt.Sprintf("upscale-flat-%d.png", time.Now().UnixNano()))
		cleanup = append(cleanup, prep)
		if err := flattenPNG(ctx, ffmpegBin, work, prep, keyColor); err != nil {
			return "", fmt.Errorf("flatten for upscale: %w", err)
		}
		work = prep
	}

	width, height, err := ImageSize(work)
	if err != nil {
		return "", fmt.Errorf("read image size: %w", err)
	}

	if preferWidthResizeOnly(width, height) {
		if err := resizeToWidth(ctx, ffmpegBin, work, outputPath, targetWidth); err != nil {
			return "", err
		}
		if err := maybeRestoreKey(ctx, ffmpegBin, outputPath, keyColor); err != nil {
			return "", err
		}
		if err := trimPNGToAlphaBounds(outputPath, 0); err != nil {
			return "", fmt.Errorf("trim upscale output: %w", err)
		}
		return UpscaleModeWidthResize, nil
	}

	const maxPasses = 5
	for pass := 0; pass < maxPasses && width < targetWidth; pass++ {
		if width*2 > targetWidth {
			break
		}
		scale := 2
		if width*4 <= targetWidth {
			scale = 4
		}

		next := filepath.Join(filepath.Dir(outputPath), fmt.Sprintf("motif-upscale-%d-%d.png", time.Now().UnixNano(), pass))
		cleanup = append(cleanup, next)
		if err := runRealESRGAN(ctx, realesrganBin, work, next, model, scale); err != nil {
			return "", err
		}
		work = next
		width, height, err = ImageSize(work)
		if err != nil {
			return "", fmt.Errorf("read upscaled size: %w", err)
		}
	}

	mode = UpscaleModeRealESRGAN
	if width == targetWidth && work == outputPath {
		// output already at target
	} else if width == targetWidth {
		if err := copyFile(work, outputPath); err != nil {
			return "", err
		}
	} else if err := resizeToWidth(ctx, ffmpegBin, work, outputPath, targetWidth); err != nil {
		return "", err
	}

	if err := maybeRestoreKey(ctx, ffmpegBin, outputPath, keyColor); err != nil {
		return "", err
	}
	if err := trimPNGToAlphaBounds(outputPath, 0); err != nil {
		return "", fmt.Errorf("trim upscale output: %w", err)
	}
	return mode, nil
}

func runRealESRGAN(ctx context.Context, bin, inputPath, outputPath, model string, scale int) error {
	if scale != 4 {
		scale = 2
	}

	installDir, err := resolveRealESRGANInstallDir(bin)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, upscaleTimeout)
	defer cancel()

	width, height, err := ImageSize(inputPath)
	if err != nil {
		return fmt.Errorf("read tile size: %w", err)
	}
	tile := realesrganTileSize(width, height)

	args := []string{
		"-i", inputPath,
		"-o", outputPath,
		"-n", model,
		"-s", strconv.Itoa(scale),
		"-m", filepath.Join(installDir, "models"),
		"-t", strconv.Itoa(tile),
	}

	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Dir = installDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("Real-ESRGAN timed out after %s", upscaleTimeout)
		}
		if errors.Is(err, exec.ErrNotFound) || isExecutableNotFound(err) {
			return fmt.Errorf(
				"Real-ESRGAN binary not found (%s). Install realesrgan-ncnn-vulkan and set REALESRGAN_BIN or pass -realesrgan. See tools/motif-editor/README.md",
				bin,
			)
		}
		msg := strings.TrimSpace(string(out))
		if msg == "" {
			msg = err.Error()
		}
		return fmt.Errorf("Real-ESRGAN failed: %s", msg)
	}

	if _, err := os.Stat(outputPath); err != nil {
		return fmt.Errorf("Real-ESRGAN did not produce an output file")
	}
	return nil
}

func resizeToWidth(ctx context.Context, ffmpegBin, inputPath, outputPath string, width int) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	args := []string{
		"-y",
		"-i", inputPath,
		"-vf", fmt.Sprintf("scale=%d:-1", width),
		outputPath,
	}
	cmd := exec.CommandContext(ctx, ffmpegBin, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) || isExecutableNotFound(err) {
			return fmt.Errorf("ffmpeg not found (%s). Install ffmpeg or set FFMPEG_BIN", ffmpegBin)
		}
		msg := strings.TrimSpace(string(out))
		if msg == "" {
			msg = err.Error()
		}
		return fmt.Errorf("ffmpeg resize failed: %s", msg)
	}
	if _, err := os.Stat(outputPath); err != nil {
		return fmt.Errorf("ffmpeg did not produce resized output")
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

func isExecutableNotFound(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "executable file not found") ||
		strings.Contains(err.Error(), "no such file or directory")
}

func resolveRealESRGANInstallDir(bin string) (string, error) {
	bin = strings.TrimSpace(bin)
	if bin == "" {
		return "", fmt.Errorf("Real-ESRGAN binary path is empty")
	}

	resolved := bin
	if strings.Contains(bin, string(os.PathSeparator)) {
		if abs, err := filepath.Abs(bin); err == nil {
			resolved = abs
		}
	} else if path, err := exec.LookPath(bin); err == nil {
		resolved = path
	}

	if link, err := filepath.EvalSymlinks(resolved); err == nil {
		resolved = link
	}

	dir := filepath.Dir(resolved)
	if dir == "" || dir == "." {
		return "", fmt.Errorf("could not resolve Real-ESRGAN install directory for %q", bin)
	}
	return dir, nil
}

func preferWidthResizeOnly(width, height int) bool {
	if width < 1 || height < 1 {
		return true
	}
	// Motif strips are wide and short; Real-ESRGAN tile padding corrupts them (black blobs, sliced edges).
	if height <= 320 && float64(width)/float64(height) >= 5 {
		return true
	}
	return false
}

func realesrganTileSize(width, height int) int {
	if width < 1 || height < 1 {
		return 0
	}
	// ncnn-vulkan tiles are square; using the long side pads short strips vertically.
	short := height
	if width < short {
		short = width
	}
	const maxTile = 512
	if short > maxTile {
		short = maxTile
	}
	if short < 32 {
		return 32
	}
	return short
}

func trimPNGToAlphaBounds(path string, padding int) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return err
	}
	bounds, ok := findAlphaBounds(img, padding)
	if !ok {
		return nil
	}
	cropped := cropImage(img, bounds)
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	return png.Encode(out, cropped)
}

func findAlphaBounds(img image.Image, padding int) (image.Rectangle, bool) {
	minX := img.Bounds().Max.X
	minY := img.Bounds().Max.Y
	maxX := img.Bounds().Min.X - 1
	maxY := img.Bounds().Min.Y - 1
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a>>8 > 1 {
				if x < minX {
					minX = x
				}
				if y < minY {
					minY = y
				}
				if x > maxX {
					maxX = x
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}
	if maxX < minX || maxY < minY {
		return image.Rectangle{}, false
	}
	if padding > 0 {
		minX -= padding
		minY -= padding
		maxX += padding
		maxY += padding
		if minX < 0 {
			minX = 0
		}
		if minY < 0 {
			minY = 0
		}
	}
	return image.Rect(minX, minY, maxX+1, maxY+1), true
}

func cropImage(img image.Image, bounds image.Rectangle) image.Image {
	if sub, ok := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}); ok {
		return sub.SubImage(bounds)
	}
	out := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			out.Set(x-bounds.Min.X, y-bounds.Min.Y, img.At(x, y))
		}
	}
	return out
}

func maybeRestoreKey(ctx context.Context, ffmpegBin, outputPath, keyColor string) error {
	keyColor = strings.TrimSpace(keyColor)
	if keyColor == "" {
		return nil
	}
	keyed := outputPath + ".keyed.png"
	defer os.Remove(keyed)
	if err := restoreKeyTransparency(ctx, ffmpegBin, outputPath, keyed, keyColor); err != nil {
		return err
	}
	return copyFile(keyed, outputPath)
}

func maybeStripPanoramaSeparators(inputPath string, slideCount, slideWidth, gapPx int) (string, error) {
	strippedPath := filepath.Join(filepath.Dir(inputPath), fmt.Sprintf("upscale-stripped-%d.png", time.Now().UnixNano()))
	stripped, err := StripPanoramaSeparators(inputPath, strippedPath, slideCount, slideWidth, gapPx)
	if err != nil {
		os.Remove(strippedPath)
		return inputPath, err
	}
	if !stripped {
		return inputPath, nil
	}
	return strippedPath, nil
}

func flattenPNG(ctx context.Context, ffmpegBin, inputPath, outputPath, hexColor string) error {
	width, height, err := ImageSize(inputPath)
	if err != nil {
		return err
	}
	color := ffmpegHexColor(hexColor)
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	args := []string{
		"-y",
		"-f", "lavfi", "-i", fmt.Sprintf("color=c=%s:s=%dx%d", color, width, height),
		"-i", inputPath,
		"-filter_complex", "[0][1]overlay=0:0:shortest=1",
		"-frames:v", "1",
		outputPath,
	}
	return runFFmpeg(ctx, ffmpegBin, args, "flatten")
}

func restoreKeyTransparency(ctx context.Context, ffmpegBin, inputPath, outputPath, hexColor string) error {
	color := ffmpegHexColor(hexColor)
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	args := []string{
		"-y",
		"-i", inputPath,
		"-vf", fmt.Sprintf("colorkey=%s:0.12:0.04", color),
		outputPath,
	}
	return runFFmpeg(ctx, ffmpegBin, args, "colorkey")
}

func ffmpegHexColor(hex string) string {
	hex = strings.TrimSpace(hex)
	hex = strings.TrimPrefix(hex, "#")
	if strings.HasPrefix(strings.ToLower(hex), "0x") {
		return strings.ToLower(hex)
	}
	return "0x" + strings.ToLower(hex)
}

func EncodeWebPFromPNG(ctx context.Context, ffmpegBin, cwebpBin, pngPath, webpPath string) error {
	var failures []string

	if BinaryAvailable(cwebpBin) {
		if err := encodeWebPWithCwebp(ctx, cwebpBin, pngPath, webpPath); err == nil {
			return nil
		} else {
			failures = append(failures, "cwebp: "+err.Error())
		}
	}

	if err := encodeWebPWithFFmpeg(ctx, ffmpegBin, pngPath, webpPath); err == nil {
		return nil
	} else {
		failures = append(failures, "ffmpeg: "+err.Error())
	}

	if len(failures) == 0 {
		return fmt.Errorf(
			"WebP encoder not found. Install google/webp for cwebp (brew install webp) or ffmpeg built with libwebp",
		)
	}
	return fmt.Errorf("WebP encode failed (%s). Install webp (brew install webp) for cwebp", strings.Join(failures, "; "))
}

func encodeWebPWithCwebp(ctx context.Context, cwebpBin, pngPath, webpPath string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	args := []string{"-lossless", pngPath, "-o", webpPath}
	cmd := exec.CommandContext(ctx, cwebpBin, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) || isExecutableNotFound(err) {
			return fmt.Errorf("cwebp not found (%s). Install webp (brew install webp) or set CWEBP_BIN", cwebpBin)
		}
		msg := strings.TrimSpace(string(out))
		if msg == "" {
			msg = err.Error()
		}
		return fmt.Errorf("%s", msg)
	}
	if _, err := os.Stat(webpPath); err != nil {
		return fmt.Errorf("cwebp did not produce output")
	}
	return nil
}

func encodeWebPWithFFmpeg(ctx context.Context, ffmpegBin, pngPath, webpPath string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	args := []string{
		"-y",
		"-i", pngPath,
		"-c:v", "libwebp",
		"-lossless", "1",
		webpPath,
	}
	return runFFmpeg(ctx, ffmpegBin, args, "webp")
}

func runFFmpeg(ctx context.Context, ffmpegBin string, args []string, step string) error {
	cmd := exec.CommandContext(ctx, ffmpegBin, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) || isExecutableNotFound(err) {
			return fmt.Errorf("ffmpeg not found (%s). Install ffmpeg or set FFMPEG_BIN", ffmpegBin)
		}
		msg := strings.TrimSpace(string(out))
		if msg == "" {
			msg = err.Error()
		}
		return fmt.Errorf("ffmpeg %s failed: %s", step, msg)
	}
	if len(args) > 0 {
		last := args[len(args)-1]
		if !strings.HasPrefix(last, "-") {
			if _, err := os.Stat(last); err != nil {
				return fmt.Errorf("ffmpeg %s did not produce output", step)
			}
		}
	}
	return nil
}
