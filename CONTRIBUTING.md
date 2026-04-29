# Contributing to fmatch

Thank you for your interest in contributing!

## Prerequisites

- [Go 1.24+](https://go.dev/dl/)
- [Docker](https://www.docker.com/) (recommended for reproducible builds)
- [Make](https://www.gnu.org/software/make/)

## Running tests locally

All tests must pass before opening a PR. Use Docker for a consistent environment:

```bash
# Build the development image (first time or after dependency changes)
docker build -t fmatch-dev .

# Run the full test suite with race detector
docker run --rm -v $(pwd):/app fmatch-dev make test

# Or if Go 1.24+ is installed locally
make test
```

## Branching strategy

| Branch | Purpose |
|--------|---------|
| `dev`  | Active development — all work happens here |
| `main` | Public/release branch — merged from `dev` via PR only |

- Create your feature branch from `dev`
- Open a pull request targeting `dev`
- `main` is updated only when a release is ready

## Commit convention

This project uses [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <short description>

Types: feat, fix, docs, test, chore, refactor
Scope: cmd, comparator, hash, ignore, output (package name)

Examples:
  feat(comparator): add hash-based directory comparison
  fix(output): handle empty directory in FormatDuplicates
  docs: update README with single-dir usage
  test(hash): add large-file test case
```

## Opening an issue

- Describe the problem clearly: what you expected vs. what actually happened
- Include the `fmatch` version (`fmatch --version`) and OS
- Attach a minimal reproducible example if possible

## Opening a pull request

- Reference the related issue in the PR description
- Include tests for any new behaviour or bug fix
- Ensure `make test` passes with no race conditions
- Keep commits atomic and well-described (Conventional Commits)
