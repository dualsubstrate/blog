package servectl

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
)

// Up runs process-compose until interrupted, then cleans up.
func Up(root string) error {
	root, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("root: %w", err)
	}
	sock := filepath.Join(root, "tmp", "process-compose.sock")
	if err := os.MkdirAll(filepath.Dir(sock), 0o755); err != nil {
		return err
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	cmd := exec.Command("process-compose", "up", "-f", "process-compose.yaml", "-U", "-u", sock)
	cmd.Dir = root
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("process-compose: %w", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		_ = Cleanup(root, false, true)
		if err != nil {
			return err
		}
		return nil
	case <-sigCh:
		_ = cmd.Process.Signal(syscall.SIGTERM)
		<-done
		return Cleanup(root, false, false)
	}
}
