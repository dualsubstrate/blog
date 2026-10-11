package substackhtml

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"strings"

	xhtml "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func normalizeHTML(fragment string, opt Options) (string, error) {
	roots, err := xhtml.ParseFragment(strings.NewReader(fragment), &xhtml.Node{
		Type:     xhtml.ElementNode,
		Data:     "div",
		DataAtom: atom.Div,
	})
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	for _, n := range roots {
		n = maybeReplaceTables(n, opt.TableMode)
		nodes, err := sanitizeTree(n, opt)
		if err != nil {
			return "", err
		}
		nodes = maybeFlattenParagraphsToBR(nodes, opt.ParagraphMode, opt.ParagraphBreakBRCount)
		nodes = normalizeBlockquotesForBR(nodes, opt.ParagraphMode, opt.ParagraphBreakBRCount)
		if opt.ParagraphMode == ParagraphBR {
			nodes = collapseBRRunsTopLevel(nodes, opt.ParagraphBreakBRCount)
		}
		for _, nn := range nodes {
			if err := renderNode(&buf, nn); err != nil {
				return "", err
			}
		}
	}
	s := strings.TrimSpace(buf.String())
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return s, nil
}

func normalizeBlockquotesForBR(nodes []*xhtml.Node, mode ParagraphMode, paraBRCount int) []*xhtml.Node {
	if mode != ParagraphBR {
		return nodes
	}
	for _, n := range nodes {
		normalizeBlockquoteNode(n, paraBRCount)
	}
	return nodes
}

func normalizeBlockquoteNode(n *xhtml.Node, paraBRCount int) {
	if n == nil {
		return
	}
	if n.Type == xhtml.ElementNode && n.Data == "blockquote" {
		// Collapse <br> runs inside blockquotes to at most 1, and drop leading/trailing breaks.
		var kids []*xhtml.Node
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			kids = append(kids, c)
		}
		kids = trimLeadingTrailingBR(kids)
		kids = collapseBRRuns(kids, 1)
		replaceChildren(n, kids)
		return
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		normalizeBlockquoteNode(c, paraBRCount)
	}
}

func collapseBRRunsTopLevel(nodes []*xhtml.Node, maxRun int) []*xhtml.Node {
	if maxRun <= 0 {
		maxRun = 1
	}
	return collapseBRRuns(nodes, maxRun)
}

func trimLeadingTrailingBR(nodes []*xhtml.Node) []*xhtml.Node {
	start := 0
	for start < len(nodes) && isBR(nodes[start]) {
		start++
	}
	end := len(nodes)
	for end > start && isBR(nodes[end-1]) {
		end--
	}
	return nodes[start:end]
}

func collapseBRRuns(nodes []*xhtml.Node, maxRun int) []*xhtml.Node {
	if maxRun <= 0 {
		maxRun = 1
	}
	var out []*xhtml.Node
	run := 0
	for _, n := range nodes {
		if isBR(n) {
			run++
			if run <= maxRun {
				out = append(out, n)
			}
			continue
		}
		run = 0
		out = append(out, n)
	}
	return out
}

func isBR(n *xhtml.Node) bool {
	return n != nil && n.Type == xhtml.ElementNode && n.Data == "br"
}

func replaceChildren(parent *xhtml.Node, kids []*xhtml.Node) {
	parent.FirstChild = nil
	parent.LastChild = nil
	for i, k := range kids {
		k.Parent = parent
		k.PrevSibling = nil
		k.NextSibling = nil
		if i == 0 {
			parent.FirstChild = k
			parent.LastChild = k
			continue
		}
		parent.LastChild.NextSibling = k
		k.PrevSibling = parent.LastChild
		parent.LastChild = k
	}
}

// normalizeParagraphTextAlignCenterStyle allows a single safe declaration so pasted HTML
// matches Substack's own centering pattern.
func normalizeParagraphTextAlignCenterStyle(style string) (normalized string, ok bool) {
	parts := strings.Split(style, ";")
	var nonEmpty []string
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			nonEmpty = append(nonEmpty, t)
		}
	}
	if len(nonEmpty) != 1 {
		return "", false
	}
	one := strings.ToLower(strings.ReplaceAll(nonEmpty[0], " ", ""))
	if one == "text-align:center" {
		return "text-align: center;", true
	}
	return "", false
}

