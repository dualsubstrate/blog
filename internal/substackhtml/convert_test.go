package substackhtml

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStripFrontMatter(t *testing.T) {
	src := []byte("---\ntitle: x\ndraft: true\n---\n\nHello\n")
	body, err := StripFrontMatter(src)
	if err != nil {
		t.Fatal(err)
	}
	if got := string(body); got != "Hello" {
		t.Fatalf("body: %q", got)
	}
}

func TestStripFrontMatterNoFM(t *testing.T) {
	src := []byte("# Hi\n")
	body, err := StripFrontMatter(src)
	if err != nil {
		t.Fatal(err)
	}
	if got := string(body); got != "# Hi" {
		t.Fatalf("body: %q", got)
	}
}

func TestPreprocessYouTubeShortcode(t *testing.T) {
	in := []byte("Before\n{{< youtube dQw4w9WgXcQ >}}\nAfter\n")
	got := string(PreprocessMarkdown(in, ""))
	if !strings.Contains(got, "https://www.youtube.com/watch?v=dQw4w9WgXcQ") {
		t.Fatalf("missing watch URL: %q", got)
	}
	if strings.Contains(got, "{{<") {
		t.Fatalf("shortcode leaked: %q", got)
	}
}

func TestPreprocessMermaidfileShortcode(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "diagram.mmd"), []byte("flowchart LR\n    A --> B"), 0o644); err != nil {
		t.Fatal(err)
	}
	idx := filepath.Join(dir, "index.md")
	if err := os.WriteFile(idx, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	src := []byte("before\n{{< mermaidfile >}}\nafter\n")
	got := string(PreprocessMarkdown(src, idx))
	if !strings.Contains(got, "```mermaid") || !strings.Contains(got, "flowchart LR") {
		t.Fatalf("expected fenced mermaid body: %q", got)
	}
	if strings.Contains(got, "{{<") {
		t.Fatalf("shortcode leaked: %q", got)
	}
}

func TestPreprocessMermaidfilePathTraversalRejected(t *testing.T) {
	dir := t.TempDir()
	secret := filepath.Join(dir, "secret.txt")
	if err := os.WriteFile(secret, []byte("nope"), 0o600); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(dir, "post")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	idx := filepath.Join(sub, "index.md")
	if err := os.WriteFile(idx, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	cases := []string{
		"../secret.txt",
		"../../secret.txt",
		"/etc/passwd",
		`..\secret.txt`,
		"diagram/../../../secret.txt",
	}
	for _, fname := range cases {
		src := []byte("{{< mermaidfile \"" + fname + "\" >}}\n")
		got := string(PreprocessMarkdown(src, idx))
		if !strings.Contains(got, "invalid filename") {
			t.Fatalf("fname %q: expected invalid filename marker, got %q", fname, got)
		}
		if strings.Contains(got, "nope") {
			t.Fatalf("fname %q: leaked secret file content", fname)
		}
	}
}

func TestPreprocessMermaidfileNamedFileInSameDir(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "custom.mmd"), []byte("graph TD\n  X"), 0o644); err != nil {
		t.Fatal(err)
	}
	idx := filepath.Join(dir, "index.md")
	if err := os.WriteFile(idx, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	src := []byte("{{< mermaidfile \"custom.mmd\" >}}\n")
	got := string(PreprocessMarkdown(src, idx))
	if !strings.Contains(got, "graph TD") {
		t.Fatalf("expected inlined mermaid: %q", got)
	}
}

