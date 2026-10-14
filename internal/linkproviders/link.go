package linkproviders

import (
	"fmt"
	"os"
	"path/filepath"
)

type linkSpec struct {
	display string
	target  string
}

// Run creates provider and skill symlinks for local development.
func Run(root string) error {
	root, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("root: %w", err)
	}
	providers := filepath.Join(root, "providers")
	if err := os.MkdirAll(providers, 0o755); err != nil {
		return fmt.Errorf("mkdir providers: %w", err)
	}

	providerLinks := []linkSpec{
		{display: "providers/content-pipelines", target: "../../../ai/content-pipelines"},
		{display: "providers/content-pipelines-mcp", target: "../../../ai/content-pipelines-mcp"},
	}
	for _, spec := range providerLinks {
		if err := ensureSymlink(filepath.Join(providers, filepath.Base(spec.display)), spec.target, spec.display); err != nil {
			return err
		}
	}

	skillsDir := filepath.Join(root, ".cursor", "skills")
	if err := os.MkdirAll(skillsDir, 0o755); err != nil {
		return fmt.Errorf("mkdir skills: %w", err)
	}
	skillLinks := []linkSpec{
		{display: ".cursor/skills/essay-command-center", target: "../../../../ai/content-pipelines-mcp/.cursor/skills/essay-command-center"},
	}
	for _, spec := range skillLinks {
		name := filepath.Base(spec.display)
		if err := ensureSymlink(filepath.Join(skillsDir, name), spec.target, spec.display); err != nil {
			return err
		}
	}
	return nil
}

func ensureSymlink(linkPath, target, display string) error {
	if err := os.MkdirAll(filepath.Dir(linkPath), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", linkPath, err)
	}

	if fi, err := os.Lstat(linkPath); err == nil {
		if fi.Mode()&os.ModeSymlink != 0 {
			cur, err := os.Readlink(linkPath)
			if err != nil {
				return fmt.Errorf("readlink %s: %w", linkPath, err)
			}
			if cur == target {
				fmt.Printf("%s -> %s (ok)\n", display, target)
				return nil
			}
			if err := os.Remove(linkPath); err != nil {
				return fmt.Errorf("remove %s: %w", linkPath, err)
			}
			if err := os.Symlink(target, linkPath); err != nil {
				return fmt.Errorf("symlink %s: %w", linkPath, err)
			}
			fmt.Printf("%s -> %s (retargeted)\n", display, target)
			return nil
		}
		return fmt.Errorf("%s exists and is not the expected symlink; remove it and retry", display)
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := os.Symlink(target, linkPath); err != nil {
		return fmt.Errorf("symlink %s -> %s: %w", linkPath, target, err)
	}
	fmt.Printf("%s -> %s\n", display, target)
	return nil
}
