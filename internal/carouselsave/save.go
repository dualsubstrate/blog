// Package carouselsave writes studio carousel.json files under content/.
package carouselsave

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const maxBodyBytes = 4 << 20

// Request is a studio save payload. Body is the pretty-printed file contents.
type Request struct {
	Source   string `json:"source"`
	Filename string `json:"filename"`
	Body     string `json:"body"`
}

// ResolvePath maps a Hugo bundle source path to an allowed carousel JSON file.
func ResolvePath(repoRoot, source, filename string) (string, error) {
	name := filepath.Base(strings.TrimSpace(filename))
	if name != "carousel.json" && name != "carousel.es.json" {
		return "", fmt.Errorf("filename must be carousel.json or carousel.es.json")
	}

	source = strings.TrimSpace(strings.ReplaceAll(source, "\\", "/"))
	if source == "" {
		return "", fmt.Errorf("source is empty")
	}
	if strings.Contains(source, "..") {
		return "", fmt.Errorf("source path is invalid")
	}
	if !strings.HasPrefix(source, "content/") {
		return "", fmt.Errorf("source must be under content/")
	}

	absRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return "", fmt.Errorf("repo root: %w", err)
	}
	bundleDir := filepath.FromSlash(filepath.Dir(source))
	absOut, err := filepath.Abs(filepath.Join(absRoot, bundleDir, name))
	if err != nil {
		return "", fmt.Errorf("output path: %w", err)
	}

	contentRoot, err := filepath.Abs(filepath.Join(absRoot, "content"))
	if err != nil {
		return "", fmt.Errorf("content root: %w", err)
	}
	rel, err := filepath.Rel(contentRoot, absOut)
	if err != nil {
		return "", fmt.Errorf("relativize: %w", err)
	}
	if rel == "." || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("refusing to write outside content/")
	}
	return absOut, nil
}

// WriteBody validates JSON and writes Body to the resolved carousel file.
func WriteBody(repoRoot string, req Request) (string, error) {
	path, err := ResolvePath(repoRoot, req.Source, req.Filename)
	if err != nil {
		return "", err
	}
	body := strings.TrimSpace(req.Body)
	if body == "" {
		return "", fmt.Errorf("body is empty")
	}
	if len(body) > maxBodyBytes {
		return "", fmt.Errorf("body exceeds %d bytes", maxBodyBytes)
	}
	if err := validateDeckJSON(body); err != nil {
		return "", err
	}
	if !strings.HasSuffix(body, "\n") {
		body += "\n"
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", fmt.Errorf("mkdir: %w", err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		return "", fmt.Errorf("write %s: %w", path, err)
	}
	return path, nil
}

func validateDeckJSON(body string) error {
	var probe any
	if err := json.Unmarshal([]byte(body), &probe); err != nil {
		return fmt.Errorf("body is not valid JSON: %w", err)
	}
	obj, ok := probe.(map[string]any)
	if !ok {
		return fmt.Errorf("body must be a JSON object")
	}
	if _, ok := obj["slides"]; !ok {
		return fmt.Errorf("body is missing slides")
	}
	if _, ok := obj["deck"]; !ok {
		return fmt.Errorf("body is missing deck")
	}
	return nil
}
