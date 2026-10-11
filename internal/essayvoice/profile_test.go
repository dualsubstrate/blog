package essayvoice

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPatientNarratorShape(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	body := "---\nid: patient-narrator\nlabel: Patient Narrator\nsections: [\"human-condition\"]\nbanned_patterns:\n  - \"Look,\"\nrequired_verb_classes: [\"collapse\"]\nmax_staccato_run: 2\n---\n\nGenerator contract.\n"
	if err := os.WriteFile(filepath.Join(dir, "patient-narrator.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := Load(context.Background(), dir, "patient-narrator")
	if err != nil {
		t.Fatal(err)
	}
	if got.Profile.ID != "patient-narrator" || len(got.Profile.BannedPatterns) != 1 || got.Profile.RequiredVerbs[0] != "collapse" {
		t.Fatalf("profile = %#v", got.Profile)
	}
}

func TestAuditFileFlagsStaccato(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	post := filepath.Join(dir, "index.md")
	raw := "---\ntitle: Demo\n---\n\nLook, people hide the real fact today. Folks avoid the hard talk at night. We all just nod and stay quiet.\n"
	if err := os.WriteFile(post, []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}
	file, err := Load(context.Background(), writeVoice(t, dir), "patient-narrator")
	if err != nil {
		t.Fatal(err)
	}
	result, err := AuditFile(context.Background(), post, file.Profile)
	if err != nil {
		t.Fatal(err)
	}
	if result.Pass {
		t.Fatal("expected staccato and banned phrase to fail")
	}
}

func TestAuditFilePassesVariedProse(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	post := filepath.Join(dir, "index.md")
	raw := "---\ntitle: Demo\n---\n\nA single sentence cannot carry an entire mind. When you speak, you collapse a dense landscape of memories, private associations, and immediate mood into a flat sequence of words. The listener rebuilds that meaning from another history.\n"
	if err := os.WriteFile(post, []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}
	file, err := Load(context.Background(), writeVoice(t, dir), "patient-narrator")
	if err != nil {
		t.Fatal(err)
	}
	result, err := AuditFile(context.Background(), post, file.Profile)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Pass {
		t.Fatalf("expected pass, violations=%v", result.Violations)
	}
}

func writeVoice(t *testing.T, dir string) string {
	t.Helper()
	voices := filepath.Join(dir, "voices")
	if err := os.Mkdir(voices, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\nid: patient-narrator\nlabel: Patient Narrator\nsections: [\"human-condition\"]\nbanned_patterns:\n  - \"Look,\"\nrequired_verb_classes: [\"collapse\", \"rebuild\"]\nmax_staccato_run: 2\n---\n\nBrief.\n"
	if err := os.WriteFile(filepath.Join(voices, "patient-narrator.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return voices
}
