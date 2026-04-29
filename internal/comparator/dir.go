package comparator

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mlabate/fmatch/internal/hash"
	"github.com/mlabate/fmatch/internal/ignore"
)

// DirOptions controls directory traversal behaviour.
type DirOptions struct {
	Matcher *ignore.Matcher // patterns to ignore (never nil; use LoadPatterns(nil) for empty)
	Depth   int             // max traversal depth; -1 = unlimited; 0 = root files only
}

// HashGroup groups files that share the same SHA-256 hash.
// Used for both two-directory comparison and single-directory duplicate detection.
type HashGroup struct {
	Hash string   // SHA-256 hex (64 chars)
	InA  []string // relative paths in A with this hash (nil in DuplicateResult.Groups)
	InB  []string // relative paths in B with this hash (nil in DuplicateResult.Groups)
}

// DirCompareResult is the result of a hash-based two-directory comparison.
type DirCompareResult struct {
	Identical bool        // true only if OnlyInA and OnlyInB are both empty
	Matched   []HashGroup // hashes present in both dirs; each entry lists all matching paths on each side
	OnlyInA   []HashGroup // hashes present only in A (no file with that content exists in B)
	OnlyInB   []HashGroup // hashes present only in B (no file with that content exists in A)
}

// DuplicateResult is the result of duplicate detection within a single directory.
type DuplicateResult struct {
	HasDuplicates bool
	Groups        []HashGroup // groups of 2+ files sharing the same hash (InA used; InB always nil)
	Unique        []string    // relative paths of files with a unique hash
}

// CompareDir compares two directories by matching files on SHA-256 hash.
// Files with identical content are matched regardless of their name or location.
// Returns an error if either root path cannot be walked.
func CompareDir(pathA, pathB string, opts DirOptions) (DirCompareResult, error) {
	mapA, err := hashDir(pathA, opts)
	if err != nil {
		return DirCompareResult{}, err
	}
	mapB, err := hashDir(pathB, opts)
	if err != nil {
		return DirCompareResult{}, err
	}

	var result DirCompareResult

	for h, pathsA := range mapA {
		if pathsB, inB := mapB[h]; inB {
			result.Matched = append(result.Matched, HashGroup{
				Hash: h,
				InA:  sortedCopy(pathsA),
				InB:  sortedCopy(pathsB),
			})
		} else {
			result.OnlyInA = append(result.OnlyInA, HashGroup{
				Hash: h,
				InA:  sortedCopy(pathsA),
			})
		}
	}

	for h, pathsB := range mapB {
		if _, inA := mapA[h]; !inA {
			result.OnlyInB = append(result.OnlyInB, HashGroup{
				Hash: h,
				InB:  sortedCopy(pathsB),
			})
		}
	}

	sortHashGroups(result.Matched)
	sortHashGroups(result.OnlyInA)
	sortHashGroups(result.OnlyInB)

	result.Identical = len(result.OnlyInA) == 0 && len(result.OnlyInB) == 0

	return result, nil
}

// FindDuplicates finds files with identical content within a single directory.
// Files sharing a hash are grouped together; files with a unique hash go to Unique.
// Returns an error if the directory cannot be walked.
func FindDuplicates(path string, opts DirOptions) (DuplicateResult, error) {
	hashMap, err := hashDir(path, opts)
	if err != nil {
		return DuplicateResult{}, err
	}

	var result DuplicateResult

	for h, paths := range hashMap {
		sorted := sortedCopy(paths)
		if len(sorted) >= 2 {
			result.Groups = append(result.Groups, HashGroup{
				Hash: h,
				InA:  sorted,
			})
		} else {
			result.Unique = append(result.Unique, sorted[0])
		}
	}

	sortHashGroups(result.Groups)
	sort.Strings(result.Unique)

	result.HasDuplicates = len(result.Groups) > 0

	return result, nil
}

// hashDir scans root and returns a map of SHA-256 hash → []relPath.
// It is the shared primitive for both CompareDir and FindDuplicates.
//
// Depth semantics: depth is the number of path separators in the relative path.
//
//	"file.txt"     → depth 0  (root file)
//	"sub/file.txt" → depth 1  (one level deep)
//
// With Depth=0 only root files are included; with Depth=-1 there is no limit.
func hashDir(root string, opts DirOptions) (map[string][]string, error) {
	result := make(map[string][]string)

	err := filepath.WalkDir(root, func(abs string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, relErr := filepath.Rel(root, abs)
		if relErr != nil {
			return relErr
		}

		// Skip the root entry itself.
		if rel == "." {
			return nil
		}

		// Depth is the number of path separators in the relative path.
		depth := strings.Count(rel, string(filepath.Separator))

		if opts.Depth >= 0 {
			if d.IsDir() {
				// A directory at depth d only contains files at depth d+1.
				// If d >= Depth, those files would exceed the limit: skip the dir entirely.
				if depth >= opts.Depth {
					return filepath.SkipDir
				}
			} else {
				if depth > opts.Depth {
					return nil
				}
			}
		}

		// Apply ignore patterns.
		if opts.Matcher != nil && opts.Matcher.Match(rel) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Compute SHA-256 for regular files only.
		if !d.IsDir() {
			h, err := hash.FileHash(abs)
			if err != nil {
				return err
			}
			result[h] = append(result[h], rel)
		}

		return nil
	})

	return result, err
}

// sortedCopy returns a sorted copy of ss.
func sortedCopy(ss []string) []string {
	cp := make([]string, len(ss))
	copy(cp, ss)
	sort.Strings(cp)
	return cp
}

// sortHashGroups sorts a slice of HashGroup by Hash for deterministic output.
func sortHashGroups(groups []HashGroup) {
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Hash < groups[j].Hash
	})
}
