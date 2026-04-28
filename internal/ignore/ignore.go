package ignore

import (
	"os"
	"strings"

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

// LoadFileAndPatterns combines patterns from a file with additional inline patterns.
// If the file does not exist, only the extra patterns are used (no error).
// Returns an error only for real I/O failures.
func LoadFileAndPatterns(path string, extra []string) (*Matcher, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return LoadPatterns(extra), nil
		}
		return nil, err
	}
	lines := append(strings.Split(string(data), "\n"), extra...)
	ig := gitignore.CompileIgnoreLines(lines...)
	return &Matcher{ig: ig}, nil
}
