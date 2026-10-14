package synccontentai

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var skills = []string{
	"site-claims-content", "site-video-content", "site-cognitive-memetics-content", "site-spanish-translation-content",
	"site-revise-post", "site-revise-flow", "site-revise-hooks", "site-revise-format", "site-revise-emphasis",
	"site-curiosity-title", "site-revise-score", "site-revise-post-es", "site-revise-spanish", "site-revise-prose",
	"site-linkedin-post", "site-facebook-post", "site-substack-post", "site-carousel-post",
	"site-tag-register", "site-tag-unify", "site-short-link-register", "site-ai-for-general-audience",
}

var rules = []string{
	"always-rules-0-ai.mdc",
	"always-rules-01-human-interaction.mdc",
	"site-always-rules-01-human-interaction.mdc",
	"site-content-placement.mdc",
	"site-content-markdown-writing.mdc",
}

var allowedSkillFiles = map[string]bool{
	"SKILL.md": true, "reference.md": true, "examples.md": true,
}

// Run copies selected skills and rules into content-ai/ for Continue-readable copies.
func Run(root string) error {
	root, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("root: %w", err)
	}
	skillsOut := filepath.Join(root, "content-ai", "skills")
	rulesOut := filepath.Join(root, "content-ai", "rules")
	if err := os.MkdirAll(skillsOut, 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(rulesOut, 0o755); err != nil {
		return err
	}

	for _, name := range skills {
		src := filepath.Join(root, ".cursor", "skills", name)
		dst := filepath.Join(skillsOut, name)
		if _, err := os.Stat(src); err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("skip missing %s\n", src)
				continue
			}
			return err
		}
		if err := syncSkillDir(src, dst); err != nil {
			return err
		}
	}

	for _, name := range rules {
		src := filepath.Join(root, ".cursor", "rules", name)
		dst := filepath.Join(rulesOut, name)
		if err := copyFile(src, dst); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("missing rule %s", src)
			}
			return err
		}
	}

	if err := rewriteContentAI(root); err != nil {
		return err
	}
	fmt.Println("Synced content-ai/ from .cursor/ (Continue-readable copies).")
	return nil
}

func syncSkillDir(src, dst string) error {
	if err := os.RemoveAll(dst); err != nil && !os.IsNotExist(err) {
		return err
	}
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(dst, rel), 0o755)
		}
		base := filepath.Base(path)
		if !allowedSkillFiles[base] {
			return nil
		}
		return copyFile(path, filepath.Join(dst, rel))
	})
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}

func rewriteContentAI(root string) error {
	contentAI := filepath.Join(root, "content-ai")
	return filepath.WalkDir(contentAI, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext != ".md" && ext != ".mdc" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		s := string(data)
		s = strings.ReplaceAll(s, "`.cursor/", "`content-ai/")
		s = strings.ReplaceAll(s, ".cursor/skills/", "content-ai/skills/")
		s = strings.ReplaceAll(s, ".cursor/rules/", "content-ai/rules/")
		return os.WriteFile(path, []byte(s), 0o644)
	})
}
