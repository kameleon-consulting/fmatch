package output_test

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mlabate/fmatch/internal/comparator"
	"github.com/mlabate/fmatch/internal/output"
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

// ── Quiet mode ───────────────────────────────────────────────────────────────

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

// ── Normal mode ──────────────────────────────────────────────────────────────

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

// ── Verbose mode ─────────────────────────────────────────────────────────────

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

// ── Very Verbose mode ────────────────────────────────────────────────────────

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

// ── Color output ─────────────────────────────────────────────────────────────

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

// ── FormatDir ────────────────────────────────────────────────────────────────

func TestFormatDir_Quiet(t *testing.T) {
	result := comparator.DirResult{Identical: true, TotalA: 3, TotalB: 3}
	opts := output.Options{Level: output.VerbosityQuiet, NoColor: true}
	got, err := output.FormatDir(result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("quiet mode: expected empty string, got %q", got)
	}
}

func TestFormatDir_Normal_Identical(t *testing.T) {
	result := comparator.DirResult{Identical: true, TotalA: 3, TotalB: 3}
	opts := output.Options{Level: output.VerbosityNormal, NoColor: true}
	got, err := output.FormatDir(result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "IDENTICAL" {
		t.Errorf("expected %q, got %q", "IDENTICAL", got)
	}
}

func TestFormatDir_Normal_Different_ContainsCounts(t *testing.T) {
	result := comparator.DirResult{
		Identical: false,
		TotalA:    5, TotalB: 5,
		Different: []string{"changed.txt"},
		OnlyInA:   []string{"extra_a.txt"},
		OnlyInB:   []string{},
	}
	opts := output.Options{Level: output.VerbosityNormal, NoColor: true}
	got, err := output.FormatDir(result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, "DIFFERENT") {
		t.Errorf("expected DIFFERENT in output, got %q", got)
	}
	if !strings.Contains(got, "1 different") {
		t.Errorf("expected count '1 different' in output, got %q", got)
	}
	if !strings.Contains(got, "1 only in A") {
		t.Errorf("expected '1 only in A' in output, got %q", got)
	}
}

func TestFormatDir_Verbose_Different_ContainsFileLists(t *testing.T) {
	result := comparator.DirResult{
		Identical: false,
		TotalA:    3, TotalB: 3,
		Different: []string{"changed.txt"},
		OnlyInA:   []string{"only_a.txt"},
		OnlyInB:   []string{"only_b.txt"},
	}
	opts := output.Options{
		Level:   output.VerbosityVerbose,
		NoColor: true,
		PathA:   "/dir/A",
		PathB:   "/dir/B",
	}
	got, err := output.FormatDir(result, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, expected := range []string{"changed.txt", "only_a.txt", "only_b.txt", "/dir/A", "/dir/B"} {
		if !strings.Contains(got, expected) {
			t.Errorf("expected %q in verbose output, got:\n%s", expected, got)
		}
	}
}
