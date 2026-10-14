package filetree

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// MirrorDir copies src into dst, replacing dst contents (rsync -a --delete semantics for files).
func MirrorDir(src, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	if err := os.MkdirAll(dst, 0o755); err != nil {
		return fmt.Errorf("mkdir dst: %w", err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("read src: %w", err)
	}

	keep := make(map[string]struct{}, len(entries))
	for _, ent := range entries {
		keep[ent.Name()] = struct{}{}
		name := ent.Name()
		from := filepath.Join(src, name)
		to := filepath.Join(dst, name)
		if ent.IsDir() {
			if err := os.RemoveAll(to); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("remove old dir %s: %w", to, err)
			}
			if err := copyDir(from, to); err != nil {
				return err
			}
			continue
		}
		if err := copyFile(from, to); err != nil {
			return err
		}
	}

	dstEntries, err := os.ReadDir(dst)
	if err != nil {
		return fmt.Errorf("read dst: %w", err)
	}
	for _, ent := range dstEntries {
		if _, ok := keep[ent.Name()]; ok {
			continue
		}
		p := filepath.Join(dst, ent.Name())
		if err := os.RemoveAll(p); err != nil {
			return fmt.Errorf("remove stale %s: %w", p, err)
		}
	}
	return nil
}

func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("readdir: %w", err)
	}
	for _, ent := range entries {
		from := filepath.Join(src, ent.Name())
		to := filepath.Join(dst, ent.Name())
		if ent.IsDir() {
			if err := copyDir(from, to); err != nil {
				return err
			}
			continue
		}
		if err := copyFile(from, to); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open src: %w", err)
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("mkdir parent: %w", err)
	}
	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create dst: %w", err)
	}
	defer func() {
		_ = out.Close()
	}()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("copy: %w", err)
	}
	if err := out.Close(); err != nil {
		return fmt.Errorf("close dst: %w", err)
	}
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("stat src: %w", err)
	}
	return os.Chmod(dst, info.Mode().Perm())
}
