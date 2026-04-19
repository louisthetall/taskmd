---
name: list-tasks
description: List tasks with optional filters. Use when the user wants to see their tasks.
allowed-tools: Bash
---

# List Tasks

List tasks using the `taskmd` CLI.

## Instructions

The user's arguments are in `$ARGUMENTS` (e.g. `--status pending`, `--format json`, a directory path).

1. Run `taskmd list $ARGUMENTS`
   - If `$ARGUMENTS` is empty, run: `taskmd list`
   - Common flags: `--status`, `--priority`, `--format`, `--filter`
   - Phase filtering: `taskmd list --filter phase=core-cli` to list tasks in a specific phase
   - Use `taskmd phases` to see all configured phases with progress stats
2. Present the output to the user