func TestPreprocessMermaidfilePrefersWebpSibling(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "diagram.mmd"), []byte("flowchart LR\n    A --> B"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "diagram.webp"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	idx := filepath.Join(dir, "index.md")
	if err := os.WriteFile(idx, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	src := []byte("before\n{{< mermaidfile >}}\nafter\n")
	got := string(PreprocessMarkdown(src, idx))
	if !strings.Contains(got, "![Diagram](diagram.webp)") {
		t.Fatalf("expected markdown image for webp sibling, got %q", got)
	}
	if strings.Contains(got, "```mermaid") {
		t.Fatalf("did not expect fenced mermaid when webp exists: %q", got)
	}
}

func TestConvertBodyRelativeImageResolvesWithPagePermalink(t *testing.T) {
	src := []byte("---\ntitle: T\n---\n\n![](chart.webp)\n")
	opt := DefaultOptions()
	opt.PagePermalink = "https://blog.dualsubstrate.com/human-condition/test-post/"
	out, err := Convert(src, opt)
	if err != nil {
		t.Fatal(err)
	}
	want := `src="https://blog.dualsubstrate.com/human-condition/test-post/chart.webp"`
	if !strings.Contains(out, want) {
		t.Fatalf("want substring %q in %q", want, out)
	}
}

func TestExtractFrontMatterMetaFromSample(t *testing.T) {
	repo := filepath.Clean(filepath.Join("..", ".."))
	mdPath := filepath.Join(repo, "content", "mind-infrastructure", "2026-04-29-free-energy-principle-hallucination-machine", "index.md")
	raw, err := os.ReadFile(mdPath)
	if err != nil {
		t.Skipf("sample article not present: %v", err)
	}
	meta := ExtractFrontMatterMeta(raw)
	if strings.TrimSpace(meta.Title) == "" {
		t.Fatalf("expected title in front matter meta")
	}
	if strings.TrimSpace(meta.Description) == "" {
		t.Fatalf("expected description in front matter meta")
	}
	if strings.TrimSpace(meta.SoWhat) == "" {
		t.Fatalf("expected sowhat in front matter meta")
	}
	if strings.TrimSpace(meta.YouTubeID) == "" {
		t.Fatalf("expected youtube_id in front matter meta")
	}
	if strings.TrimSpace(meta.Type) == "" {
		t.Fatalf("expected type in front matter meta")
	}
	if strings.TrimSpace(meta.ImageURL) == "" {
		t.Fatalf("expected images[0] in front matter meta")
	}
	if len(meta.Tags) == 0 {
		t.Fatalf("expected tags in front matter meta")
	}
}

func TestExtractFrontMatterMetaListTags(t *testing.T) {
	src := []byte(strings.Join([]string{
		"---",
		"title: x",
		"type: sayings",
		"tags:",
		"  - A",
		"  - B",
		"---",
		"",
		"Body",
		"",
	}, "\n"))
	meta := ExtractFrontMatterMeta(src)
	if len(meta.Tags) != 2 || meta.Tags[0] != "A" || meta.Tags[1] != "B" {
		t.Fatalf("tags: %#v", meta.Tags)
	}
}

func TestExtractFrontMatterMetaInlineListWithComma(t *testing.T) {
	src := []byte(strings.Join([]string{
		"---",
		"title: x",
		"type: sayings",
		`tags: ["tag, with comma", "tag2"]`,
		`categories: ["Cognitive-Memetics", "Sm-Art"]`,
		"---",
		"",
		"Body",
		"",
	}, "\n"))
	meta := ExtractFrontMatterMeta(src)
	if len(meta.Tags) != 2 || meta.Tags[0] != "tag, with comma" || meta.Tags[1] != "tag2" {
		t.Fatalf("tags: %#v", meta.Tags)
	}
	if len(meta.Categories) != 2 || meta.Categories[0] != "Cognitive-Memetics" || meta.Categories[1] != "Sm-Art" {
		t.Fatalf("categories: %#v", meta.Categories)
	}
}

func TestExtractFrontMatterMetaLeadingWhitespace(t *testing.T) {
	src := []byte("\n\n  \t---\n" +
		"title: x\n" +
		"tags: [a]\n" +
		"---\n\nBody\n")
	meta := ExtractFrontMatterMeta(src)
	if meta.Title != "x" {
		t.Fatalf("title: %#v", meta.Title)
	}
	if len(meta.Tags) != 1 || meta.Tags[0] != "a" {
		t.Fatalf("tags: %#v", meta.Tags)
	}
}

func TestConvertVideoLeadUsesCenterParagraphsInBRAndTableListMode(t *testing.T) {
	src := []byte(strings.Join([]string{
		"---",
		"title: T",
		"type: video",
		"youtube_id: dQw4w9WgXcQ",
		"images:",
		"  - https://example.com/thumb.jpg",
		"---",
		"",
		"Body text",
	}, "\n"))
	opt := DefaultOptions()
	opt.IncludeFrontMatterLead = true
	opt.TableMode = TableList
	opt.ParagraphMode = ParagraphBR
	out, err := Convert(src, opt)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `style="text-align: center;"`) {
		t.Fatalf("expected Substack-style centered p: %s", out)
	}
	if strings.Count(out, `<p style="text-align: center;">`) < 2 {
		t.Fatalf("expected thumb and CTA each in a centered p: %s", out)
	}
	if !strings.Contains(out, "<img") || !strings.Contains(out, "youtube.com/watch?v=") {
		t.Fatalf("expected img and watch link in lead: %s", out)
	}
	// Multi-cell markdown tables must still become lists.
	src2 := append(src, []byte("\n\n| a | b |\n| - | - |\n| x | y |\n")...)
	out2, err := Convert(src2, opt)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out2, "<ol>") {
		t.Fatalf("expected markdown table as list: %s", out2)
	}
}

