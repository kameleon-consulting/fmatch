package ignore

import (
	"os"

	gitignore "github.com/sabhiram/go-gitignore"
)

// Matcher holds compiled ignore patterns and reports whether a path should be ignored.
type Matcher struct {
	ig *gitignore.GitIgnore
}

// empty is a no-op Matcher that never matches anything.
var empty = &Matcher{}

// LoadFile loads patterns from a .fmatchignore-style file.
// If the file does not exist, returns an empty Matcher (no error — missing file = no patterns).
// Returns an error only for I/O problems other than "file not found".
func LoadFile(path string) (*Matcher, error) {
	ig, err := gitignore.CompileIgnoreFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return empty, nil
		}
		return nil, err
	}
	return &Matcher{ig: ig}, nil
}

// LoadPatterns compiles a Matcher from a list of pattern strings (e.g., from -i flag).
// Returns an empty Matcher if patterns is empty.
func LoadPatterns(patterns []string) *Matcher {
	if len(patterns) == 0 {
		return empty
	}
	ig := gitignore.CompileIgnoreLines(patterns...)
	return &Matcher{ig: ig}
}

// Match reports whether the given path matches any ignore pattern.
// Returns false if the Matcher was created with no patterns.
func (m *Matcher) Match(path string) bool {
	if m.ig == nil {
		return false
	}
	return m.ig.MatchesPath(path)
}
