package substackbrowser

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestScheduleDebugSnapshotJSIsValidJSONShape(t *testing.T) {
	js := ScheduleDebugSnapshotJS()
	if !strings.Contains(js, "cmdkItems") || !strings.Contains(js, "publishDialogOuterHTML") {
		head := js
		if len(head) > 400 {
			head = head[:400] + "..."
		}
		t.Fatalf("unexpected snapshot script: %s", head)
	}
}

func TestDefaultScheduleDebugSnapshotPathUnderTmp(t *testing.T) {
	p := DefaultScheduleDebugSnapshotPath()
	if got := filepath.ToSlash(filepath.Dir(filepath.Clean(p))); got != "tmp" {
		t.Fatalf("expected parent dir tmp, got %q (full %q)", got, p)
	}
}

func TestWriteScheduleDebugSnapshotOutputShape(t *testing.T) {
	// Document expected wrapper JSON without a browser.
	raw := `{"capturedAtISO":"2026-05-01T12:00:00.000Z","href":"https://example.com","title":"t","activeElement":null,"publishDialogFound":false,"publishDialogOuterHTML":"","cmdkItems":[],"radixPortalsSample":[],"cmdkListsSample":[],"tagInputsInDialog":[]}`
	var inner map[string]json.RawMessage
	if err := json.Unmarshal([]byte(raw), &inner); err != nil {
		t.Fatal(err)
	}
	out := struct {
		WrittenAt             string          `json:"writtenAt"`
		ScheduleFailureReason string          `json:"scheduleFailureReason"`
		Snapshot              json.RawMessage `json:"snapshot"`
	}{
		WrittenAt:             "2026-05-01T12:00:01Z",
		ScheduleFailureReason: "tag create UI not confirmed",
		Snapshot:              json.RawMessage(raw),
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "scheduleFailureReason") {
		t.Fatalf("missing reason in %s", string(b))
	}
}
