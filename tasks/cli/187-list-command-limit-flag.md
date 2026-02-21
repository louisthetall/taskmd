---
id: "187"
title: "Add --limit flag to list command"
status: completed
priority: medium
effort: small
tags: [cli, list]
created: 2026-02-21
---

# Add --limit flag to list command

## Objective

Add a `--limit` flag to the `list` CLI command, matching the behavior of the `next` command's existing `--limit` flag. The limit should be applied **after sorting**, so users get the top N results according to the current sort order.

## Tasks

- [x] Add `--limit` flag definition to the `list` command (integer, default 0 = unlimited)
- [x] Apply the limit after sorting the task list
- [x] Add unit tests for the `--limit` flag
- [x] Verify the flag works with all output formats (table, json, yaml)

## Acceptance Criteria

- `taskmd list --limit 5` returns at most 5 tasks
- The limit is applied after sorting (e.g., `--sort priority --limit 3` returns the top 3 by priority)
- `--limit 0` or omitting the flag returns all tasks (no limit)
- The flag works consistently across all output formats
- Tests cover happy path, edge cases (limit > total tasks, limit = 0), and interaction with sorting
