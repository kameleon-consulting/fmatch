package comparator_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/kameleon-consulting/fmatch/internal/comparator"
)

// writeTemp creates a temporary file with the given content.
// The file is removed automatically when the test ends.
func writeTemp(t *testing.T, content []byte) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "fmatch-*")
	if err != nil {
		t.Fatalf("writeTemp: CreateTemp: %v", err)
	}
	defer f.Close()
	if _, err := f.Write(content); err != nil {
		t.Fatalf("writeTemp: Write: %v", err)
	}
	return f.Name()
}

// chunk64KB returns a byte slice of exactly 64*1024 bytes filled with b.
func chunk64KB(b byte) []byte {
	return bytes.Repeat([]byte{b}, 64*1024)
}

// ── Identical files ──────────────────────────────────────────────────────────

func TestCompareFiles_Identical_Small(t *testing.T) {
	content := []byte("hello, fmatch!")
	a, b := writeTemp(t, content), writeTemp(t, content)

	got, err := comparator.CompareFiles(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Identical {
		t.Error("expected Identical=true")
	}
	if got.DiffOffset != -1 {
		t.Errorf("expected DiffOffset=-1, got %d", got.DiffOffset)
	}
	if got.SizeA != int64(len(content)) {
		t.Errorf("expected SizeA=%d, got %d", len(content), got.SizeA)
	}
}

func TestCompareFiles_Identical_Large(t *testing.T) {
	// 3 full 64 KB chunks → 192 KB
	content := bytes.Repeat(chunk64KB(0xAB), 3)
	a, b := writeTemp(t, content), writeTemp(t, content)

	got, err := comparator.CompareFiles(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Identical {
		t.Error("expected Identical=true for large identical files")
	}
}

// ── Different size ───────────────────────────────────────────────────────────

func TestCompareFiles_DifferentSize(t *testing.T) {
	a := writeTemp(t, []byte("short"))
	b := writeTemp(t, []byte("much longer content here"))

	got, err := comparator.CompareFiles(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Identical {
		t.Error("expected Identical=false for files with different sizes")
	}
	if got.SizeA == got.SizeB {
		t.Errorf("expected SizeA(%d) != SizeB(%d)", got.SizeA, got.SizeB)
	}
}

// ── Different content, same size ─────────────────────────────────────────────

func TestCompareFiles_DifferentContent_SameSize(t *testing.T) {
	a := writeTemp(t, []byte("hello world!"))
	b := writeTemp(t, []byte("hello WORLD!")) // diff at byte 6

	got, err := comparator.CompareFiles(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Identical {
		t.Error("expected Identical=false")
	}
	if got.DiffOffset != 6 {
		t.Errorf("expected DiffOffset=6, got %d", got.DiffOffset)
	}
}

func TestCompareFiles_DifferentContent_DiffInSecondChunk(t *testing.T) {
	// First chunk identical, second chunk differs at byte 0 of that chunk.
	contentA := append(chunk64KB(0x00), chunk64KB(0xAA)...)
	contentB := append(chunk64KB(0x00), chunk64KB(0xBB)...)
	a, b := writeTemp(t, contentA), writeTemp(t, contentB)

	got, err := comparator.CompareFiles(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Identical {
		t.Error("expected Identical=false")
	}
	// Difference starts at the first byte of the second chunk.
	if got.DiffOffset != 64*1024 {
		t.Errorf("expected DiffOffset=%d, got %d", 64*1024, got.DiffOffset)
	}
}

// ── Empty files ──────────────────────────────────────────────────────────────

func TestCompareFiles_BothEmpty(t *testing.T) {
	a, b := writeTemp(t, []byte{}), writeTemp(t, []byte{})

	got, err := comparator.CompareFiles(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Identical {
		t.Error("expected Identical=true for two empty files")
	}
}

func TestCompareFiles_OneEmpty(t *testing.T) {
	a := writeTemp(t, []byte{})
	b := writeTemp(t, []byte("not empty"))

	got, err := comparator.CompareFiles(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Identical {
		t.Error("expected Identical=false when one file is empty")
	}
}

// ── Binary files ─────────────────────────────────────────────────────────────

func TestCompareFiles_BinaryIdentical(t *testing.T) {
	content := []byte{0x00, 0xFF, 0x0A, 0x1B, 0xDE, 0xAD, 0xBE, 0xEF}
	a, b := writeTemp(t, content), writeTemp(t, content)

	got, err := comparator.CompareFiles(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Identical {
		t.Error("expected Identical=true for identical binary files")
	}
}

func TestCompareFiles_BinaryDifferent(t *testing.T) {
	a := writeTemp(t, []byte{0x00, 0xFF, 0xAA})
	b := writeTemp(t, []byte{0x00, 0xFF, 0xBB}) // diff at byte 2

	got, err := comparator.CompareFiles(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Identical {
		t.Error("expected Identical=false")
	}
	if got.DiffOffset != 2 {
		t.Errorf("expected DiffOffset=2, got %d", got.DiffOffset)
	}
}

// ── Error cases ──────────────────────────────────────────────────────────────

func TestCompareFiles_FileNotFound(t *testing.T) {
	_, err := comparator.CompareFiles("/nonexistent/a", "/nonexistent/b")
	if err == nil {
		t.Error("expected error for nonexistent files, got nil")
	}
}
