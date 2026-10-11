package substackpublishstate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBundleReadyForPublishDraftAndFutureDateAllowed(t *testing.T) {
	root := t.TempDir()
	b := filepath.Join(root, "bundle")
	if err := os.MkdirAll(b, 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(body string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(b, "index.md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("---\ntitle: T\ndraft: true\n---\n")
	ok, err := BundleReadyForPublish(b, LegacyDefaultTarget)
	if err != nil || ok {
		t.Fatalf("draft: ok=%v err=%v", ok, err)
	}
	substack := filepath.Join(b, "substack.md")
	writeSubstack := func() {
		t.Helper()
		if err := os.WriteFile(substack, []byte("## Post\n\nNewsletter body.\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("---\ntitle: T\ndate: 2099-01-01\n---\n")
	writeSubstack()
	ok, err = BundleReadyForPublish(b, LegacyDefaultTarget)
	if err != nil || !ok {
		t.Fatalf("future date should still be ready: ok=%v err=%v", ok, err)
	}
	write("---\ntitle: T\ndate: 2026-04-01\n---\n")
	ok, err = BundleReadyForPublish(b, LegacyDefaultTarget)
	if err != nil || !ok {
		t.Fatalf("past date: ok=%v err=%v", ok, err)
	}
	write("---\ntitle: T\n---\n")
	if err := os.Remove(substack); err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	ok, err = BundleReadyForPublish(b, LegacyDefaultTarget)
	if err != nil || ok {
		t.Fatalf("without substack.md: ok=%v err=%v", ok, err)
	}
	writeSubstack()
	ok, err = BundleReadyForPublish(b, LegacyDefaultTarget)
	if err != nil || !ok {
		t.Fatalf("with substack.md: ok=%v err=%v", ok, err)
	}
}

func TestBundleReadyForPublishSubstackESRequiresSpanishFile(t *testing.T) {
	root := t.TempDir()
	b := filepath.Join(root, "bundle")
	if err := os.MkdirAll(b, 0o755); err != nil {
		t.Fatal(err)
	}
	index := "---\ntitle: T\n---\nbody\n"
	if err := os.WriteFile(filepath.Join(b, "index.md"), []byte(index), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err := BundleReadyForPublish(b, TargetSubstackES)
	if err != nil || ok {
		t.Fatalf("without index.es.md: ok=%v err=%v", ok, err)
	}
	if err := os.WriteFile(filepath.Join(b, "index.es.md"), []byte("   \n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err = BundleReadyForPublish(b, TargetSubstackES)
	if err != nil || ok {
		t.Fatalf("whitespace-only index.es.md: ok=%v err=%v", ok, err)
	}
	if err := os.WriteFile(filepath.Join(b, "index.es.md"), []byte("---\n---\nhola\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err = BundleReadyForPublish(b, TargetSubstackES)
	if err != nil || ok {
		t.Fatalf("without substack.es.md: ok=%v err=%v", ok, err)
	}
	if err := os.WriteFile(filepath.Join(b, "substack.es.md"), []byte("## Hola\n\nCuerpo.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err = BundleReadyForPublish(b, TargetSubstackES)
	if err != nil || !ok {
		t.Fatalf("with index.es.md and substack.es.md: ok=%v err=%v", ok, err)
	}
}

func TestBundleReadyForPublishSubstackESRespectsSpanishDraft(t *testing.T) {
	root := t.TempDir()
	b := filepath.Join(root, "bundle")
	if err := os.MkdirAll(b, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(b, "index.md"), []byte("---\ntitle: T\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(b, "index.es.md"), []byte("---\ntitle: T\ndraft: true\n---\nhola\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err := BundleReadyForPublish(b, TargetSubstackES)
	if err != nil || ok {
		t.Fatalf("spanish draft should block ES: ok=%v err=%v", ok, err)
	}
	if err := os.WriteFile(filepath.Join(b, "substack.md"), []byte("## EN\n\nBody.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err = BundleReadyForPublish(b, LegacyDefaultTarget)
	if err != nil || !ok {
		t.Fatalf("english target ignores spanish draft: ok=%v err=%v", ok, err)
	}
}

func TestBundleReadyForPublishLinkedInRequiresFile(t *testing.T) {
	root := t.TempDir()
	b := filepath.Join(root, "bundle")
	if err := os.MkdirAll(b, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(b, "index.md"), []byte("---\ntitle: T\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err := BundleReadyForPublish(b, TargetLinkedIn)
	if err != nil || ok {
		t.Fatalf("no linkedin.txt: ok=%v err=%v", ok, err)
	}
	if err := os.WriteFile(filepath.Join(b, "linkedin.txt"), []byte("\n\t\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err = BundleReadyForPublish(b, TargetLinkedIn)
	if err != nil || ok {
		t.Fatalf("empty linkedin.txt: ok=%v err=%v", ok, err)
	}
	if err := os.WriteFile(filepath.Join(b, "linkedin.txt"), []byte("Hello\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err = BundleReadyForPublish(b, TargetLinkedIn)
	if err != nil || !ok {
		t.Fatalf("with linkedin.txt: ok=%v err=%v", ok, err)
	}
}

func TestBundleReadyForPublishSiteENRequiresBody(t *testing.T) {
	root := t.TempDir()
	b := filepath.Join(root, "bundle")
	if err := os.MkdirAll(b, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(b, "index.md"), []byte("---\ntitle: T\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err := BundleReadyForPublish(b, TargetSiteEN)
	if err != nil || ok {
		t.Fatalf("empty body: ok=%v err=%v", ok, err)
	}
	if err := os.WriteFile(filepath.Join(b, "index.md"), []byte("---\ntitle: T\n---\nBody.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err = BundleReadyForPublish(b, TargetSiteEN)
	if err != nil || !ok {
		t.Fatalf("with body: ok=%v err=%v", ok, err)
	}
}

func TestBundleReadyForPublishSiteESRequiresSpanishPage(t *testing.T) {
	root := t.TempDir()
	b := filepath.Join(root, "bundle")
	if err := os.MkdirAll(b, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(b, "index.md"), []byte("---\ntitle: T\n---\nBody.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err := BundleReadyForPublish(b, TargetSiteES)
	if err != nil || ok {
		t.Fatalf("without index.es.md: ok=%v err=%v", ok, err)
	}
	if err := os.WriteFile(filepath.Join(b, "index.es.md"), []byte("---\ntitle: T\n---\nHola.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err = BundleReadyForPublish(b, TargetSiteES)
	if err != nil || !ok {
		t.Fatalf("with index.es.md: ok=%v err=%v", ok, err)
	}
}

func TestBundleReadyForPublishFacebookRequiresSpanishCopy(t *testing.T) {
	root := t.TempDir()
	b := filepath.Join(root, "bundle")
	if err := os.MkdirAll(b, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(b, "index.md"), []byte("---\ntitle: T\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err := BundleReadyForPublish(b, TargetFacebook)
	if err != nil || ok {
		t.Fatalf("no facebook-es.txt: ok=%v err=%v", ok, err)
	}
	if err := os.WriteFile(filepath.Join(b, "facebook-es.txt"), []byte("Hola amigos\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err = BundleReadyForPublish(b, TargetFacebook)
	if err != nil || !ok {
		t.Fatalf("with facebook-es.txt: ok=%v err=%v", ok, err)
	}
}

func TestListUnpublishedExcludesNotReady(t *testing.T) {
	root := t.TempDir()
	content := filepath.Join(root, "content")
	sec := filepath.Join(content, "human-condition")
	for _, name := range []string{"draft-post", "scheduled-post", "ready-post"} {
		if err := os.MkdirAll(filepath.Join(sec, name), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(sec, "draft-post", "index.md"), []byte("---\ntitle: D\ndraft: true\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sec, "scheduled-post", "index.md"), []byte("---\ntitle: F\ndate: 2099-01-01\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sec, "ready-post", "index.md"), []byte("---\ntitle: R\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sec, "ready-post", "substack.md"), []byte("## Ready\n\nBody.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sec, "scheduled-post", "substack.md"), []byte("## Scheduled\n\nBody.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	list, err := ListUnpublished(content, DefaultMarker, LegacyDefaultTarget)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Fatalf("got %v want 2 entries", list)
	}
	if string(list[0]) != "human-condition/ready-post" || string(list[1]) != "human-condition/scheduled-post" {
		t.Fatalf("got %v want [human-condition/ready-post human-condition/scheduled-post]", list)
	}
}
