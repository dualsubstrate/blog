package essayvoice

import (
	"context"
	"fmt"

	"github.com/behaviorengineering/strop/pkg/evaluation/voice"
)

// AuditFile runs the deterministic voice audit on a Hugo post body.
func AuditFile(ctx context.Context, postPath string, profile voice.Profile) (voice.AuditResult, error) {
	if err := ctx.Err(); err != nil {
		return voice.AuditResult{}, err
	}
	prose, err := BodyProse(ctx, postPath)
	if err != nil {
		return voice.AuditResult{}, err
	}
	if stringsEmpty(prose) {
		return voice.AuditResult{}, fmt.Errorf("post %s has no body prose", postPath)
	}
	return voice.HeuristicAudit(profile, prose), nil
}

func stringsEmpty(s string) bool {
	for _, r := range s {
		if r != ' ' && r != '\n' && r != '\t' && r != '\r' {
			return false
		}
	}
	return true
}
