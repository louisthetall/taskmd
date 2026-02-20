---
id: "171"
title: "Document missing CLI flags and commands"
status: completed
priority: medium
effort: medium
type: docs
tags:
  - docs
  - cli
created: 2026-02-20
---

# Document missing CLI flags and commands

## Objective

Fill documentation gaps for CLI commands and flags that exist in the CLI but are missing from documentation. Covers both `apps/docs/guide/cli.md` (VitePress) and `docs/guides/cli-guide.md` (standalone).

## Tasks

### VitePress CLI guide (`apps/docs/guide/cli.md`)

- [x] Add `todos` / `todos list` command section with all flags (`--dir`, `--marker`, `--include`, `--exclude`, `--rich`, `--raw-text`, `--format`)
- [x] Add `todos` to Quick Reference table
- [x] Add `next --quick-wins` and `next --critical` flags to `next` command section
- [x] Add `web export` subcommand section with flags (`--output`, `--base-path`)
- [x] Add `web start --readonly` flag

### Standalone CLI guide (`docs/guides/cli-guide.md`)

- [x] Add `add` command section
- [x] Add `search` command section
- [x] Add `verify` command section
- [x] Add `status` command section
- [x] Add `context` command section
- [x] Add `worklog` command section
- [x] Add `import` command section
- [x] Add `spec` command section
- [x] Add `commit-msg` command section
- [x] Add `todos` / `todos list` command section
- [x] Add `completion` to Quick Reference table
- [x] Add `web export` subcommand section
- [x] Add `next --quick-wins` and `next --critical` flags
- [x] Add `web start --readonly` flag
- [x] Add `get --context` flag
- [x] Add `graph --filter` flag

## Acceptance Criteria

- Every CLI command from `taskmd --help` has a corresponding section in both CLI guides
- All flags shown by `<command> --help` are documented in the respective command section
- New documentation follows existing style (description, flags table, examples)
