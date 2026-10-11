package substackbrowser

import (
	"errors"
	"testing"
)

func TestIsChromeTargetGone(t *testing.T) {
	if isChromeTargetGone(nil) {
		t.Fatal("nil should be false")
	}
	if !isChromeTargetGone(errors.New("Inspected target navigated or closed (-32000)")) {
		t.Fatal("expected navigated message")
	}
	if !isChromeTargetGone(errors.New("rpc: execution context was destroyed")) {
		t.Fatal("expected destroyed context")
	}
	if isChromeTargetGone(errors.New("some other failure")) {
		t.Fatal("unrelated error should be false")
	}
}
