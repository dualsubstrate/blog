package facebookautopost

import (
	"encoding/json"
	"fmt"
	"os"
)

// State tracks GUIDs that were successfully posted to Facebook.
type State struct {
	PostedGUIDs []string `json:"posted_guids"`
}

// Posted reports whether guid was already posted.
func (s *State) Posted(guid string) bool {
	for _, g := range s.PostedGUIDs {
		if g == guid {
			return true
		}
	}
	return false
}

// MarkPosted records a successful post. It does not write to disk.
func (s *State) MarkPosted(guid string) {
	if s.Posted(guid) {
		return
	}
	s.PostedGUIDs = append(s.PostedGUIDs, guid)
}

// LoadState reads state from path. Missing file yields an empty state.
func LoadState(path string) (*State, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &State{PostedGUIDs: nil}, nil
		}
		return nil, err
	}
	var st State
	if len(b) == 0 {
		return &State{}, nil
	}
	if err := json.Unmarshal(b, &st); err != nil {
		return nil, fmt.Errorf("state json: %w", err)
	}
	return &st, nil
}

// SaveState writes state to path (typically only after a confirmed Facebook success).
func SaveState(path string, st *State) error {
	b, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}
