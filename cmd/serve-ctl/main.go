// Command serve-ctl manages the local process-compose dev stack.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dualsubstrate/blog/internal/servectl"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}

	sub := os.Args[1]
	switch sub {
	case "cleanup":
		fs := flag.NewFlagSet("cleanup", flag.ExitOnError)
		rootFlag := fs.String("root", ".", "blog repository root")
		preflight := fs.Bool("preflight", false, "exit after cleanup without status message")
		_ = fs.Parse(os.Args[2:])
		if err := servectl.Cleanup(*rootFlag, *preflight, false); err != nil {
			log.Fatalf("serve-ctl cleanup: %v", err)
		}
	case "up":
		fs := flag.NewFlagSet("up", flag.ExitOnError)
		rootFlag := fs.String("root", ".", "blog repository root")
		_ = fs.Parse(os.Args[2:])
		if err := servectl.Up(*rootFlag); err != nil {
			log.Fatalf("serve-ctl up: %v", err)
		}
	case "wrap":
		if len(os.Args) < 3 {
			log.Fatal("serve-ctl wrap: missing command")
		}
		if err := servectl.Wrap(os.Args[2:]); err != nil {
			os.Exit(1)
		}
	case "mcp":
		fs := flag.NewFlagSet("mcp", flag.ExitOnError)
		rootFlag := fs.String("root", ".", "blog repository root")
		_ = fs.Parse(os.Args[2:])
		if err := servectl.RunContentPipelinesMCP(*rootFlag, fs.Args()); err != nil {
			log.Fatalf("serve-ctl mcp: %v", err)
		}
	case "help", "-h", "--help":
		printUsage()
	default:
		log.Fatalf("serve-ctl: unknown subcommand %q", sub)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `serve-ctl — local dev stack helpers

Usage:
  serve-ctl cleanup [--preflight]
  serve-ctl up
  serve-ctl wrap <command> [args...]
  serve-ctl mcp [args...]

`)
}
