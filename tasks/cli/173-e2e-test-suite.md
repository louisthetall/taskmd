---
id: "173"
title: "Build e2e test suite for CLI"
status: pending
priority: medium
effort: large
type: improvement
tags:
  - testing
  - cli
created: 2026-02-20
---

# Build e2e test suite for CLI

## Objective

Create an end-to-end test suite that invokes the compiled `taskmd` binary as a subprocess (via `exec.Command`) rather than calling `run*` functions directly. This tests the full pipeline: argument parsing, config loading, command execution, output formatting, and exit codes.

## Tasks

- [ ] Create test package at `apps/cli/e2e_test.go` or `apps/cli/internal/e2e/`
- [ ] Implement `TestMain` to build the binary once before all tests
- [ ] Implement test helpers: `run()`, `mustRun()`, `writeTask()`
- [ ] Isolate tests with `t.TempDir()` and overridden `HOME` env var
- [ ] Test full workflows: add -> list -> set -> next chains
- [ ] Test JSON and YAML output parseability across commands
- [ ] Test flag wiring (flags actually affect behavior end-to-end)
- [ ] Test `.taskmd.yaml` config loading and precedence
- [ ] Test stdin/pipe behavior (e.g. `taskmd validate --stdin`)
- [ ] Test dependency resolution across commands (graph, next)
- [ ] Test error cases: unknown commands, malformed task files, missing args
- [ ] Test exit codes for success and failure scenarios
- [ ] Add `make e2e` target to run e2e tests separately from unit tests

## Sub-tasks

- **174** — Set up e2e test foundation and helpers
- **175** — E2e tests for command workflows (depends on 174)
- **176** — E2e tests for error handling and edge cases (depends on 174)
- **177** — E2e tests for config loading and precedence (depends on 174)

## Acceptance Criteria

- E2e tests invoke the compiled binary, not Go functions directly
- Binary is built once per test run via `TestMain`
- Each test is fully isolated (temp dirs, no shared state, no user config leakage)
- Workflow tests cover add, list, set, get, next, graph, validate, board commands
- Error cases verify non-zero exit codes and meaningful stderr output
- JSON/YAML output from commands is parsed and structurally validated
- `make e2e` runs only e2e tests; `make test` continues to run unit/integration tests
- Tests pass in CI with no external dependencies beyond Go toolchain
