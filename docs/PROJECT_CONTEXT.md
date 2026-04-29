# fmatch — Project Context

## Initial Request

The user needs a program that:
- Compares any two files and determines if they are exactly identical
- Accepts directories (with file list output) or single files as input
- Is cross-platform (compiled binary)

## Collected Requirements

1. **Language**: Go (confirmed as best choice)
2. **Output**: configurable verbosity (quiet by default)
3. **Directory depth**: configurable via `--depth` flag (default: unlimited)
4. **Exclusion patterns**: `.gitignore`-style (`.fmatchignore` file)
5. **Scope**: open-source, public distribution

## Confirmed Technology Choices

- **Go** as language
- **Cobra** as CLI framework
- **go-gitignore** (`github.com/sabhiram/go-gitignore`) for pattern matching
- File comparison **byte-by-byte with early exit** (no external libs, pure stdlib)
- **4 verbosity levels**: quiet (-q), normal (default), verbose (-v), very verbose (-vv)
- **Exit codes**: 0 (identical), 1 (different), 2 (error) — Unix convention
- **Docker** for development environment
- **GoReleaser** for cross-compilation and release publishing (active — runs via Docker, `make release`)
- **Mandatory TDD**

## Roadmap

- **v1.0**: Full CLI (released 2026-04-28)
- **v1.1**: Hash-based directory comparison + duplicate detection (released 2026-04-29)
- **v2.0**: Embedded Web UI (`fmatch --ui`) — local HTTP server + vanilla HTML/CSS/JS frontend embedded in the binary via Go `embed`

## Closed Decisions

| # | Question | Decision |
|---|----------|-----------|
| 1 | Binary name | ✅ **`fmatch`** |
| 2 | Documentation language | ✅ **All EN** (open-source project) |
| 3 | License | ✅ **GPL v3** |
| 4 | Default directory depth | ✅ **-1 (unlimited)** |
| 5 | Symlink handling | ✅ **Follow by default + `--no-follow-symlinks` to disable** |
| 6 | Task tracking / backlog | ✅ **No separate backlog file.** Implementation order lives in `IMPLEMENTATION_PLAN.md`. Current state (step in progress, next step) lives in `CURRENT_STATUS.md`. Public task tracking via GitHub Issues when the repo goes live. |
| 7 | Versioning strategy | ✅ **Semantic versioning + Git tags, managed via GoReleaser.** Released versions: `v1.0.0` (2026-04-28), `v1.1.0` (2026-04-29). Next planned: `v2.0.0` (Web UI). Release command: `make release` (GoReleaser via Docker). Requires `GITHUB_TOKEN` in `.env`. |
| 8 | Go dependencies (Cobra, go-gitignore) | ✅ **Added in Step 1 (Scaffolding).** Deps are already closed decisions; `go.sum` must be committed with scaffolding so the Dockerfile `go mod download` layer is deterministic. |
| 9 | Empty internal directories | ✅ **No `.gitkeep` files.** Directories (`internal/comparator/`, `internal/output/`, `internal/ignore/`) are created naturally when the first Go file (test) is added in their respective steps. Structure is documented in `IMPLEMENTATION_PLAN.md`. |
| 10 | `cmd/root.go` scope in Step 1 | ✅ **Minimal Cobra root command only.** `Use`, `Short`, `Long`, `RunE` stub. No flags. Flags are added in Step 4 with their integration tests. |

## Branching Strategy

- `dev` — working branch. All code, architecture docs, `CURRENT_STATUS.md`, and analysis files are committed here.
- `main` — public/release branch. Only approved, working code and polished docs (README, CONTRIBUTING, CHANGELOG). Merged from `dev` via PR.



## Next Steps

> Per lo stato corrente e il prossimo step, vedere `CURRENT_STATUS.md`.


## Design Bug — Directory Comparison (detected 2026-04-28)

The original directory comparison matched files by **relative path** — incorrect behavior relative to the actual requirement.

**Confirmed real requirement:**
> "I want to know if two files are identical regardless of their name."

Comparison must be hash-based (SHA-256), not name/path-based.

## Additional Confirmed Requirements (2026-04-28)

| # | Requirement | Detail |
|---|-------------|--------|
| A | `fmatch dirA dirB` — hash-based | Files are matched by content (SHA-256), not by relative path. Files with identical content but different names are considered matched. |
| B | `fmatch dirA` — find duplicates | With a single directory argument: find all files with identical content within that directory, grouped by hash. |
| C | `fmatch fileA` — error | A single file argument has no semantic meaning: exit 2 with explanatory message. |
| D | `fmatch fileA fileB` — unchanged | Single file comparison remains unchanged (already correct: content-based, name irrelevant). |

## Additional Closed Decisions (2026-04-28)

| # | Question | Decision |
|---|----------|---------|
| 11 | CLI args | ✅ `RangeArgs(1, 2)` — 1 or 2 arguments. 1 arg = directory only (duplicates); 2 args = files or directories. |
| 12 | Directory data structure | ✅ `HashGroup{Hash, InA []string, InB []string}` as common primitive for both cases (2-dir compare and 1-dir duplicates). |
| 13 | `fileHash` refactoring | ✅ Move `fileHash` from `output/formatter.go` to a shared package (`internal/hash`) to avoid circular dependencies. |

