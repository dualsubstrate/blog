package cliout

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintPublishQueueEmpty(t *testing.T) {
	var buf bytes.Buffer
	PrintPublishQueueEmpty(&buf, "content", "substack-en", "social-published", PublishQueueHints{
		LangLabel: "English Substack",
		ListCmd:   "make sb-list-unpublished",
		PasteCmd:  "make sb-en section/slug",
	})
	out := buf.String()
	if !strings.Contains(out, "✅") || !strings.Contains(out, "queue: empty") {
		t.Fatalf("got %q", out)
	}
}

func TestPrintSubstackDryRun(t *testing.T) {
	var buf bytes.Buffer
	PrintSubstackDryRun(&buf, map[string]string{"url": "https://example.com", "html_bytes": "1234"})
	if !strings.Contains(buf.String(), "🧪") {
		t.Fatalf("got %q", buf.String())
	}
}
