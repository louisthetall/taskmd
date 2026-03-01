---
title: "Implement minimal feed command"
id: "01kjmg6sc"
status: in-progress
priority: medium
type: feature
tags: ["cli", "git"]
created: "2026-03-01"
---

# Implement minimal feed command

## Objective

Add a `taskmd feed` command that shows a chronological activity feed of recent changes to task files. The minimal version uses `git log` on the tasks directory to detect task creation, status changes, and modifications, presenting them as a time-ordered feed.

This gives users a quick way to answer "what happened recently?" — complementing `list` (spatial) with a temporal view.

## Tasks

- [ ] Create `internal/cli/feed.go` with the `feed` cobra command
- [ ] Implement git log parsing scoped to the tasks directory (`git log --name-status --diff-filter=ACMR -- 'tasks/**/*.md'`)
- [ ] Detect event types from git diffs: created (A), modified (M), renamed (R)
- [ ] Parse commit messages and timestamps into structured feed entries
- [ ] Support `--since` flag for time-based filtering (e.g. `--since 2d`, `--since 2026-02-28`)
- [ ] Support `--limit` flag to cap the number of entries (default: 20)
- [ ] Support `--scope` flag to filter to a specific tasks subdirectory
- [ ] Implement plain text output format (default) with one-line-per-event style
- [ ] Support `--format json` for machine-readable output
- [ ] Add comprehensive tests in `internal/cli/feed_test.go`
- [ ] Register the command in the CLI

## Acceptance Criteria

- `taskmd feed` outputs a chronological list of recent task file changes with timestamps, event type, and task title
- `taskmd feed --since 7d` filters to the last 7 days
- `taskmd feed --limit 10` caps output at 10 entries
- `taskmd feed --scope cli` filters to the `tasks/cli/` subdirectory
- `taskmd feed --format json` outputs structured JSON
- Works correctly in repos with no task changes (empty feed, no errors)
- Gracefully handles non-git directories with a clear error message
- Tests cover happy path, flags, empty results, and error cases
