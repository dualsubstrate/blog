package substackbrowser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadLocalConfigEnvOverridesPublicationSubdomain(t *testing.T) {
	t.Setenv("SUBSTACK_PUBLICATION_SUBDOMAIN", "fromenv")
	dir := t.TempDir()
	p := filepath.Join(dir, "c.json")
	body := `{"substack_browser":{"publication_subdomain":"infile","step_delay_milliseconds":3}}`
	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, found, err := LoadLocalConfig(p)
	if err != nil || !found {
		t.Fatalf("found=%v err=%v", found, err)
	}
	if cfg.Pub != "fromenv" {
		t.Fatalf("pub: %q", cfg.Pub)
	}
	if cfg.NavigationDelayMS != 3 {
		t.Fatalf("delay: %d", cfg.NavigationDelayMS)
	}
}

func TestLoadLocalConfigEnvInvalidBoolIgnored(t *testing.T) {
	t.Setenv("SUBSTACK_PUBLISH_SCHEDULE_DEBUG_DOM_ON_FAILURE", "maybe")
	dir := t.TempDir()
	p := filepath.Join(dir, "c.json")
	body := `{"substack_publish_schedule_debug_dom_on_failure": true}`
	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, found, err := LoadLocalConfig(p)
	if err != nil || !found {
		t.Fatalf("found=%v err=%v", found, err)
	}
	if !cfg.ScheduleDebugDOMOnFailure {
		t.Fatal("expected JSON true to remain when env bool invalid")
	}
}

func TestLoadLocalConfigMissingFileStillAppliesEnvScheduleDebugDom(t *testing.T) {
	t.Setenv("SUBSTACK_PUBLISH_SCHEDULE_DEBUG_DOM_ON_FAILURE", "1")
	dir := t.TempDir()
	p := filepath.Join(dir, "no-such-substack.json")
	cfg, found, err := LoadLocalConfig(p)
	if err != nil {
		t.Fatal(err)
	}
	if found {
		t.Fatal("expected found=false when config file is missing")
	}
	if !cfg.ScheduleDebugDOMOnFailure {
		t.Fatal("expected SUBSTACK_PUBLISH_SCHEDULE_DEBUG_DOM_ON_FAILURE=1 to apply even without a JSON file")
	}
}

func TestLoadLocalConfigWithGlobalEnvOverridesOverlay(t *testing.T) {
	t.Setenv("SUBSTACK_PUBLICATION_SUBDOMAIN", "envwins")
	dir := t.TempDir()
	globalPath := filepath.Join(dir, "g.json")
	overlayPath := filepath.Join(dir, "o.json")
	if err := os.WriteFile(globalPath, []byte(`{"substack_browser":{"publication_subdomain":"g","step_delay_milliseconds":9}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(overlayPath, []byte(`{"substack_browser":{"publication_subdomain":"o"}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, found, err := LoadLocalConfigWithGlobal(globalPath, overlayPath)
	if err != nil || !found {
		t.Fatalf("found=%v err=%v", found, err)
	}
	if cfg.Pub != "envwins" {
		t.Fatalf("pub: %q", cfg.Pub)
	}
	if cfg.NavigationDelayMS != 9 {
		t.Fatalf("delay: %d", cfg.NavigationDelayMS)
	}
}
