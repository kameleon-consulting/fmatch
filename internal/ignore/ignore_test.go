package ignore_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kameleon-consulting/fmatch/internal/ignore"
)

// writeIgnoreFile creates a temporary .fmatchignore-style file with the given content.
func writeIgnoreFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), ".fmatchignore-*")
	if err != nil {
		t.Fatalf("writeIgnoreFile: %v", err)
	}
	defer f.Close()
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writeIgnoreFile write: %v", err)
	}
	return f.Name()
}

// ── LoadFile ─────────────────────────────────────────────────────────────────

func TestLoadFile_SimplePattern(t *testing.T) {
	path := writeIgnoreFile(t, "*.log\n")
	m, err := ignore.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.Match("error.log") {
		t.Error("expected error.log to match *.log")
	}
	if m.Match("main.go") {
		t.Error("expected main.go NOT to match *.log")
	}
}

func TestLoadFile_FileNotExist_ReturnsEmpty(t *testing.T) {
	m, err := ignore.LoadFile("/nonexistent/.fmatchignore")
	if err != nil {
		t.Fatalf("missing file should return empty Matcher, not error: %v", err)
	}
	if m.Match("anything.log") {
		t.Error("empty Matcher should not match anything")
	}
}

func TestLoadFile_Comments_Ignored(t *testing.T) {
	path := writeIgnoreFile(t, "# this is a comment\n*.tmp\n")
	m, err := ignore.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.Match("file.tmp") {
		t.Error("expected file.tmp to match *.tmp")
	}
}

func TestLoadFile_EmptyLines_Ignored(t *testing.T) {
	path := writeIgnoreFile(t, "\n\n*.bak\n\n")
	m, err := ignore.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.Match("old.bak") {
		t.Error("expected old.bak to match *.bak")
	}
}

func TestLoadFile_Negation(t *testing.T) {
	path := writeIgnoreFile(t, "*.log\n!important.log\n")
	m, err := ignore.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.Match("debug.log") {
		t.Error("expected debug.log to match *.log")
	}
	if m.Match("important.log") {
		t.Error("expected important.log NOT to match (negation with !)")
	}
}

func TestLoadFile_GlobStar(t *testing.T) {
	path := writeIgnoreFile(t, "**/*.log\n")
	m, err := ignore.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.Match(filepath.Join("subdir", "error.log")) {
		t.Errorf("expected subdir/error.log to match **/*.log")
	}
}

// ── LoadPatterns ─────────────────────────────────────────────────────────────

func TestLoadPatterns_SimplePatterns(t *testing.T) {
	m := ignore.LoadPatterns([]string{"*.swp", "*.tmp"})
	if !m.Match("file.swp") {
		t.Error("expected file.swp to match *.swp")
	}
	if !m.Match("temp.tmp") {
		t.Error("expected temp.tmp to match *.tmp")
	}
	if m.Match("main.go") {
		t.Error("expected main.go NOT to match")
	}
}

func TestLoadPatterns_Empty_MatchesNothing(t *testing.T) {
	m := ignore.LoadPatterns([]string{})
	if m.Match("anything.log") {
		t.Error("empty patterns should match nothing")
	}
}

// ── LoadFileAndPatterns ──────────────────────────────────────────────────────

func TestLoadFileAndPatterns_CombinesFileAndExtra(t *testing.T) {
	path := writeIgnoreFile(t, "*.log\n")
	m, err := ignore.LoadFileAndPatterns(path, []string{"*.tmp"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.Match("error.log") {
		t.Error("expected error.log to match *.log from file")
	}
	if !m.Match("temp.tmp") {
		t.Error("expected temp.tmp to match *.tmp from extra patterns")
	}
	if m.Match("main.go") {
		t.Error("expected main.go NOT to match")
	}
}

func TestLoadFileAndPatterns_FileNotExist_UsesExtraOnly(t *testing.T) {
	m, err := ignore.LoadFileAndPatterns("/nonexistent/.fmatchignore", []string{"*.swp"})
	if err != nil {
		t.Fatalf("missing file should not return error: %v", err)
	}
	if !m.Match("vim.swp") {
		t.Error("expected vim.swp to match *.swp from extra patterns")
	}
	if m.Match("main.go") {
		t.Error("expected main.go NOT to match")
	}
}

func TestLoadFileAndPatterns_NoExtra_BehavesLikeLoadFile(t *testing.T) {
	path := writeIgnoreFile(t, "*.bak\n")
	m, err := ignore.LoadFileAndPatterns(path, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.Match("old.bak") {
		t.Error("expected old.bak to match *.bak")
	}
}

