// Command hugo-build runs a production-like Hugo build and post-processes carousel previews.
package main

import (
	"flag"
	"log"

	"github.com/dualsubstrate/blog/internal/hugobuild"
)

func main() {
	log.SetFlags(0)
	root := flag.String("root", ".", "blog repository root")
	hugo := flag.String("hugo", "", "hugo binary (default: $HUGO or hugo)")
	preview := flag.Bool("preview", false, "build drafts/future with preview environment (also when HUGO_PREVIEW is set)")
	baseURL := flag.String("baseURL", "", "optional base URL (default: $BASE_URL)")
	flag.Parse()

	if err := hugobuild.Run(hugobuild.Options{
		Root:    *root,
		Hugo:    *hugo,
		Preview: *preview,
		BaseURL: *baseURL,
	}); err != nil {
		log.Fatalf("hugo-build: %v", err)
	}
}
