// Command tag-register scans Hugo content front matter and writes data/tag-register.txt.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/dualsubstrate/blog/internal/cliout"
	"github.com/dualsubstrate/blog/internal/tagregister"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	contentDir := flag.String("content", "content", "directory of Markdown content (recursive *.md)")
	depPath := flag.String("deprecations", "data/tag-deprecations.toml", "TOML file with [[deprecated]] from/to rows (optional if missing)")
	outPath := flag.String("out", "data/tag-register.txt", "output register path")
	flag.Parse()

	absContent, err := filepath.Abs(*contentDir)
	if err != nil {
		log.Fatalf("content path: %v", err)
	}
	st, err := os.Stat(absContent)
	if err != nil || !st.IsDir() {
		log.Fatalf("content %q must be an existing directory", *contentDir)
	}

	counts, err := tagregister.ScanMarkdownTags(absContent)
	if err != nil {
		log.Fatalf("scan: %v", err)
	}

	var deps []tagregister.Deprecation
	if *depPath != "" {
		absDep, err := filepath.Abs(*depPath)
		if err != nil {
			log.Fatalf("deprecations path: %v", err)
		}
		deps, err = tagregister.LoadDeprecations(absDep)
		if err != nil {
			log.Fatalf("deprecations: %v", err)
		}
	}

	body := tagregister.RenderRegister(counts, deps)
	absOut, err := filepath.Abs(*outPath)
	if err != nil {
		log.Fatalf("out path: %v", err)
	}
	if err := writeFileAtomically(absOut, []byte(body)); err != nil {
		log.Fatalf("write: %v", err)
	}
	cliout.PrintFileWritten(os.Stdout, "tag-register", absOut, strconv.Itoa(len(counts))+" distinct tags")
}

func writeFileAtomically(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	f, err := os.CreateTemp(dir, ".tag-register-*.tmp")
	if err != nil {
		return fmt.Errorf("temp: %w", err)
	}
	tmpName := f.Name()
	abort := func(err error) error {
		_ = f.Close()
		_ = os.Remove(tmpName)
		return err
	}

	if _, err := f.Write(data); err != nil {
		return abort(fmt.Errorf("write temp: %w", err))
	}
	if err := f.Sync(); err != nil {
		return abort(fmt.Errorf("sync: %w", err))
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("close temp: %w", err)
	}
	if err := os.Chmod(tmpName, 0o644); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("chmod: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("rename: %w", err)
	}
	return nil
}
