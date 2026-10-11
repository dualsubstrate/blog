package substackbrowser

import (
	"strings"
	"testing"
)

func TestScheduleAfterContinueJSMarshals(t *testing.T) {
	js, err := ScheduleAfterContinueJS(ScheduleAfterContinueOptions{
		Tags:              []string{"A", "B"},
		SectionLabel:      "Mind Infrastructure",
		DateTimeLocal:     "2026-05-01T09:00",
		DateDisplay:       "29/04/2026, 08:40 am",
		TickEmailSubstack: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(js, `"tags"`) || !strings.Contains(js, `"dateDisplay"`) || !strings.Contains(js, `"Mind Infrastructure"`) {
		head := js
		if len(head) > 200 {
			head = head[:200]
		}
		t.Fatalf("unexpected script head: %s", head)
	}
}

func TestScheduleAfterContinueJSUsesNativeTagTyping(t *testing.T) {
	js, err := ScheduleAfterContinueJS(ScheduleAfterContinueOptions{
		Tags:         []string{"RulesForTheRest"},
		SectionLabel: "Human Condition",
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, needle := range []string{
		"function typeTagQuery(",
		"function setNativeInputValue(",
		"function tagQueryLooksApplied(",
		"tag query did not stick in combobox for:",
		"highlightedTagRowMatchesNeedle(",
		"HTMLInputElement.prototype",
		"execCommand('insertText'",
	} {
		if !strings.Contains(js, needle) {
			t.Fatalf("schedule JS missing %q", needle)
		}
	}
	// Blind Enter on an unfiltered list used to commit wrong chips; keep the gate.
	if !strings.Contains(js, "highlightedTagRowMatchesNeedle(t)") {
		t.Fatal("schedule JS missing Enter gate on addTags")
	}
}

func TestScheduleAfterContinueJSCommitsExactTagOption(t *testing.T) {
	js, err := ScheduleAfterContinueJS(ScheduleAfterContinueOptions{
		Tags:         []string{"LetsDefineBad"},
		SectionLabel: "Human Condition",
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, needle := range []string{
		"function pressTagComboboxEnter(",
		"function exactTagOptionVisible(",
		"function highlightedTagListRow(preferTagText)",
		"exactTagOptionVisible(t)",
		"pressTagComboboxEnter(inp)",
		// List checkmark is not a chip: only succeed when hasCommittedTagPill is true.
		"List checkmark is not a chip",
		// aria-selected highlight is not a chip; still click when no pill.
		"aria-selected highlight is not a chip",
	} {
		if !strings.Contains(js, needle) {
			t.Fatalf("schedule JS missing %q", needle)
		}
	}
	// False success: checkmark branch must not return true without a pill.
	bad := "if (listOptionShowsAlreadyApplied(el)) {\n          await sleep(200);\n          if (hasCommittedTagPill(tagText)) return true;\n          return true;\n        }"
	if strings.Contains(js, bad) {
		t.Fatal("schedule JS still returns true on list checkmark without hasCommittedTagPill")
	}
	if strings.Contains(js, "aria-selected') === 'true') {\n          return false;") {
		t.Fatal("schedule JS still skips click when aria-selected without pill")
	}
}
