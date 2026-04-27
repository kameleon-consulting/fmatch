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
- **GoReleaser** for cross-compilation (future activation)
- **Mandatory TDD**

## Roadmap

- **v1.0**: Full CLI (working, distributable binary)
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
| 7 | Versioning strategy | ✅ **No tags during development.** First public release = `v1.0.0` on merge to `main` after Step 9 (all v1.0 steps complete). `v2.0.0` = Web UI (Step 10). SemVer + Git tags, managed via GoReleaser. |
| 8 | Go dependencies (Cobra, go-gitignore) | ✅ **Added in Step 1 (Scaffolding).** Deps are already closed decisions; `go.sum` must be committed with scaffolding so the Dockerfile `go mod download` layer is deterministic. |
| 9 | Empty internal directories | ✅ **No `.gitkeep` files.** Directories (`internal/comparator/`, `internal/output/`, `internal/ignore/`) are created naturally when the first Go file (test) is added in their respective steps. Structure is documented in `IMPLEMENTATION_PLAN.md`. |
| 10 | `cmd/root.go` scope in Step 1 | ✅ **Minimal Cobra root command only.** `Use`, `Short`, `Long`, `RunE` stub. No flags. Flags are added in Step 4 with their integration tests. |

## Branching Strategy

- `dev` — working branch. All code, architecture docs, `CURRENT_STATUS.md`, and analysis files are committed here.
- `main` — public/release branch. Only approved, working code and polished docs (README, CONTRIBUTING, CHANGELOG). Merged from `dev` via PR.

## Workspace

- Path: `/Users/mario/repos/001-sandbox/files_diff`
- State: Step 1 partial. Files committed on `dev`: `go.mod`, `go.sum`, `main.go`, `cmd/root.go`
- Next commit: `Dockerfile` + `Makefile` → Step 1 complete

## Next Steps

1. Create `Dockerfile` (golang:1.24-alpine, with `make` and `git`) — Step 1 pending
2. Create `Makefile` (targets: build, test, lint, cross-compile) — Step 1 pending
3. Commit → Step 1 ✅ complete → Step 2: TDD `internal/comparator` (file comparison)
