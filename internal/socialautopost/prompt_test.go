package socialautopost

import (
	"bytes"
	"strings"
	"testing"
)

func TestPromptEnabled(t *testing.T) {
	t.Setenv("CI", "")
	t.Setenv("GITHUB_ACTIONS", "")
	t.Setenv("SOCIAL_AUTOPOST_ASK", "")
	t.Setenv("SOCIAL_AUTOPOST_NO_ASK", "")

	if PromptEnabled(false, true) {
		t.Fatal("forceNoAsk should disable")
	}
	if !PromptEnabled(true, false) {
		t.Fatal("forceAsk should enable")
	}

	t.Setenv("SOCIAL_AUTOPOST_ASK", "1")
	if !PromptEnabled(false, false) {
		t.Fatal("SOCIAL_AUTOPOST_ASK should enable")
	}
	t.Setenv("SOCIAL_AUTOPOST_ASK", "")
	t.Setenv("SOCIAL_AUTOPOST_NO_ASK", "yes")
	if PromptEnabled(false, false) {
		t.Fatal("SOCIAL_AUTOPOST_NO_ASK should disable")
	}
}

func TestIsCI(t *testing.T) {
	t.Setenv("CI", "")
	t.Setenv("GITHUB_ACTIONS", "")
	if isCI() {
		t.Fatal("expected false with empty env")
	}
	t.Setenv("GITHUB_ACTIONS", "true")
	if !isCI() {
		t.Fatal("GITHUB_ACTIONS=true")
	}
	t.Setenv("GITHUB_ACTIONS", "")
	t.Setenv("CI", "true")
	if !isCI() {
		t.Fatal("CI=true")
	}
}

func TestConfirmPublishItemPlain(t *testing.T) {
	p := ItemPrompt{
		Index:           0,
		Total:           2,
		RelUnderContent: "cognitive-memetics/foo",
		PostURL:         "https://blog.dualsubstrate.com/foo/",
		Network:         "LinkedIn",
		WithImage:       true,
	}

	var out bytes.Buffer
	choice, err := confirmPublishItemPlain(p, &out, strings.NewReader("p\n"))
	if err != nil || choice != ChoicePublish {
		t.Fatalf("publish: choice=%v err=%v", choice, err)
	}

	choice, err = confirmPublishItemPlain(p, &out, strings.NewReader("t\n"))
	if err != nil || choice != ChoiceTagAsPublished {
		t.Fatalf("tag: choice=%v err=%v", choice, err)
	}

	choice, err = confirmPublishItemPlain(p, &out, strings.NewReader("q\n"))
	if err != nil || choice != ChoiceQuit {
		t.Fatalf("quit: choice=%v err=%v", choice, err)
	}
}

func TestConfirmPublishItem_idempotencySkipped(t *testing.T) {
	choice, err := ConfirmPublishItem(ItemPrompt{IdempotencySkipped: true})
	if err != nil || choice != ChoiceTagAsPublished {
		t.Fatalf("got choice=%v err=%v", choice, err)
	}
}