func TestConvertIncludeFrontMatterLeadEmitsVideoIntroBits(t *testing.T) {
	repo := filepath.Clean(filepath.Join("..", ".."))
	mdPath := filepath.Join(repo, "content", "mind-infrastructure", "2026-04-29-free-energy-principle-hallucination-machine", "index.md")
	raw, err := os.ReadFile(mdPath)
	if err != nil {
		t.Skipf("sample article not present: %v", err)
	}
	opt := DefaultOptions()
	opt.IncludeFrontMatterLead = true
	opt.ParagraphMode = ParagraphBR
	opt.ParagraphBreakBRCount = 1
	opt.SourcePath = mdPath
	out, err := Convert(raw, opt)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "What you probably do not know yet") {
		t.Fatalf("expected video intro heading: %s", out)
	}
	if !strings.Contains(out, "<img") {
		t.Fatalf("expected lead image: %s", out)
	}
	if !strings.Contains(out, "youtube.com/watch?v=") {
		t.Fatalf("expected youtube link: %s", out)
	}
}

func TestConvertIncludeFrontMatterLeadSpanishVideoHeadings(t *testing.T) {
	repo := filepath.Clean(filepath.Join("..", ".."))
	mdPath := filepath.Join(repo, "content", "mind-infrastructure", "2026-04-29-free-energy-principle-hallucination-machine", "index.es.md")
	raw, err := os.ReadFile(mdPath)
	if err != nil {
		t.Skipf("sample article not present: %v", err)
	}
	opt := DefaultOptions()
	opt.IncludeFrontMatterLead = true
	opt.ParagraphMode = ParagraphBR
	opt.ParagraphBreakBRCount = 1
	opt.SourcePath = mdPath
	out, err := Convert(raw, opt)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "Lo que probablemente aún no sabes") {
		t.Fatalf("expected Spanish video intro heading: %s", out)
	}
	if !strings.Contains(out, "Lo que sabrás después") {
		t.Fatalf("expected Spanish sowhat heading: %s", out)
	}
	if !strings.Contains(out, "Ver en ") || !strings.Contains(out, "YouTube") {
		t.Fatalf("expected Spanish watch CTA: %s", out)
	}
}