func paragraphIsCenterPreservedInBRMode(n *xhtml.Node) bool {
	if n == nil || n.Type != xhtml.ElementNode || n.Data != "p" {
		return false
	}
	for _, a := range n.Attr {
		if !strings.EqualFold(strings.TrimSpace(a.Key), "style") {
			continue
		}
		_, ok := normalizeParagraphTextAlignCenterStyle(a.Val)
		return ok
	}
	return false
}

func maybeFlattenParagraphsToBR(nodes []*xhtml.Node, mode ParagraphMode, paraBRCount int) []*xhtml.Node {
	if mode != ParagraphBR {
		return nodes
	}
	if paraBRCount <= 0 {
		paraBRCount = 2
	}
	var out []*xhtml.Node
	for _, n := range nodes {
		if n == nil {
			continue
		}
		if n.Type == xhtml.TextNode && strings.Contains(n.Data, "\n") {
			// In Markdown, a single newline inside a paragraph is usually a soft wrap, not a paragraph break.
			// Treat blank lines as paragraph breaks, and treat single newlines as spaces.
			out = append(out, normalizeBRModeTextNode(n.Data, paraBRCount)...)
			continue
		}
		if n.Type == xhtml.TextNode && strings.TrimSpace(n.Data) != "" {
			// Adjacent top-level text nodes are often created when inline elements are stripped during
			// sanitization (e.g. word1 <span>word2</span> word3). Treat these as inline continuation,
			// not as paragraph breaks.
			if len(out) > 0 && out[len(out)-1].Type == xhtml.TextNode {
				out[len(out)-1].Data = joinInlineText(out[len(out)-1].Data, n.Data)
			} else {
				out = append(out, n)
			}
			continue
		}
		if n.Type == xhtml.ElementNode && n.Data == "p" {
			if paragraphIsCenterPreservedInBRMode(n) {
				out = append(out, n)
				continue
			}
			kids := nonEmptyInlineChildren(n)
			if len(kids) == 0 {
				continue
			}
			out = append(out, kids...)
			// paragraph break
			out = append(out, repeatBR(paraBRCount)...)
			continue
		}
		out = append(out, n)
	}
	return out
}

func joinInlineText(a, b string) string {
	a = strings.Trim(a, "\r\n")
	b = strings.Trim(b, "\r\n")
	if a == "" {
		return b
	}
	if b == "" {
		return a
	}
	aEnd := a[len(a)-1]
	bStart := b[0]
	aSpace := aEnd == ' ' || aEnd == '\t' || aEnd == '\n'
	bSpace := bStart == ' ' || bStart == '\t' || bStart == '\n'
	if aSpace || bSpace {
		return a + b
	}
	return a + " " + b
}

// normalizeBRModeTextNode converts a text node that contains newlines into a slice of nodes:
// - blank lines become paragraph breaks (represented as paraBRCount <br> tags)
// - single newlines become spaces (Markdown soft line breaks)
//
// It trims newline characters at segment boundaries, but preserves spaces so we do not glue
// words to inline markup when the next sibling begins with <strong>, etc.
func normalizeBRModeTextNode(s string, paraBRCount int) []*xhtml.Node {
	if paraBRCount <= 0 {
		paraBRCount = 2
	}
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")

	var out []*xhtml.Node
	var buf strings.Builder
	needBreak := false

	flush := func() {
		if buf.Len() == 0 {
			needBreak = false
			return
		}
		t := strings.Trim(buf.String(), "\n")
		if strings.TrimSpace(t) == "" {
			buf.Reset()
			needBreak = false
			return
		}
		if needBreak && len(out) > 0 {
			out = append(out, repeatBR(paraBRCount)...)
		}
		out = append(out, &xhtml.Node{Type: xhtml.TextNode, Data: t})
		buf.Reset()
		needBreak = false
	}

	i := 0
	for i < len(s) {
		ch := s[i]
		if ch != '\n' {
			buf.WriteByte(ch)
			i++
			continue
		}

		// Count consecutive newlines, allowing spaces/tabs between them (blank lines).
		j := i
		newlines := 0
		for j < len(s) {
			if s[j] == '\n' {
				newlines++
				j++
				continue
			}
			if s[j] == ' ' || s[j] == '\t' {
				j++
				continue
			}
			break
		}

		if newlines >= 2 {
			flush()
			needBreak = true
			i = j
			continue
		}

		// Single newline: treat as space if it joins two non-space runs.
		prevIsSpace := buf.Len() == 0
		if buf.Len() > 0 {
			b := buf.String()
			prev := b[len(b)-1]
			prevIsSpace = prev == ' ' || prev == '\t' || prev == '\n'
		}
		nextIsSpace := false
		if j < len(s) {
			nxt := s[j]
			nextIsSpace = nxt == ' ' || nxt == '\t' || nxt == '\n'
		}
		if !prevIsSpace && !nextIsSpace {
			buf.WriteByte(' ')
		}
		i = j
	}
	flush()
	return out
}

