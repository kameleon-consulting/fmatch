# fmatch — Implementation Plan

Go CLI tool to compare files and directories, determining exact equality. Cross-platform (macOS/Linux/Windows).

---

## Roadmap

- **v1.0 — CLI**: full terminal tool, distributable as a single binary
- **v2.0 — Embedded Web UI**: `fmatch --ui` flag starts a local HTTP server, opens the browser, vanilla HTML/CSS/JS frontend embedded in the binary via Go `embed` package. Same comparison core as v1.

---

## Technology Choices

### Language: Go
- Single binary, zero runtime dependencies
- Native cross-compilation (`GOOS`/`GOARCH`)
- Excellent stdlib for file I/O
- Industry standard for CLI tools

### External Dependencies

| Library | Purpose | Rationale |
|---------|---------|----------|
| `github.com/spf13/cobra` | CLI framework | De-facto standard for Go CLIs. Auto-help, shell completion, POSIX flags. No subcommands needed but provides professional UX |
| `github.com/sabhiram/go-gitignore` | .gitignore-style pattern matching | Lightweight, simple API, supports `*`, `**`, `!` (negation) patterns. More than sufficient for our use case (single `.fmatchignore` file) |

> **Note**: No external dependency for file comparison. Custom byte-by-byte implementation using stdlib (`os`, `io`, `bufio`) for maximum control and performance.

### Build & Release
- **Makefile** for local builds and cross-compilation
- **GoReleaser** (configuration prepared, future activation with GitHub Actions)
- **Docker** for reproducible development environment

---

## Architecture

### Project Structure

```
fmatch/
├── cmd/
│   └── root.go                 # Cobra command definition, flag parsing
├── internal/
│   ├── comparator/
│   │   ├── file.go             # Single file comparison logic
│   │   ├── file_test.go        # File comparison tests
│   │   ├── dir.go              # Directory comparison (hash-based, v1.1)
│   │   └── dir_test.go         # Directory comparison tests
│   ├── hash/
│   │   ├── hash.go             # FileHash(path) — SHA-256 primitive (v1.1)
│   │   └── hash_test.go
│   ├── ignore/
│   │   ├── ignore.go           # .fmatchignore pattern loading and matching
│   │   └── ignore_test.go      # Pattern matching tests
│   └── output/
│       ├── formatter.go        # Output formatting per verbosity level
│       └── formatter_test.go   # Formatter tests
├── main.go                     # Minimal entry point
├── go.mod
├── go.sum
├── .env.example                # Template for local environment secrets
├── .fmatchignore.example       # Example exclusion pattern file
├── Makefile                    # Build targets
├── .goreleaser.yaml            # GoReleaser configuration
├── Dockerfile                  # Development environment
├── LICENSE
├── README.md
├── CHANGELOG.md
└── CONTRIBUTING.md
```

### Execution Flow

```
Input: path_a [path_b]
  │
  ├─ 1 argument?
  │   ├─ Is file → Exit 2: "single file argument requires a second path"
  │   └─ Is dir  → FindDuplicates → FormatDuplicates → exit 0/1
  │
  └─ 2 arguments:
      ├─ Both exist? ──No──► Exit 2: error
      ├─ Type mismatch? ──Yes──► Exit 2: "cannot compare file with directory"
      ├─ File vs File ──► CompareFiles → FormatFile
      └─ Dir vs Dir   ──► CompareDir (hash-based) → FormatDirCompare

  compareFiles:
      ├─ Different size? ──► DIFFERENT (exit 1)
      └─ Byte-by-byte comparison
          ├─ Match    ──► IDENTICAL (exit 0)
          └─ Mismatch ──► DIFFERENT (exit 1)

  compareDir (hash-based):
      ├─ hashDir(A) + hashDir(B) → map[sha256][]relPath
      ├─ Hash in both → Matched
      ├─ Hash only in A → OnlyInA
      ├─ Hash only in B → OnlyInB
      └─ Report results ──► exit 0 or exit 1
```

---

## Functional Specifications

### 1. File Comparison (core)

**Algorithm** (in order):
1. `os.Stat()` on both files
2. If sizes differ → **DIFFERENT** (immediate exit, zero I/O)
3. Open both files with `bufio.Reader`
4. Read in **64 KB** chunks (optimal for disk I/O and CPU cache)
5. `bytes.Equal()` comparison per chunk
6. First difference → **DIFFERENT** (early exit, reports offset in verbose mode)
7. All chunks equal → **IDENTICAL**

### 2. Directory Comparison

**Algorithm**:
1. Recursive scan of both directories (respecting `--depth` and `.fmatchignore`)
2. Compute relative paths from each directory root
3. Set difference:
   - Files **only in A**
   - Files **only in B**
   - Files **in both**
4. For common files → `compareFiles` on each pair
5. Aggregated report

