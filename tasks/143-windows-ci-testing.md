---
id: "143"
title: "Add Windows CI testing"
status: pending
priority: low
effort: small
tags: [ci, windows, testing]
created: 2026-02-17
---

# Add Windows CI testing

## Objective

Ensure taskmd's Go tests pass on Windows by adding `windows-latest` to the CI matrix. This catches platform-specific issues like path separators, file permissions, and line endings.

## Tasks

- [ ] Add `windows-latest` to the test job matrix in `.github/workflows/ci.yml` using `runs-on` strategy
- [ ] Add `windows-latest` to the build-cli job matrix
- [ ] Keep the lint job on `ubuntu-latest` only (golangci-lint action works best on Linux)
- [ ] Fix any test failures caused by Windows-specific behavior (path separators, file operations)

## Acceptance Criteria

- CI runs Go tests on both `ubuntu-latest` and `windows-latest`
- All tests pass on Windows
- Lint job remains Linux-only
