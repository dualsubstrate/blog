package substackbrowser

import (
	"strings"
	"testing"
)

func TestInsertSubscribeButtonJSEmbedsAsync(t *testing.T) {
	js, err := InsertSubscribeButton()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(js, "Subscribe") || !strings.Contains(js, "async function") {
		head := js
		if len(head) > 120 {
			head = head[:120]
		}
		t.Fatalf("unexpected script: %s", head)
	}
}