func TestConvertVideoLeadSpanishFromLangFrontMatter(t *testing.T) {
	src := []byte(strings.Join([]string{
		"---",
		"title: T",
		"type: video",
		"lang: es",
		"description: |",
		"  - Uno",
		"sowhat: |",
		"  Dos",
		"---",
		"",
		"Body",
	}, "\n"))
	opt := DefaultOptions()
	opt.IncludeFrontMatterLead = true
	out, err := Convert(src, opt)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "Lo que probablemente aún no sabes") {
		t.Fatalf("expected Spanish heading from lang: %s", out)
	}
}

func TestConvertClaimsLeadThoughtsHeadingAfterFeaturedImageHTML(t *testing.T) {
	// Regression: one newline between </p> from featured_image and "## 💭 Thoughts"
	// lets goldmark treat "##" as paragraph text, not an ATX heading (CommonMark HTML block).
	src := []byte(strings.Join([]string{
		"---",
		"title: T",
		"type: claims",
		"description: |",
		"  Short claim.",
		"images:",
		"  - pic.png",
		"---",
		"",
		"### Body heading",
		"",
		"Body text.",
	}, "\n"))
	opt := DefaultOptions()
	opt.IncludeFrontMatterLead = true
	opt.PagePermalink = "https://blog.dualsubstrate.com/human-condition/slug/"
	out, err := Convert(src, opt)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "<p>## ") {
		t.Fatalf("ATX heading leaked as paragraph after featured HTML: %s", out)
	}
	if !strings.Contains(out, "Thoughts") {
		t.Fatalf("expected Thoughts section in output: %s", out)
	}
}

func TestConvertSayingsLeadRendersTLDRAndContextAsHeadings(t *testing.T) {
	// Regression: sayings_prefix used {{- end -}} after {{.TLDR}}, which trims the
	// blank lines before ## 💬 Context so goldmark does not see ATX headings at line
	// start (TLDR and Context glue together or show as literal ## in pasted output).
	src := []byte(strings.Join([]string{
		"---",
		"title: T",
		"type: sayings",
		"tldr: |",
		"  First line of tldr.",
		"fluff: |",
		"  Context paragraph here.",
		"---",
		"",
		"## Essay body",
		"",
		"Hello.",
	}, "\n"))
	opt := DefaultOptions()
	opt.IncludeFrontMatterLead = true
	out, err := Convert(src, opt)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "First line of tldr.##") {
		t.Fatalf("TLDR text glued to Context heading marker: %s", out)
	}
	if strings.Contains(out, "<p>## ") {
		t.Fatalf("ATX heading leaked as paragraph text: %s", out)
	}
	if !strings.Contains(out, "TLDR") || !strings.Contains(out, "Context") {
		t.Fatalf("expected TLDR and Context labels in output: %s", out)
	}
}

func TestConvertCowsIncludesFeaturedImageWithPermalink(t *testing.T) {
	src := []byte(strings.Join([]string{
		"---",
		"title: T",
		"type: panel",
		"images:",
		"  - cow.png",
		"description: Teaser **here**.",
		"---",
		"",
		"Body line.",
	}, "\n"))
	opt := DefaultOptions()
	opt.IncludeFrontMatterLead = true
	opt.PagePermalink = "https://blog.dualsubstrate.com/cognitive-memetics/cows/slug/"
	out, err := Convert(src, opt)
	if err != nil {
		t.Fatal(err)
	}
	wantSrc := `src="https://blog.dualsubstrate.com/cognitive-memetics/cows/slug/cow.png"`
	if !strings.Contains(out, wantSrc) {
		t.Fatalf("expected absolute img in output, want substring %q: %s", wantSrc, out)
	}
	if !strings.Contains(out, "Teaser") {
		t.Fatalf("expected teaser heading: %s", out)
	}
}