### 3. Verbosity Levels

| Flag | Level | Output |
|------|-------|--------|
| `-q` / `--quiet` | Quiet | Exit code only (0/1/2). No stdout output |
| *(default)* | Normal | One line: `IDENTICAL` or `DIFFERENT` per file pair. For directories: summary count |
| `-v` / `--verbose` | Verbose | Details: file sizes, full paths. For directories: full file list with status |
| `-vv` | Very Verbose | All verbose output + SHA-256 hash of each file + offset of first difference (if different) |

### 4. CLI Flags

```
fmatch [flags] <path_a> [path_b]

Modes:
  fmatch <file_a> <file_b>   # byte-by-byte file comparison
  fmatch <dir_a> <dir_b>     # hash-based directory comparison
  fmatch <dir>               # find duplicate files within a directory

Flags:
  -q, --quiet              Quiet mode: exit code only
  -v, --verbose            Verbose output (repeatable: -vv for extra detail)
  -d, --depth int          Maximum depth for directories (-1 = unlimited) (default -1)
  -i, --ignore string      Additional patterns to ignore (repeatable)
      --ignore-file string  Path to pattern ignore file (default ".fmatchignore")
      --no-ignore           Disable .fmatchignore file
      --no-follow-symlinks  Do not follow symlinks (default: follow) [declared, not yet implemented]
      --no-color            Disable colored output
  -h, --help               Help
      --version             Version
```

### 5. Exit Codes

| Code | Meaning |
|------|--------|
| `0` | Files/directories are identical |
| `1` | Differences found |
| `2` | Error (file not found, permissions, type mismatch, etc.) |

Aligned with Unix conventions (`diff`, `cmp`).

### 6. `.fmatchignore` File

Same rules as `.gitignore`:
- One pattern per line
- `#` for comments
- `!` for negation
- `*`, `**`, `?` for glob
- Trailing `/` to match directories only

**Example file** (`.fmatchignore.example`):
```
# OS-generated files
.DS_Store
Thumbs.db

# Version control
.git/
.svn/

# IDE
.idea/
.vscode/
*.swp

# Build artifacts
node_modules/
__pycache__/
*.pyc
```

---

## TDD Plan

For each package, tests are written **before** the implementation.

### Tests `internal/comparator`
- Identical files (small and large > 64KB)
- Files differing by size
- Files differing by content (same size)
- Empty files (both empty, one empty)
- Binary files
- Directories with identical files
- Directories with different files
- Directories with missing files (only in A or only in B)
- `--depth` flag respected
- Symlink handling

### Tests `internal/ignore`
- Simple patterns (`*.log`)
- Patterns with `**`
- Negation (`!important.log`)
- Comments ignored
- Empty lines ignored

### Tests `internal/output`
- Correct format for each verbosity level
- Colored vs no-color output

---

## Docker Setup

```dockerfile
FROM golang:1.24-alpine
WORKDIR /app
RUN apk add --no-cache make git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
```

Used for:
- Consistent local development
- Running tests
- Cross-compilation

---

## Implementation Order (v1.0)

1. **Scaffolding**: `go mod init`, directory structure, Dockerfile, Makefile
2. **Package `internal/comparator` — File**: TDD for single file comparison
3. **Package `internal/output`**: TDD for output formatting
4. **Package `cmd`**: Cobra command with flags
5. **Integration**: `main.go` wiring
6. **Package `internal/ignore`**: TDD for pattern matching
7. **Package `internal/comparator` — Directory**: TDD for directory comparison
8. **Polish**: colored output, `.fmatchignore.example`, README
9. **Release**: `.goreleaser.yaml`, cross-platform build

## Implementation Order (v2.0)

10. **Web UI**: embedded HTTP server, vanilla HTML/CSS/JS frontend, Go `embed` package

---

## Verification Plan

### Automated Tests
- `go test ./...` — all unit tests
- `go vet ./...` — static analysis
- `go build ./...` — compilation check

### Manual Verification
- Manual test with real files on macOS
- Cross-compilation check: `GOOS=linux go build`, `GOOS=windows go build`
- Performance test with large directories

---

## Directory Comparison Redesign v1.1 (confirmed 2026-04-28)

### Problem to fix

The v1.0 directory comparison matched files by **relative path** — incorrect.
The actual requirement is: match by **content (SHA-256)**, regardless of name.

---

### New data structures (`internal/comparator/dir.go`)

