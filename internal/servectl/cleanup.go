package servectl

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var servePorts = []int{3848, 3849, 28492}

// Cleanup stops process-compose and frees loopback serve ports.
func Cleanup(root string, preflight bool, quiet bool) error {
	root, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("root: %w", err)
	}
	sock := filepath.Join(root, "tmp", "process-compose.sock")

	_ = exec.Command("process-compose", "down", "-U", "-u", sock).Run()
	_ = os.Remove(sock)

	for _, port := range servePorts {
		if err := freePort(port); err != nil {
			return err
		}
	}

	if preflight {
		return nil
	}
	if !quiet {
		fmt.Printf("serve: stopped process-compose and freed ports %v\n", servePorts)
	}
	return nil
}

func freePort(port int) error {
	pids, err := listenersOnPort(port)
	if err != nil || len(pids) == 0 {
		return err
	}
	for _, pid := range pids {
		_ = signalPID(pid, "TERM")
	}
	time.Sleep(300 * time.Millisecond)
	pids, err = listenersOnPort(port)
	if err != nil {
		return err
	}
	for _, pid := range pids {
		_ = signalPID(pid, "KILL")
	}
	return nil
}

func listenersOnPort(port int) ([]int, error) {
	out, err := exec.Command("lsof", "-nP", fmt.Sprintf("-iTCP:%d", port), "-sTCP:LISTEN", "-t").Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return nil, nil
		}
		if _, lookErr := exec.LookPath("lsof"); lookErr != nil {
			return nil, nil
		}
		return nil, fmt.Errorf("lsof port %d: %w", port, err)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var pids []int
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		pid, err := strconv.Atoi(line)
		if err != nil {
			continue
		}
		pids = append(pids, pid)
	}
	return pids, nil
}

func signalPID(pid int, sig string) error {
	return exec.Command("kill", "-"+sig, strconv.Itoa(pid)).Run()
}
