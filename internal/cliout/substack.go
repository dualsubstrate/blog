package cliout

import (
	"io"
	"strconv"
)

// PrintSubstackMarkPublished reports a successful marker write.
func PrintSubstackMarkPublished(out io.Writer, bundleDir, marker, target string) {
	Heading(out, "📌", "marked published ("+target+")")
	Meta(out, "📂", "bundle", bundleDir)
	Meta(out, "📄", "marker", marker)
	Blank(out)
}

// PrintSubstackDryRun summarizes a dry-run build without opening Chrome.
func PrintSubstackDryRun(out io.Writer, fields map[string]string) {
	Heading(out, "🧪", "Substack dry-run (no Chrome)")
	for _, key := range []string{"url", "html_bytes", "headless", "insert_subscribe_button", "include_project_about", "config", "schedule"} {
		if v, ok := fields[key]; ok && v != "" {
			icon := dryRunIcon(key)
			Meta(out, icon, dryRunLabel(key), v)
		}
	}
	Blank(out)
}

// PrintSubstackBrowserOpen reports navigation to Substack or a fixture.
func PrintSubstackBrowserOpen(out io.Writer, action, url string) {
	icon := "🌐"
	title := "opening Substack"
	switch action {
	case "login":
		icon = "🔐"
		title = "Substack login"
	case "cc":
		icon = "🎛️"
		title = "Substack command center"
	case "paste", "paste-schedule":
		icon = "📝"
		title = "Substack editor"
	}
	Heading(out, icon, title)
	Meta(out, "🔗", "url", url)
	Blank(out)
}

// PrintSubstackPasteFinished reports body paste / schedule automation completion.
func PrintSubstackPasteFinished(out io.Writer, action string, marked bool, bundleDir, marker, target string) {
	if action == "paste-schedule" {
		Heading(out, "✅", "paste-schedule finished")
		Body(out,
			"Body and publish settings were filled in Substack.",
			"If the publish modal is still open, click Schedule or Publish when ready.",
		)
	} else {
		Heading(out, "✅", "paste finished")
		Body(out,
			"HTML was inserted into the editor.",
			"Review the draft in the browser. This tool did not publish.",
		)
	}
	if marked {
		Divider(out)
		Meta(out, "📌", "recorded in", bundleDir+"/"+marker+" ("+target+")")
	}
	Blank(out)
}

// PrintSubstackAborted reports user-cancelled flow before publish.
func PrintSubstackAborted(out io.Writer) {
	Heading(out, "🛑", "aborted before publish")
	Body(out, "No publish-settings step ran. social-published was not updated.")
	Blank(out)
}

// PrintSubstackPickMissingSpanish explains a missing index.es.md after picker selection.
func PrintSubstackPickMissingSpanish(out io.Writer, path string) {
	Heading(out, "❌", "Spanish bundle missing index.es.md")
	Meta(out, "📄", "expected", path)
	Divider(out)
	Hint(out, "Add Spanish page", "author index.es.md in the bundle, then retry make sb-es")
	Blank(out)
}

// PrintSubstackMarkPublishedUsage explains mark-published needs -in.
func PrintSubstackMarkPublishedUsage(out io.Writer) {
	Heading(out, "❌", "mark-published needs a bundle path")
	Hint(out, "Example", "make sb-mark-published POST=human-condition/2026-05-01-ego-as-game")
	Blank(out)
}

// PrintSubstackHTMLWritten reports HTML export to a file.
func PrintSubstackHTMLWritten(out io.Writer, inPath, outPath string, bytes int) {
	Heading(out, "✅", "Substack HTML written")
	Meta(out, "📥", "source", inPath)
	Meta(out, "📤", "output", outPath)
	Meta(out, "📊", "bytes", strconv.Itoa(bytes))
	Blank(out)
}

func dryRunIcon(key string) string {
	switch key {
	case "url":
		return "🔗"
	case "html_bytes":
		return "📊"
	case "schedule":
		return "📅"
	case "config":
		return "⚙️"
	default:
		return "•"
	}
}

func dryRunLabel(key string) string {
	switch key {
	case "url":
		return "target URL"
	case "html_bytes":
		return "html bytes"
	case "headless":
		return "headless"
	case "insert_subscribe_button":
		return "subscribe button"
	case "include_project_about":
		return "cognitive-memetics footer"
	case "config":
		return "config"
	case "schedule":
		return "schedule"
	default:
		return key
	}
}
