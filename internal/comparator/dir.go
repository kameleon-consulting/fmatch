package comparator

import (
	"io/fs"
	"path/filepath"
	"sort"

	"github.com/mlabate/fmatch/internal/ignore"
)

// DirOptions controls directory comparison behaviour.
type DirOptions struct {
	Matcher *ignore.Matcher // patterns to ignore (never nil; use LoadPatterns(nil) for empty)
	Depth   int             // max traversal depth; -1 = unlimited; 0 = root files only
}

// DirResult holds the aggregated outcome of a directory comparison.
type DirResult struct {
	Identical bool     // true only if OnlyInA, OnlyInB and Different are all empty
	TotalA    int      // total files scanned in A (after ignore)
	TotalB    int      // total files scanned in B (after ignore)
	OnlyInA   []string // relative paths present only in A
	OnlyInB   []string // relative paths present only in B
	Different []string // relative paths present in both but with different content
	Common    []string // relative paths present in both and identical
}

// CompareDir recursively compares two directories and returns an aggregated DirResult.
// Returns an error if either root path cannot be walked.
func CompareDir(pathA, pathB string, opts DirOptions) (DirResult, error) {
	filesA, err := walkDir(pathA, opts)
	if err != nil {
		return DirResult{}, err
	}
	filesB, err := walkDir(pathB, opts)
	if err != nil {
		return DirResult{}, err
	}

	var result DirResult
	result.TotalA = len(filesA)
	result.TotalB = len(filesB)

	// Set difference using sorted keys.
	for rel := range filesA {
		if _, inB := filesB[rel]; !inB {
			result.OnlyInA = append(result.OnlyInA, rel)
		}
	}
	for rel := range filesB {
		if _, inA := filesA[rel]; !inA {
			result.OnlyInB = append(result.OnlyInB, rel)
		}
	}

	// Compare files present in both directories.
	for rel, absA := range filesA {
		absB, inB := filesB[rel]
		if !inB {
			continue
		}
		fileResult, err := CompareFiles(absA, absB)
		if err != nil {
			// Treat unreadable files as different.
			result.Different = append(result.Different, rel)
			continue
		}
		if fileResult.Identical {
			result.Common = append(result.Common, rel)
		} else {
			result.Different = append(result.Different, rel)
		}
	}

	sort.Strings(result.OnlyInA)
	sort.Strings(result.OnlyInB)
	sort.Strings(result.Different)
	sort.Strings(result.Common)

	result.Identical = len(result.OnlyInA) == 0 &&
		len(result.OnlyInB) == 0 &&
		len(result.Different) == 0

	return result, nil
}

// walkDir returns a map of relative_path → absolute_path for all files under root,
// respecting depth limit and ignore patterns.
func walkDir(root string, opts DirOptions) (map[string]string, error) {
	files := make(map[string]string)

	err := filepath.WalkDir(root, func(abs string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, relErr := filepath.Rel(root, abs)
		if relErr != nil {
			return relErr
		}

		// Skip the root itself.
		if rel == "." {
			return nil
		}

		// Enforce depth limit.
		if opts.Depth >= 0 {
			depth := len(filepath.SplitList(filepath.ToSlash(rel)))
			// filepath.SplitList splits on os.PathListSeparator, not path separator.
			// Use a manual depth count instead.
			depth = depthOf(rel)
			if depth > opts.Depth {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Apply ignore patterns.
		if opts.Matcher != nil && opts.Matcher.Match(rel) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Record files only (not directories themselves).
		if !d.IsDir() {
			files[rel] = abs
		}

		return nil
	})

	return files, err
}

// depthOf returns the depth of a relative path (number of path segments).
// "file.txt" → 1, "sub/file.txt" → 2, "a/b/c.txt" → 3.
func depthOf(rel string) int {
	count := 0
	for _, part := range filepath.SplitList(filepath.Clean(rel)) {
		if part != "" {
			count++
		}
	}
	// filepath.SplitList splits on os.PathListSeparator (':' on Unix).
	// We need to split on the path separator instead.
	count = 1
	for _, c := range rel {
		if c == filepath.Separator {
			count++
		}
	}
	return count
}
