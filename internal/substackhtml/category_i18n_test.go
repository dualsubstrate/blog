package substackhtml

import "testing"

func TestCategorySubstackSectionLabelSpanish(t *testing.T) {
	got := CategorySubstackSectionLabel("Social-Protocols", true)
	want := "Protocolos-sociales"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestCategorySubstackSectionLabelSpanishHumanCondition(t *testing.T) {
	got := CategorySubstackSectionLabel("Human-Condition", true)
	want := "Condición-Humana"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestCategorySubstackSectionLabelEnglishUnchanged(t *testing.T) {
	got := CategorySubstackSectionLabel("Mind-Infrastructure", false)
	if got != "Mind-Infrastructure" {
		t.Fatalf("got %q want Mind-Infrastructure", got)
	}
}

func TestCategorySubstackSectionLabelEmpty(t *testing.T) {
	if got := CategorySubstackSectionLabel("", true); got != "" {
		t.Fatalf("got %q want empty", got)
	}
}