```go
// HashGroup groups files that share the same SHA-256 hash.
// Used for both 2-dir comparison and 1-dir duplicate detection.
type HashGroup struct {
    Hash string   // SHA-256 hex (64 chars)
    InA  []string // relative paths in A with this hash (1+ entries)
    InB  []string // relative paths in B with this hash (1+ entries; nil in DuplicateResult)
}

// DirCompareResult: result of a hash-based two-directory comparison.
type DirCompareResult struct {
    Identical bool        // true only if both OnlyInA and OnlyInB are empty
    Matched   []HashGroup // hashes present in both dirs; InA and InB list ALL files with that hash
    OnlyInA   []HashGroup // hashes present only in A (no file with that content in B)
    OnlyInB   []HashGroup // hashes present only in B (no file with that content in A)
}

// DuplicateResult: result of duplicate detection within a single directory.
type DuplicateResult struct {
    HasDuplicates bool
    Groups        []HashGroup // groups with 2+ files (InA used; InB always nil)
    Unique        []string    // relative paths of files with a unique hash
}
```

**Confirmed matching rule:** all files with the same hash from either side are grouped together. Example: A has `f1.txt` and `f2.txt` with hash X, B has `f3.txt` with hash X → HashGroup{Hash: X, InA: ["f1.txt","f2.txt"], InB: ["f3.txt"]} → **Matched**.

---

### New functions (`internal/comparator/dir.go`)

```go
// CompareDir compares two directories by matching files on SHA-256 hash.
func CompareDir(pathA, pathB string, opts DirOptions) (DirCompareResult, error)

// FindDuplicates finds files with identical content within a single directory.
func FindDuplicates(path string, opts DirOptions) (DuplicateResult, error)

// hashDir scans a directory and returns a map of hash → []relPath.
func hashDir(root string, opts DirOptions) (map[string][]string, error)
```

`hashDir` is the shared primitive for both. Replaces `walkDir` + path-based matching.

---

### `fileHash` refactoring (circular dependency)

`fileHash` currently lives in `output/formatter.go`. To be used by `comparator/dir.go`
without creating a circular dependency (`comparator` → `output` → `comparator`), it must be moved.

**Decision:** create `internal/hash/hash.go` with a single exported function:
```go
// FileHash computes the SHA-256 hash of a file and returns it as a lowercase hex string.
func FileHash(path string) (string, error)
```
Both `comparator` and `output` import `internal/hash`. No circular dependency.

---

### CLI changes (`cmd/root.go`)

- `cobra.ExactArgs(2)` → `cobra.RangeArgs(1, 2)`
- Routing in `runE`:
  - 1 argument → must be a directory → `FindDuplicates` → `FormatDuplicates`
  - 1 argument → is a file → exit 2: `"fmatch: single file argument requires a second path to compare against"`
  - 2 arguments, both files → `CompareFiles` (unchanged)
  - 2 arguments, both dirs → `CompareDir` hash-based → `FormatDirCompare`
  - 2 arguments, mixed file/dir → exit 2: `"fmatch: cannot compare file with directory"`

---

### New output functions (`internal/output/formatter.go`)

`FormatDir` is **removed** and replaced by:

```go
// FormatDirCompare formats a DirCompareResult for all verbosity levels.
func FormatDirCompare(result comparator.DirCompareResult, opts Options) (string, error)

// FormatDuplicates formats a DuplicateResult for all verbosity levels.
func FormatDuplicates(result comparator.DuplicateResult, opts Options) (string, error)
```

**FormatDirCompare output:**

| Level | Output |
|-------|--------|
| Quiet | `""` |
| Normal | `IDENTICAL` or `DIFFERENT\n  N matched · N only in A · N only in B` |
| Verbose/VV | Label + path_a/path_b (with file count) + matched list with abbreviated hash + only in A + only in B |

Verbose example:
```
DIFFERENT
  path_a: /dirA (5 files)
  path_b: /dirB (6 files)
  matched (2):
    [a3f9c1d2] dirA/report.pdf, dirA/report_copy.pdf  ↔  dirB/report_final.pdf
    [b12c3e4f] dirA/img.png                           ↔  dirB/foto.png
  only in A (1):
    [ff001122] dirA/old.txt
  only in B (3):
    [cc112233] dirB/new1.txt, dirB/new2.txt
    [dd223344] dirB/extra.log
```

**FormatDuplicates output:**

| Level | Output |
|-------|--------|
| Quiet | `""` |
| Normal | `N duplicate groups found` or `No duplicates found` |
| Verbose/VV | Label + for each group: abbreviated hash, count, file list |

Verbose example:
```
3 duplicate groups found
  [a3f9c1d2] (2 files):
    report.pdf
    report_final.pdf
  [b12c3e4f] (3 files):
    img.png
    foto.png
    backup.png
```

---

### Files to modify

