package socialautopost

import (
	"fmt"
	"path/filepath"

	"github.com/dualsubstrate/blog/internal/substackpublishstate"
)

// RecordPublished appends or updates targetKey in the bundle's social-published file.
// No-op when dryRun or noMark is true.
func RecordPublished(bundleDir, targetKey string, dryRun, noMark bool) error {
	if dryRun || noMark {
		return nil
	}
	if err := substackpublishstate.MarkPublished(bundleDir, substackpublishstate.DefaultMarker, targetKey); err != nil {
		return fmt.Errorf("social-published: %w", err)
	}
	return nil
}

// MarkerHasTarget reports whether bundleDir/social-published records targetKey.
func MarkerHasTarget(bundleDir, targetKey string) (bool, error) {
	p := filepath.Join(bundleDir, substackpublishstate.DefaultMarker)
	return substackpublishstate.MarkerHasTarget(p, targetKey)
}

// PublishTargetForNetwork maps autopost network label to social-published target id.
func PublishTargetForNetwork(network string) string {
	switch network {
	case "Facebook":
		return substackpublishstate.TargetFacebook
	default:
		return substackpublishstate.TargetLinkedIn
	}
}
