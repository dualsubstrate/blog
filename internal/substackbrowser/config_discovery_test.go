package substackbrowser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveSubstackConfigPath_prefersUserConfig(t *testing.T) {
	root := t.TempDir()
	userDir := filepath.Join(root, "behaviour-engineering")
	require.NoError(t, os.MkdirAll(userDir, 0o700))
	userFile := filepath.Join(userDir, "substack.json")
	require.NoError(t, os.WriteFile(userFile, []byte("{}"), 0o600))
	t.Setenv("XDG_CONFIG_HOME", root)

	path, err := ResolveSubstackConfigPath("")
	require.NoError(t, err)
	assert.Equal(t, userFile, path)
}

func TestResolveSubstackConfigPath_SUBSTACK_CONFIGOverride(t *testing.T) {
	custom := filepath.Join(t.TempDir(), "custom.json")
	require.NoError(t, os.WriteFile(custom, []byte("{}"), 0o600))
	t.Setenv("SUBSTACK_CONFIG", custom)

	path, err := ResolveSubstackConfigPath("")
	require.NoError(t, err)
	assert.Equal(t, custom, path)
}

func TestDefaultLocalConfigPath_fallsBackToRepoName(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	assert.Equal(t, "substack.json", DefaultLocalConfigPath())
}
