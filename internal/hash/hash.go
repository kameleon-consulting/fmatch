package hash

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// FileHash computes the SHA-256 hash of the file at path and returns it as a
// lowercase hex string (64 characters). Returns a non-nil error if the file
// cannot be opened or read.
func FileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("hash %s: %w", path, err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