func repeatBR(n int) []*xhtml.Node {
	if n <= 0 {
		n = 1
	}
	out := make([]*xhtml.Node, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, &xhtml.Node{Type: xhtml.ElementNode, Data: "br", DataAtom: atom.Br})
	}
	return out
}

func nonEmptyInlineChildren(p *xhtml.Node) []*xhtml.Node {
	var out []*xhtml.Node
	for c := p.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == xhtml.TextNode {
			// In BR paragraph mode we flatten <p> into inline nodes. Preserve Markdown semantics:
			// soft-wrapped newlines inside a paragraph become spaces, not paragraph breaks.
			if strings.Contains(c.Data, "\n") || strings.Contains(c.Data, "\r") {
				c.Data = normalizeInlineSoftNewlines(c.Data)
			}
			if strings.TrimSpace(c.Data) == "" {
				continue
			}
		}
		out = append(out, c)
	}
	return out
}

func normalizeInlineSoftNewlines(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	// Replace remaining newlines with spaces.
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}

func maybeReplaceTables(n *xhtml.Node, mode TableMode) *xhtml.Node {
	if mode != TableList {
		return n
	}
	return replaceTablesWithLists(n)
}

// isSingleCellLayoutTable is true when t has exactly one row and one cell (td or th).
// We keep those tables in TableList mode so paste-friendly centering (td align="center") is not flattened to a list.
func isSingleCellLayoutTable(t *xhtml.Node) bool {
	if t == nil || t.Type != xhtml.ElementNode || t.Data != "table" {
		return false
	}
	trs := collectTRs(t)
	if len(trs) != 1 {
		return false
	}
	nCells := 0
	for c := trs[0].FirstChild; c != nil; c = c.NextSibling {
		if c.Type != xhtml.ElementNode {
			continue
		}
		if c.Data == "td" || c.Data == "th" {
			nCells++
		}
	}
	return nCells == 1
}

func replaceTablesWithLists(root *xhtml.Node) *xhtml.Node {
	if root == nil {
		return nil
	}
	if root.Type == xhtml.ElementNode && root.Data == "table" {
		if isSingleCellLayoutTable(root) {
			clone := &xhtml.Node{
				Type:     root.Type,
				Data:     root.Data,
				DataAtom: root.DataAtom,
				Attr:     append([]xhtml.Attribute(nil), root.Attr...),
			}
			for c := root.FirstChild; c != nil; c = c.NextSibling {
				if nc := replaceTablesWithLists(c); nc != nil {
					appendChild(clone, nc)
				}
			}
			return clone
		}
		ol := &xhtml.Node{Type: xhtml.ElementNode, Data: "ol", DataAtom: atom.Ol}
		for _, tr := range collectTRs(root) {
			type cell struct {
				tag      string
				children []*xhtml.Node
			}
			var cells []cell
			allHeader := true
			for c := tr.FirstChild; c != nil; c = c.NextSibling {
				if c.Type != xhtml.ElementNode {
					continue
				}
				if c.Data == "th" || c.Data == "td" {
					if c.Data != "th" {
						allHeader = false
					}
					var kids []*xhtml.Node
					for k := c.FirstChild; k != nil; k = k.NextSibling {
						kids = append(kids, cloneDeep(k))
					}
					cells = append(cells, cell{tag: c.Data, children: kids})
				}
			}
			if len(cells) == 0 || allHeader {
				continue
			}
			li := &xhtml.Node{Type: xhtml.ElementNode, Data: "li", DataAtom: atom.Li}
			for i, cc := range cells {
				if i > 0 {
					appendChild(li, &xhtml.Node{Type: xhtml.TextNode, Data: " — "})
				}
				// For Chapter Guide style tables, the first cell is a time link. Make it bold.
				if i == 0 {
					strong := &xhtml.Node{Type: xhtml.ElementNode, Data: "strong", DataAtom: atom.Strong}
					for _, k := range cc.children {
						appendChild(strong, k)
					}
					appendChild(li, strong)
					continue
				}
				// If the cell rendered as a single <p>, unwrap it for list compactness.
				if len(cc.children) == 1 && cc.children[0].Type == xhtml.ElementNode && cc.children[0].Data == "p" {
					for k := cc.children[0].FirstChild; k != nil; k = k.NextSibling {
						appendChild(li, cloneDeep(k))
					}
					continue
				}
				for _, k := range cc.children {
					appendChild(li, k)
				}
			}
			appendChild(ol, li)
		}
		return ol
	}
	if root.Type != xhtml.ElementNode {
		return cloneShallowText(root)
	}
	clone := &xhtml.Node{
		Type:     root.Type,
		Data:     root.Data,
		DataAtom: root.DataAtom,
		Attr:     append([]xhtml.Attribute(nil), root.Attr...),
	}
	for c := root.FirstChild; c != nil; c = c.NextSibling {
		if nc := replaceTablesWithLists(c); nc != nil {
			appendChild(clone, nc)
		}
	}
	return clone
}

