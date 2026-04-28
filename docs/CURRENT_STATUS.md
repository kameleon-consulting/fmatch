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
> **Current version**: pre-v1.0 (Step 5 complete — Step 6 starting)
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
| 6 | Package `internal/ignore` (TDD) | ⬜ |
| 7 | Package `internal/comparator` — Directory (TDD) | ⬜ |
| 8 | Polish: colored output, `.fmatchignore.example`, README | ⬜ |
| 9 | Release: `.goreleaser.yaml`, cross-platform build | ⬜ |
| 10 | v2.0 — Embedded Web UI | ⬜ |

---

## Last Completed

**Step 5 — Integration `RunE` ✅**
- `ExitError` type (exported): porta exit code, testabile senza `os.Exit`
- `runE`: flag reading, stat, type mismatch, `CompareFiles`, `Format`, exit 0/1/2
- `resolveVerbosity`: mappa `-q`/`-v`/`-vv` → `output.Verbosity`
- `Execute()`: ispeziona `ExitError` per `os.Exit(Code)` corretto
- 5 test integrazione: identical, different/exit1, not-found/exit2, type-mismatch/exit2, quiet

---

## In Progress

**Step 6 — Package `internal/ignore` (TDD)**

---

## Next Step

**Step 6 — TDD `internal/ignore` (pattern matching)**

1. Write `internal/ignore/ignore_test.go` (tests first — TDD)
2. Write `internal/ignore/ignore.go` (wraps go-gitignore)
3. Run `docker run --rm -v $(pwd):/app fmatch-dev make test` → all green
4. Commit → move to Step 7

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
