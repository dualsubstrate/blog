// Command publish-calendar scans Hugo content and writes a JSON snapshot for the calendar UI.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dualsubstrate/blog/internal/cliout"
	"github.com/dualsubstrate/blog/internal/publishcalendar"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	contentRoot := flag.String("content", "content", "Hugo content directory")
	outPath := flag.String("out", "static/calendar/publish-calendar.json", "output JSON path")
	channelsFlag := flag.String("channels", strings.Join(publishcalendar.DefaultChannels, ","), "comma-separated channel ids")
	flag.Parse()

	channels := splitCSV(*channelsFlag)
	cal, err := publishcalendar.Build(*contentRoot, channels)
	if err != nil {
		log.Fatalf("build: %v", err)
	}

	absOut, err := filepath.Abs(*outPath)
	if err != nil {
		log.Fatalf("out path: %v", err)
	}
	if err := publishcalendar.WriteJSON(absOut, cal); err != nil {
		log.Fatalf("write: %v", err)
	}
	cliout.PrintFileWritten(os.Stdout, "publish-calendar", absOut, fmt.Sprintf("%d bundles", len(cal.Entries)))
}

func splitCSV(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}