func cloneDeep(n *xhtml.Node) *xhtml.Node {
	if n == nil {
		return nil
	}
	switch n.Type {
	case xhtml.ErrorNode, xhtml.DocumentNode, xhtml.CommentNode, xhtml.DoctypeNode, xhtml.RawNode:
		return nil
	case xhtml.TextNode:
		return &xhtml.Node{Type: xhtml.TextNode, Data: n.Data}
	case xhtml.ElementNode:
		clone := &xhtml.Node{Type: xhtml.ElementNode, Data: n.Data, DataAtom: n.DataAtom, Attr: append([]xhtml.Attribute(nil), n.Attr...)}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if cc := cloneDeep(c); cc != nil {
				appendChild(clone, cc)
			}
		}
		return clone
	default:
		// Drop comments/doctype in conversion output.
		return nil
	}
}

func cloneShallowText(root *xhtml.Node) *xhtml.Node {
	if root == nil {
		return nil
	}
	if root.Type == xhtml.TextNode {
		return &xhtml.Node{Type: xhtml.TextNode, Data: root.Data}
	}
	return root
}

func collectTRs(t *xhtml.Node) []*xhtml.Node {
	var out []*xhtml.Node
	var walk func(*xhtml.Node)
	walk = func(n *xhtml.Node) {
		if n.Type == xhtml.ElementNode && n.Data == "tr" {
			out = append(out, n)
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(t)
	return out
}

func appendChild(parent, child *xhtml.Node) {
	child.Parent = parent
	if parent.LastChild != nil {
		parent.LastChild.NextSibling = child
		child.PrevSibling = parent.LastChild
		parent.LastChild = child
	} else {
		parent.FirstChild = child
		parent.LastChild = child
	}
}

func linkText(parent *xhtml.Node, text string) {
	tn := &xhtml.Node{Type: xhtml.TextNode, Data: text}
	tn.Parent = parent
	parent.FirstChild = tn
	parent.LastChild = tn
}

func collectText(n *xhtml.Node) string {
	var b strings.Builder
	var walk func(*xhtml.Node)
	walk = func(x *xhtml.Node) {
		if x.Type == xhtml.TextNode {
			b.WriteString(x.Data)
		}
		for c := x.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return strings.TrimSpace(b.String())
}

var allowedTags = map[string]bool{
	"p": true, "br": true, "hr": true,
	"h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true,
	"strong": true, "b": true, "em": true, "i": true, "u": true,
	"s": true, "del": true,
	"a":  true,
	"ul": true, "ol": true, "li": true,
	"blockquote": true,
	"pre":        true, "code": true,
	"center": true,
	"div":    true,
	"table":  true, "thead": true, "tbody": true, "tr": true, "th": true, "td": true,
	"img": true,
}

func sanitizeTree(n *xhtml.Node, opt Options) ([]*xhtml.Node, error) {
	if n == nil {
		return nil, nil
	}
	switch n.Type {
	case xhtml.ErrorNode, xhtml.DocumentNode, xhtml.DoctypeNode, xhtml.RawNode:
		return nil, nil
	case xhtml.CommentNode:
		return nil, nil
	case xhtml.TextNode:
		return []*xhtml.Node{{Type: xhtml.TextNode, Data: n.Data}}, nil
	case xhtml.ElementNode:
		tag := strings.ToLower(n.Data)
		if tag == "script" || tag == "style" {
			return nil, nil
		}
		if tag == "iframe" || tag == "object" || tag == "embed" {
			return linkFallbackNodes(n)
		}
		if tag == "blockquote" && (opt.QuoteMode == QuoteMonospace || opt.QuoteMode == QuotePullquoteMonospace) {
			txt := normalizeQuoteText(collectTextWithBR(n))
			if strings.TrimSpace(txt) == "" {
				return nil, nil
			}
			code := &xhtml.Node{Type: xhtml.ElementNode, Data: "code", DataAtom: atom.Code}
			appendChild(code, &xhtml.Node{Type: xhtml.TextNode, Data: txt})
			if opt.QuoteMode == QuotePullquoteMonospace {
				bq := &xhtml.Node{Type: xhtml.ElementNode, Data: "blockquote", DataAtom: atom.Blockquote}
				appendChild(bq, code)
				return []*xhtml.Node{bq}, nil
			}
			pre := &xhtml.Node{Type: xhtml.ElementNode, Data: "pre", DataAtom: atom.Pre}
			appendChild(pre, code)
			return []*xhtml.Node{pre}, nil
		}
		if tag == "input" {
			checked := false
			for _, a := range n.Attr {
				if strings.EqualFold(a.Key, "checked") {
					checked = true
				}
			}
			mark := "[ ] "
			if checked {
				mark = "[x] "
			}
			return []*xhtml.Node{{Type: xhtml.TextNode, Data: mark}}, nil
		}
		if !allowedTags[tag] {
			var out []*xhtml.Node
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				kids, err := sanitizeTree(c, opt)
				if err != nil {
					return nil, err
				}
				out = append(out, kids...)
			}
			return out, nil
		}
		outTag := tag
		if opt.DemoteHeadings && len(tag) == 2 && tag[0] == 'h' && tag[1] >= '1' && tag[1] <= '6' {
			if tag[1] < '6' {
				outTag = "h" + string(rune(tag[1]+1))
			}
		}
		attrs := filterAttrs(tag, n.Attr, opt)
		if tag == "a" && len(attrs) == 0 {
			var out []*xhtml.Node
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				kids, err := sanitizeTree(c, opt)
				if err != nil {
					return nil, err
				}
				out = append(out, kids...)
			}
			return out, nil
		}
		if tag == "img" && len(attrs) == 0 {
			alt := ""
			for _, a := range n.Attr {
				if strings.EqualFold(a.Key, "alt") {
					alt = a.Val
				}
			}
			if strings.TrimSpace(alt) == "" {
				return nil, nil
			}
			return []*xhtml.Node{{Type: xhtml.TextNode, Data: alt}}, nil
		}
		clone := &xhtml.Node{Type: xhtml.ElementNode, Data: outTag}
		clone.Attr = attrs
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			kids, err := sanitizeTree(c, opt)
			if err != nil {
				return nil, err
			}
			for _, k := range kids {
				appendChild(clone, k)
			}
		}
		return []*xhtml.Node{clone}, nil
	default:
		return nil, nil
	}
}

