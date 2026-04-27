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

> **Last updated**: 2026-04-27
> **Current version**: pre-v1.0 (no code yet)
> **Active branch**: `dev`

---

## Implementation Progress

Reference: `docs/IMPLEMENTATION_PLAN.md` — v1.0 Implementation Order

| Step | Description | Status |
|------|-------------|--------|
| 1 | Scaffolding: `go mod init`, directory structure, Dockerfile, Makefile | ⬜ |
| 2 | Package `internal/comparator` — File (TDD) | ⬜ |
| 3 | Package `internal/output` (TDD) | ⬜ |
| 4 | Package `cmd` — Cobra command + flags | ⬜ |
| 5 | Integration: `main.go` wiring | ⬜ |
| 6 | Package `internal/ignore` (TDD) | ⬜ |
| 7 | Package `internal/comparator` — Directory (TDD) | ⬜ |
| 8 | Polish: colored output, `.fmatchignore.example`, README | ⬜ |
| 9 | Release: `.goreleaser.yaml`, cross-platform build | ⬜ |
| 10 | v2.0 — Embedded Web UI | ⬜ |

---

## Last Completed

- Defined and finalized the complete documentation structure for the project.
- Aligned `PROJECT_CONTEXT.md`: translated to English, removed stale entries, added branching strategy and all closed decisions (including #6: no separate backlog).
- Decision: no separate backlog/roadmap file. Task order in `IMPLEMENTATION_PLAN.md`, current state here.

---

## In Progress

_Nothing in progress. Ready to start Step 1: Scaffolding._

---

## Next Step

**Step 1 — Scaffolding**

1. `go mod init github.com/<user>/fmatch`
2. Create directory structure: `cmd/`, `internal/comparator/`, `internal/output/`, `internal/ignore/`
3. Create `Dockerfile` (golang:1.24-alpine)
4. Create `Makefile` with targets: `build`, `test`, `lint`, `cross-compile`
5. Create `main.go` (minimal entry point)

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
