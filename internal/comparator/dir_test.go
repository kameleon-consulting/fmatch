package comparator_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/mlabate/fmatch/internal/comparator"
	"github.com/mlabate/fmatch/internal/ignore"
)

// noMatcher is an empty ignore.Matcher that never ignores anything.
var noMatcher = ignore.LoadPatterns([]string{})

// makeDirTree creates a temporary directory tree from a map of relative path → content.
// Parent directories are created automatically. Cleaned up by t.TempDir().
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

// sortedStrings returns a sorted copy of a string slice.
func sortedStrings(ss []string) []string {
	cp := make([]string, len(ss))
	copy(cp, ss)
	sort.Strings(cp)
	return cp
}

// ── CompareDir (hash-based, 2-dir) ───────────────────────────────────────────

// TestCompareDir_BothEmpty: two empty directories are identical.
func TestCompareDir_BothEmpty(t *testing.T) {
	a, b := t.TempDir(), t.TempDir()
	result, err := comparator.CompareDir(a, b, comparator.DirOptions{Matcher: noMatcher, Depth: -1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Identical {
		t.Error("expected two empty dirs to be Identical=true")
	}
	if len(result.Matched) != 0 || len(result.OnlyInA) != 0 || len(result.OnlyInB) != 0 {
		t.Errorf("expected all slices empty; got Matched=%v OnlyInA=%v OnlyInB=%v",
			result.Matched, result.OnlyInA, result.OnlyInB)
	}
}

// TestCompareDir_SameContent_SameNames: same files, same names → Identical=true, all in Matched.
func TestCompareDir_SameContent_SameNames(t *testing.T) {
	files := map[string]string{
		"a.txt":     "hello",
		"sub/b.txt": "world",
	}
	a := makeDirTree(t, files)
	b := makeDirTree(t, files)

	result, err := comparator.CompareDir(a, b, comparator.DirOptions{Matcher: noMatcher, Depth: -1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Identical {
		t.Errorf("expected Identical=true; got: %+v", result)
	}
	if len(result.Matched) != 2 {
		t.Errorf("expected 2 Matched groups (one per unique hash); got %d: %+v", len(result.Matched), result.Matched)
	}
	if len(result.OnlyInA) != 0 || len(result.OnlyInB) != 0 {
		t.Errorf("expected no unmatched groups; OnlyInA=%v OnlyInB=%v", result.OnlyInA, result.OnlyInB)
	}
}

// TestCompareDir_SameContent_DifferentNames is the KEY v1.1 test:
// files with identical content but different names must be considered matched.
func TestCompareDir_SameContent_DifferentNames(t *testing.T) {
	a := makeDirTree(t, map[string]string{"report.pdf": "document content"})
	b := makeDirTree(t, map[string]string{"report_copy.pdf": "document content"})

	result, err := comparator.CompareDir(a, b, comparator.DirOptions{Matcher: noMatcher, Depth: -1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Identical {
		t.Errorf("same content, different names: expected Identical=true; got: %+v", result)
	}
	if len(result.Matched) != 1 {
		t.Fatalf("expected 1 Matched group; got %d: %+v", len(result.Matched), result.Matched)
	}
	g := result.Matched[0]
	if len(g.InA) != 1 || g.InA[0] != "report.pdf" {
		t.Errorf("expected InA=[report.pdf]; got %v", g.InA)
	}
	if len(g.InB) != 1 || g.InB[0] != "report_copy.pdf" {
		t.Errorf("expected InB=[report_copy.pdf]; got %v", g.InB)
	}
}

// TestCompareDir_OnlyInA: content present in A with no match in B → OnlyInA.
func TestCompareDir_OnlyInA(t *testing.T) {
	a := makeDirTree(t, map[string]string{
		"common.txt": "shared",
		"extra.txt":  "unique to a",
	})
	b := makeDirTree(t, map[string]string{
		"common.txt": "shared",
	})

	result, err := comparator.CompareDir(a, b, comparator.DirOptions{Matcher: noMatcher, Depth: -1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Identical {
		t.Error("expected Identical=false")
	}
	if len(result.OnlyInA) != 1 {
		t.Fatalf("expected 1 OnlyInA group; got %d: %+v", len(result.OnlyInA), result.OnlyInA)
	}
	if len(result.OnlyInA[0].InA) != 1 || result.OnlyInA[0].InA[0] != "extra.txt" {
		t.Errorf("expected OnlyInA[0].InA=[extra.txt]; got %v", result.OnlyInA[0].InA)
	}
	if len(result.OnlyInB) != 0 {
		t.Errorf("expected no OnlyInB; got %v", result.OnlyInB)
	}
}

// TestCompareDir_OnlyInB: content present in B with no match in A → OnlyInB.
func TestCompareDir_OnlyInB(t *testing.T) {
	a := makeDirTree(t, map[string]string{
		"common.txt": "shared",
	})
	b := makeDirTree(t, map[string]string{
		"common.txt": "shared",
		"extra.txt":  "unique to b",
	})

	result, err := comparator.CompareDir(a, b, comparator.DirOptions{Matcher: noMatcher, Depth: -1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Identical {
		t.Error("expected Identical=false")
	}
	if len(result.OnlyInB) != 1 {
		t.Fatalf("expected 1 OnlyInB group; got %d: %+v", len(result.OnlyInB), result.OnlyInB)
	}
	if len(result.OnlyInB[0].InB) != 1 || result.OnlyInB[0].InB[0] != "extra.txt" {
		t.Errorf("expected OnlyInB[0].InB=[extra.txt]; got %v", result.OnlyInB[0].InB)
	}
	if len(result.OnlyInA) != 0 {
		t.Errorf("expected no OnlyInA; got %v", result.OnlyInA)
	}
}

// TestCompareDir_ManyToMany: A has f1+f2 (same hash), B has f3 (same hash)
// → single HashGroup in Matched with InA=[f1,f2], InB=[f3].
func TestCompareDir_ManyToMany(t *testing.T) {
	a := makeDirTree(t, map[string]string{
		"f1.txt": "same content",
		"f2.txt": "same content",
	})
	b := makeDirTree(t, map[string]string{
		"f3.txt": "same content",
	})

	result, err := comparator.CompareDir(a, b, comparator.DirOptions{Matcher: noMatcher, Depth: -1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Identical {
		t.Errorf("all content is matched: expected Identical=true; got: %+v", result)
	}
	if len(result.Matched) != 1 {
		t.Fatalf("expected 1 Matched group; got %d: %+v", len(result.Matched), result.Matched)
	}
	g := result.Matched[0]
	if len(g.InA) != 2 {
		t.Errorf("expected 2 files in InA; got %v", g.InA)
	}
	got := sortedStrings(g.InA)
	if got[0] != "f1.txt" || got[1] != "f2.txt" {
		t.Errorf("expected InA=[f1.txt, f2.txt]; got %v", got)
	}
	if len(g.InB) != 1 || g.InB[0] != "f3.txt" {
		t.Errorf("expected InB=[f3.txt]; got %v", g.InB)
	}
}

// TestCompareDir_MixedResult: Matched + OnlyInA + OnlyInB all present.
func TestCompareDir_MixedResult(t *testing.T) {
	// file1 ↔ file3 (same content), file2 only in A, file4 only in B.
	a := makeDirTree(t, map[string]string{
		"file1.txt": "shared content",
		"file2.txt": "unique to a",
	})
	b := makeDirTree(t, map[string]string{
		"file3.txt": "shared content",
		"file4.txt": "unique to b",
	})

	result, err := comparator.CompareDir(a, b, comparator.DirOptions{Matcher: noMatcher, Depth: -1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Identical {
		t.Error("expected Identical=false")
	}
	if len(result.Matched) != 1 {
		t.Errorf("expected 1 Matched group; got %d", len(result.Matched))
	}
	if len(result.OnlyInA) != 1 {
		t.Errorf("expected 1 OnlyInA group; got %d", len(result.OnlyInA))
	}
	if len(result.OnlyInB) != 1 {
		t.Errorf("expected 1 OnlyInB group; got %d", len(result.OnlyInB))
	}
}

// TestCompareDir_DepthLimit: with Depth=0, subdirectory differences are invisible.
func TestCompareDir_DepthLimit(t *testing.T) {
	a := makeDirTree(t, map[string]string{
		"root.txt":     "same",
		"sub/deep.txt": "deep content A",
	})
	b := makeDirTree(t, map[string]string{
		"root.txt":     "same",
		"sub/deep.txt": "deep content B",
	})

	result, err := comparator.CompareDir(a, b, comparator.DirOptions{Matcher: noMatcher, Depth: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Identical {
		t.Errorf("depth=0: sub differences should be invisible; got: %+v", result)
	}
}

// TestCompareDir_IgnorePatterns: files matching ignore patterns are excluded.
func TestCompareDir_IgnorePatterns(t *testing.T) {
	a := makeDirTree(t, map[string]string{
		"main.go":   "same code",
		"debug.log": "log A",
	})
	b := makeDirTree(t, map[string]string{
		"main.go":   "same code",
		"debug.log": "log B",
	})

	m := ignore.LoadPatterns([]string{"*.log"})
	result, err := comparator.CompareDir(a, b, comparator.DirOptions{Matcher: m, Depth: -1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Identical {
		t.Errorf("*.log should be ignored; got: %+v", result)
	}
}

// TestCompareDir_ErrorNonExistent: non-existent path must return an error.
func TestCompareDir_ErrorNonExistent(t *testing.T) {
	_, err := comparator.CompareDir("/nonexistent/a", "/nonexistent/b",
		comparator.DirOptions{Matcher: noMatcher, Depth: -1})
	if err == nil {
		t.Error("expected error for nonexistent path, got nil")
	}
}

// ── FindDuplicates (1-dir) ────────────────────────────────────────────────────

// TestFindDuplicates_EmptyDir: empty directory → no duplicates, no files.
func TestFindDuplicates_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	result, err := comparator.FindDuplicates(dir, comparator.DirOptions{Matcher: noMatcher, Depth: -1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasDuplicates {
		t.Error("expected HasDuplicates=false for empty dir")
	}
	if len(result.Groups) != 0 {
		t.Errorf("expected no Groups; got %v", result.Groups)
	}
	if len(result.Unique) != 0 {
		t.Errorf("expected no Unique; got %v", result.Unique)
	}
}

// TestFindDuplicates_NoDuplicates: all files have unique content.
func TestFindDuplicates_NoDuplicates(t *testing.T) {
	dir := makeDirTree(t, map[string]string{
		"a.txt": "content a",
		"b.txt": "content b",
		"c.txt": "content c",
	})

	result, err := comparator.FindDuplicates(dir, comparator.DirOptions{Matcher: noMatcher, Depth: -1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.HasDuplicates {
		t.Error("expected HasDuplicates=false")
	}
	if len(result.Groups) != 0 {
		t.Errorf("expected no Groups; got %v", result.Groups)
	}
	if len(result.Unique) != 3 {
		t.Errorf("expected 3 Unique files; got %d: %v", len(result.Unique), result.Unique)
	}
}

// TestFindDuplicates_WithDuplicates: two files share content → 1 group.
func TestFindDuplicates_WithDuplicates(t *testing.T) {
	dir := makeDirTree(t, map[string]string{
		"a.txt": "duplicate content",
		"b.txt": "duplicate content",
		"c.txt": "unique content",
	})

	result, err := comparator.FindDuplicates(dir, comparator.DirOptions{Matcher: noMatcher, Depth: -1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasDuplicates {
		t.Error("expected HasDuplicates=true")
	}
	if len(result.Groups) != 1 {
		t.Fatalf("expected 1 Group; got %d: %+v", len(result.Groups), result.Groups)
	}
	if len(result.Groups[0].InA) != 2 {
		t.Errorf("expected 2 files in group; got %v", result.Groups[0].InA)
	}
	if len(result.Unique) != 1 {
		t.Errorf("expected 1 Unique file; got %v", result.Unique)
	}
}

// TestFindDuplicates_MultipleGroups: two distinct duplicate groups.
func TestFindDuplicates_MultipleGroups(t *testing.T) {
	dir := makeDirTree(t, map[string]string{
		"g1a.txt":    "group one",
		"g1b.txt":    "group one",
		"g2a.txt":    "group two",
		"g2b.txt":    "group two",
		"unique.txt": "something else",
	})

	result, err := comparator.FindDuplicates(dir, comparator.DirOptions{Matcher: noMatcher, Depth: -1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasDuplicates {
		t.Error("expected HasDuplicates=true")
	}
	if len(result.Groups) != 2 {
		t.Errorf("expected 2 Groups; got %d: %+v", len(result.Groups), result.Groups)
	}
	if len(result.Unique) != 1 {
		t.Errorf("expected 1 Unique file; got %v", result.Unique)
	}
}

// TestFindDuplicates_DepthLimit: with Depth=0, files in subdirs are excluded.
func TestFindDuplicates_DepthLimit(t *testing.T) {
	// Root has 2 duplicates; a third copy lives in a subdir (excluded by depth).
	dir := makeDirTree(t, map[string]string{
		"dup1.txt":     "duplicate",
		"dup2.txt":     "duplicate",
		"sub/dup3.txt": "duplicate",
	})

	result, err := comparator.FindDuplicates(dir, comparator.DirOptions{Matcher: noMatcher, Depth: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasDuplicates {
		t.Error("expected duplicates at root level")
	}
	if len(result.Groups) != 1 {
		t.Fatalf("expected 1 group; got %d", len(result.Groups))
	}
	// sub/dup3.txt must NOT be in the group (depth=0 excludes subdirs).
	if len(result.Groups[0].InA) != 2 {
		t.Errorf("expected 2 files in group (sub/ excluded); got %v", result.Groups[0].InA)
	}
}

// TestFindDuplicates_IgnorePatterns: files matching ignore patterns are excluded.
func TestFindDuplicates_IgnorePatterns(t *testing.T) {
	// a.txt and b.txt are duplicates; c.log has same content but must be ignored.
	dir := makeDirTree(t, map[string]string{
		"a.txt": "same content",
		"b.txt": "same content",
		"c.log": "same content",
	})

	m := ignore.LoadPatterns([]string{"*.log"})
	result, err := comparator.FindDuplicates(dir, comparator.DirOptions{Matcher: m, Depth: -1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.HasDuplicates {
		t.Error("expected HasDuplicates=true (a.txt and b.txt are duplicates)")
	}
	if len(result.Groups) != 1 {
		t.Fatalf("expected 1 Group; got %d: %+v", len(result.Groups), result.Groups)
	}
	for _, f := range result.Groups[0].InA {
		if f == "c.log" {
			t.Errorf("c.log should be ignored but appears in group: %v", result.Groups[0].InA)
		}
	}
	if len(result.Groups[0].InA) != 2 {
		t.Errorf("expected 2 files in group (c.log ignored); got %v", result.Groups[0].InA)
	}
}
