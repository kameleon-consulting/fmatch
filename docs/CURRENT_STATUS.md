# fmatch — Current Status

<!-- MAINTENANCE RULES
PURPOSE: Session memory. Read this FIRST at the start of every working session.
         Update this LAST at the end of every working session.
UPDATE WHEN: end of every session, no exceptions.
HOW TO UPDATE:
  1. Move "In Progress" → "Last Completed" (summarize what was done)
  2. Move "Next Step" → "In Progress"
  3. Set a new "Next Step" based on IMPLEMENTATION_PLAN.md
  4. Update "Open Points" with any unresolved decisions or blockers
  5. Update "Last Updated" date
NEVER: leave this file reflecting a state older than the last session.
-->

> **Last updated**: 2026-04-29
> **Current version**: v1.1.0 — **RELEASED** (published on GitHub Releases)
> **Next version**: v2.0.0 — Embedded Web UI (deferred — no timeline)
> **Active branch**: `dev`

---

## Implementation Progress

Reference: `docs/IMPLEMENTATION_PLAN.md`

**v1.0 — CLI (complete)**

| Step | Description | Status |
|------|-------------|--------|
| 1 | Scaffolding: `go mod init`, directory structure, Dockerfile, Makefile | ✅ |
| 2 | Package `internal/comparator` — File (TDD) | ✅ |
| 3 | Package `internal/output` (TDD) | ✅ |
| 4 | Package `cmd` — Cobra command + flags | ✅ |
| 5 | Integration: `main.go` wiring | ✅ |
| 6 | Package `internal/ignore` (TDD) | ✅ |
| 7 | Package `internal/comparator` — Directory (TDD) | ✅ |
| 8 | Polish: colored output, `.fmatchignore.example`, README | ✅ |
| 9 | Release: `.goreleaser.yaml`, cross-platform build | ✅ |

**v1.1 — Directory comparison redesign (complete ✅)**

| Step | Description | Status |
|------|-------------|--------|
| v1.1-1 | `internal/hash` — new package `FileHash` (TDD) | ✅ |
| v1.1-2 | `internal/comparator/dir.go` — hash-based rewrite + `FindDuplicates` (TDD) | ✅ |
| v1.1-3 | `internal/output/formatter.go` — `FormatDirCompare` + `FormatDuplicates` | ✅ |
| v1.1-4 | `cmd/root.go` — `RangeArgs(1,2)`, 1-arg vs 2-arg routing | ✅ |
| v1.1-5 | Verify, README, CHANGELOG, CONTRIBUTING, CI, tag `v1.1.0` | ✅ |

**v2.0 — Embedded Web UI (future)**

| Step | Description | Status |
|------|-------------|--------|
| 10 | Embedded Web UI (`fmatch --ui`) | ⬜ |

---

## Last Completed

**v1.1 — Directory comparison redesign ✅**
- `internal/hash`: new package `FileHash` (SHA-256, TDD — 5 tests)
- `internal/comparator/dir.go`: full hash-based rewrite (`CompareDir` + `FindDuplicates` + `hashDir`)
- `internal/output/formatter.go`: `FormatDirCompare` + `FormatDuplicates`; `fileHash` moved to `hash.FileHash`
- `cmd/root.go`: `RangeArgs(1,2)`, routing 1-arg (single-dir duplicate detection) vs 2-arg
- Open-source artifacts: `CHANGELOG.md`, `CONTRIBUTING.md`, `.github/workflows/ci.yml`
- `README.md` updated; license corrected MIT → GPL v3
- Full test suite (`go test -race ./...`): all packages green ✅
- Module path migrated from `github.com/mlabate/fmatch` to `github.com/kameleon-consulting/fmatch`
- Goreleaser moved to Docker execution (`goreleaser/goreleaser:latest`); `.env` support added
- `v1.1.0` tagged on `main` HEAD and published via `make release` (GoReleaser → GitHub Releases)

---

## In Progress

*No active task — v1.1.0 released.*

---

## Next Step

No next step planned. v1.1.0 is released. Options for future work (not yet decided):
- Implement `--no-follow-symlinks` (flag declared but not functional in walk logic)
- Plan v2.0 (Embedded Web UI) — requires separate architectural decisions

---

## Open Points

- `--no-follow-symlinks` flag declared but not implemented in walk logic (deferred post-v1.1)
- v2.0 (Web UI) requires separate architectural decisions — deferred
- **FIXED in v1.1**: v1.0 directory comparison was path-based → now hash-based (SHA-256)

---

## Key Context

- Binary name: `fmatch`
- Go module path: `github.com/kameleon-consulting/fmatch`
- GitHub remote: `git@github.com:kameleon-consulting/fmatch.git`
- GitHub Releases: `https://github.com/kameleon-consulting/fmatch/releases`
- License: GPL v3
- Go version: 1.24
- All documentation and code comments: English only
- Working branch: `dev` — merge to `main` only when step is approved and working
- Release process: `make release` (GoReleaser via Docker) — requires `GITHUB_TOKEN` in `.env`
- Dev environment: Docker (`fmatch-dev` image) — `docker build -t fmatch-dev .`
