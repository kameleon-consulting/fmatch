package hash_test

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mlabate/fmatch/internal/hash"
)

// expectedHash computes the SHA-256 hash of the given bytes using stdlib directly.
// Used as the reference implementation in tests.
func expectedHash(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// writeTempFile creates a temporary file with the given content and returns its path.
// The file is placed in t.TempDir() and cleaned up automatically.
func writeTempFile(t *testing.T, content []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "testfile")
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return path
}

// TestFileHash_KnownContent verifies that FileHash returns the correct SHA-256
// for a file with known content, cross-checked against stdlib crypto/sha256.
func TestFileHash_KnownContent(t *testing.T) {
	content := []byte("the quick brown fox jumps over the lazy dog\n")
	path := writeTempFile(t, content)

	got, err := hash.FileHash(path)
	if err != nil {
		t.Fatalf("FileHash returned unexpected error: %v", err)
	}

	want := expectedHash(content)
	if got != want {
		t.Errorf("FileHash(%q) = %q, want %q", path, got, want)
	}
}

// TestFileHash_EmptyFile verifies that FileHash returns the correct SHA-256
// for an empty file (SHA-256 of zero bytes).
func TestFileHash_EmptyFile(t *testing.T) {
	path := writeTempFile(t, []byte{})

	got, err := hash.FileHash(path)
	if err != nil {
		t.Fatalf("FileHash returned unexpected error: %v", err)
	}

	// SHA-256 of empty input is a well-known constant; verify via stdlib.
	h := sha256.New()
	want := fmt.Sprintf("%x", h.Sum(nil))

	if got != want {
		t.Errorf("FileHash(empty) = %q, want %q", got, want)
	}
}

// TestFileHash_LargeFile verifies that FileHash handles files larger than a single
// read buffer (> 64 KB) and still returns the correct hash.
func TestFileHash_LargeFile(t *testing.T) {
	// 128 KB of repeated bytes — larger than a typical 64 KB read buffer.
	content := make([]byte, 128*1024)
	for i := range content {
		content[i] = byte(i % 251) // arbitrary non-trivial pattern
	}
	path := writeTempFile(t, content)

	got, err := hash.FileHash(path)
	if err != nil {
		t.Fatalf("FileHash returned unexpected error: %v", err)
	}

	want := expectedHash(content)
	if got != want {
		t.Errorf("FileHash(largeFile) = %q, want %q", got, want)
	}
}

// TestFileHash_NonExistent verifies that FileHash returns a non-nil error
// when the target file does not exist.
func TestFileHash_NonExistent(t *testing.T) {
	_, err := hash.FileHash("/non/existent/path/that/cannot/exist.txt")
	if err == nil {
		t.Error("FileHash(nonExistent): expected error, got nil")
	}
}

// TestFileHash_OutputFormat verifies that the returned hash is a lowercase
// hex string of exactly 64 characters (SHA-256 produces 32 bytes = 64 hex chars).
func TestFileHash_OutputFormat(t *testing.T) {
	path := writeTempFile(t, []byte("format check"))

	got, err := hash.FileHash(path)
	if err != nil {
		t.Fatalf("FileHash returned unexpected error: %v", err)
	}

	if len(got) != 64 {
		t.Errorf("FileHash output length = %d, want 64", len(got))
	}

	for i, c := range got {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("FileHash output[%d] = %q: not lowercase hex", i, c)
			break
		}
	}
}


