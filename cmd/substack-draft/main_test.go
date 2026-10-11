package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dualsubstrate/blog/internal/substackbrowser"
	"github.com/dualsubstrate/blog/internal/substackhtml"
)

func writeTestSubstackMD(bundleDir string) error {
	body := "## Newsletter\n\nSidecar body for Substack paste.\n"
	return os.WriteFile(filepath.Join(bundleDir, "substack.md"), []byte(body), 0o644)
}

func TestBundleDirFromMarkdownPath(t *testing.T) {
	dir := t.TempDir()
	bundle := filepath.Join(dir, "human-condition", "2026-01-01-test")
	if err := os.MkdirAll(bundle, 0o755); err != nil {
		t.Fatal(err)
	}
	md := filepath.Join(bundle, "index.md")
	if err := os.WriteFile(md, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := bundleDirFromMarkdownPath(md)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Clean(got) != filepath.Clean(bundle) {
		t.Fatalf("got %q want %q", got, bundle)
	}
	got2, err := bundleDirFromMarkdownPath(bundle)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Clean(got2) != filepath.Clean(bundle) {
		t.Fatalf("dir: got %q want %q", got2, bundle)
	}
}

func TestSubtitleTypeCategoryLine(t *testing.T) {
	got := subtitleTypeCategoryLine([]string{"Mind-Infrastructure"}, 3, "video")
	want := "Type: mind-infrastructure:video"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestSubtitleTypeCategoryLineMultipleCategories(t *testing.T) {
	got := subtitleTypeCategoryLine([]string{"Mind-Infrastructure", "Human Condition"}, 2, "claims")
	want := "Type: mind-infrastructure:human-condition:claims"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestSubtitleTypeCategoryLineTypeOnlyTrailingWhenPresent(t *testing.T) {
	got := subtitleTypeCategoryLine([]string{"Social-Protocols"}, 3, "")
	want := "Type: social-protocols"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestBuildHTMLPrefersTypeLineWhenDescriptionAndCategoriesPresent(t *testing.T) {
	dir := t.TempDir()
	md := filepath.Join(dir, "index.md")
	fm := "---\ntitle: Cow\ndescription: \"Hello **world** tail\"\ncategories: [\"Cognitive-Memetics\", \"Cube-Cows\"]\ntype: panel\n---\n\nBody.\n"
	if err := os.WriteFile(md, []byte(fm), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := writeTestSubstackMD(dir); err != nil {
		t.Fatal(err)
	}
	lc := substackbrowser.LocalConfig{
		SubtitleIncludeCategories: true,
		SubtitleCategoriesMax:     3,
		SubtitleMaxChars:          200,
		IncludeFrontMatterLead:    false,
	}
	_, cfg, _, err := buildHTMLAndURL(buildOptions{
		Action: "paste", MDPath: md, HTMLPath: "", Fixture: true, Pub: "", TargetURL: "",
		Tables: "html", LocalCfg: lc, LocalCfgFound: true, NoDemote: false, TitleOverride: "", SubtitleOverride: "",
	})
	if err != nil {
		t.Fatal(err)
	}
	want := "Type: cognitive-memetics:cube-cows:panel"
	if cfg.Subtitle != want {
		t.Fatalf("subtitle: got %q want %q", cfg.Subtitle, want)
	}
}

func TestBuildHTMLUsesTypeLineWhenDescriptionEmpty(t *testing.T) {
	dir := t.TempDir()
	md := filepath.Join(dir, "index.md")
	fm := "---\ntitle: X\ncategories: [\"Mind-Infrastructure\"]\ntype: video\n---\n\nBody.\n"
	if err := os.WriteFile(md, []byte(fm), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := writeTestSubstackMD(dir); err != nil {
		t.Fatal(err)
	}
	lc := substackbrowser.LocalConfig{
		SubtitleIncludeCategories: true,
		SubtitleCategoriesMax:     3,
		SubtitleMaxChars:          200,
		IncludeFrontMatterLead:    false,
	}
	_, cfg, _, err := buildHTMLAndURL(buildOptions{
		Action: "paste", MDPath: md, HTMLPath: "", Fixture: true, Pub: "", TargetURL: "",
		Tables: "html", LocalCfg: lc, LocalCfgFound: true, NoDemote: false, TitleOverride: "", SubtitleOverride: "",
	})
	if err != nil {
		t.Fatal(err)
	}
	want := "Type: mind-infrastructure:video"
	if cfg.Subtitle != want {
		t.Fatalf("subtitle: got %q want %q", cfg.Subtitle, want)
	}
}

func TestScheduleDateTimeLocalFromYAMLDate(t *testing.T) {
	got := scheduleDateTimeLocalFromYAMLDate(`2026-05-01T01:00:00+11:00`)
	if got == "" {
		t.Fatal("expected non-empty datetime-local fragment")
	}
	if !strings.Contains(got, "T") || len(got) < 15 {
		t.Fatalf("unexpected format: %q", got)
	}
	_, err := time.ParseInLocation("2006-01-02T15:04", got, time.Local)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCategoryBrowseButtonURL(t *testing.T) {
	lc := substackbrowser.LocalConfig{SiteBaseURL: "https://blog.dualsubstrate.com/"}
	meta := substackhtml.FrontMatterMeta{Categories: []string{"Mind-Infrastructure"}}
	got := categoryBrowseButtonURL(lc, meta, "/repo/content/x/index.md")
	want := "https://blog.dualsubstrate.com/categories/mind-infrastructure/"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestCategoryBrowseButtonURLSpanish(t *testing.T) {
	lc := substackbrowser.LocalConfig{SiteBaseURL: "https://blog.dualsubstrate.com/"}
	meta := substackhtml.FrontMatterMeta{Categories: []string{"Mind-Infrastructure"}, Lang: "es"}
	got := categoryBrowseButtonURL(lc, meta, "")
	want := "https://blog.dualsubstrate.com/es/categories/mind-infrastructure/"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestCategoryBrowseButtonURLSpanishFromEsMdPath(t *testing.T) {
	lc := substackbrowser.LocalConfig{SiteBaseURL: "https://blog.dualsubstrate.com/"}
	meta := substackhtml.FrontMatterMeta{Categories: []string{"Mind-Infrastructure"}}
	got := categoryBrowseButtonURL(lc, meta, "/repo/content/x/index.es.md")
	want := "https://blog.dualsubstrate.com/es/categories/mind-infrastructure/"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestCategoryBrowseButtonURLSpanishBaseAlreadyUnderEs(t *testing.T) {
	lc := substackbrowser.LocalConfig{SiteBaseURL: "https://blog.dualsubstrate.com/es/"}
	meta := substackhtml.FrontMatterMeta{Categories: []string{"Mind-Infrastructure"}, Lang: "es"}
	got := categoryBrowseButtonURL(lc, meta, "")
	want := "https://blog.dualsubstrate.com/es/categories/mind-infrastructure/"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestScheduleSectionLabelUsesFirstCategory(t *testing.T) {
	meta := substackhtml.FrontMatterMeta{Categories: []string{"Mind-Infrastructure"}}
	got := scheduleSectionLabel(meta, "content/mind-infrastructure/x/index.md")
	if got != "Mind-Infrastructure" {
		t.Fatalf("got %q want %q", got, "Mind-Infrastructure")
	}
}

func TestScheduleSectionLabelSpanishUsesI18n(t *testing.T) {
	meta := substackhtml.FrontMatterMeta{Categories: []string{"Social-Protocols"}}
	got := scheduleSectionLabel(meta, "content/social-protocols/x/index.es.md")
	want := "Protocolos-sociales"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestScheduleSectionLabelSpanishCognitiveLane(t *testing.T) {
	meta := substackhtml.FrontMatterMeta{Categories: []string{"Cognitive-Memetics", "Por-Estas-Calles"}}
	p := "content/cognitive-memetics/sayings/2026-01-26-saying-03/index.es.md"
	got := scheduleSectionLabel(meta, p)
	// Lane name unchanged when i18n matches Hugo id.
	if got != "Por-Estas-Calles" {
		t.Fatalf("got %q want Por-Estas-Calles", got)
	}
}

func TestScheduleSectionLabelHumanConditionKeepsFirstWhenThemeSecond(t *testing.T) {
	meta := substackhtml.FrontMatterMeta{Categories: []string{"Human-Condition", "Mental-Processes"}}
	got := scheduleSectionLabel(meta, "content/human-condition/x/index.md")
	if got != "Human-Condition" {
		t.Fatalf("got %q want %q (pillar, not theme hub)", got, "Human-Condition")
	}
}

func TestScheduleSectionLabelCognitiveMemeticsUsesSecondCategory(t *testing.T) {
	meta := substackhtml.FrontMatterMeta{Categories: []string{"Cognitive-Memetics", "Cube-Cows"}}
	got := scheduleSectionLabel(meta, "content/mind-infrastructure/x/index.md")
	if got != "Cube-Cows" {
		t.Fatalf("got %q want %q", got, "Cube-Cows")
	}
}

func TestScheduleSectionLabelCognitivePathUsesSecondWhenUmbrellaFirst(t *testing.T) {
	meta := substackhtml.FrontMatterMeta{Categories: []string{"Cognitive-Memetics", "Por-Estas-Calles"}}
	p := "content/cognitive-memetics/sayings/2026-01-26-saying-03/index.md"
	got := scheduleSectionLabel(meta, p)
	if got != "Por-Estas-Calles" {
		t.Fatalf("got %q want %q", got, "Por-Estas-Calles")
	}
}

func TestScheduleSectionLabelCognitiveMemeticsFallsBackWhenNoSecond(t *testing.T) {
	meta := substackhtml.FrontMatterMeta{Categories: []string{"Cognitive-Memetics"}}
	got := scheduleSectionLabel(meta, "content/social-protocols/x/index.md")
	if got != "Cognitive-Memetics" {
		t.Fatalf("got %q want %q", got, "Cognitive-Memetics")
	}
}

func TestScheduleSectionLabelCognitivePathSingleCategory(t *testing.T) {
	meta := substackhtml.FrontMatterMeta{Categories: []string{"Reptilocracy"}}
	p := "content/cognitive-memetics/reptilocracy/2026-04-12-not-in-our-term/index.md"
	got := scheduleSectionLabel(meta, p)
	if got != "Reptilocracy" {
		t.Fatalf("got %q want %q", got, "Reptilocracy")
	}
}

func TestBuildHTMLIncludesCognitiveMemeticsProjectAbout(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "content", "cognitive-memetics", "cows", "2026-02-26-cow-w01")
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatal(err)
	}
	md := filepath.Join(subdir, "index.md")
	fm := "---\ntitle: Cow\ncategories: [\"Cognitive-Memetics\", \"Cube-Cows\"]\ntype: panel\n---\n\nBody.\n"
	if err := os.WriteFile(md, []byte(fm), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := writeTestSubstackMD(subdir); err != nil {
		t.Fatal(err)
	}
	lc := substackbrowser.LocalConfig{
		IncludeFrontMatterLead: false,
	}
	html, _, _, err := buildHTMLAndURL(buildOptions{
		Action: "paste", MDPath: md, HTMLPath: "", Fixture: true, Pub: "", TargetURL: "",
		Tables: "html", LocalCfg: lc, LocalCfgFound: true, NoDemote: false, TitleOverride: "", SubtitleOverride: "",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "Tales from the Cube Farm") {
		t.Fatalf("expected cube-cows project about in HTML, got: %s", html)
	}
}

func TestBuildHTMLSkipsCognitiveMemeticsProjectAboutWhenDisabled(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "content", "cognitive-memetics", "cows", "2026-02-26-cow-w01")
	if err := os.MkdirAll(subdir, 0o755); err != nil {
		t.Fatal(err)
	}
	md := filepath.Join(subdir, "index.md")
	fm := "---\ntitle: Cow\ncategories: [\"Cognitive-Memetics\", \"Cube-Cows\"]\ntype: panel\n---\n\nBody.\n"
	if err := os.WriteFile(md, []byte(fm), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := writeTestSubstackMD(subdir); err != nil {
		t.Fatal(err)
	}
	off := false
	lc := substackbrowser.LocalConfig{
		IncludeFrontMatterLead:               false,
		IncludeCognitiveMemeticsProjectAbout: &off,
	}
	html, _, _, err := buildHTMLAndURL(buildOptions{
		Action: "paste", MDPath: md, HTMLPath: "", Fixture: true, Pub: "", TargetURL: "",
		Tables: "html", LocalCfg: lc, LocalCfgFound: true, NoDemote: false, TitleOverride: "", SubtitleOverride: "",
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(html, "Tales from the Cube Farm") {
		t.Fatalf("did not expect project about when disabled, got: %s", html)
	}
}

func TestAppendFooterHTMLUsesEsPrefixForSpanish(t *testing.T) {
	meta := substackhtml.FrontMatterMeta{
		Categories: []string{"Mind-Infrastructure"},
		Tags:       []string{"FreeEnergyPrinciple"},
	}
	lc := substackbrowser.LocalConfig{
		SiteBaseURL:               "https://blog.dualsubstrate.com/",
		FooterIncludeSiteLink:     true,
		FooterIncludeCategoryLink: true,
		FooterIncludeTags:         true,
		FooterCategoryLinkIndex:   0,
	}
	md := "/repo/content/mind-infrastructure/2026-04-29-free-energy-principle-hallucination-machine/index.es.md"
	got := appendFooterHTML("<p>x</p>", md, meta, lc)
	if !strings.Contains(got, `href="https://blog.dualsubstrate.com/es/mind-infrastructure/2026-04-29-free-energy-principle-hallucination-machine/"`) {
		t.Fatalf("expected Spanish post URL in footer, got: %s", got)
	}
	if !strings.Contains(got, `href="https://blog.dualsubstrate.com/es/categories/mind-infrastructure/"`) {
		t.Fatalf("expected Spanish category URL in footer, got: %s", got)
	}
}

func TestFormatSubstackScheduleTextField(t *testing.T) {
	tt := time.Date(2026, 4, 29, 8, 40, 0, 0, time.Local)
	if g := formatSubstackScheduleTextField(tt); g != "29/04/2026, 08:40 am" {
		t.Fatalf("got %q", g)
	}
	tt2 := time.Date(2026, 4, 29, 13, 5, 0, 0, time.Local)
	if g := formatSubstackScheduleTextField(tt2); g != "29/04/2026, 01:05 pm" {
		t.Fatalf("got %q", g)
	}
}
