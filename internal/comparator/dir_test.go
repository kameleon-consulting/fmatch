package comparator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mlabate/fmatch/internal/comparator"
	"github.com/mlabate/fmatch/internal/ignore"
)

// noMatcher is an empty ignore.Matcher that never ignores anything.
var noMatcher = ignore.LoadPatterns([]string{})

// makeDirTree creates a directory tree from a map of relative path → content.
// Directories are created automatically.
func makeDirTree(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, content := range files {
		abs := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			t.Fatalf("makeDirTree mkdir: %v", err)
		}
		if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
			t.Fatalf("makeDirTree write %s: %v", rel, err)
		}
	}
	return root
}

// ── Identical directories ────────────────────────────────────────────────────

func TestCompareDir_Identical_Empty(t *testing.T) {
	a, b := t.TempDir(), t.TempDir()
	result, err := comparator.CompareDir(a, b, comparator.DirOptions{Matcher: noMatcher, Depth: -1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Identical {
		t.Error("expected two empty directories to be identical")
	}
	if len(result.OnlyInA) != 0 || len(result.OnlyInB) != 0 || len(result.Different) != 0 {
		t.Errorf("expected no diffs, got OnlyInA=%v OnlyInB=%v Different=%v",
			result.OnlyInA, result.OnlyInB, result.Different)
	}
}

func TestCompareDir_Identical_SameFiles(t *testing.T) {
	files := map[string]string{
		"a.txt":      "hello",
		"sub/b.txt":  "world",
	}
	a := makeDirTree(t, files)
	b := makeDirTree(t, files)

	result, err := comparator.CompareDir(a, b, comparator.DirOptions{Matcher: noMatcher, Depth: -1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Identical {
		t.Errorf("expected identical dirs, got diffs: %+v", result)
	}
}

// ── Files only in one side ───────────────────────────────────────────────────

func TestCompareDir_OnlyInA(t *testing.T) {
	a := makeDirTree(t, map[string]string{
		"a.txt":  "hello",
		"extra.txt": "only in a",
	})
	b := makeDirTree(t, map[string]string{
		"a.txt": "hello",
	})

	result, err := comparator.CompareDir(a, b, comparator.DirOptions{Matcher: noMatcher, Depth: -1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Identical {
		t.Error("expected dirs to differ")
	}
	if len(result.OnlyInA) != 1 || result.OnlyInA[0] != "extra.txt" {
		t.Errorf("expected OnlyInA=[extra.txt], got %v", result.OnlyInA)
	}
}

func TestCompareDir_OnlyInB(t *testing.T) {
	a := makeDirTree(t, map[string]string{
		"a.txt": "hello",
	})
	b := makeDirTree(t, map[string]string{
		"a.txt":  "hello",
		"extra.txt": "only in b",
	})

	result, err := comparator.CompareDir(a, b, comparator.DirOptions{Matcher: noMatcher, Depth: -1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Identical {
		t.Error("expected dirs to differ")
	}
	if len(result.OnlyInB) != 1 || result.OnlyInB[0] != "extra.txt" {
		t.Errorf("expected OnlyInB=[extra.txt], got %v", result.OnlyInB)
	}
}

// ── Different content ────────────────────────────────────────────────────────

func TestCompareDir_DifferentContent(t *testing.T) {
	a := makeDirTree(t, map[string]string{"file.txt": "version A"})
	b := makeDirTree(t, map[string]string{"file.txt": "version B"})

	result, err := comparator.CompareDir(a, b, comparator.DirOptions{Matcher: noMatcher, Depth: -1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Identical {
		t.Error("expected dirs to differ")
	}
	if len(result.Different) != 1 || result.Different[0] != "file.txt" {
		t.Errorf("expected Different=[file.txt], got %v", result.Different)
	}
}

// ── Depth limit ──────────────────────────────────────────────────────────────

func TestCompareDir_Depth0_IgnoresSubdirs(t *testing.T) {
	// Both dirs have same root files but different subdir content.
	a := makeDirTree(t, map[string]string{
		"root.txt":    "same",
		"sub/deep.txt": "only in A",
	})
	b := makeDirTree(t, map[string]string{
		"root.txt": "same",
	})

	// Depth 0 = root files only, subdirs not scanned.
	result, err := comparator.CompareDir(a, b, comparator.DirOptions{Matcher: noMatcher, Depth: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Identical {
		t.Errorf("with depth=0, sub differences should be invisible; got diffs: %+v", result)
	}
}

// ── Ignore patterns ──────────────────────────────────────────────────────────

func TestCompareDir_IgnoredFiles_NotCompared(t *testing.T) {
	a := makeDirTree(t, map[string]string{
		"main.go":   "same",
		"debug.log": "differs in A",
	})
	b := makeDirTree(t, map[string]string{
		"main.go":   "same",
		"debug.log": "differs in B",
	})

	m := ignore.LoadPatterns([]string{"*.log"})
	result, err := comparator.CompareDir(a, b, comparator.DirOptions{Matcher: m, Depth: -1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Identical {
		t.Errorf("*.log should be ignored; got diffs: %+v", result)
	}
}

// ── Aggregated counters ──────────────────────────────────────────────────────

func TestCompareDir_Counters(t *testing.T) {
	a := makeDirTree(t, map[string]string{
		"common_same.txt":  "same",
		"common_diff.txt":  "version A",
		"only_a.txt":       "only A",
	})
	b := makeDirTree(t, map[string]string{
		"common_same.txt": "same",
		"common_diff.txt": "version B",
		"only_b.txt":      "only B",
	})

	result, err := comparator.CompareDir(a, b, comparator.DirOptions{Matcher: noMatcher, Depth: -1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalA != 3 {
		t.Errorf("expected TotalA=3, got %d", result.TotalA)
	}
	if result.TotalB != 3 {
		t.Errorf("expected TotalB=3, got %d", result.TotalB)
	}
	if len(result.OnlyInA) != 1 {
		t.Errorf("expected 1 file only in A, got %v", result.OnlyInA)
	}
	if len(result.OnlyInB) != 1 {
		t.Errorf("expected 1 file only in B, got %v", result.OnlyInB)
	}
	if len(result.Different) != 1 {
		t.Errorf("expected 1 different file, got %v", result.Different)
	}
}

// ── Error cases ──────────────────────────────────────────────────────────────

func TestCompareDir_PathNotFound_Error(t *testing.T) {
	_, err := comparator.CompareDir("/nonexistent/a", "/nonexistent/b",
		comparator.DirOptions{Matcher: noMatcher, Depth: -1})
	if err == nil {
		t.Error("expected error for nonexistent paths, got nil")
	}
}
