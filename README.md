# fmatch

**fmatch** is a fast, cross-platform CLI tool to compare two files or directories for exact equality.

## Features

- **Byte-exact comparison** — no false positives
- **Directory comparison** — recursive, with set difference (only-in-A, only-in-B, different)
- **Ignore patterns** — `.fmatchignore` file (`.gitignore` syntax) + inline `-i` flags
- **Multiple verbosity levels** — quiet, normal, verbose, very-verbose (with SHA-256)
- **Colored output** — ANSI colors, disable with `--no-color`
- **Unix exit codes** — `0` identical, `1` different, `2` error

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
# Binary: ./bin/fmatch
```

---

## Usage

```
fmatch [flags] <path_a> <path_b>
```

### Examples

```bash
# Compare two files (default: normal output)
fmatch file_a.txt file_b.txt

# Compare two directories
fmatch dir_a/ dir_b/

# Quiet mode — exit code only
fmatch -q file_a.txt file_b.txt; echo $?

# Verbose — show paths and sizes
fmatch -v file_a.txt file_b.txt

# Very verbose — show SHA-256 and diff offset
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

# Save output to file (--no-color avoids ANSI codes in the file)
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
| `0` | Files/directories are identical |
| `1` | Differences found |
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

MIT — see [LICENSE](LICENSE).
