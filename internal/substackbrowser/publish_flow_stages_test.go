package substackbrowser

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestPublishFlowStageSnapshotJSONDecode(t *testing.T) {
	raw := `{"editorProseMirror":true,"continueButton":true,"publishSettingsVisible":false,"publishModalLikely":true,"dialogCount":2}`
	var s PublishFlowStageSnapshot
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatal(err)
	}
	if !s.EditorProseMirror || !s.ContinueButton || s.DialogCount != 2 {
		t.Fatalf("%+v", s)
	}
}

func TestPublishFlowStageSnapshotJSHasKeys(t *testing.T) {
	js := PublishFlowStageSnapshotJS()
	if !strings.Contains(js, "editorProseMirror") || !strings.Contains(js, "continueButton") {
		head := js
		if len(head) > 120 {
			head = head[:120]
		}
		t.Fatalf("unexpected script head: %s", head)
	}
}

func TestRecoverPublishFlowJSHasEscape(t *testing.T) {
	js := RecoverPublishFlowSyncJS()
	if !strings.Contains(js, "Escape") {
		t.Fatal("expected Escape handling in recovery script")
	}
	if !strings.Contains(js, "document.documentElement") {
		t.Fatal("expected documentElement fallback when body is null during navigation")
	}
	if !strings.Contains(js, "hdrSelectors") {
		t.Fatal("full recovery script should still offer header Publish/Next reopen path")
	}
}

func TestRecoverPublishFlowScheduleRetryJSNoHeaderClicks(t *testing.T) {
	js := RecoverPublishFlowScheduleRetrySyncJS()
	if !strings.Contains(js, "Escape") {
		t.Fatal("expected Escape in schedule-retry recovery")
	}
	if strings.Contains(js, "hdrSelectors") {
		t.Fatal("schedule retry recovery must not click header Publish/Next")
	}
}
