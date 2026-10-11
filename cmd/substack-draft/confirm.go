package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/mattn/go-isatty"

	"github.com/dualsubstrate/blog/internal/substackbrowser"
)

// runHuhForm runs a huh form with:
//   - optional parent context (browser cancel still tears down the form),
//   - OS signals cancelling the bubbletea program (raw TTY mode often eats SIGINT as a key; this covers the signal path),
//   - on Unix, tea.WithInputTTY so keyboard input works when stdin is not the controlling TTY (e.g. after make).
func runHuhForm(form *huh.Form) error {
	return runHuhFormContext(context.Background(), form)
}

func runHuhFormContext(parent context.Context, form *huh.Form) error {
	ctx, stop := notifyContextInterrupt(parent)
	defer stop()
	f := form
	if runtime.GOOS != "windows" {
		f = f.WithProgramOptions(tea.WithInputTTY())
	}
	err := f.RunWithContext(ctx)
	if err == nil {
		return nil
	}
	if errors.Is(err, huh.ErrUserAborted) {
		return huh.ErrUserAborted
	}
	// Context cancelled (signal) ends tea as killed; huh maps that to ErrTimeout.
	if ctx.Err() != nil {
		return huh.ErrUserAborted
	}
	if errors.Is(err, huh.ErrTimeout) {
		return huh.ErrUserAborted
	}
	return err
}

func notifyContextInterrupt(parent context.Context) (context.Context, context.CancelFunc) {
	if runtime.GOOS == "windows" {
		return signal.NotifyContext(parent, os.Interrupt)
	}
	return signal.NotifyContext(parent, os.Interrupt, syscall.SIGTERM)
}

// confirmScheduleAfterPaste asks after HTML paste whether to run Substack Continue plus publish settings.
// Returns true to run section/tags/schedule automation, false to close without that step.
func confirmScheduleAfterPaste() (bool, error) {
	if isatty.IsTerminal(uintptr(os.Stdin.Fd())) && isatty.IsTerminal(uintptr(os.Stdout.Fd())) {
		return confirmScheduleAfterPasteHuh()
	}
	return confirmScheduleAfterPastePlain()
}

func confirmScheduleAfterPasteHuh() (bool, error) {
	proceed := true
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Body pasted. Continue to publish settings?").
				Description("Continue clicks Substack Continue and fills section, tags, delivery, and schedule from your Markdown. Close leaves the draft as-is in the editor and ends this browser session.").
				Affirmative("Continue").
				Negative("Close").
				Value(&proceed),
		),
	)
	if err := runHuhForm(form); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return false, substackbrowser.ErrAbortedBeforePublish
		}
		return false, err
	}
	if !proceed {
		return false, substackbrowser.ErrAbortedBeforePublish
	}
	return true, nil
}

func confirmScheduleAfterPastePlain() (bool, error) {
	rd := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "⏸️  Body pasted. Run publish settings automation (Continue in Substack, section, tags, schedule)?")
		fmt.Fprintln(os.Stderr, "  c + Enter  continue with automation")
		fmt.Fprintln(os.Stderr, "  q + Enter  close browser, skip automation")
		fmt.Fprint(os.Stderr, "> ")
		line, err := rd.ReadString('\n')
		if err != nil {
			return false, err
		}
		s := strings.TrimSpace(strings.ToLower(line))
		switch s {
		case "", "c", "continue", "y", "yes":
			return true, nil
		case "q", "quit", "close", "n", "no":
			return false, substackbrowser.ErrAbortedBeforePublish
		default:
			fmt.Fprintln(os.Stderr, "Please type c (continue) or q (close).")
		}
	}
}

// confirmDraftReview blocks until the user is done with the automated Chrome session (paste-only flow).
// Uses charmbracelet/huh when stdin and stdout are both TTYs; otherwise plain y/n on stderr (IDE-friendly).
func confirmDraftReview(ctx context.Context) error {
	_ = ctx
	if isatty.IsTerminal(uintptr(os.Stdin.Fd())) && isatty.IsTerminal(uintptr(os.Stdout.Fd())) {
		return confirmDraftReviewHuh()
	}
	return confirmDraftReviewPlain()
}

func confirmDraftReviewHuh() error {
	for {
		done := true
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Done with this automated Chrome session?").
					Description("Continue closes this Chrome window. Your draft remains in Substack. Not yet keeps the window open so you can keep editing in that tab.").
					Affirmative("Continue").
					Negative("Not yet").
					Value(&done),
			),
		)
		if err := runHuhForm(form); err != nil {
			if errors.Is(err, huh.ErrUserAborted) {
				return substackbrowser.ErrAbortedBeforePublish
			}
			return err
		}
		if done {
			return nil
		}
	}
}

