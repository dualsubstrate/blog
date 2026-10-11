package substackbrowser

import (
	"net/url"
)

// FixtureDataURL returns a data: URL that loads a tiny page with a ProseMirror
// style contenteditable for local paste testing without Substack credentials.
func FixtureDataURL() string {
	const page = `<!DOCTYPE html><html><head><meta charset="utf-8"><title>paste fixture</title></head><body>
<p style="font-family:system-ui">Local paste fixture (not Substack). Body editor is below.</p>
<div class="ProseMirror" contenteditable="true" style="border:1px solid #444;min-height:220px;padding:10px;font-family:system-ui">
<p>Start here. Automation will append converted HTML at the end of this field.</p>
</div>
</body></html>`
	return "data:text/html;charset=utf-8," + url.PathEscape(page)
}
