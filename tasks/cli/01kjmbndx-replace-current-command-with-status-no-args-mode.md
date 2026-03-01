---
title: "Replace current command with status no-args mode"
id: "01kjmbndx"
status: completed
priority: medium
type: feature
tags: []
created: "2026-03-01"
---

# Replace current command with status no-args mode

## Objective

Remove the `current` command and consolidate its functionality into the `status` command. When `status` is called with no arguments, it should display all in-progress tasks using the existing rich format. This replaces the single-task compact output of `current` with a more useful multi-task view, while preserving the compact statusline format via a `--statusline` flag.

## Tasks

- [ ] Remove `current` command (code, tests, docs)
  - [ ] Delete `apps/cli/internal/cli/current.go`
  - [ ] Delete `apps/cli/internal/cli/current_test.go`
  - [ ] Remove any e2e tests for `current`
  - [ ] Remove `current` references from documentation (`docs/`, `apps/docs/`)
- [ ] Update `status` command to handle no-args mode
  - [ ] Make the `<query>` argument optional (change from `ExactArgs(1)` to `MaximumNArgs(1)`)
  - [ ] When no args: scan for all in-progress tasks and display them using the rich status format
  - [ ] Respect `--scope` flag to filter by group/directory
  - [ ] Support existing output formats (text, json, yaml) for no-args mode
- [ ] Add `--statusline` flag to `status` command
  - [ ] Output compact `#ID title` format (matching old `current` behavior)
  - [ ] Truncate titles longer than 30 chars with "..."
  - [ ] For multiple in-progress tasks, show first task with `(+N more)` suffix
- [ ] Update tests for `status` command
  - [ ] Test no-args mode with zero, one, and multiple in-progress tasks
  - [ ] Test `--statusline` flag output
  - [ ] Test `--scope` filtering in no-args mode
  - [ ] Test all output formats in no-args mode

## Acceptance Criteria

- `current` command is fully removed (code, tests, docs)
- `taskmd status` (no args) shows all in-progress tasks in rich format
- `taskmd status <query>` continues to work as before
- `taskmd status --statusline` outputs compact format for shell integrations
- `taskmd status --scope cli` filters to only cli group tasks
- All output formats (text, json, yaml) work in no-args mode
- Tests cover all new functionality
