// Command audit-essay-voice runs the deterministic voice audit on one Hugo post.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/dualsubstrate/blog/internal/essayvoice"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("audit-essay-voice", flag.ContinueOnError)
	fs.SetOutput(stderr)
	post := fs.String("post", "", "path to a Hugo Markdown file (index.md)")
	voiceID := fs.String("voice", "", "voice id under -voices-dir (for example patient-narrator)")
	voicesDir := fs.String("voices-dir", "data/voices", "directory of voice Markdown files")
	asJSON := fs.Bool("json", false, "print a JSON report")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *post == "" || *voiceID == "" {
		fmt.Fprintln(stderr, "usage: audit-essay-voice -post content/.../index.md -voice patient-narrator [-voices-dir data/voices] [-json]")
		return 2
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	file, err := essayvoice.Load(ctx, *voicesDir, *voiceID)
	if err != nil {
		fmt.Fprintf(stderr, "voice: %v\n", err)
		return 1
	}
	result, err := essayvoice.AuditFile(ctx, *post, file.Profile)
	if err != nil {
		fmt.Fprintf(stderr, "audit: %v\n", err)
		return 1
	}
	if *asJSON {
		payload := struct {
			Pass                   bool   `json:"pass"`
			Voice                  string `json:"voice"`
			LongestStaccatoRun     int    `json:"longest_staccato_run"`
			ParagraphsMissingVerbs int    `json:"paragraphs_missing_verbs"`
			Violations             any    `json:"violations"`
		}{
			Pass:                   result.Pass,
			Voice:                  file.Profile.ID,
			LongestStaccatoRun:     result.LongestStaccatoRun,
			ParagraphsMissingVerbs: result.ParagraphsMissingVerbs,
			Violations:             result.Violations,
		}
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(payload); err != nil {
			fmt.Fprintf(stderr, "json: %v\n", err)
			return 1
		}
	} else if result.Pass {
		fmt.Fprintf(stdout, "pass %s %s\n", file.Profile.ID, *post)
	} else {
		fmt.Fprintf(stdout, "fail %s %s (%d)\n", file.Profile.ID, *post, len(result.Violations))
		for _, v := range result.Violations {
			fmt.Fprintf(stdout, "- %s: %s\n", v.Kind, v.Detail)
		}
	}
	if result.Pass {
		return 0
	}
	return 1
}
