package main

import (
	"fmt"
	"os"

	"github.com/dualsubstrate/blog/internal/substackbrowser"
)

func main() {
	force := false
	for _, arg := range os.Args[1:] {
		if arg == "--force" {
			force = true
		}
	}
	example, err := os.ReadFile("docs/substack-html/substack-config.example.json")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	written, err := substackbrowser.InitUserSubstackConfig(example, force)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if written {
		fmt.Fprintln(os.Stderr, "Wrote user substack.json under ~/.config/behaviour-engineering/")
	} else {
		fmt.Fprintln(os.Stderr, "User substack.json already exists (use --force to overwrite)")
	}
	if err := substackbrowser.InitRepoSubstackConfig(example); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if _, err := os.Stat("substack.json"); err == nil {
		fmt.Fprintln(os.Stderr, "Repo substack.json is present (created if it was missing)")
	}
	fmt.Fprintln(os.Stderr, "Chrome profile: use direnv .envrc + SUBSTACK_CHROMIUM_USER_DATA_DIRECTORY")
}
