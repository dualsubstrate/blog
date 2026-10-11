package socialautopost

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/mattn/go-isatty"
)

// PublishChoice is the user's decision for one bundle in the autopost loop.
type PublishChoice int

const (
	ChoicePublish PublishChoice = iota
	ChoiceTagAsPublished
	ChoiceQuit
)

// ErrAutopostQuit means the user chose to stop; remaining bundles are not processed.
var ErrAutopostQuit = errors.New("socialautopost: user stopped remaining bundles")

// ItemPrompt describes one bundle before publish or skip.
type ItemPrompt struct {
	Index              int // 0-based
	Total              int
	RelUnderContent    string
	PostURL            string
	Network            string // e.g. "LinkedIn", "Facebook"
	WithImage          bool
	DryRun             bool
	IdempotencySkipped bool // true when API idempotency already skipped this item
}

// PromptEnabled reports whether to ask publish or skip per bundle.
// ForceAsk and forceNoAsk override auto-detection (-ask / -no-ask).
func PromptEnabled(forceAsk, forceNoAsk bool) bool {
	if forceNoAsk {
		return false
	}
	if forceAsk || envTruthy("SOCIAL_AUTOPOST_ASK") {
		return true
	}
	if envTruthy("SOCIAL_AUTOPOST_NO_ASK") {
		return false
	}
	if isCI() {
		return false
	}
	return isatty.IsTerminal(uintptr(os.Stdin.Fd()))
}

func envTruthy(key string) bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(key))) {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

func isCI() bool {
	if strings.TrimSpace(os.Getenv("GITHUB_ACTIONS")) == "true" {
		return true
	}
	v := strings.TrimSpace(os.Getenv("CI"))
	return v != "" && !strings.EqualFold(v, "false") && !strings.EqualFold(v, "0")
}

// ConfirmPublishItem asks whether to publish, tag-as-published, or quit for this bundle.
// When IdempotencySkipped is true, returns ChoiceTagAsPublished without prompting.
func ConfirmPublishItem(p ItemPrompt) (PublishChoice, error) {
	if p.IdempotencySkipped {
		return ChoiceTagAsPublished, nil
	}
	if isatty.IsTerminal(uintptr(os.Stdin.Fd())) && isatty.IsTerminal(uintptr(os.Stdout.Fd())) {
		return confirmPublishItemHuh(p)
	}
	return confirmPublishItemPlain(p, os.Stderr, os.Stdin)
}

func confirmPublishItemPlain(p ItemPrompt, w io.Writer, rd io.Reader) (PublishChoice, error) {
	r := bufio.NewReader(rd)
	for {
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, formatItemPromptHeader(p))
		if p.DryRun {
			fmt.Fprintln(w, "  (dry-run: no API calls, no social-published writes)")
		}
		fmt.Fprintln(w, "  p + Enter  publish (API + social-published; dry-run: preview only)")
		fmt.Fprintln(w, "  t + Enter  tag-as-published (file only; dry-run: no file write)")
		fmt.Fprintln(w, "  q + Enter  quit (leave remaining bundles untouched)")
		fmt.Fprint(w, "> ")
		line, err := r.ReadString('\n')
		if err != nil {
			return ChoiceTagAsPublished, err
		}
		switch strings.TrimSpace(strings.ToLower(line)) {
		case "p", "publish":
			return ChoicePublish, nil
		case "q", "quit", "stop", "exit":
			return ChoiceQuit, nil
		case "", "t", "tag", "tag-as-published":
			return ChoiceTagAsPublished, nil
		default:
			fmt.Fprintln(w, "Please type p (publish), t (tag-as-published), or q (quit).")
		}
	}
}

func confirmPublishItemHuh(p ItemPrompt) (PublishChoice, error) {
	var choice string
	desc := "publish: API + social-published. tag-as-published: social-published only."
	if p.DryRun {
		desc = "Dry-run: no API and no social-published writes. publish: preview only. tag-as-published: no file write."
	}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(formatItemPromptHeader(p)).
				Description(desc).
				Options(
					huh.NewOption("publish", "publish"),
					huh.NewOption("tag-as-published", "tag-as-published"),
					huh.NewOption("quit", "quit"),
				).
				Value(&choice),
		),
	)
	if err := runHuhForm(form); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return ChoiceQuit, ErrAutopostQuit
		}
		return ChoiceTagAsPublished, err
	}
	switch choice {
	case "publish":
		return ChoicePublish, nil
	case "quit":
		return ChoiceQuit, nil
	default:
		return ChoiceTagAsPublished, nil
	}
}

func formatItemPromptHeader(p ItemPrompt) string {
	n := p.Index + 1
	img := "text or link"
	if p.WithImage {
		img = "image + caption"
	}
	prefix := ""
	if p.DryRun {
		prefix = "[dry-run] "
	}
	return fmt.Sprintf("%s%s item %d/%d: content/%s (%s)\nURL: %s",
		prefix, p.Network, n, p.Total, p.RelUnderContent, img, p.PostURL)
}

func runHuhForm(form *huh.Form) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	if runtime.GOOS != "windows" {
		ctx, stop = signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	}
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
	if ctx.Err() != nil {
		return huh.ErrUserAborted
	}
	if errors.Is(err, huh.ErrTimeout) {
		return huh.ErrUserAborted
	}
	return err
}
