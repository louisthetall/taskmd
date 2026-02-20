---
id: "174"
title: "Set up e2e test foundation and helpers"
status: pending
priority: medium
effort: medium
type: improvement
tags:
  - testing
  - cli
parent: "173"
created: 2026-02-20
---

# Set up e2e test foundation and helpers

## Objective

Create the e2e test package, build infrastructure, and shared helpers that all subsequent e2e tests will depend on. This is the foundation — no command-specific tests yet, just the scaffolding and a smoke test to prove it works.

## Tasks

- [ ] Create test package (e.g. `apps/cli/internal/e2e/e2e_test.go`)
- [ ] Implement `TestMain` that builds the `taskmd` binary once into a temp directory
- [ ] Implement `run(t, dir, args...) (stdout, stderr, error)` helper that invokes the binary as a subprocess
- [ ] Implement `mustRun(t, dir, args...)` helper that fails the test on non-zero exit
- [ ] Implement `writeTask(t, dir, filename, id, title, status, deps)` helper for creating test task files
- [ ] Isolate tests from user config by overriding `HOME` env var to a temp directory
- [ ] Set `NO_COLOR=1` in subprocess env for deterministic output
- [ ] Add `make e2e` target to Makefile that runs only the e2e test package
- [ ] Ensure `make test` still runs unit/integration tests and does not include e2e tests
- [ ] Add a basic smoke test: `taskmd --help` returns exit code 0 and includes expected output

## Acceptance Criteria

- `make e2e` builds the binary and runs the e2e test package
- `make test` is unaffected and does not run e2e tests
- Smoke test passes: `taskmd --help` exits 0 with usage text
- Each test gets an isolated temp dir and clean environment
- Helpers are reusable and well-documented for subsequent test tasks
