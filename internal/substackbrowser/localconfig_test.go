package substackbrowser

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadLocalConfigFlatLongKeys(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "config.json")
	body := `{"substack_publication_subdomain":"flatpub","markdown_table_mode":"list"}`
	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, found, err := LoadLocalConfig(p)
	if err != nil || !found {
		t.Fatal(err)
	}
	if cfg.Pub != "flatpub" || cfg.TableMode != "list" {
		t.Fatalf("%+v", cfg)
	}
}

func TestLoadLocalConfigFallbackLegacyDotSubstackFile(t *testing.T) {
	dir := t.TempDir()
	legacyDir := filepath.Join(dir, ".substack")
	if err := os.MkdirAll(legacyDir, 0o700); err != nil {
		t.Fatal(err)
	}
	legacyPath := filepath.Join(legacyDir, "config.json")
	if err := os.WriteFile(legacyPath, []byte(`{"pub":"fromlegacy"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	oldWd, wdErr := os.Getwd()
	if wdErr != nil {
		t.Fatal(wdErr)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	cfg, found, err := LoadLocalConfig("substack.json")
	if err != nil || !found {
		t.Fatalf("found=%v err=%v", found, err)
	}
	if cfg.Pub != "fromlegacy" {
		t.Fatalf("pub: %q", cfg.Pub)
	}
}

func TestLoadLocalConfigFallbackSubstackConfigWhenJSONMissing(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "substack.config"), []byte(`{"substack_browser":{"publication_subdomain":"fromoldfile"}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	cfg, found, err := LoadLocalConfig("substack.json")
	if err != nil || !found {
		t.Fatalf("found=%v err=%v", found, err)
	}
	if cfg.Pub != "fromoldfile" {
		t.Fatalf("pub: %q", cfg.Pub)
	}
}

func TestLoadLocalConfigLegacyKeys(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "config.json")
	old := `{"pub":"mytestpub","url":"https://example.com/edit"}`
	if err := os.WriteFile(p, []byte(old), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, found, err := LoadLocalConfig(p)
	if err != nil || !found {
		t.Fatalf("found=%v err=%v", found, err)
	}
	if cfg.Pub != "mytestpub" {
		t.Fatalf("pub: %q", cfg.Pub)
	}
	if cfg.URL != "https://example.com/edit" {
		t.Fatalf("url: %q", cfg.URL)
	}
}

func TestLoadLocalConfigGrouped(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "config.json")
	body := `{
	  "substack_browser": {"publication_subdomain": "gpub", "step_delay_milliseconds": 99},
	  "markdown_export": {"table_mode": "html"}
	}`
	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, found, err := LoadLocalConfig(p)
	if err != nil || !found {
		t.Fatalf("found=%v err=%v", found, err)
	}
	if cfg.Pub != "gpub" {
		t.Fatalf("pub: %q", cfg.Pub)
	}
	if cfg.NavigationDelayMS != 99 {
		t.Fatalf("delay: %d", cfg.NavigationDelayMS)
	}
	if cfg.TableMode != "html" {
		t.Fatalf("table: %q", cfg.TableMode)
	}
}

func TestLoadLocalConfigNewKeysWinOverLegacy(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "config.json")
	body := `{"pub":"old","substack_publication_subdomain":"newsub"}`
	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, found, err := LoadLocalConfig(p)
	if err != nil || !found {
		t.Fatal(err)
	}
	if cfg.Pub != "newsub" {
		t.Fatalf("want new key to win, got %q", cfg.Pub)
	}
}

func TestLoadLocalConfigFlatRenamedSchedulePushDeliveriesLegacyLongKey(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "config.json")
	body := `{"substack_publish_schedule_enable_email_and_app_delivery": false}`
	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, found, err := LoadLocalConfig(p)
	if err != nil || !found {
		t.Fatalf("found=%v err=%v", found, err)
	}
	if cfg.SchedulePushDeliveries == nil || *cfg.SchedulePushDeliveries != false {
		t.Fatalf("SchedulePushDeliveries: %+v", cfg.SchedulePushDeliveries)
	}
}

