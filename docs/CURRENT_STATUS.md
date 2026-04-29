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

> **Last updated**: 2026-04-28
> **Current version**: v1.0.0 (Steps 1-9 complete — v1.0 released)
> **Next version**: v1.1.0 — directory comparison redesign (hash-based)
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

**v1.1 — Directory comparison redesign (next)**

| Step | Description | Status |
|------|-------------|--------|
| v1.1-1 | `internal/hash` — new package `FileHash` (TDD) | ⬜ |
| v1.1-2 | `internal/comparator/dir.go` — hash-based rewrite + `FindDuplicates` (TDD) | ⬜ |
| v1.1-3 | `internal/output/formatter.go` — `FormatDirCompare` + `FormatDuplicates` | ⬜ |
| v1.1-4 | `cmd/root.go` — `RangeArgs(1,2)`, 1-arg vs 2-arg routing | ⬜ |
| v1.1-5 | Verify, README, CHANGELOG, tag `v1.1.0` | ⬜ |

**v2.0 — Embedded Web UI (future)**

| Step | Description | Status |
|------|-------------|--------|
| 10 | Embedded Web UI (`fmatch --ui`) | ⬜ |

---

## Last Completed

**Step 9 — Release pipeline ✅**
- `.goreleaser.yaml`: 5 piattaforme (linux/darwin/windows × amd64/arm64 meno windows/arm64)
- `Makefile`: target `snapshot` e `release` con goreleaser
- Cross-compile verificato: tutti i binari compilano (`make cross-compile`)
- Merge `dev` → `main`, tag `v1.0.0`

---

## In Progress

*Nessun task attivo — v0.1.0 rilasciata.*

---

## Next Step

**v1.1-1 — Package `internal/hash`** (TDD): create `internal/hash/hash.go` with `FileHash(path string) (string, error)`. This is the building block for all of v1.1 — must be done first.

---

## Open Points

- `--no-follow-symlinks` flag declared but not implemented in walk logic (TODO v0.2.0 — deferred after v1.1)
- v2.0 (Web UI) requires separate architectural decisions — deferred after v1.1
- **KNOWN BUG**: v1.0 compares directories by relative path instead of hash — fixed in v1.1 (see `IMPLEMENTATION_PLAN.md` section "Directory Comparison Redesign v1.1")

---

## Key Context

- Binary name: `fmatch`
- Go module path: `github.com/mlabate/fmatch`
- GitHub remote: `git@github.com:mlabate/fmatch.git`
- License: GPL v3
- Go version: 1.24
- All documentation and code: English only
- Working branch: `dev` — merge to `main` only when step is approved and working
