// Command normalize-post prints a normalized section/slug path for Task and publishing scripts.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dualsubstrate/blog/internal/content"
)

func main() {
	flag.Parse()
	post := flag.Arg(0)
	if post == "" {
		post = os.Getenv("POST")
	}
	if post == "" {
		fmt.Fprintln(os.Stderr, "usage: normalize-post <section/slug|content/.../index.md>")
		os.Exit(2)
	}
	fmt.Print(content.NormalizePostSlug(post))
}
