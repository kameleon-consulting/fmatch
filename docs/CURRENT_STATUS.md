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
> **Current version**: v1.1.0 (Steps v1.1-1 → v1.1-5 complete — v1.1 ready for release)
> **Next version**: v2.0.0 — Embedded Web UI (future)
> **Active branch**: `dev` — ready to merge to `main`

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

**Step v1.1 — Directory comparison redesign ✅**
- `internal/hash`: nuovo package `FileHash` (SHA-256, TDD)
- `internal/comparator/dir.go`: rewrite hash-based (`CompareDir` + `FindDuplicates` + `hashDir`)
- `internal/output/formatter.go`: `FormatDirCompare` + `FormatDuplicates`; `fileHash` → `hash.FileHash`
- `cmd/root.go`: `RangeArgs(1,2)`, routing 1-arg (single-dir) vs 2-arg
- Open-source artifacts: `CHANGELOG.md`, `CONTRIBUTING.md`, `.github/workflows/ci.yml`
- `README.md` aggiornato; fix licenza MIT → GPL v3
- Suite completa (`go test -race ./...`): 50+ test passati, tutti i package ✅

---

## In Progress

*Nessun task attivo — v0.1.0 rilasciata.*

---

## Next Step

**v1.1-5g — Merge `dev` → `main`, tag `v1.1.0`, `make release`** (ultimo step — da eseguire ora).

---

## Open Points

- `--no-follow-symlinks` flag declared but not implemented in walk logic (TODO — deferred post-v1.1)
- v2.0 (Web UI) requires separate architectural decisions — deferred
- **FIXED in v1.1**: v1.0 directory comparison by relative path → now hash-based

---

## Key Context

- Binary name: `fmatch`
- Go module path: `github.com/kameleon-consulting/fmatch`
- GitHub remote: `git@github.com:kameleon-consulting/fmatch.git`
- License: GPL v3
- Go version: 1.24
- All documentation and code: English only
- Working branch: `dev` — merge to `main` only when step is approved and working
