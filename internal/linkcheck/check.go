package linkcheck

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/dualsubstrate/blog/internal/hugobuild"
)

const defaultPort = 8765

// Options configures a muffet crawl against a local Hugo build.
type Options struct {
	Root    string
	Port    int
	Hugo    string
	Stdout  io.Writer
	Stderr  io.Writer
	SkipHugo bool // when true, assumes public/ already exists with correct baseURL
}

// Run builds the site (unless SkipHugo), serves public/ on loopback, and runs muffet.
func Run(ctx context.Context, opts Options) error {
	root, err := filepath.Abs(opts.Root)
	if err != nil {
		return fmt.Errorf("root: %w", err)
	}
	port := opts.Port
	if port == 0 {
		port = defaultPort
	}
	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = os.Stderr
	}

	if !opts.SkipHugo {
		if err := hugobuild.CleanPublic(root); err != nil {
			return err
		}
		dep := filepath.Join(root, "data", "tag-deprecations.toml")
		out := filepath.Join(root, "data", "tag-register.txt")
		tagCmd := exec.CommandContext(ctx, "go", "run", "./cmd/tag-register",
			"-content", "content",
			"-deprecations", dep,
			"-out", out,
		)
		tagCmd.Dir = root
		tagCmd.Stdout = stdout
		tagCmd.Stderr = stderr
		if err := tagCmd.Run(); err != nil {
			return fmt.Errorf("tag-register: %w", err)
		}

		base := hugobuild.LoopbackBaseURL(port)
		if err := hugobuild.Run(hugobuild.Options{
			Root:    root,
			Hugo:    opts.Hugo,
			BaseURL: base,
			Stdout:  stdout,
			Stderr:  stderr,
		}); err != nil {
			return err
		}
	}

	public := hugobuild.PublicDir(root)
	baseURL := hugobuild.LoopbackBaseURL(port)

	srv := &http.Server{
		Addr:              fmt.Sprintf("127.0.0.1:%d", port),
		Handler:           http.FileServer(http.Dir(public)),
		ReadHeaderTimeout: 10 * time.Second,
	}

	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	defer ln.Close()

	go func() {
		_ = srv.Serve(ln)
	}()

	if err := waitReady(ctx, baseURL); err != nil {
		return err
	}

	muffetArgs := []string{
		"tool", "github.com/raviqqe/muffet/v2",
		baseURL,
		"--include=^http://127\\.0\\.0\\.1:",
		"--max-connections=32",
		"--timeout=30",
		"--buffer-size=65536",
		"--exclude=^mailto:",
		"--exclude=^https://(www\\.)?facebook\\.com",
		"--exclude=/apple-touch-icon\\.png$",
		"--exclude=/safari-pinned-tab\\.svg$",
		"--exclude=livereload",
	}
	muffet := exec.CommandContext(ctx, "go", muffetArgs...)
	muffet.Dir = root
	muffet.Stdout = stdout
	muffet.Stderr = stderr
	if err := muffet.Run(); err != nil {
		return fmt.Errorf("muffet: %w", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}

func waitReady(ctx context.Context, baseURL string) error {
	client := &http.Client{Timeout: 2 * time.Second}
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(30 * time.Second)
	}
	for time.Now().Before(deadline) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
		if err != nil {
			return err
		}
		resp, err := client.Do(req)
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 500 {
				return nil
			}
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(1 * time.Second):
		}
	}
	return fmt.Errorf("static server did not become ready at %s", baseURL)
}
