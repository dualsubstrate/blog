package substackbrowser

import (
	"strings"
	"testing"
)

func TestTagOptionCenterJSTargetsHeadlessListbox(t *testing.T) {
	js := tagOptionCenterJS("LetsDefineBad")
	for _, needle := range []string{
		`[role="listbox"] [role="option"]`,
		`li[role="option"]`,
		`findExactTagOption`,
		`findCreateTagOption`,
		`looksLikeCreateTagLine`,
		`hasCommittedTagPill`,
		`LetsDefineBad`,
	} {
		if !strings.Contains(js, needle) {
			t.Fatalf("tagOptionCenterJS missing %q", needle)
		}
	}
}

func TestTagInputCenterJSFindsHeadlessCombobox(t *testing.T) {
	js := tagInputCenterJS()
	for _, needle := range []string{
		`input[role="combobox"]`,
		`headlessui-combobox-input`,
		`inputBox`,
		`hasChips`,
	} {
		if !strings.Contains(js, needle) {
			t.Fatalf("tagInputCenterJS missing %q", needle)
		}
	}
}

func TestTagClearQueryJSClearsNativeValue(t *testing.T) {
	js := tagClearQueryJS()
	if !strings.Contains(js, "HTMLInputElement.prototype") {
		t.Fatal("tagClearQueryJS missing native value setter")
	}
	if !strings.Contains(js, "deleteContentBackward") {
		t.Fatal("tagClearQueryJS missing input clear event")
	}
}

func TestTagCapacityJSChecksReadonly(t *testing.T) {
	js := tagCapacityJS()
	for _, needle := range []string{
		`readOnly`,
		`disabled`,
		`capacity`,
		`tagInputEl`,
	} {
		if !strings.Contains(js, needle) {
			t.Fatalf("tagCapacityJS missing %q", needle)
		}
	}
}
