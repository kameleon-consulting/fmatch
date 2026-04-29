package output

import (
	"fmt"
	"strings"

	"github.com/kameleon-consulting/fmatch/internal/comparator"
	"github.com/kameleon-consulting/fmatch/internal/hash"
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

// coloredLabel returns the label with or without ANSI color codes.
func (o Options) coloredLabel(label, ansiColor string) string {
	if o.NoColor {
		return label
	}
	return ansiColor + label + colorReset
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
	cl := opts.coloredLabel(label, labelColor)

	switch opts.Level {
	case VerbosityNormal:
		return cl, nil

	case VerbosityVerbose:
		var b strings.Builder
		b.WriteString(cl)
		b.WriteString(formatPaths(result, opts))
		return b.String(), nil

	case VerbosityVV:
		hashA, err := hash.FileHash(opts.PathA)
		if err != nil {
			return "", fmt.Errorf("sha256 %s: %w", opts.PathA, err)
		}
		hashB, err := hash.FileHash(opts.PathB)
		if err != nil {
			return "", fmt.Errorf("sha256 %s: %w", opts.PathB, err)
		}

		var b strings.Builder
		b.WriteString(cl)
		b.WriteString(formatPaths(result, opts))
		b.WriteString(fmt.Sprintf("\n  sha256(a): %s", hashA))
		b.WriteString(fmt.Sprintf("\n  sha256(b): %s", hashB))
		if !result.Identical && result.DiffOffset >= 0 {
			b.WriteString(fmt.Sprintf("\n  first difference at byte: %d", result.DiffOffset))
		}
		return b.String(), nil
	}

	// Fallback (should not be reached with defined constants).
	return cl, nil
}

// formatPaths formats path and size lines for verbose file-comparison output.
func formatPaths(result comparator.Result, opts Options) string {
	return fmt.Sprintf(
		"\n  path_a: %s (%d bytes)\n  path_b: %s (%d bytes)",
		opts.PathA, result.SizeA,
		opts.PathB, result.SizeB,
	)
}

// FormatDirCompare returns the formatted output string for a hash-based
// directory comparison result. Returns ("", nil) for VerbosityQuiet.
func FormatDirCompare(result comparator.DirCompareResult, opts Options) (string, error) {
	if opts.Level == VerbosityQuiet {
		return "", nil
	}

	label, labelColor := "IDENTICAL", colorGreen
	if !result.Identical {
		label, labelColor = "DIFFERENT", colorRed
	}
	cl := opts.coloredLabel(label, labelColor)

	switch opts.Level {
	case VerbosityNormal:
		if result.Identical {
			return cl, nil
		}
		return fmt.Sprintf("%s\n  %d matched · %d only in A · %d only in B",
			cl,
			len(result.Matched),
			len(result.OnlyInA),
			len(result.OnlyInB),
		), nil

	case VerbosityVerbose, VerbosityVV:
		var b strings.Builder
		b.WriteString(cl)

		// Count total files on each side.
		totalA, totalB := 0, 0
		for _, g := range result.Matched {
			totalA += len(g.InA)
			totalB += len(g.InB)
		}
		for _, g := range result.OnlyInA {
			totalA += len(g.InA)
		}
		for _, g := range result.OnlyInB {
			totalB += len(g.InB)
		}

		b.WriteString(fmt.Sprintf("\n  path_a: %s (%d files)", opts.PathA, totalA))
		b.WriteString(fmt.Sprintf("\n  path_b: %s (%d files)", opts.PathB, totalB))

		if len(result.Matched) > 0 {
			b.WriteString(fmt.Sprintf("\n  matched (%d):", len(result.Matched)))
			for _, g := range result.Matched {
				b.WriteString(fmt.Sprintf("\n    [%s] %s  ↔  %s",
					g.Hash[:8],
					strings.Join(g.InA, ", "),
					strings.Join(g.InB, ", "),
				))
			}
		}

		if len(result.OnlyInA) > 0 {
			b.WriteString(fmt.Sprintf("\n  only in A (%d):", len(result.OnlyInA)))
			for _, g := range result.OnlyInA {
				b.WriteString(fmt.Sprintf("\n    [%s] %s",
					g.Hash[:8],
					strings.Join(g.InA, ", "),
				))
			}
		}

		if len(result.OnlyInB) > 0 {
			b.WriteString(fmt.Sprintf("\n  only in B (%d):", len(result.OnlyInB)))
			for _, g := range result.OnlyInB {
				b.WriteString(fmt.Sprintf("\n    [%s] %s",
					g.Hash[:8],
					strings.Join(g.InB, ", "),
				))
			}
		}

		return b.String(), nil
	}

	return cl, nil
}

// FormatDuplicates returns the formatted output string for a duplicate detection
// result. Returns ("", nil) for VerbosityQuiet.
func FormatDuplicates(result comparator.DuplicateResult, opts Options) (string, error) {
	if opts.Level == VerbosityQuiet {
		return "", nil
	}

	switch opts.Level {
	case VerbosityNormal:
		if !result.HasDuplicates {
			return "No duplicates found", nil
		}
		n := len(result.Groups)
		if n == 1 {
			return "1 duplicate group found", nil
		}
		return fmt.Sprintf("%d duplicate groups found", n), nil

	case VerbosityVerbose, VerbosityVV:
		var b strings.Builder

		if !result.HasDuplicates {
			b.WriteString("No duplicates found")
			return b.String(), nil
		}

		n := len(result.Groups)
		if n == 1 {
			b.WriteString("1 duplicate group found")
		} else {
			b.WriteString(fmt.Sprintf("%d duplicate groups found", n))
		}

		for _, g := range result.Groups {
			b.WriteString(fmt.Sprintf("\n  [%s] (%d files):", g.Hash[:8], len(g.InA)))
			for _, f := range g.InA {
				b.WriteString(fmt.Sprintf("\n    %s", f))
			}
		}

		return b.String(), nil
	}

	return "", nil
}
