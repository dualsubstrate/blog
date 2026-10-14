package canvasbuild

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/dualsubstrate/blog/internal/filetree"
)

const (
	relDestCanvas = "web/dss-hops/src/DSS-core-topology.canvas.tsx"
	relHopsDir    = "web/dss-hops"
	relDist       = "web/dss-hops/dist"
	relStaticOut  = "static/js/dss-hops"
)

// Options configures the DSS hops canvas build.
type Options struct {
	Root string
	// CanvasSource overrides the default Cursor canvas path (env DSS_HOPS_CANVAS).
	CanvasSource string
	// CI when true runs npm ci instead of npm install.
	CI bool
	Stdout io.Writer
	Stderr io.Writer
}

// Run syncs the canvas source, runs npm build, and mirrors dist into static/js/dss-hops.
func Run(opts Options) error {
	root, err := filepath.Abs(opts.Root)
	if err != nil {
		return fmt.Errorf("root: %w", err)
	}
	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = os.Stderr
	}

	dest := filepath.Join(root, relDestCanvas)
	src := opts.CanvasSource
	if src == "" {
		src = os.Getenv("DSS_HOPS_CANVAS")
	}
	if src == "" {
		src = defaultCursorCanvasPath()
	}

	if err := syncCanvas(src, dest, stderr); err != nil {
		return err
	}

	hopsDir := filepath.Join(root, relHopsDir)
	if opts.CI || os.Getenv("CI") != "" {
		if err := runCmd(stdout, stderr, hopsDir, "npm", "ci"); err != nil {
			return err
		}
	} else {
		if err := runCmd(stdout, stderr, hopsDir, "npm", "install"); err != nil {
			return err
		}
	}
	if err := runCmd(stdout, stderr, hopsDir, "npm", "run", "build"); err != nil {
		return err
	}

	staticOut := filepath.Join(root, relStaticOut)
	if err := os.MkdirAll(staticOut, 0o755); err != nil {
		return fmt.Errorf("mkdir static out: %w", err)
	}
	keep := filepath.Join(staticOut, ".gitkeep")
	if err := touchFile(keep); err != nil {
		return err
	}
	dist := filepath.Join(root, relDist)
	if err := filetree.MirrorDir(dist, staticOut); err != nil {
		return fmt.Errorf("mirror dist: %w", err)
	}
	if err := touchFile(keep); err != nil {
		return err
	}
	return nil
}

func defaultCursorCanvasPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".cursor", "projects", "Users-hector-Xynova-dss", "canvases", "DSS-core-topology.canvas.tsx")
}

func syncCanvas(src, dest string, stderr io.Writer) error {
	if src != "" {
		st, err := os.Stat(src)
		if err == nil && !st.IsDir() {
			if err := copyCanvasFile(src, dest); err != nil {
				return err
			}
			fmt.Fprintf(stderr, "canvas-sync: copied %s -> %s\n", src, dest)
			return nil
		}
	}
	st, err := os.Stat(dest)
	if err == nil && !st.IsDir() {
		fmt.Fprintf(stderr, "canvas-sync: using committed canvas at %s\n", dest)
		return nil
	}
	return fmt.Errorf("canvas-sync: missing source canvas %q and no committed file at %s", src, dest)
}

func copyCanvasFile(src, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open src: %w", err)
	}
	defer in.Close()
	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("create dest: %w", err)
	}
	defer func() {
		_ = out.Close()
	}()
	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("copy: %w", err)
	}
	return out.Close()
}

func touchFile(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return fmt.Errorf("touch %s: %w", path, err)
	}
	return f.Close()
}

func runCmd(stdout, stderr io.Writer, dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Env = os.Environ()
	if runtime.GOOS == "windows" {
		cmd.Env = append(cmd.Env, "npm_config_script_shell=cmd")
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s %v in %s: %w", name, args, dir, err)
	}
	return nil
}
