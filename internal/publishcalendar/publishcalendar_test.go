package publishcalendar_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/dualsubstrate/blog/internal/publishcalendar"
	"github.com/dualsubstrate/blog/internal/substackpublishstate"
)

func TestBuildBundleWithMarker(t *testing.T) {
	root := t.TempDir()
	content := filepath.Join(root, "content")
	bundle := filepath.Join(content, "human-condition", "2026-05-28-test-post")
	if err := os.MkdirAll(bundle, 0o755); err != nil {
		t.Fatal(err)
	}
	index := `---
title: Test Post
date: '2026-05-28T01:00:00+11:00'
draft: false
type: claims
---
Body.
`
	if err := os.WriteFile(filepath.Join(bundle, "index.md"), []byte(index), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bundle, "substack.md"), []byte("Newsletter body\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	marker := "# <target> <RFC3339-UTC>\n" + substackpublishstate.TargetSubstackEN + " 2026-05-28T10:00:00Z\n"
	if err := os.WriteFile(filepath.Join(bundle, substackpublishstate.DefaultMarker), []byte(marker), 0o644); err != nil {
		t.Fatal(err)
	}

	cal, err := publishcalendar.Build(content, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(cal.Entries) != 1 {
		t.Fatalf("entries: got %d want 1", len(cal.Entries))
	}
	e := cal.Entries[0]
	if e.Planned != "2026-05-28" {
		t.Fatalf("planned: got %q", e.Planned)
	}
	if e.Channels[substackpublishstate.TargetSiteEN].Status != publishcalendar.StatusPresent {
		t.Fatalf("site-en: %+v", e.Channels[substackpublishstate.TargetSiteEN])
	}
	if e.Channels[substackpublishstate.TargetSiteES].Status != publishcalendar.StatusMissing {
		t.Fatalf("site-es: %+v", e.Channels[substackpublishstate.TargetSiteES])
	}
	if e.Channels[substackpublishstate.TargetLinkedIn].Status != publishcalendar.StatusMissing {
		t.Fatalf("linkedin: %+v", e.Channels[substackpublishstate.TargetLinkedIn])
	}
	if e.Channels[substackpublishstate.TargetSubstackEN].Status != publishcalendar.StatusPresent {
		t.Fatalf("substack-en: %+v", e.Channels[substackpublishstate.TargetSubstackEN])
	}
}

func TestBuildPanelFrontMatterOnly(t *testing.T) {
	root := t.TempDir()
	content := filepath.Join(root, "content")
	bundle := filepath.Join(content, "cognitive-memetics", "cows", "2026-06-10-cow-w16")
	if err := os.MkdirAll(bundle, 0o755); err != nil {
		t.Fatal(err)
	}
	index := `---
title: Collaboration days
date: '2026-06-10T01:00:00+11:00'
draft: false
type: panel
description: "Teaser in front matter only."
---
`
	indexES := `---
title: Días de colaboración
date: '2026-06-10T01:00:00+11:00'
draft: false
type: panel
description: "Teaser en front matter solamente."
---
`
	if err := os.WriteFile(filepath.Join(bundle, "index.md"), []byte(index), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bundle, "index.es.md"), []byte(indexES), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bundle, "linkedin.txt"), []byte("LinkedIn copy\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bundle, "facebook-es.txt"), []byte("Facebook copy\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cal, err := publishcalendar.Build(content, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(cal.Entries) != 1 {
		t.Fatalf("entries: got %d want 1", len(cal.Entries))
	}
	e := cal.Entries[0]
	for _, key := range []string{
		substackpublishstate.TargetSiteEN,
		substackpublishstate.TargetSiteES,
		substackpublishstate.TargetLinkedIn,
		substackpublishstate.TargetFacebook,
	} {
		if e.Channels[key].Status != publishcalendar.StatusPresent {
			t.Fatalf("%s: got %+v want present", key, e.Channels[key])
		}
	}
	if e.Channels[substackpublishstate.TargetSubstackEN].Status != publishcalendar.StatusMissing {
		t.Fatalf("substack-en: %+v", e.Channels[substackpublishstate.TargetSubstackEN])
	}
}

func TestWriteJSON(t *testing.T) {
	root := t.TempDir()
	out := filepath.Join(root, "calendar.json")
	cal := &publishcalendar.Calendar{
		GeneratedAt: "2026-05-30T00:00:00Z",
		Channels:    []string{"substack-en"},
		Entries:     []publishcalendar.Entry{},
	}
	if err := publishcalendar.WriteJSON(out, cal); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	var decoded publishcalendar.Calendar
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.GeneratedAt != cal.GeneratedAt {
		t.Fatalf("generatedAt: got %q", decoded.GeneratedAt)
	}
}
