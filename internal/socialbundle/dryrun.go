package socialbundle

import (
	"fmt"
	"io"
	"strings"

	"github.com/dualsubstrate/blog/internal/cliout"
)

// PrintDryRunItem prints one bundle preview in the same shape as facebook-autopost dry-run.
// ImagePath is optional; when empty the Image line is omitted.
func PrintDryRunItem(w io.Writer, index, total int, mode, postURL, imagePath, message string) {
	PrintDryRunLinkedInItem(w, index, total, mode, postURL, imagePath, "", "", message)
}

// PrintDryRunLinkedInItem prints LinkedIn dry-run details including optional YouTube article card fields.
func PrintDryRunLinkedInItem(w io.Writer, index, total int, mode, postURL, imagePath, youtubeURL, thumbURL, message string) {
	cliout.PrintDryRunBanner(w, index, total, "social")
	cliout.KV(w, "Mode", mode)
	cliout.KV(w, "URL", postURL)
	if strings.TrimSpace(imagePath) != "" {
		cliout.KV(w, "Image", imagePath)
	}
	if strings.TrimSpace(youtubeURL) != "" {
		cliout.KV(w, "YouTube", youtubeURL)
	}
	if strings.TrimSpace(thumbURL) != "" {
		cliout.KV(w, "Thumbnail", thumbURL)
	}
	fmt.Fprint(w, "\n   ----- BEGIN POST TEXT -----\n\n")
	fmt.Fprint(w, message)
	fmt.Fprint(w, "\n\n   ----- END POST TEXT -----\n\n")
}

// PrintDryRunLinkedInIdempotencySkipped prints a one-line note when LinkedIn would skip the recent-post scan.
func PrintDryRunLinkedInIdempotencySkipped(w io.Writer) {
	cliout.PrintDryRunNote(w, "recent-post idempotency scan would be skipped (-disable-idempotency or LINKEDIN_DISABLE_IDEMPOTENCY)")
}
