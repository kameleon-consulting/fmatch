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
> **Current version**: pre-v1.0 (Step 8 complete — Step 9 starting)
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
| 7 | Package `internal/comparator` — Directory (TDD) | ✅ |
| 8 | Polish: colored output, `.fmatchignore.example`, README | ✅ |
| 9 | Release: `.goreleaser.yaml`, cross-platform build | ⬜ |
| 10 | v2.0 — Embedded Web UI | ⬜ |

---

## Last Completed

**Step 8 — Polish ✅**
- `ignore.LoadFileAndPatterns`: combina file + pattern `-i` in un unico Matcher (+ 3 test)
- `output.FormatDir`: formatta `DirResult` per tutti i livelli verbosità (+ 4 test)
- `cmd/root.go`: `loadMatcher` helper, directory stub → `CompareDir` + `FormatDir`
- `.fmatchignore.example`: template con pattern comuni
- `README.md`: installazione, esempi, flag, exit codes

---

## In Progress

**Step 9 — Release: `.goreleaser.yaml`, cross-platform build**

---

## Next Step

**Step 9 — Release pipeline**

1. Creare `.goreleaser.yaml` (cross-platform: linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64)
2. Configurare ldflags per iniettare `Version` al build time
3. Testare build locale con `goreleaser build --snapshot --clean`
4. Aggiornare `Makefile` con target `release`
5. Commit → tag `v0.1.0` → `git push`

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
