package socialautopost

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dualsubstrate/blog/internal/substackpublishstate"
)

func TestRecordPublished_andMarkerHasTarget(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.md"), []byte("---\ndraft: false\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "linkedin.txt"), []byte("https://blog.dualsubstrate.com/x/\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := RecordPublished(dir, substackpublishstate.TargetLinkedIn, false, false); err != nil {
		t.Fatal(err)
	}
	has, err := MarkerHasTarget(dir, substackpublishstate.TargetLinkedIn)
	if err != nil || !has {
		t.Fatalf("has=%v err=%v", has, err)
	}
	if err := RecordPublished(dir, substackpublishstate.TargetLinkedIn, true, false); err != nil {
		t.Fatal(err)
	}
	if err := RecordPublished(dir, substackpublishstate.TargetFacebook, false, true); err != nil {
		t.Fatal(err)
	}
	hasFB, err := MarkerHasTarget(dir, substackpublishstate.TargetFacebook)
	if err != nil || hasFB {
		t.Fatalf("facebook should not be marked: has=%v err=%v", hasFB, err)
	}
}

func TestPublishTargetForNetwork(t *testing.T) {
	if got := PublishTargetForNetwork("LinkedIn"); got != substackpublishstate.TargetLinkedIn {
		t.Fatalf("LinkedIn: %q", got)
	}
	if got := PublishTargetForNetwork("Facebook"); got != substackpublishstate.TargetFacebook {
		t.Fatalf("Facebook: %q", got)
	}
}