func TestLoadLocalConfigGroupedSchedulePushDeliveriesLegacyKey(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "config.json")
	body := `{"substack_publish": {"schedule_enable_email_and_app_delivery": false}}`
	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, found, err := LoadLocalConfig(p)
	if err != nil || !found {
		t.Fatalf("found=%v err=%v", found, err)
	}
	if cfg.SchedulePushDeliveries == nil || *cfg.SchedulePushDeliveries != false {
		t.Fatalf("SchedulePushDeliveries: %+v", cfg.SchedulePushDeliveries)
	}
}

func TestLoadLocalConfigGroupedSchedulePushDeliveriesNewKeyWins(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "config.json")
	body := `{"substack_publish": {
	  "schedule_push_deliveries": false,
	  "schedule_enable_email_and_app_delivery": true
	}}`
	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, found, err := LoadLocalConfig(p)
	if err != nil || !found {
		t.Fatalf("found=%v err=%v", found, err)
	}
	if cfg.SchedulePushDeliveries == nil || *cfg.SchedulePushDeliveries != false {
		t.Fatalf("want new key to win, got %+v", cfg.SchedulePushDeliveries)
	}
}

func TestLoadLocalConfigWithGlobalMergeGrouped(t *testing.T) {
	dir := t.TempDir()
	globalPath := filepath.Join(dir, "global.json")
	overlayPath := filepath.Join(dir, "overlay.json")
	globalBody := `{
	  "substack_browser": {"publication_subdomain": "englishpub", "step_delay_milliseconds": 42},
	  "markdown_export": {"table_mode": "list"}
	}`
	overlayBody := `{
	  "substack_browser": {"publication_subdomain": "spanishpub"}
	}`
	if err := os.WriteFile(globalPath, []byte(globalBody), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(overlayPath, []byte(overlayBody), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, found, err := LoadLocalConfigWithGlobal(globalPath, overlayPath)
	if err != nil || !found {
		t.Fatalf("found=%v err=%v", found, err)
	}
	if cfg.Pub != "spanishpub" {
		t.Fatalf("overlay pub: %q", cfg.Pub)
	}
	if cfg.NavigationDelayMS != 42 {
		t.Fatalf("global delay: %d", cfg.NavigationDelayMS)
	}
	if cfg.TableMode != "list" {
		t.Fatalf("global table: %q", cfg.TableMode)
	}
}

func TestLoadLocalConfigWithGlobalExplicitOverlayMissing(t *testing.T) {
	dir := t.TempDir()
	globalPath := filepath.Join(dir, "global.json")
	overlayPath := filepath.Join(dir, "nope.json")
	if err := os.WriteFile(globalPath, []byte(`{"substack_browser":{"publication_subdomain":"x"}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	_, _, err := LoadLocalConfigWithGlobal(globalPath, overlayPath)
	if err == nil {
		t.Fatal("expected error for missing explicit overlay path")
	}
}

func TestMergeJSONConfigBytesDeepMerge(t *testing.T) {
	base := []byte(`{"a":{"x":1},"b":2}`)
	over := []byte(`{"a":{"y":3},"b":4}`)
	got, err := mergeJSONConfigBytes(base, over)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(got, &m); err != nil {
		t.Fatal(err)
	}
	var inner map[string]int
	if err := json.Unmarshal(m["a"], &inner); err != nil {
		t.Fatal(err)
	}
	if inner["x"] != 1 || inner["y"] != 3 {
		t.Fatalf("deep merge inner: %+v", inner)
	}
	var b int
	if err := json.Unmarshal(m["b"], &b); err != nil || b != 4 {
		t.Fatalf("scalar overlay: %v err=%v", b, err)
	}
}

func TestEffectiveIncludeCognitiveMemeticsProjectAboutDefaultTrue(t *testing.T) {
	if !EffectiveIncludeCognitiveMemeticsProjectAbout(LocalConfig{}) {
		t.Fatal("nil pointer should mean true")
	}
	off := false
	if EffectiveIncludeCognitiveMemeticsProjectAbout(LocalConfig{IncludeCognitiveMemeticsProjectAbout: &off}) {
		t.Fatal("explicit false should mean false")
	}
	on := true
	if !EffectiveIncludeCognitiveMemeticsProjectAbout(LocalConfig{IncludeCognitiveMemeticsProjectAbout: &on}) {
		t.Fatal("explicit true should mean true")
	}
}
