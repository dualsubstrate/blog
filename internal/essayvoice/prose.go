package essayvoice

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// BodyProse reads a Hugo Markdown file, drops front matter, and returns paragraph text.
func BodyProse(ctx context.Context, path string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read post: %w", err)
	}
	_, body, err := splitFrontMatter(raw)
	if err != nil {
		return "", fmt.Errorf("post front matter: %w", err)
	}
	prose, err := markdownToProse(body)
	if err != nil {
		return "", fmt.Errorf("post prose: %w", err)
	}
	return prose, nil
}

func markdownToProse(src []byte) (string, error) {
	doc := goldmark.New().Parser().Parse(text.NewReader(src))
	var b strings.Builder
	var pending strings.Builder
	flush := func() {
		paragraph := strings.TrimSpace(pending.String())
		pending.Reset()
		if paragraph == "" {
			return
		}
		if b.Len() > 0 {
			b.WriteString("\n\n")
		}
		b.WriteString(paragraph)
	}
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			switch n.Kind() {
			case ast.KindParagraph, ast.KindListItem:
				flush()
			}
			return ast.WalkContinue, nil
		}
		switch n.Kind() {
		case ast.KindHeading, ast.KindCodeBlock, ast.KindFencedCodeBlock, ast.KindHTMLBlock:
			return ast.WalkSkipChildren, nil
		}
		if t, ok := n.(*ast.Text); ok {
			pending.Write(t.Segment.Value(src))
			if t.SoftLineBreak() || t.HardLineBreak() {
				pending.WriteByte(' ')
			}
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		return "", err
	}
	flush()
	return strings.TrimSpace(b.String()), nil
}
