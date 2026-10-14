// Command sync-content-ai copies selected Cursor skills and rules into content-ai/.
package main

import (
	"flag"
	"log"

	"github.com/dualsubstrate/blog/internal/synccontentai"
)

func main() {
	log.SetFlags(0)
	root := flag.String("root", ".", "blog repository root")
	flag.Parse()

	if err := synccontentai.Run(*root); err != nil {
		log.Fatalf("sync-content-ai: %v", err)
	}
}
