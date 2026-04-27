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
│   │   ├── dir.go              # Directory comparison logic
│   │   └── dir_test.go         # Directory comparison tests
│   ├── ignore/
│   │   ├── ignore.go           # .fmatchignore pattern loading and matching
│   │   └── ignore_test.go      # Pattern matching tests
│   └── output/
│       ├── formatter.go        # Output formatting per verbosity level
│       └── formatter_test.go   # Formatter tests
├── main.go                     # Minimal entry point
├── go.mod
├── go.sum
├── .fmatchignore.example       # Example exclusion pattern file
├── Makefile                    # Build targets
├── .goreleaser.yaml            # GoReleaser configuration (future)
├── Dockerfile                  # Development environment
├── LICENSE
└── README.md
```

### Execution Flow

```
Input: path_a, path_b
  │
  ├─ Both exist? ──No──► Exit 2: error
  │
  ├─ Path type?
  │   ├─ File vs File ──► compareFiles
  │   ├─ Dir vs Dir   ──► compareDirs
  │   └─ Mismatch     ──► Exit 2: type mismatch error
  │
  ├─ compareFiles:
  │   ├─ Different size? ──► DIFFERENT (exit 1)
  │   └─ Byte-by-byte comparison
  │       ├─ Match    ──► IDENTICAL (exit 0)
  │       └─ Mismatch ──► DIFFERENT (exit 1)
  │
  └─ compareDirs:
      ├─ Scan files (depth + ignore)
      ├─ Set difference (only in A, only in B, in both)
      ├─ Compare common files
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
fmatch [flags] <path_a> <path_b>

Flags:
  -q, --quiet              Quiet mode: exit code only
  -v, --verbose            Verbose output (repeatable: -vv for extra detail)
  -d, --depth int          Maximum depth for directories (-1 = unlimited) (default -1)
  -i, --ignore string      Additional patterns to ignore (repeatable)
      --ignore-file string  Path to pattern ignore file (default ".fmatchignore")
      --no-ignore           Disable .fmatchignore file
      --no-follow-symlinks  Do not follow symlinks (default: follow)
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
