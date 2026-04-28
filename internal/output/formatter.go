package output

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mlabate/fmatch/internal/comparator"
)


// Verbosity represents the output verbosity level.
type Verbosity int

const (
	VerbosityQuiet   Verbosity = -1 // -q: no output, exit code only
	VerbosityNormal  Verbosity = 0  // default: one-line result
	VerbosityVerbose Verbosity = 1  // -v: paths and sizes
	VerbosityVV      Verbosity = 2  // -vv: + SHA-256 hashes and diff offset
)

// ANSI color escape codes.
const (
	colorGreen = "\x1b[32m"
	colorRed   = "\x1b[31m"
	colorReset = "\x1b[0m"
)

// Options controls the output format.
type Options struct {
	Level   Verbosity // verbosity level
	NoColor bool      // disable ANSI color codes
	PathA   string    // path to file/dir A (required for Verbose and VV)
	PathB   string    // path to file/dir B (required for Verbose and VV)
}

// Format returns the formatted output string for a file comparison result.
// Returns ("", nil) for VerbosityQuiet.
// Returns an error only in VerbosityVV mode if SHA-256 computation fails.
func Format(result comparator.Result, opts Options) (string, error) {
	if opts.Level == VerbosityQuiet {
		return "", nil
	}

	label, labelColor := "IDENTICAL", colorGreen
	if !result.Identical {
		label, labelColor = "DIFFERENT", colorRed
	}

	coloredLabel := label
	if !opts.NoColor {
		coloredLabel = labelColor + label + colorReset
	}

	switch opts.Level {
	case VerbosityNormal:
		return coloredLabel, nil

	case VerbosityVerbose:
		var b strings.Builder
		b.WriteString(coloredLabel)
		b.WriteString(formatPaths(result, opts))
		return b.String(), nil

	case VerbosityVV:
		hashA, err := fileHash(opts.PathA)
		if err != nil {
			return "", fmt.Errorf("sha256 %s: %w", opts.PathA, err)
		}
		hashB, err := fileHash(opts.PathB)
		if err != nil {
			return "", fmt.Errorf("sha256 %s: %w", opts.PathB, err)
		}

		var b strings.Builder
		b.WriteString(coloredLabel)
		b.WriteString(formatPaths(result, opts))
		b.WriteString(fmt.Sprintf("\n  sha256(a): %s", hashA))
		b.WriteString(fmt.Sprintf("\n  sha256(b): %s", hashB))
		if !result.Identical && result.DiffOffset >= 0 {
			b.WriteString(fmt.Sprintf("\n  first difference at byte: %d", result.DiffOffset))
		}
		return b.String(), nil
	}

	// Fallback (should not be reached with defined constants).
	return coloredLabel, nil
}

// formatPaths formats path and size lines for verbose output.
func formatPaths(result comparator.Result, opts Options) string {
	return fmt.Sprintf(
		"\n  path_a: %s (%d bytes)\n  path_b: %s (%d bytes)",
		opts.PathA, result.SizeA,
		opts.PathB, result.SizeB,
	)
}

// fileHash computes the SHA-256 hash of a file and returns it as a lowercase hex string.
func fileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// FormatDir returns the formatted output string for a directory comparison result.
// Returns ("", nil) for VerbosityQuiet.
func FormatDir(result comparator.DirResult, opts Options) (string, error) {
	if opts.Level == VerbosityQuiet {
		return "", nil
	}

	label, labelColor := "IDENTICAL", colorGreen
	if !result.Identical {
		label, labelColor = "DIFFERENT", colorRed
	}

	coloredLabel := label
	if !opts.NoColor {
		coloredLabel = labelColor + label + colorReset
	}

	switch opts.Level {
	case VerbosityNormal:
		var b strings.Builder
		b.WriteString(coloredLabel)
		if !result.Identical {
			b.WriteString(fmt.Sprintf(
				"\n  %d different · %d only in A · %d only in B",
				len(result.Different), len(result.OnlyInA), len(result.OnlyInB),
			))
		}
		return b.String(), nil

	case VerbosityVerbose, VerbosityVV:
		var b strings.Builder
		b.WriteString(coloredLabel)
		b.WriteString(fmt.Sprintf("\n  path_a: %s (%d files)", opts.PathA, result.TotalA))
		b.WriteString(fmt.Sprintf("\n  path_b: %s (%d files)", opts.PathB, result.TotalB))

		if result.Identical {
			b.WriteString(fmt.Sprintf("\n  all %d files are identical", result.TotalA))
		} else {
			if len(result.OnlyInA) > 0 {
				b.WriteString(fmt.Sprintf("\n  only in A (%d):", len(result.OnlyInA)))
				for _, f := range result.OnlyInA {
					b.WriteString(fmt.Sprintf("\n    %s", f))
				}
			}
			if len(result.OnlyInB) > 0 {
				b.WriteString(fmt.Sprintf("\n  only in B (%d):", len(result.OnlyInB)))
				for _, f := range result.OnlyInB {
					b.WriteString(fmt.Sprintf("\n    %s", f))
				}
			}
			if len(result.Different) > 0 {
				b.WriteString(fmt.Sprintf("\n  different  (%d):", len(result.Different)))
				for _, f := range result.Different {
					b.WriteString(fmt.Sprintf("\n    %s", f))
				}
			}
		}
		return b.String(), nil
	}

	return coloredLabel, nil
}
