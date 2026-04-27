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

## Branching Strategy

- `dev` — working branch. All code, architecture docs, `CURRENT_STATUS.md`, and analysis files are committed here.
- `main` — public/release branch. Only approved, working code and polished docs (README, CONTRIBUTING, CHANGELOG). Merged from `dev` via PR.

## Workspace

- Path: `/Users/mario/repos/001-sandbox/files_diff`
- State: documentation only (`docs/`), no code yet

## Next Steps

1. Complete documentation structure (`CURRENT_STATUS.md`, `docs/working/` convention)
2. Scaffolding: `go mod init`, directory structure, `Dockerfile`, `Makefile`
3. Proceed TDD package by package per implementation order in `IMPLEMENTATION_PLAN.md`
