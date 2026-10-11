// Package cliout formats user-facing terminal panels with consistent icons and spacing.
package cliout

import (
	"fmt"
	"io"
	"strings"
)

// Blank prints one empty line.
func Blank(out io.Writer) {
	fmt.Fprintln(out)
}

// Heading prints a titled panel header: "icon  title".
func Heading(out io.Writer, icon, title string) {
	Blank(out)
	fmt.Fprintf(out, "%s  %s\n", icon, title)
}

// Subheading prints a secondary line under a heading.
func Subheading(out io.Writer, text string) {
	fmt.Fprintf(out, "\n   %s\n", text)
}

// Meta prints a labeled metadata row with a leading icon.
func Meta(out io.Writer, icon, label, value string) {
	fmt.Fprintf(out, "   %s  %s: %s\n", icon, label, value)
}

// Body prints wrapped explanatory text indented under a panel.
func Body(out io.Writer, lines ...string) {
	for _, line := range lines {
		fmt.Fprintf(out, "   %s\n", line)
	}
}

// Hint prints a next-step command with the standard pointer icon.
func Hint(out io.Writer, label, command string) {
	if strings.TrimSpace(label) != "" {
		fmt.Fprintf(out, "   👉  %s:\n", label)
	} else {
		fmt.Fprintln(out, "   👉")
	}
	fmt.Fprintf(out, "       %s\n", strings.TrimSpace(command))
}

// Hints groups hint blocks with a blank line between them.
func Hints(out io.Writer, pairs ...[2]string) {
	for i, p := range pairs {
		if i > 0 {
			Blank(out)
		}
		Hint(out, p[0], p[1])
	}
}

// KV prints a simple key: value row (no icon).
func KV(out io.Writer, key, value string) {
	fmt.Fprintf(out, "   %s: %s\n", key, value)
}

// Divider prints a light section break inside a panel.
func Divider(out io.Writer) {
	Blank(out)
}

// Status prints a one-line status with icon (for log-style messages).
func Status(out io.Writer, icon, message string) {
	fmt.Fprintf(out, "%s  %s\n", icon, message)
}

// Item prints a numbered list row.
func Item(out io.Writer, index int, text string) {
	fmt.Fprintf(out, "   %d.  %s\n", index, text)
}
