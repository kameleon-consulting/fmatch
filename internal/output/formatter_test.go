package output_test

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/kameleon-consulting/fmatch/internal/comparator"
	"github.com/kameleon-consulting/fmatch/internal/output"
)

// writeTemp creates a temporary file with the given content.
func writeTemp(t *testing.T, content []byte) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "fmatch-out-*")
	if err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	defer f.Close()
	if _, err := f.Write(content); err != nil {
		t.Fatalf("writeTemp write: %v", err)
	}
	return f.Name()
}

// sha256hex returns the SHA-256 hex digest of data (for expected value in tests).
func sha256hex(data []byte) string {
	h := sha256.Sum256(data)
	return fmt.Sprintf("%x", h)
}

// Fake 64-char SHA-256 hex strings for constructing test fixtures.
// Abbreviated (first 8 chars) used in expected output checks.
const (
	hashAlpha = "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2" // abbrev: a1b2c3d4
	hashBeta  = "b1c2d3e4f5a6b1c2d3e4f5a6b1c2d3e4f5a6b1c2d3e4f5a6b1c2d3e4f5a6b1c2" // abbrev: b1c2d3e4
	hashGamma = "c1d2e3f4a5b6c1d2e3f4a5b6c1d2e3f4a5b6c1d2e3f4a5b6c1d2e3f4a5b6c1d2" // abbrev: c1d2e3f4
)

// ── Format (file comparison) ──────────────────────────────────────────────────

func TestFormat_Quiet_Identical(t *testing.T) {
	result := comparator.Result{Identical: true, DiffOffset: -1, SizeA: 10, SizeB: 10}
	opts := output.Options{Level: output.VerbosityQuiet, NoColor: true}

	got, err := output.Format(result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("quiet mode: expected empty string, got %q", got)
	}
}

func TestFormat_Quiet_Different(t *testing.T) {
	result := comparator.Result{Identical: false, DiffOffset: -1, SizeA: 10, SizeB: 20}
	opts := output.Options{Level: output.VerbosityQuiet, NoColor: true}

	got, err := output.Format(result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("quiet mode: expected empty string, got %q", got)
	}
}

