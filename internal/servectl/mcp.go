package servectl

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// RunContentPipelinesMCP starts content-pipelines-mcp via go tool wgo.
func RunContentPipelinesMCP(root string, extraArgs []string) error {
	root, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("root: %w", err)
	}
	mcp := filepath.Join(root, "providers", "content-pipelines-mcp")
	cfg := filepath.Join(root, "content-pipelines-mcp.yaml")

	if _, err := os.Stat(mcp); err != nil {
		return fmt.Errorf("providers/content-pipelines-mcp missing; run link-providers")
	}
	if _, err := os.Stat(cfg); err != nil {
		return fmt.Errorf("content-pipelines-mcp.yaml missing")
	}

	args := []string{
		"tool", "wgo",
		"-cd", mcp,
		"-file", ".go", "-file", ".html", "-file", ".css",
		"go", "run", "./cmd/content-pipelines-mcp",
		"--config", cfg,
	}
	args = append(args, extraArgs...)

	cmd := exec.Command("go", args...)
	cmd.Dir = root
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	return cmd.Run()
}
