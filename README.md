# fmatch

**fmatch** is a fast, cross-platform CLI tool to compare files and directories for exact equality, or find duplicate files within a directory.

## Features

- **Byte-exact file comparison** — no false positives, early exit on first difference
- **Hash-based directory comparison** — matches files by content (SHA-256), regardless of name or path
- **Duplicate detection** — single-directory mode finds all files with identical content
- **Ignore patterns** — `.fmatchignore` file (`.gitignore` syntax) + inline `-i` flags
- **Multiple verbosity levels** — quiet, normal, verbose, very-verbose (with SHA-256 hashes)
- **Colored output** — ANSI colors, disable with `--no-color`
- **Unix exit codes** — `0` identical/no duplicates, `1` different/duplicates found, `2` error

---

## Installation

### Download pre-built binary

Download the latest release from the [Releases](https://github.com/mlabate/fmatch/releases) page.

### Build from source

Requirements: Go 1.24+, Docker (for reproducible builds)

```bash
git clone https://github.com/mlabate/fmatch.git
cd fmatch
docker build -t fmatch-dev .
docker run --rm -v $(pwd):/app fmatch-dev make build
# Binary: ./fmatch
```

---

## Usage

```
fmatch [flags] <path_a> [path_b]
```

| Invocation | Behaviour |
|-----------|-----------|
| `fmatch <file_a> <file_b>` | Byte-by-byte file comparison |
| `fmatch <dir_a> <dir_b>` | Hash-based directory comparison (content, not names) |
| `fmatch <dir>` | Find duplicate files within a single directory |
| `fmatch <file>` | Error — a single file argument requires a second path |

### Examples

```bash
# Compare two files
fmatch file_a.txt file_b.txt

# Compare two directories (files matched by content, not by name)
fmatch dir_a/ dir_b/

# Find duplicate files within a directory
fmatch /path/to/photos/

# Quiet mode — exit code only
fmatch -q file_a.txt file_b.txt; echo $?

# Verbose — show paths, sizes, matched groups
fmatch -v dir_a/ dir_b/

# Very verbose — also show SHA-256 hashes and diff offset
fmatch -vv file_a.txt file_b.txt

# Limit directory depth
fmatch -d 2 dir_a/ dir_b/

# Ignore patterns
fmatch -i "*.log" -i "*.tmp" dir_a/ dir_b/

# Use a custom ignore file
fmatch --ignore-file /path/to/patterns.txt dir_a/ dir_b/

# Disable ignore file
fmatch --no-ignore dir_a/ dir_b/

# No colored output
fmatch --no-color file_a.txt file_b.txt

# Save output to file (--no-color avoids ANSI codes)
fmatch --no-color -v dir_a/ dir_b/ > result.txt
```

---

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-q`, `--quiet` | false | Exit code only, no output |
| `-v`, `--verbose` | — | Verbose output (repeat for `-vv`) |
| `-d`, `--depth int` | -1 | Max directory depth (-1 = unlimited) |
| `-i`, `--ignore string` | — | Pattern to ignore (repeatable) |
| `--ignore-file string` | `.fmatchignore` | Path to ignore file |
| `--no-ignore` | false | Disable `.fmatchignore` file |
| `--no-follow-symlinks` | false | Do not follow symlinks |
| `--no-color` | false | Disable colored output |
| `--version` | — | Print version and exit |

---

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Files/directories are identical, or no duplicates found |
| `1` | Differences found, or duplicate files found |
| `2` | Error (file not found, permission denied, type mismatch) |

Same convention as `diff` and `cmp` — suitable for use in scripts.

---

## `.fmatchignore`

Place a `.fmatchignore` file in your working directory. Syntax follows `.gitignore` rules:

```
# comment
*.log
*.tmp
node_modules/
!important.log   # negation: never ignore this file
```

See `.fmatchignore.example` for a ready-to-use template.

---

## License

GPL v3 — see [LICENSE](LICENSE).
