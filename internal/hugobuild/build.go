package hugobuild

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Options configures a production-like Hugo build.
type Options struct {
	Root string
	// Hugo binary name or path (default hugo or env HUGO).
	Hugo string
	// Preview enables draft/future/expired and preview environment.
	Preview bool
	// BaseURL when non-empty is passed to hugo --baseURL.
	BaseURL string
	Stdout io.Writer
	Stderr io.Writer
}

// Run invokes Hugo and renames carousel.preview files to .html siblings.
func Run(opts Options) error {
	root, err := filepath.Abs(opts.Root)
	if err != nil {
		return fmt.Errorf("root: %w", err)
	}
	hugo := opts.Hugo
	if hugo == "" {
		hugo = os.Getenv("HUGO")
	}
	if hugo == "" {
		hugo = "hugo"
	}
	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = os.Stderr
	}

	preview := opts.Preview
	if !preview && os.Getenv("HUGO_PREVIEW") != "" {
		preview = true
	}
	baseURL := opts.BaseURL
	if baseURL == "" {
		baseURL = os.Getenv("BASE_URL")
	}

	args := []string{"--minify", "--gc"}
	if preview {
		args = append(args, "--environment", "preview", "--buildDrafts", "--buildFuture", "--buildExpired")
	}
	if baseURL != "" {
		args = append(args, "--baseURL", baseURL)
	}

	cmd := exec.Command(hugo, args...)
	cmd.Dir = root
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("hugo: %w", err)
	}

	public := filepath.Join(root, "public")
	return renameCarouselPreviews(public)
}

func renameCarouselPreviews(publicDir string) error {
	var renamed int
	err := filepath.WalkDir(publicDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() != "carousel.preview" {
			return nil
		}
		dest := path + ".html"
		if err := os.Rename(path, dest); err != nil {
			return fmt.Errorf("rename %s: %w", path, err)
		}
		renamed++
		return nil
	})
	if err != nil {
		return err
	}
	if renamed > 0 {
		fmt.Fprintf(os.Stderr, "hugo-build: renamed %d carousel.preview file(s)\n", renamed)
	}
	return nil
}

// PublicDir returns the Hugo output directory under root.
func PublicDir(root string) string {
	return filepath.Join(root, "public")
}

// CleanPublic removes the public output directory.
func CleanPublic(root string) error {
	public := PublicDir(root)
	if err := os.RemoveAll(public); err != nil {
		return fmt.Errorf("remove public: %w", err)
	}
	return nil
}

// LoopbackBaseURL formats a link-check base URL for a port.
func LoopbackBaseURL(port int) string {
	return fmt.Sprintf("http://127.0.0.1:%d/", port)
}

// IsLoopbackBase checks whether a base URL targets local link checks.
func IsLoopbackBase(base string) bool {
	return strings.HasPrefix(base, "http://127.0.0.1:")
}
