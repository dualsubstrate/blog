package substackbrowser

import (
	"strings"
	"testing"
)

func TestPasteExpressionEmbedsJSONString(t *testing.T) {
	html := `<p>He said "hi" & 'bye'</p>`
	expr, err := pasteExpression(html)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(expr, `const html = `) || !strings.Contains(expr, `\u003cp\u003e`) {
		t.Fatalf("expected escaped html payload in script: %s", expr)
	}
}

func TestParsePasteResult(t *testing.T) {
	r, err := ParsePasteResult(`{"ok":true,"reason":""}`)
	if err != nil {
		t.Fatal(err)
	}
	if !r.OK || r.Reason != "" {
		t.Fatalf("unexpected: %+v", r)
	}
	r2, err := ParsePasteResult(`{"ok":true,"reason":"note"}`)
	if err != nil {
		t.Fatal(err)
	}
	if !r2.OK || r2.Reason != "note" {
		t.Fatalf("unexpected: %+v", r2)
	}
}
