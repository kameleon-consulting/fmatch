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
> **Current version**: pre-v1.0 (Step 6 complete — Step 7 starting)
> **Active branch**: `dev`

---

## Implementation Progress

Reference: `docs/IMPLEMENTATION_PLAN.md` — v1.0 Implementation Order

| Step | Description | Status |
|------|-------------|--------|
| 1 | Scaffolding: `go mod init`, directory structure, Dockerfile, Makefile | ✅ |
| 2 | Package `internal/comparator` — File (TDD) | ✅ |
| 3 | Package `internal/output` (TDD) | ✅ |
| 4 | Package `cmd` — Cobra command + flags | ✅ |
| 5 | Integration: `main.go` wiring | ✅ |
| 6 | Package `internal/ignore` (TDD) | ✅ |
| 7 | Package `internal/comparator` — Directory (TDD) | ⬜ |
| 8 | Polish: colored output, `.fmatchignore.example`, README | ⬜ |
| 9 | Release: `.goreleaser.yaml`, cross-platform build | ⬜ |
| 10 | v2.0 — Embedded Web UI | ⬜ |

---

## Last Completed

**Step 6 — Package `internal/ignore` ✅**
- `go.mod`/`go.sum`: aggiunto `github.com/sabhiram/go-gitignore v0.0.0-20210923224102`
- `internal/ignore/ignore.go`: `Matcher` — `LoadFile` (fallback su file mancante), `LoadPatterns`, `Match()`
- `internal/ignore/ignore_test.go`: 8 test (pattern semplici, file mancante, commenti, righe vuote, negazione `!`, `**`, LoadPatterns, lista vuota)
- All tests pass with `-race` detector

---

## In Progress

**Step 7 — Package `internal/comparator` — Directory (TDD)**

---

## Next Step

**Step 7 — TDD `internal/comparator` — Directory**

1. Write `internal/comparator/dir_test.go` (tests first — TDD)
2. Write `internal/comparator/dir.go` (recursive walk, ignore integration)
3. Run `docker run --rm -v $(pwd):/app fmatch-dev make test` → all green
4. Wire directory support into `cmd/root.go` `runE`
5. Commit → move to Step 8

---

## Open Points

_None. All decisions closed. See `PROJECT_CONTEXT.md` — Closed Decisions._

---

## Key Context

- Binary name: `fmatch`
- Go module path: `github.com/mlabate/fmatch`
- GitHub remote: `git@github.com:mlabate/fmatch.git`
- License: GPL v3
- Go version: 1.24
- All documentation and code: English only
- Working branch: `dev` — merge to `main` only when step is approved and working
