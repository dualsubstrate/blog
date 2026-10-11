package tagregister

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Deprecation maps a legacy tag string to the preferred replacement for new content.
type Deprecation struct {
	From string `toml:"from"`
	To   string `toml:"to"`
}

type deprecationsFile struct {
	Deprecated []Deprecation `toml:"deprecated"`
}

// LoadDeprecations reads TOML deprecations. Missing file returns nil, nil.
func LoadDeprecations(path string) ([]Deprecation, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var f deprecationsFile
	if _, err := toml.Decode(string(b), &f); err != nil {
		return nil, fmt.Errorf("deprecations toml: %w", err)
	}
	seenFrom := make(map[string]struct{})
	for i := range f.Deprecated {
		from := f.Deprecated[i].From
		to := f.Deprecated[i].To
		if from == "" || to == "" {
			return nil, fmt.Errorf("deprecation row %d: from and to must be non-empty", i+1)
		}
		if _, dup := seenFrom[from]; dup {
			return nil, fmt.Errorf("duplicate deprecated from %q", from)
		}
		seenFrom[from] = struct{}{}
	}
	return f.Deprecated, nil
}
