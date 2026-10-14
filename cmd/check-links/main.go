// Command check-links builds the Hugo site and crawls it with muffet.
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dualsubstrate/blog/internal/linkcheck"
)

func main() {
	log.SetFlags(0)
	root := flag.String("root", ".", "blog repository root")
	port := flag.Int("port", 0, "loopback port for static server (default: LINK_CHECK_PORT or 8765)")
	hugo := flag.String("hugo", "", "hugo binary (default: $HUGO or hugo)")
	flag.Parse()

	if *port == 0 {
		if v := os.Getenv("LINK_CHECK_PORT"); v != "" {
			p, err := strconv.Atoi(v)
			if err != nil {
				log.Fatalf("check-links: invalid LINK_CHECK_PORT: %v", err)
			}
			*port = p
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()

	if err := linkcheck.Run(ctx, linkcheck.Options{
		Root:  *root,
		Port:  *port,
		Hugo:  *hugo,
	}); err != nil {
		log.Fatalf("check-links: %v", err)
	}
}