func collectTextWithBR(n *xhtml.Node) string {
	var b strings.Builder
	var walk func(*xhtml.Node)
	walk = func(x *xhtml.Node) {
		if x == nil {
			return
		}
		if x.Type == xhtml.TextNode {
			b.WriteString(x.Data)
		}
		if x.Type == xhtml.ElementNode && x.Data == "br" {
			b.WriteString("\n")
		}
		for c := x.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return strings.TrimSpace(b.String())
}

func normalizeQuoteText(s string) string {
	// Substack-friendly: strip any accidental leading '>' markers per line and
	// drop empty marker-only lines so quote-as-code looks clean.
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	for _, ln := range lines {
		t := strings.TrimSpace(ln)
		if t == ">" {
			continue
		}
		if strings.HasPrefix(t, ">") {
			t = strings.TrimSpace(strings.TrimPrefix(t, ">"))
		}
		if t == "" {
			continue
		}
		out = append(out, t)
	}
	return strings.TrimSpace(strings.Join(out, "\n"))
}

func linkFallbackNodes(n *xhtml.Node) ([]*xhtml.Node, error) {
	href := ""
	for _, a := range n.Attr {
		if strings.EqualFold(a.Key, "src") {
			href = strings.TrimSpace(a.Val)
			break
		}
	}
	if !safeHref(href) {
		return nil, nil
	}
	aNode := &xhtml.Node{Type: xhtml.ElementNode, Data: "a", DataAtom: atom.A}
	aNode.Attr = []xhtml.Attribute{{Key: "href", Val: href}}
	linkText(aNode, href)
	return []*xhtml.Node{aNode}, nil
}

func filterAttrs(tag string, attrs []xhtml.Attribute, opt Options) []xhtml.Attribute {
	switch tag {
	case "p":
		for _, a := range attrs {
			if !strings.EqualFold(strings.TrimSpace(a.Key), "style") {
				continue
			}
			if normalized, ok := normalizeParagraphTextAlignCenterStyle(a.Val); ok {
				return []xhtml.Attribute{{Key: "style", Val: normalized}}
			}
		}
		return nil
	case "a":
		var href string
		for _, a := range attrs {
			if strings.EqualFold(a.Key, "href") {
				href = strings.TrimSpace(a.Val)
				break
			}
		}
		if !safeHref(href) {
			return nil
		}
		return []xhtml.Attribute{{Key: "href", Val: href}}
	case "img":
		var src, alt string
		for _, a := range attrs {
			switch strings.ToLower(a.Key) {
			case "src":
				src = strings.TrimSpace(a.Val)
			case "alt":
				alt = a.Val
			}
		}
		if !safeImgSrc(src) {
			base := strings.TrimSpace(opt.bundleImageBase)
			ref := strings.TrimSpace(strings.TrimPrefix(src, "./"))
			if base != "" && ref != "" && !isAbsHTTPURL(ref) {
				src = ResolveImageReference(ref, base)
			}
		}
		if !safeImgSrc(src) {
			return nil
		}
		out := []xhtml.Attribute{{Key: "src", Val: src}}
		if alt != "" {
			out = append(out, xhtml.Attribute{Key: "alt", Val: alt})
		}
		return out
	case "div":
		align := ""
		for _, a := range attrs {
			if strings.EqualFold(a.Key, "align") {
				align = strings.TrimSpace(a.Val)
				break
			}
		}
		if strings.EqualFold(align, "center") {
			return []xhtml.Attribute{{Key: "align", Val: "center"}}
		}
		return nil
	case "td", "th":
		align := ""
		for _, a := range attrs {
			if strings.EqualFold(a.Key, "align") {
				align = strings.TrimSpace(a.Val)
				break
			}
		}
		if strings.EqualFold(align, "center") {
			return []xhtml.Attribute{{Key: "align", Val: "center"}}
		}
		return nil
	case "table":
		var out []xhtml.Attribute
		for _, a := range attrs {
			switch strings.ToLower(strings.TrimSpace(a.Key)) {
			case "align":
				if strings.EqualFold(strings.TrimSpace(a.Val), "center") {
					out = append(out, xhtml.Attribute{Key: "align", Val: "center"})
				}
			case "width":
				w := strings.TrimSpace(strings.ToLower(a.Val))
				if w == "100%" {
					out = append(out, xhtml.Attribute{Key: "width", Val: "100%"})
				}
			}
		}
		return out
	default:
		return nil
	}
}

func safeHref(h string) bool {
	h = strings.TrimSpace(h)
	ls := strings.ToLower(h)
	return strings.HasPrefix(ls, "http://") || strings.HasPrefix(ls, "https://") || strings.HasPrefix(ls, "mailto:")
}

func safeImgSrc(s string) bool {
	ls := strings.ToLower(strings.TrimSpace(s))
	return strings.HasPrefix(ls, "https://") || strings.HasPrefix(ls, "http://")
}

func renderNode(w io.Writer, n *xhtml.Node) error {
	if n == nil {
		return nil
	}
	switch n.Type {
	case xhtml.ErrorNode, xhtml.DocumentNode, xhtml.CommentNode, xhtml.DoctypeNode, xhtml.RawNode:
		return nil
	case xhtml.TextNode:
		_, err := io.WriteString(w, html.EscapeString(n.Data))
		return err
	case xhtml.ElementNode:
		tag := n.Data
		void := tag == "br" || tag == "hr" || tag == "img"
		if _, err := fmt.Fprintf(w, "<%s", tag); err != nil {
			return err
		}
		for _, a := range n.Attr {
			if err := writeAttr(w, a); err != nil {
				return err
			}
		}
		if void {
			_, err := w.Write([]byte(">"))
			return err
		}
		if _, err := fmt.Fprintf(w, ">"); err != nil {
			return err
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if err := renderNode(w, c); err != nil {
				return err
			}
		}
		_, err := fmt.Fprintf(w, "</%s>", tag)
		return err
	default:
		return nil
	}
}

func writeAttr(w io.Writer, a xhtml.Attribute) error {
	key := a.Key
	if a.Namespace != "" {
		key = a.Namespace + ":" + a.Key
	}
	_, err := fmt.Fprintf(w, " %s=\"%s\"", key, html.EscapeString(a.Val))
	return err
}
