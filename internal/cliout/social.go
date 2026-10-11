package cliout

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// PrintSocialAutopostEmpty reports no bundles matched the requested publish date (warn / skip, exit 0).
func PrintSocialAutopostEmpty(out io.Writer, network, date, postPath string) {
	Heading(out, "⚠️", network+" autopost: nothing scheduled for "+date+" (skip)")
	if postPath != "" {
		Meta(out, "📌", "post filter", postPath)
	}
	Body(out,
		"No bundle under content/ has that publish date in index.md front matter,",
		"or the bundle is missing linkedin.txt. This is not a failure.",
	)
	Divider(out)
	cmd := strings.ToLower(network) + "-autopost"
	Hints(out,
		[2]string{"Check the publish calendar", "make calendar  (then open /calendar/ under make serve)"},
		[2]string{"Target one bundle by path", "make " + cmd + " DATE=" + date + " SOCIAL_POST=section/slug"},
		[2]string{"List Hugo permalinks", "make list"},
	)
	Blank(out)
}

// PrintSocialAutopostStart announces a live posting run.
func PrintSocialAutopostStart(out io.Writer, network string, count, httpRetries int) {
	Status(out, "🚀", network+" autopost: posting "+strconv.Itoa(count)+" item(s) (http-retries="+strconv.Itoa(httpRetries)+")")
}

// PrintSocialSkipAlreadyMarked reports a bundle skipped because social-published already records the target.
func PrintSocialSkipAlreadyMarked(out io.Writer, target, rel string) {
	Status(out, "⏭️", "skip: already in social-published ("+target+"): content/"+rel)
}

// PrintSocialPosted reports a successful network post.
func PrintSocialPosted(out io.Writer, network, detail string) {
	Status(out, "✅", network+" posted: "+detail)
}

// PrintSocialMarkerRecorded reports social-published was updated.
func PrintSocialMarkerRecorded(out io.Writer, target, reason string) {
	Status(out, "📌", "social-published: recorded "+target+" ("+reason+")")
}

// PrintSocialFailures summarizes a failed autopost run.
func PrintSocialFailures(out io.Writer, network string, count int) {
	Status(out, "❌", network+" autopost: "+strconv.Itoa(count)+" bundle(s) failed (publish or idempotency check)")
}

// PrintDryRunBanner opens a dry-run preview block for one bundle.
func PrintDryRunBanner(out io.Writer, index, total int, network string) {
	if index > 0 {
		Blank(out)
		fmt.Fprintln(out, "------------------------------------------------------------")
		Blank(out)
	}
	Heading(out, "🧪", "DRY-RUN item "+strconv.Itoa(index+1)+"/"+strconv.Itoa(total)+" ("+network+")")
}

// PrintDryRunNote prints a secondary dry-run note.
func PrintDryRunNote(out io.Writer, note string) {
	Status(out, "ℹ️", note)
}

// PrintFileWritten reports a generated output file.
func PrintFileWritten(out io.Writer, tool, path, detail string) {
	Heading(out, "✅", tool+": wrote "+path)
	if detail != "" {
		Meta(out, "📊", "detail", detail)
	}
	Blank(out)
}

// PrintConfigExists reports an init command found existing config.
func PrintConfigExists(out io.Writer, path string) {
	Status(out, "ℹ️", path+" already exists")
}

// PrintConfigCreated reports a new config file was copied.
func PrintConfigCreated(out io.Writer, path, note string) {
	Heading(out, "✅", "created "+path)
	if note != "" {
		Body(out, note)
	}
	Blank(out)
}

// PrintMissingBundle explains a Makefile bundle path error.
func PrintMissingBundle(out io.Writer, relPath, makeTarget string) {
	Heading(out, "❌", "missing bundle: content/"+relPath)
	Body(out, "Set POST=section/slug (path under content/).")
	Divider(out)
	Hint(out, "Example", makeTarget+" human-condition/2026-05-01-ego-as-game")
	Blank(out)
}
