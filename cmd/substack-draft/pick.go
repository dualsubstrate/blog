package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/mattn/go-isatty"
	"github.com/dualsubstrate/blog/internal/substackbrowser"
	"github.com/dualsubstrate/blog/internal/substackpublishstate"
)

// errPickTryNext means open /dev/tty failed or similar; fall back to another picker mode.
var errPickTryNext = errors.New("pick: try next method")

// pickUnpublishedBundleRel shows unpublished bundles and returns the chosen path
// relative to contentRoot (e.g. "human-condition/2026-05-01-ego-as-game").
func pickUnpublishedBundleRel(contentRoot, marker, publishTarget string) (string, error) {
	list, err := substackpublishstate.ListUnpublished(contentRoot, marker, publishTarget)
	if err != nil {
		return "", err
	}
	if len(list) == 0 {
		return "", newPublishQueueEmptyError(contentRoot, marker, publishTarget)
	}

	if runtime.GOOS != "windows" {
		s, err := runHuhOnDevTTY(list, publishTarget)
		if err == nil {
			return s, nil
		}
		if !errors.Is(err, errPickTryNext) {
			return "", err
		}
	}

	if isatty.IsTerminal(uintptr(os.Stdin.Fd())) && isatty.IsTerminal(uintptr(os.Stdout.Fd())) {
		return pickBundleHuh(list, publishTarget)
	}

	return pickBundlePlain(list, publishTarget)
}

// runHuhOnDevTTY runs the huh form with stdin/stdout/stderr on /dev/tty so the UI works
// when a parent process (e.g. make, IDE task) does not attach a TTY to the child.
func runHuhOnDevTTY(list []substackpublishstate.BundlePostPath, publishTarget string) (string, error) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return "", errPickTryNext
	}
	defer tty.Close()

	oldIn, oldOut, oldErr := os.Stdin, os.Stdout, os.Stderr
	os.Stdin = tty
	os.Stdout = tty
	os.Stderr = tty
	defer func() {
		os.Stdin = oldIn
		os.Stdout = oldOut
		os.Stderr = oldErr
	}()

	return pickBundleHuh(list, publishTarget)
}

func pickBundleHuh(list []substackpublishstate.BundlePostPath, publishTarget string) (string, error) {
	strs := make([]string, len(list))
	for i, p := range list {
		strs[i] = string(p)
	}
	chosen := strs[0]
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose post to open in Substack").
				Description(fmt.Sprintf("Target %q not recorded in marker file yet. Type / to filter.", publishTarget)).
				Options(huh.NewOptions(strs...)...).
				Value(&chosen).
				Height(min(14, max(6, len(strs)+2))),
		),
	)
	if err := runHuhForm(form); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return "", substackbrowser.ErrAbortedBeforePublish
		}
		return "", err
	}
	return chosen, nil
}

func pickBundlePlain(list []substackpublishstate.BundlePostPath, publishTarget string) (string, error) {
	rd := os.Stdin
	out := io.Writer(os.Stderr)
	if runtime.GOOS != "windows" {
		if tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0); err == nil {
			defer tty.Close()
			rd = tty
			out = tty
		}
	}
	reader := bufio.NewReader(rd)
	for {
		fmt.Fprintln(out, "")
		fmt.Fprintf(out, "📋  Bundles not yet recorded for target %q:\n", publishTarget)
		for i, p := range list {
			fmt.Fprintf(out, "   %d  %s\n", i+1, p)
		}
		fmt.Fprintf(out, "🔢  Enter number 1-%d (q to quit): ", len(list))
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		s := strings.TrimSpace(strings.ToLower(line))
		if s == "q" || s == "quit" {
			return "", substackbrowser.ErrAbortedBeforePublish
		}
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 || n > len(list) {
			fmt.Fprintln(out, "Invalid number.")
			continue
		}
		return string(list[n-1]), nil
	}
}
