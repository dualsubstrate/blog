package cliout

import (
	"io"
	"strconv"
	"strings"
)

// PublishQueueHints lists Makefile commands shown after queue panels.
type PublishQueueHints struct {
	LangLabel string
	ListCmd   string
	PickCmd   string
	PasteCmd  string
	MarkHint  string
}

// PrintPublishQueueEmpty explains that no bundles are waiting for a publish target.
func PrintPublishQueueEmpty(out io.Writer, contentRoot, target, marker string, hints PublishQueueHints) {
	Heading(out, "✅", hints.LangLabel+" queue: empty ("+target+")")
	Subheading(out, "Every ready bundle already has a line for this target in the marker file,")
	Body(out, "or nothing is ready to publish yet (missing sidecar, draft, etc.).")
	Divider(out)
	Meta(out, "📂", "content root", contentRoot)
	Meta(out, "🏷️", "target", target)
	Meta(out, "📌", "marker", marker)
	Divider(out)
	Hints(out,
		[2]string{"Paste one post anyway (skips the picker)", hints.PasteCmd},
		[2]string{"Inspect the queue again", hints.ListCmd},
		[2]string{"Put a post back in the picker", "remove the " + target + " line from content/section/slug/" + marker},
	)
	Blank(out)
}

// PrintPublishQueueList shows unpublished bundles and next-step commands.
func PrintPublishQueueList(out io.Writer, paths []string, target, marker string, hints PublishQueueHints) {
	first := paths[0]
	Heading(out, "📋", hints.LangLabel+" queue: "+strconv.Itoa(len(paths))+" bundle(s) waiting ("+target+")")
	for i, p := range paths {
		Item(out, i+1, p)
	}
	Divider(out)
	directPaste := strings.TrimSuffix(hints.PasteCmd, " section/slug") + " " + first
	markCmd := strings.Replace(hints.MarkHint, "section/slug", first, 1)
	Hints(out,
		[2]string{"Interactive picker", hints.PickCmd},
		[2]string{"Paste the first one directly", directPaste},
		[2]string{"After it is live on Substack", markCmd},
	)
	Blank(out)
}