func confirmDraftReviewPlain() error {
	rd := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "🪟  Done with the automated Chrome session? Draft is saved in Substack.")
		fmt.Fprintln(os.Stderr, "  y + Enter  close this Chrome window")
		fmt.Fprintln(os.Stderr, "  n + Enter  keep Chrome open longer")
		fmt.Fprint(os.Stderr, "> ")
		line, err := rd.ReadString('\n')
		if err != nil {
			return err
		}
		s := strings.TrimSpace(strings.ToLower(line))
		switch s {
		case "", "y", "yes":
			return nil
		case "n", "no":
			continue
		default:
			fmt.Fprintln(os.Stderr, "Please type y or n and press Enter.")
		}
	}
}

// confirmAfterPublishSettingsReview runs after paste-schedule automation. Affirmative "Publish" clicks
// Substack's purple footer button in the open tab, then ends the session. "Keep open" loops without clicking.
func confirmAfterPublishSettingsReview(ctx context.Context) error {
	if isatty.IsTerminal(uintptr(os.Stdin.Fd())) && isatty.IsTerminal(uintptr(os.Stdout.Fd())) {
		return confirmAfterPublishSettingsReviewHuh(ctx)
	}
	return confirmAfterPublishSettingsReviewPlain(ctx)
}

func confirmAfterPublishSettingsReviewHuh(ctx context.Context) error {
	for {
		publish := true
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Automation finished: review Substack's Publish modal.").
					Description("Publish (default) clicks Substack's purple footer (send or schedule), then Chrome exits. Keep open leaves the window so you can finish in Substack yourself.").
					Affirmative("Publish").
					Negative("Keep open").
					Value(&publish),
			),
		)
		if err := runHuhFormContext(ctx, form); err != nil {
			if errors.Is(err, huh.ErrUserAborted) {
				return substackbrowser.ErrAbortedBeforePublish
			}
			return err
		}
		if publish {
			if err := substackbrowser.ClickPublishModalFooter(ctx); err != nil {
				return fmt.Errorf("substack publish click: %w", err)
			}
			return nil
		}
	}
}

func confirmAfterPublishSettingsReviewPlain(ctx context.Context) error {
	rd := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "📤  Automation finished. Review Substack's Publish modal.")
		fmt.Fprintln(os.Stderr, "  p + Enter  click Publish in Substack (purple footer), then close Chrome (default if you press Enter on an empty line)")
		fmt.Fprintln(os.Stderr, "  k + Enter  keep Chrome open (finish in Substack yourself)")
		fmt.Fprint(os.Stderr, "> ")
		line, err := rd.ReadString('\n')
		if err != nil {
			return err
		}
		s := strings.TrimSpace(strings.ToLower(line))
		switch s {
		case "", "p", "publish", "y", "yes":
			if err := substackbrowser.ClickPublishModalFooter(ctx); err != nil {
				return fmt.Errorf("substack publish click: %w", err)
			}
			return nil
		case "k", "keep", "n", "no":
			continue
		default:
			fmt.Fprintln(os.Stderr, "Please type p (Publish in Substack, default on empty line) or k (keep open).")
		}
	}
}

// confirmMarkPublishedOptional asks whether to append an entry to social-published for targetKey.
// Default is skip (No in huh; empty line in plain mode). Choose Yes only after the post is live for that channel.
func confirmMarkPublishedOptional(bundleDirDisplay string, targetKey string) (bool, error) {
	if isatty.IsTerminal(uintptr(os.Stdin.Fd())) && isatty.IsTerminal(uintptr(os.Stdout.Fd())) {
		return confirmMarkPublishedOptionalHuh(bundleDirDisplay, targetKey)
	}
	return confirmMarkPublishedOptionalPlain(bundleDirDisplay, targetKey)
}

func confirmMarkPublishedOptionalHuh(bundleDirDisplay, targetKey string) (bool, error) {
	mark := false
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Record this post as published in the repo?").
				Description("Yes appends " + targetKey + " and a UTC timestamp to social-published under " + bundleDirDisplay + " so list-unpublished and pick skip this bundle. Choose Yes only after the post is live for that channel. Default: No.").
				Affirmative("Yes, mark published").
				Negative("No, skip").
				Value(&mark),
		),
	)
	if err := runHuhForm(form); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return false, substackbrowser.ErrAbortedBeforePublish
		}
		return false, err
	}
	return mark, nil
}

func confirmMarkPublishedOptionalPlain(bundleDirDisplay, targetKey string) (bool, error) {
	rd := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "📌  Record as published in the repo (append "+targetKey+" to social-published in "+bundleDirDisplay+")?")
		fmt.Fprintln(os.Stderr, "  y + Enter  yes, write marker")
		fmt.Fprintln(os.Stderr, "  n + Enter  skip (default on empty line)")
		fmt.Fprint(os.Stderr, "> ")
		line, err := rd.ReadString('\n')
		if err != nil {
			return false, nil
		}
		s := strings.TrimSpace(strings.ToLower(line))
		switch s {
		case "y", "yes":
			return true, nil
		case "", "n", "no":
			return false, nil
		default:
			fmt.Fprintln(os.Stderr, "Please type y (yes, write marker) or n (skip); empty line skips.")
		}
	}
}
