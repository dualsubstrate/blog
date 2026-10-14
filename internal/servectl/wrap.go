package servectl

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

// Wrap runs a child command and terminates its process group on exit signals.
func Wrap(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("serve-ctl wrap: missing command")
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	setProcessGroup(cmd)

	if err := cmd.Start(); err != nil {
		return err
	}

	sigCh := make(chan os.Signal, 2)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for range sigCh {
			killProcessGroup(cmd.Process.Pid, syscall.SIGTERM)
		}
	}()

	waitErr := cmd.Wait()
	signal.Stop(sigCh)
	killProcessGroup(cmd.Process.Pid, syscall.SIGTERM)
	time.Sleep(500 * time.Millisecond)
	killProcessGroup(cmd.Process.Pid, syscall.SIGKILL)
	return waitErr
}

func setProcessGroup(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setpgid = true
}

func killProcessGroup(pid int, sig syscall.Signal) {
	if pid <= 0 {
		return
	}
	_ = syscall.Kill(-pid, sig)
}