func TestIsSpanishSiteLocale(t *testing.T) {
	if !IsSpanishSiteLocale(FrontMatterMeta{Lang: "es"}, "") {
		t.Fatal("lang es should imply Spanish site locale")
	}
	if !IsSpanishSiteLocale(FrontMatterMeta{}, "/tmp/foo/index.es.md") {
		t.Fatal(".es.md basename should imply Spanish when lang is unset")
	}
	if IsSpanishSiteLocale(FrontMatterMeta{}, "/tmp/foo/index.md") {
		t.Fatal("index.md should not imply Spanish")
	}
	if IsSpanishSiteLocale(FrontMatterMeta{}, "") {
		t.Fatal("empty path and no lang should be English")
	}
}

func TestConvertSpaceBeforeStrongBRModePullquote(t *testing.T) {
	src := []byte("---\nx: 1\n---\n\nEs como si los sentidos **sacudieran** al cerebro.\n\n> quote\n\nEsto: la **señal** fin.\n")
	opt := DefaultOptions()
	opt.ParagraphMode = ParagraphBR
	opt.ParagraphBreakBRCount = 1
	opt.QuoteMode = QuotePullquoteMonospace
	opt.DemoteHeadings = false
	out, err := Convert(src, opt)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "sentidos<strong>") || strings.Contains(out, "la<strong>") {
		t.Fatalf("expected a space before <strong>, got fragment: %s", out)
	}
}

func TestConvertSoftNewlinesDoNotBecomeParagraphBreaksBRMode(t *testing.T) {
	src := []byte("---\nx: 1\n---\n\nFirst line\nsecond line\n\nThird paragraph.\n")
	opt := DefaultOptions()
	opt.ParagraphMode = ParagraphBR
	opt.ParagraphBreakBRCount = 2
	out, err := Convert(src, opt)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "First line<br><br>second line") {
		t.Fatalf("expected soft newline to be space, got: %q", out)
	}
	if !strings.Contains(out, "First line second line") {
		t.Fatalf("expected soft newline to become space, got: %q", out)
	}
	// Blank line should still create a paragraph break.
	if !strings.Contains(out, "second line<br><br>Third paragraph") {
		t.Fatalf("expected blank line to become paragraph break, got: %q", out)
	}
}

func TestConvertBRModeDoesNotInsertParagraphBreaksForStrippedInlineElements(t *testing.T) {
	src := []byte("---\nx: 1\n---\n\nword1 <span>word2</span> word3\n")
	opt := DefaultOptions()
	opt.ParagraphMode = ParagraphBR
	out, err := Convert(src, opt)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "word1<br><br>word2") || strings.Contains(out, "word2<br><br>word3") {
		t.Fatalf("expected no paragraph breaks between inline words, got: %q", out)
	}
	if !strings.Contains(out, "word1 word2 word3") {
		t.Fatalf("expected inline text preserved, got: %q", out)
	}
}

func TestMarkdownToHTMLSpaceBeforeStrong(t *testing.T) {
	h, err := markdownToHTML([]byte("Es como si los sentidos **sacudieran** al.\n"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(h, "sentidos<strong>") {
		t.Fatalf("goldmark dropped space: %q", h)
	}
}

func TestConvertSmoke(t *testing.T) {
	src := []byte("---\na: 1\n---\n\n## Title\n\n**Bold** and *italic*.\n\n- one\n- two\n")
	out, err := Convert(src, DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "<h3>") {
		t.Fatalf("expected demoted h2 to h3, got: %s", out)
	}
	if !strings.Contains(out, "<strong>") || !strings.Contains(out, "<em>") {
		t.Fatalf("expected emphasis: %s", out)
	}
	if !strings.Contains(out, "<ul>") {
		t.Fatalf("expected list: %s", out)
	}
}

func TestConvertTableListMode(t *testing.T) {
	src := []byte("---\nx: 1\n---\n\n| a | b |\n| - | - |\n| [0:00](https://example.com) | **Mask** rest |\n")
	opt := DefaultOptions()
	opt.TableMode = TableList
	out, err := Convert(src, opt)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "<ol>") || !strings.Contains(out, "<li>") {
		t.Fatalf("expected list replacement: %s", out)
	}
	if strings.Contains(out, "<table") {
		t.Fatalf("did not expect table: %s", out)
	}
	if !strings.Contains(out, `<strong><a href="https://example.com">0:00</a></strong>`) {
		t.Fatalf("expected link preserved: %s", out)
	}
	if !strings.Contains(out, "<strong>Mask</strong>") {
		t.Fatalf("expected strong preserved: %s", out)
	}
}

func TestConvertParagraphBRMode(t *testing.T) {
	src := []byte("---\nx: 1\n---\n\nFirst paragraph.\n\nSecond paragraph.\n")
	opt := DefaultOptions()
	opt.ParagraphMode = ParagraphBR
	out, err := Convert(src, opt)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "<p>") {
		t.Fatalf("did not expect <p> in br mode: %s", out)
	}
	if !strings.Contains(out, "<br>") {
		t.Fatalf("expected <br> in br mode: %q", out)
	}
}

