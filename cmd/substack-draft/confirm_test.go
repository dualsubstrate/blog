package main

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/dualsubstrate/blog/internal/substackbrowser"
)

func TestConfirmScheduleAfterPastePlainQuitAborts(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	oldIn := os.Stdin
	os.Stdin = r
	defer func() {
		os.Stdin = oldIn
		_ = w.Close()
		_ = r.Close()
	}()
	go func() {
		_, _ = io.WriteString(w, "q\n")
		_ = w.Close()
	}()

	proceed, err := confirmScheduleAfterPastePlain()
	if proceed {
		t.Fatal("expected proceed false")
	}
	if !errors.Is(err, substackbrowser.ErrAbortedBeforePublish) {
		t.Fatalf("err=%v want ErrAbortedBeforePublish", err)
	}
}

func TestErrAbortedBeforePublishSentinel(t *testing.T) {
	if !errors.Is(substackbrowser.ErrAbortedBeforePublish, substackbrowser.ErrAbortedBeforePublish) {
		t.Fatal("sentinel mismatch")
	}
	_ = strings.TrimSpace("substack-draft: aborted before publish")
}
