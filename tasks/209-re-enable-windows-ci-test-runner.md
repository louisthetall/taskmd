---
id: "209"
title: "Fix CI: re-enable Windows test runner"
status: pending
priority: medium
effort: medium
tags: [ci, windows]
created: 2026-02-25
---
# Fix CI: re-enable Windows test runner

## Objective

The Windows test job in CI (`.github/workflows/ci.yml`) was disabled because `go test -race` with `-coverprofile=coverage.txt` either hangs or takes excessively long on `windows-latest`. The test matrix was reduced to `ubuntu-latest` only. Re-enable Windows testing once the root cause is resolved.

## Tasks

- [ ] Investigate why the test step hangs/times out on Windows (PowerShell argument parsing, race detector performance, or both)
- [ ] Try `shell: bash` to bypass PowerShell argument mangling of `-coverprofile=coverage.txt`
- [ ] Consider disabling `-race` on Windows if that's the performance bottleneck
- [ ] Re-enable `windows-latest` in the test matrix in `.github/workflows/ci.yml`
- [ ] Verify the CI run passes end-to-end on Windows

## Acceptance Criteria

- `Test (windows-latest)` job is present and passing in CI
- No regressions on `Test (ubuntu-latest)`
