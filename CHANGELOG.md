# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.1.0] - 2026-04-29

### Added

- Single-directory mode: `fmatch <dir>` finds duplicate files within a directory,
  grouped by content hash. Exit code 1 if duplicates are found, 0 otherwise.
- New package `internal/hash` with shared `FileHash(path string) (string, error)` function
  (SHA-256, lowercase hex). Used by both `comparator` and `output` packages.

### Changed

- **BREAKING**: directory comparison (`fmatch <dir_a> <dir_b>`) now matches files by
  content (SHA-256 hash) instead of by relative path. Files with identical content but
  different names or locations are now considered matched.
- `FormatDir` removed and replaced by `FormatDirCompare` and `FormatDuplicates` (internal API).
- CLI now accepts 1 or 2 positional arguments (`RangeArgs(1, 2)`). A single file argument
  returns exit code 2 with an explanatory message.

## [1.0.0] - 2026-04-28

### Added

- Initial release.
- File comparison: byte-by-byte with early exit on first difference. 4 verbosity levels
  (`-q` / default / `-v` / `-vv`).
- Directory comparison: recursive walk, configurable depth (`--depth`), `.fmatchignore`
  pattern exclusions.
- Exit codes: `0` (identical), `1` (different), `2` (error) — aligned with Unix conventions.
- Cross-platform binaries: linux/darwin (amd64/arm64), windows/amd64.
- GoReleaser pipeline for automated release builds.

> **Note on distributions**: `dist/` is excluded from git. Versioned binary archives are
> published automatically by GoReleaser to GitHub Releases when `make release` is run
> with a valid `GITHUB_TOKEN` and a Git tag.
