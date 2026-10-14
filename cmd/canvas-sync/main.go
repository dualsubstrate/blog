// Command canvas-sync builds the DSS hops Vite bundle into static/js/dss-hops.
package main

import (
	"flag"
	"log"
	"os"

	"github.com/dualsubstrate/blog/internal/canvasbuild"
)

func main() {
	log.SetFlags(0)
	root := flag.String("root", ".", "blog repository root")
	canvas := flag.String("canvas", "", "optional canvas source path (default: DSS_HOPS_CANVAS or Cursor project path)")
	flag.Parse()

	if err := canvasbuild.Run(canvasbuild.Options{
		Root:         *root,
		CanvasSource: *canvas,
		CI:           os.Getenv("CI") != "",
	}); err != nil {
		log.Fatalf("canvas-sync: %v", err)
	}
}
