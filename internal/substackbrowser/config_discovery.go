package substackbrowser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/behaviorengineering/operatorconfig/pkg/operatorconfig"
)

const substackConfigApp = "behaviour-engineering"

func substackOperatorConfigOptions(configFlagPath string) operatorconfig.Options {
	opts := operatorconfig.Options{
		App:        substackConfigApp,
		ConfigEnv:  "SUBSTACK_CONFIG",
		Filename:   "substack.json",
		ExtraPaths: substackConfigExtraPaths(),
	}
	opts.ConfigFlagPath = strings.TrimSpace(configFlagPath)
	return opts
}

func substackConfigExtraPaths() []string {
	return []string{
		"substack.json",
		"substack.config",
		filepath.Join(".substack", "config.json"),
		filepath.Join(".substack", "substack.json"),
	}
}

// ResolveSubstackConfigPath discovers the Substack JSON config file.
// Order: SUBSTACK_CONFIG, ~/.config/behaviour-engineering/substack.json, then repo fallbacks.
func ResolveSubstackConfigPath(configFlagPath string) (string, error) {
	return operatorconfig.ResolveConfigPath(substackOperatorConfigOptions(configFlagPath))
}

// DefaultLocalConfigPath returns the discovered config path when present, else substack.json.
func DefaultLocalConfigPath() string {
	path, err := ResolveSubstackConfigPath("")
	if err == nil && strings.TrimSpace(path) != "" {
		return path
	}
	return "substack.json"
}

// InitUserSubstackConfig writes ~/.config/behaviour-engineering/substack.json when missing.
func InitUserSubstackConfig(example []byte, force bool) (bool, error) {
	return operatorconfig.InitUserConfig(substackOperatorConfigOptions(""), example, force)
}

// InitRepoSubstackConfig creates ./substack.json from the example when missing (repo-local workflow).
func InitRepoSubstackConfig(example []byte) error {
	const repoPath = "substack.json"
	if _, err := os.Stat(repoPath); err == nil {
		return nil
	}
	if err := os.WriteFile(repoPath, example, 0o600); err != nil {
		return fmt.Errorf("write %s: %w", repoPath, err)
	}
	return nil
}