| File | Change |
|------|--------|
| `internal/hash/hash.go` | **New** — `FileHash(path) (string, error)` |
| `internal/hash/hash_test.go` | **New** — TDD |
| `internal/comparator/dir.go` | **Rewrite** — `CompareDir` + `FindDuplicates` + `hashDir` |
| `internal/comparator/dir_test.go` | **Rewrite** — TDD (written before implementation) |
| `internal/output/formatter.go` | **Modify** — remove `FormatDir`, add `FormatDirCompare` + `FormatDuplicates`; `fileHash` → use `hash.FileHash` |
| `internal/output/formatter_test.go` | **Modify** — update/add tests |
| `cmd/root.go` | **Modify** — `RangeArgs(1,2)`, 1-arg routing, update output calls |
| `CHANGELOG.md` | **New** — retroactive v0.1.0 entry + v1.1.0 entry (BREAKING noted) |
| `CONTRIBUTING.md` | **New** — contribution guidelines |
| `.github/workflows/ci.yml` | **New** — CI pipeline (test + vet on push/PR to `main` and `dev`) |

---

### Implementation order (v1.1)

1. `internal/hash` — TDD + implementation (foundation for everything else)
2. `internal/comparator/dir.go` — TDD + rewrite (uses `hash.FileHash`)
3. `internal/output/formatter.go` — update dir formatting functions
4. `cmd/root.go` — update routing
5. Verify: `go test ./...`, `go build ./...`, manual test
6. Open-source artifacts (before merge to `main`):
   - Create `CHANGELOG.md` (see CHANGELOG notes section)
   - Create `CONTRIBUTING.md`
   - Create `.github/workflows/ci.yml` (see CI spec section)
7. Release:
   - Update `README.md` (new directory comparison behavior, new single-dir usage)
   - Update `CURRENT_STATUS.md`
   - Merge `dev` → `main`
   - Tag `v1.1.0`
   - Run `make release` (GoReleaser publishes versioned archives to GitHub Releases)

---

### CHANGELOG notes

`CHANGELOG.md` **does not exist yet** — it was not created at v0.1.0 release. It must be created before v1.1.0.

Format: [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) style. Content:

```markdown
# Changelog

## [Unreleased]

## [1.1.0] - YYYY-MM-DD
### Changed
- **BREAKING**: directory comparison now matches files by content (SHA-256 hash)
  instead of by relative path. Files with identical content but different names
  or locations are now considered matched.
- `FormatDir` replaced by `FormatDirCompare` and `FormatDuplicates` (internal).

### Added
- Single-directory mode: `fmatch <dir>` finds duplicate files within a directory,
  grouped by content hash.
- New package `internal/hash` with shared `FileHash` function.

## [0.1.0] - 2026-04-28
### Added
- Initial release.
- File comparison: byte-by-byte with early exit, 4 verbosity levels (-q/-v/-vv).
- Directory comparison: recursive walk, depth limit, `.fmatchignore` patterns.
- Exit codes: 0 (identical), 1 (different), 2 (error).
- Cross-platform binaries: linux/darwin (amd64/arm64), windows/amd64.
- GoReleaser pipeline.
```

> **Note on distributions**: `dist/` is excluded from git (`.gitignore`). Versioned
> binary archives are published automatically by GoReleaser to GitHub Releases when
> `make release` is run with a valid `GITHUB_TOKEN` and a Git tag.

---

### CI spec (`.github/workflows/ci.yml`)

Triggers: push and pull_request targeting `main` or `dev`.

Steps:
1. `actions/checkout@v4`
2. `actions/setup-go@v5` — Go version from `go.mod`
3. `go vet ./...`
4. `go test -race ./...`

No build artifact upload needed (GoReleaser handles releases separately).

```yaml
name: CI
on:
  push:
    branches: [main, dev]
  pull_request:
    branches: [main, dev]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - run: go vet ./...
      - run: go test -race ./...
```

---

### CONTRIBUTING spec (`CONTRIBUTING.md`)

Minimal content required:
- Prerequisites (Go 1.24+, Docker, Make)
- How to run tests locally (`docker run ... make test`)
- Branching strategy (`dev` → PR → `main`)
- Commit convention (Conventional Commits)
- How to open an issue or PR

---

### TDD — Tests to write before implementation


**`internal/hash/hash_test.go`:**
- Known file → expected SHA-256 value
- Empty file → correct hash (SHA-256 of empty string)
- Non-existent file → error

**`internal/comparator/dir_test.go`:**
- Identical dirs, same names → Identical=true, Matched full
- Identical dirs, different names → Identical=true, Matched full (key case)
- Files only in A → OnlyInA
- Files only in B → OnlyInB
- Same hash, many-to-many (A: f1+f2, B: f3) → single HashGroup in Matched
- Empty dirs → Identical=true
- Single dir, no duplicates → HasDuplicates=false
- Single dir, with duplicates → correct groups
- Depth limit respected
- Ignore patterns respected

**`internal/output/formatter_test.go`:**
- `FormatDirCompare` for all verbosity levels (Quiet/Normal/Verbose/VV)
- `FormatDuplicates` for all verbosity levels

