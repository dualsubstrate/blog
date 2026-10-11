package main

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/dualsubstrate/blog/internal/substackpublishstate"
)

func TestPrintPublishQueueEmptyEnglish(t *testing.T) {
	var buf bytes.Buffer
	printPublishQueueEmpty(&buf, "content", substackpublishstate.TargetSubstackEN, substackpublishstate.DefaultMarker)
	out := buf.String()
	for _, want := range []string{
		"✅  English Substack queue: empty (substack-en)",
		"make sb-en section/slug",
		"make sb-list-unpublished",
		"social-published",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in:\n%s", want, out)
		}
	}
}

func TestPrintPublishQueueListEnglish(t *testing.T) {
	var buf bytes.Buffer
	list := []substackpublishstate.BundlePostPath{
		"human-condition/2026-05-01-ego-as-game",
		"social-protocols/2026-06-09-the-prediction-business",
	}
	printPublishQueueList(&buf, list, substackpublishstate.TargetSubstackEN, substackpublishstate.DefaultMarker)
	out := buf.String()
	for _, want := range []string{
		"📋  English Substack queue: 2 bundle(s) waiting (substack-en)",
		"1.  human-condition/2026-05-01-ego-as-game",
		"make sb-en-pick",
		"make sb-en human-condition/2026-05-01-ego-as-game",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in:\n%s", want, out)
		}
	}
}

func TestHandlePublishQueueError(t *testing.T) {
	var buf bytes.Buffer
	err := newPublishQueueEmptyError("content", substackpublishstate.DefaultMarker, substackpublishstate.TargetSubstackEN)
	if !handlePublishQueueError(&buf, err) {
		t.Fatal("expected handlePublishQueueError true")
	}
	if !strings.Contains(buf.String(), "queue: empty") {
		t.Fatalf("got %q", buf.String())
	}
	if handlePublishQueueError(&buf, io.EOF) {
		t.Fatal("expected false for unrelated error")
	}
}
