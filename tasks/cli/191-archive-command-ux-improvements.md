---
id: "191"
title: "Improve archive command UX: positional task ID and interactive confirmation"
status: completed
priority: medium
effort: small
type: improvement
tags:
  - cli
  - ux
touches:
  - cli/archive
created: 2026-02-21
---

# Improve archive command UX: positional task ID and interactive confirmation

## Objective

Improve the `archive` command's user experience with two changes:

1. **Positional task ID** — Accept the task ID as an optional positional argument so users can write `taskmd archive 042` instead of `taskmd archive --id 042`.
2. **Interactive confirmation** — Replace the current behavior (which requires `--yes`/`-y` flag or errors out) with an actual interactive confirmation prompt. The `--yes` flag should remain for non-interactive/scripted use.

## Context

Currently the archive command requires `--id <task-id>` for single-task archiving and errors with "use --yes (-y) to confirm" instead of actually prompting the user. This is similar to what was done for the `set` command in task 154.

## Tasks

- [x] Change `Args` from `cobra.NoArgs` to `cobra.MaximumNArgs(1)`
- [x] Update `Use` field to `archive [task-id]` and add positional arg to help/examples
- [x] In `runArchive`, resolve task ID from positional arg and append to `archiveIDs`
- [x] Replace the error-based confirmation (`"use --yes (-y) to confirm"`) with an interactive `y/N` prompt using `fmt.Scanln` or `bufio.Scanner`
  - Show task list preview, then ask "Proceed? [y/N]:"
  - Only skip prompt when `--yes`/`-y` flag is set
  - For `--delete`, keep the extra confirmation (or use `--force`)
- [x] Update tests to cover positional arg and interactive confirmation flow
- [x] Update command examples in help text

## Acceptance Criteria

- `taskmd archive 042` shows the task and prompts "Archive these tasks? [y/N]:" interactively
- `taskmd archive 042 -y` archives without prompting
- `taskmd archive --id 042` continues to work (backward compatible)
- `taskmd archive --all-completed` prompts interactively before archiving
- Interactive prompt reads from stdin and accepts `y`/`Y` to confirm
- All existing archive tests continue to pass
- New tests cover positional arg resolution and prompt bypass with `--yes`
