package substackhtml

import (
	"strings"
	"testing"
)

func TestBuildVideoLeadMarkdownSeparatesHeadingsFromLists(t *testing.T) {
	meta := FrontMatterMeta{
		Type:        "video",
		Description: "- **A** one\n- **B** two\n",
		SoWhat:      "Culture opens the paragraph.\n",
	}
	opt := DefaultOptions()
	lead, err := buildLeadMarkdown(meta, opt)
	if err != nil {
		t.Fatal(err)
	}
	s := string(lead)
	if strings.Contains(s, "not know yet- **") {
		t.Fatalf("first ATX heading glued to list marker: %q", s)
	}
	if strings.Contains(s, "know afterCulture") {
		t.Fatalf("second heading glued to sowhat (trim marker regression): %q", s)
	}
	if !strings.Contains(s, "## 🎯 What you will know after\n") {
		t.Fatalf("expected heading line break before sowhat body: %q", s)
	}
}
