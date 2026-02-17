---
id: "141"
title: "Add `sync down` subcommand"
status: completed
priority: medium
effort: small
tags: [cli, sync]
created: 2026-02-17
---

# Add `sync down` subcommand

## Objective

Restructure `taskmd sync` from a standalone command into a parent command with subcommands, starting with `sync down`. The current `taskmd sync` behavior moves entirely into `taskmd sync down`. Running `taskmd sync` alone should display usage/help listing available subcommands.

## Tasks

- [x] Convert `syncCmd` from a runnable command to a parent command (remove `RunE`, set `Args` to allow subcommands)
- [x] Create `syncDownCmd` as a subcommand of `syncCmd` with `Use: "down"`
- [x] Move all existing flags (`--dry-run`, `--source`, `--conflict`) to `syncDownCmd`
- [x] Move `runSync` logic to the `down` subcommand's `RunE`
- [x] Update tests to invoke the `down` subcommand
- [x] Update help text and examples to reference `taskmd sync down`

## Acceptance Criteria

- `taskmd sync down` performs the same operation as the current `taskmd sync`
- `taskmd sync down --dry-run`, `--source`, and `--conflict` flags work as before
- `taskmd sync` (with no subcommand) prints help/usage showing `down` as an available subcommand
- All existing sync tests pass against the new subcommand structure
- No breaking references to the old `sync` invocation remain in code or docs
