package substackbrowser

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// mergeJSONConfigBytes deep-merges two JSON objects at the root. When both values
// for a key are JSON objects, merge recurses; otherwise the overlay value wins
// (including arrays and scalars). Use the same top-level shape in both files
// (grouped sections vs flat long keys) so unmarshaling sees a coherent document.
func mergeJSONConfigBytes(base, overlay []byte) ([]byte, error) {
	var bm, om map[string]json.RawMessage
	if err := json.Unmarshal(base, &bm); err != nil {
		return nil, fmt.Errorf("substackbrowser: merge config-global: parse base: %w", err)
	}
	if err := json.Unmarshal(overlay, &om); err != nil {
		return nil, fmt.Errorf("substackbrowser: merge config: parse overlay: %w", err)
	}
	merged := deepMergeRawObjects(bm, om)
	out, err := json.Marshal(merged)
	if err != nil {
		return nil, fmt.Errorf("substackbrowser: merge config: marshal: %w", err)
	}
	return out, nil
}

func deepMergeRawObjects(base, overlay map[string]json.RawMessage) map[string]json.RawMessage {
	out := make(map[string]json.RawMessage, len(base)+len(overlay))
	for k, v := range base {
		out[k] = v
	}
	for k, vb := range overlay {
		va, ok := out[k]
		if !ok || !isJSONObjectRaw(va) || !isJSONObjectRaw(vb) {
			out[k] = vb
			continue
		}
		var innerA, innerB map[string]json.RawMessage
		if err := json.Unmarshal(va, &innerA); err != nil {
			out[k] = vb
			continue
		}
		if err := json.Unmarshal(vb, &innerB); err != nil {
			out[k] = vb
			continue
		}
		mergedInner := deepMergeRawObjects(innerA, innerB)
		b, err := json.Marshal(mergedInner)
		if err != nil {
			out[k] = vb
			continue
		}
		out[k] = json.RawMessage(b)
	}
	return out
}

func isJSONObjectRaw(raw json.RawMessage) bool {
	if len(strings.TrimSpace(string(raw))) == 0 {
		return false
	}
	var m map[string]json.RawMessage
	return json.Unmarshal(raw, &m) == nil
}

// readLocalConfigFileBytes reads path, or when path is the default substack.json and it
// is missing, tries substack.config (previous default), then .substack/config.json,
// then .substack/substack.json for migration.
func readLocalConfigFileBytes(path string) ([]byte, error) {
	path = strings.TrimSpace(path)
	b, err := os.ReadFile(path)
	if err == nil {
		return b, nil
	}
	if !os.IsNotExist(err) {
		return nil, err
	}
	if filepath.Clean(path) != "substack.json" {
		return nil, err
	}
	for _, alt := range []string{"substack.config", filepath.Join(".substack", "config.json"), filepath.Join(".substack", "substack.json")} {
		lb, lerr := os.ReadFile(alt)
		if lerr == nil {
			return lb, nil
		}
		if !os.IsNotExist(lerr) {
			return nil, fmt.Errorf("substackbrowser: read legacy config %q: %w", alt, lerr)
		}
	}
	return nil, err
}

// decodeLocalConfigFromBytes parses grouped or flat local config JSON after any merge.
func decodeLocalConfigFromBytes(b []byte) (LocalConfig, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return LocalConfig{}, fmt.Errorf("substackbrowser: parse config: %w", err)
	}
	if rootJSONUsesGroupedSections(raw) {
		cfg, err := decodeGroupedLocalConfig(b)
		if err != nil {
			return LocalConfig{}, err
		}
		return cfg, nil
	}
	mergeLegacyLocalConfigJSON(raw)
	merged, err := json.Marshal(raw)
	if err != nil {
		return LocalConfig{}, fmt.Errorf("substackbrowser: merge config: %w", err)
	}
	var cfg LocalConfig
	if err := json.Unmarshal(merged, &cfg); err != nil {
		return LocalConfig{}, fmt.Errorf("substackbrowser: decode config: %w", err)
	}
	return cfg, nil
}

// LoadLocalConfigWithGlobal loads overlayPath (default substack.json) and merges
// it over globalPath when globalPath is non-empty. Overlay keys win on conflicts.
// When globalPath is empty, behavior matches LoadLocalConfig(overlayPath).
func LoadLocalConfigWithGlobal(globalPath, overlayPath string) (LocalConfig, bool, error) {
	overlayPath = strings.TrimSpace(overlayPath)
	if overlayPath == "" {
		overlayPath = DefaultLocalConfigPath()
	}
	globalPath = strings.TrimSpace(globalPath)
	if globalPath == "" {
		return loadLocalConfigSingleFile(overlayPath)
	}

	bGlobal, err := os.ReadFile(globalPath)
	if err != nil {
		return LocalConfig{}, false, fmt.Errorf("substackbrowser: read config-global %q: %w", globalPath, err)
	}

	var bOverlay []byte
	ob, err := readLocalConfigFileBytes(overlayPath)
	if err == nil {
		bOverlay = ob
	} else if os.IsNotExist(err) {
		if filepath.Clean(overlayPath) != DefaultLocalConfigPath() {
			return LocalConfig{}, false, fmt.Errorf("substackbrowser: read config %q: %w", overlayPath, err)
		}
	} else {
		return LocalConfig{}, false, fmt.Errorf("substackbrowser: read config %q: %w", overlayPath, err)
	}

	var merged []byte
	if len(bOverlay) == 0 {
		merged = bGlobal
	} else {
		merged, err = mergeJSONConfigBytes(bGlobal, bOverlay)
		if err != nil {
			return LocalConfig{}, false, err
		}
	}

	cfg, err := decodeLocalConfigFromBytes(merged)
	if err != nil {
		return LocalConfig{}, false, err
	}
	finalizeLoadedLocalConfig(&cfg)
	return cfg, true, nil
}

func loadLocalConfigSingleFile(p string) (LocalConfig, bool, error) {
	p = strings.TrimSpace(p)
	if p == "" {
		return LocalConfig{}, false, errors.New("substackbrowser: empty config path")
	}
	b, err := readLocalConfigFileBytes(p)
	if err != nil {
		if os.IsNotExist(err) {
			var cfg LocalConfig
			finalizeLoadedLocalConfig(&cfg)
			return cfg, false, nil
		}
		return LocalConfig{}, false, fmt.Errorf("substackbrowser: read config: %w", err)
	}
	cfg, err := decodeLocalConfigFromBytes(b)
	if err != nil {
		return LocalConfig{}, false, err
	}
	finalizeLoadedLocalConfig(&cfg)
	return cfg, true, nil
}
