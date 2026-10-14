// Command link-providers creates local provider and skill symlinks.
package main

import (
	"flag"
	"log"

	"github.com/dualsubstrate/blog/internal/linkproviders"
)

func main() {
	log.SetFlags(0)
	root := flag.String("root", ".", "blog repository root")
	flag.Parse()

	if err := linkproviders.Run(*root); err != nil {
		log.Fatalf("link-providers: %v", err)
	}
}