func TestFormat_Normal_Identical(t *testing.T) {
	result := comparator.Result{Identical: true, DiffOffset: -1, SizeA: 10, SizeB: 10}
	opts := output.Options{Level: output.VerbosityNormal, NoColor: true}

	got, err := output.Format(result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "IDENTICAL" {
		t.Errorf("expected %q, got %q", "IDENTICAL", got)
	}
}

func TestFormat_Normal_Different(t *testing.T) {
	result := comparator.Result{Identical: false, DiffOffset: 42, SizeA: 10, SizeB: 10}
	opts := output.Options{Level: output.VerbosityNormal, NoColor: true}

	got, err := output.Format(result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "DIFFERENT" {
		t.Errorf("expected %q, got %q", "DIFFERENT", got)
	}
}

func TestFormat_Verbose_Identical(t *testing.T) {
	content := []byte("hello")
	pathA, pathB := writeTemp(t, content), writeTemp(t, content)
	result := comparator.Result{Identical: true, DiffOffset: -1, SizeA: 5, SizeB: 5}
	opts := output.Options{Level: output.VerbosityVerbose, NoColor: true, PathA: pathA, PathB: pathB}

	got, err := output.Format(result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, "IDENTICAL") {
		t.Errorf("expected IDENTICAL in output, got %q", got)
	}
	if !strings.Contains(got, pathA) {
		t.Errorf("expected pathA in output, got %q", got)
	}
	if !strings.Contains(got, "5 bytes") {
		t.Errorf("expected size in output, got %q", got)
	}
}

func TestFormat_Verbose_Different(t *testing.T) {
	pathA := writeTemp(t, []byte("aaaaa"))
	pathB := writeTemp(t, []byte("bbbbb"))
	result := comparator.Result{Identical: false, DiffOffset: 0, SizeA: 5, SizeB: 5}
	opts := output.Options{Level: output.VerbosityVerbose, NoColor: true, PathA: pathA, PathB: pathB}

	got, err := output.Format(result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, "DIFFERENT") {
		t.Errorf("expected DIFFERENT in output, got %q", got)
	}
	if !strings.Contains(got, pathB) {
		t.Errorf("expected pathB in output, got %q", got)
	}
}

func TestFormat_VV_Identical_ContainsSHA256(t *testing.T) {
	content := []byte("hello, fmatch!")
	pathA, pathB := writeTemp(t, content), writeTemp(t, content)
	result := comparator.Result{Identical: true, DiffOffset: -1, SizeA: int64(len(content)), SizeB: int64(len(content))}
	opts := output.Options{Level: output.VerbosityVV, NoColor: true, PathA: pathA, PathB: pathB}

	got, err := output.Format(result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := sha256hex(content)
	if !strings.Contains(got, expected) {
		t.Errorf("expected SHA-256 %s in output, got %q", expected, got)
	}
}

func TestFormat_VV_Different_ContainsDiffOffset(t *testing.T) {
	pathA := writeTemp(t, []byte("hello world!"))
	pathB := writeTemp(t, []byte("hello WORLD!"))
	result := comparator.Result{Identical: false, DiffOffset: 6, SizeA: 12, SizeB: 12}
	opts := output.Options{Level: output.VerbosityVV, NoColor: true, PathA: pathA, PathB: pathB}

	got, err := output.Format(result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, "6") {
		t.Errorf("expected diff offset 6 in output, got %q", got)
	}
}

func TestFormat_NoColor_NoANSI(t *testing.T) {
	result := comparator.Result{Identical: true, DiffOffset: -1, SizeA: 5, SizeB: 5}
	opts := output.Options{Level: output.VerbosityNormal, NoColor: true}

	got, err := output.Format(result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(got, "\x1b[") {
		t.Errorf("expected no ANSI codes with NoColor=true, got %q", got)
	}
}

func TestFormat_WithColor_ContainsANSI(t *testing.T) {
	result := comparator.Result{Identical: true, DiffOffset: -1, SizeA: 5, SizeB: 5}
	opts := output.Options{Level: output.VerbosityNormal, NoColor: false}

	got, err := output.Format(result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, "\x1b[") {
		t.Errorf("expected ANSI codes with NoColor=false, got %q", got)
	}
}

// ── FormatDirCompare ──────────────────────────────────────────────────────────

// TestFormatDirCompare_Quiet: quiet mode always returns empty string.
func TestFormatDirCompare_Quiet(t *testing.T) {
	result := comparator.DirCompareResult{Identical: true}
	opts := output.Options{Level: output.VerbosityQuiet, NoColor: true}

	got, err := output.FormatDirCompare(result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("quiet mode: expected empty string, got %q", got)
	}
}

// TestFormatDirCompare_Normal_Identical: normal mode, identical dirs → "IDENTICAL".
func TestFormatDirCompare_Normal_Identical(t *testing.T) {
	result := comparator.DirCompareResult{
		Identical: true,
		Matched: []comparator.HashGroup{
			{Hash: hashAlpha, InA: []string{"a.txt"}, InB: []string{"a.txt"}},
		},
	}
	opts := output.Options{Level: output.VerbosityNormal, NoColor: true}

	got, err := output.FormatDirCompare(result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "IDENTICAL" {
		t.Errorf("expected %q, got %q", "IDENTICAL", got)
	}
}

// TestFormatDirCompare_Normal_Different: normal mode, different dirs → label + counts.
func TestFormatDirCompare_Normal_Different(t *testing.T) {
	result := comparator.DirCompareResult{
		Identical: false,
		Matched:   []comparator.HashGroup{{Hash: hashAlpha, InA: []string{"common.txt"}, InB: []string{"common.txt"}}},
		OnlyInA:   []comparator.HashGroup{{Hash: hashBeta, InA: []string{"extra_a.txt"}}},
		OnlyInB:   []comparator.HashGroup{{Hash: hashGamma, InB: []string{"extra_b.txt"}}},
	}
	opts := output.Options{Level: output.VerbosityNormal, NoColor: true}

	got, err := output.FormatDirCompare(result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"DIFFERENT", "1 matched", "1 only in A", "1 only in B"} {
		if !strings.Contains(got, want) {
			t.Errorf("expected %q in output, got:\n%s", want, got)
		}
	}
}

// TestFormatDirCompare_Verbose_Identical: verbose, identical → paths, file count, matched section.
func TestFormatDirCompare_Verbose_Identical(t *testing.T) {
	result := comparator.DirCompareResult{
		Identical: true,
		Matched: []comparator.HashGroup{
			{Hash: hashAlpha, InA: []string{"report.pdf"}, InB: []string{"report_copy.pdf"}},
		},
	}
	opts := output.Options{
		Level:   output.VerbosityVerbose,
		NoColor: true,
		PathA:   "/dirA",
		PathB:   "/dirB",
	}

	got, err := output.FormatDirCompare(result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"IDENTICAL", "/dirA", "/dirB", "matched", "report.pdf", "report_copy.pdf", "↔"} {
		if !strings.Contains(got, want) {
			t.Errorf("expected %q in verbose output, got:\n%s", want, got)
		}
	}
	// Hash must be shown abbreviated (first 8 chars).
	if !strings.Contains(got, hashAlpha[:8]) {
		t.Errorf("expected abbreviated hash %q in output, got:\n%s", hashAlpha[:8], got)
	}
}

// TestFormatDirCompare_Verbose_Different: verbose, different → all sections present.
func TestFormatDirCompare_Verbose_Different(t *testing.T) {
	result := comparator.DirCompareResult{
		Identical: false,
		Matched:   []comparator.HashGroup{{Hash: hashAlpha, InA: []string{"common.txt"}, InB: []string{"common_b.txt"}}},
		OnlyInA:   []comparator.HashGroup{{Hash: hashBeta, InA: []string{"only_a.txt"}}},
		OnlyInB:   []comparator.HashGroup{{Hash: hashGamma, InB: []string{"only_b.txt"}}},
	}
	opts := output.Options{
		Level:   output.VerbosityVerbose,
		NoColor: true,
		PathA:   "/dirA",
		PathB:   "/dirB",
	}

	got, err := output.FormatDirCompare(result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{
		"DIFFERENT", "/dirA", "/dirB",
		"matched", "common.txt", "common_b.txt", "↔",
		"only in A", "only_a.txt",
		"only in B", "only_b.txt",
		hashAlpha[:8], hashBeta[:8], hashGamma[:8],
	} {
		if !strings.Contains(got, want) {
			t.Errorf("expected %q in verbose output, got:\n%s", want, got)
		}
	}
}

// ── FormatDuplicates ──────────────────────────────────────────────────────────

// TestFormatDuplicates_Quiet: quiet mode always returns empty string.
func TestFormatDuplicates_Quiet(t *testing.T) {
	result := comparator.DuplicateResult{HasDuplicates: false}
	opts := output.Options{Level: output.VerbosityQuiet, NoColor: true}

	got, err := output.FormatDuplicates(result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("quiet mode: expected empty string, got %q", got)
	}
}

// TestFormatDuplicates_Normal_NoDuplicates: normal mode, no duplicates → "No duplicates found".
func TestFormatDuplicates_Normal_NoDuplicates(t *testing.T) {
	result := comparator.DuplicateResult{
		HasDuplicates: false,
		Unique:        []string{"a.txt", "b.txt"},
	}
	opts := output.Options{Level: output.VerbosityNormal, NoColor: true}

	got, err := output.FormatDuplicates(result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, "No duplicates found") {
		t.Errorf("expected %q, got %q", "No duplicates found", got)
	}
}

// TestFormatDuplicates_Normal_WithDuplicates: normal mode → "N duplicate groups found".
func TestFormatDuplicates_Normal_WithDuplicates(t *testing.T) {
	result := comparator.DuplicateResult{
		HasDuplicates: true,
		Groups: []comparator.HashGroup{
			{Hash: hashAlpha, InA: []string{"dup1.txt", "dup2.txt"}},
			{Hash: hashBeta, InA: []string{"dup3.txt", "dup4.txt"}},
		},
	}
	opts := output.Options{Level: output.VerbosityNormal, NoColor: true}

	got, err := output.FormatDuplicates(result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, "2 duplicate groups found") {
		t.Errorf("expected %q, got %q", "2 duplicate groups found", got)
	}
}

// TestFormatDuplicates_Verbose_WithDuplicates: verbose → hash, file count, file list.
func TestFormatDuplicates_Verbose_WithDuplicates(t *testing.T) {
	result := comparator.DuplicateResult{
		HasDuplicates: true,
		Groups: []comparator.HashGroup{
			{Hash: hashAlpha, InA: []string{"report.pdf", "report_copy.pdf"}},
		},
		Unique: []string{"unique.txt"},
	}
	opts := output.Options{Level: output.VerbosityVerbose, NoColor: true}

	got, err := output.FormatDuplicates(result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{
		"1 duplicate group found",
		hashAlpha[:8],
		"(2 files)",
		"report.pdf",
		"report_copy.pdf",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("expected %q in verbose output, got:\n%s", want, got)
		}
	}
}