func TestConvertBlockquoteBRModeCollapsesBreaks(t *testing.T) {
	src := []byte("---\nx: 1\n---\n\n> Quote line\n\nAfter.\n")
	opt := DefaultOptions()
	opt.ParagraphMode = ParagraphBR
	out, err := Convert(src, opt)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "<blockquote>") {
		t.Fatalf("expected blockquote: %q", out)
	}
	// No trailing <br><br> inside the blockquote.
	if strings.Contains(out, "<blockquote>Quote line<br><br>") || strings.Contains(out, "<blockquote>\nQuote line<br><br>") {
		t.Fatalf("expected collapsed breaks inside blockquote: %q", out)
	}
}

func TestConvertQuoteMonospaceReplacesBlockquote(t *testing.T) {
	src := []byte("---\nx: 1\n---\n\n> Quote line\n\nAfter.\n")
	opt := DefaultOptions()
	opt.QuoteMode = QuoteMonospace
	out, err := Convert(src, opt)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "<blockquote>") {
		t.Fatalf("did not expect blockquote in monospace mode: %q", out)
	}
	if !strings.Contains(out, "<pre><code>") {
		t.Fatalf("expected pre/code in monospace mode: %q", out)
	}
	if !strings.Contains(out, "Quote line") {
		t.Fatalf("expected quote content in monospace mode: %q", out)
	}
}

func TestConvertQuotePullquoteMonospaceKeepsBlockquoteAndCode(t *testing.T) {
	src := []byte("---\nx: 1\n---\n\n> Quote line\n\nAfter.\n")
	opt := DefaultOptions()
	opt.QuoteMode = QuotePullquoteMonospace
	out, err := Convert(src, opt)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "<blockquote>") {
		t.Fatalf("expected blockquote in pullquote monospace mode: %q", out)
	}
	if strings.Contains(out, "<pre><code>") {
		t.Fatalf("did not expect pre/code wrapper in pullquote monospace mode: %q", out)
	}
	if !strings.Contains(out, "<blockquote><code>") {
		t.Fatalf("expected blockquote/code combo in pullquote monospace mode: %q", out)
	}
	if !strings.Contains(out, "Quote line") {
		t.Fatalf("expected quote content in pullquote monospace mode: %q", out)
	}
}

func TestSampleArticleGolden(t *testing.T) {
	repo := filepath.Clean(filepath.Join("..", ".."))
	mdPath := filepath.Join(repo, "content", "mind-infrastructure", "2026-04-29-free-energy-principle-hallucination-machine", "index.md")
	raw, err := os.ReadFile(mdPath)
	if err != nil {
		t.Skipf("sample article not present: %v", err)
	}
	out, err := Convert(raw, DefaultOptions())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "<blockquote>") {
		t.Fatalf("expected blockquotes in sample: %s", out)
	}
	if !strings.Contains(out, "<table>") {
		t.Fatalf("expected default table html: %s", out)
	}
	if strings.Contains(out, "translationKey") {
		t.Fatalf("front matter leaked: %s", out)
	}
}
