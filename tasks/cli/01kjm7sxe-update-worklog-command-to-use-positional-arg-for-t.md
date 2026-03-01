---
id: "01kjm7sxe"
title: "Update worklog command to use positional arg for task ID"
status: completed
priority: medium
dependencies: []
tags: []
created: 2026-03-01
---

# Update worklog command to use positional arg for task ID

## Objective

Change the `taskmd worklog` command to accept the task ID as a positional argument instead of the `--task-id` flag. This aligns the worklog command with the convention used by other commands like `taskmd set <id>`.

**Before:** `taskmd worklog --task-id 015 --add "message"`
**After:** `taskmd worklog 015 --add "message"`

## Tasks

- [ ] Update `worklog.go` to accept task ID as a positional arg (`Args: cobra.ExactArgs(1)`)
- [ ] Remove the `--task-id` flag
- [ ] Update help text and usage examples
- [ ] Update the CLAUDE.md template worklog CLI examples
- [ ] Update the taskmd specification if it references `--task-id`
- [ ] Add/update tests for the new argument style
- [ ] Run `make lint` and `make test`

## Acceptance Criteria

- `taskmd worklog 015` displays the worklog for task 015
- `taskmd worklog 015 --add "message"` appends an entry
- `taskmd worklog` with no argument prints a usage error
- CLAUDE.md template and spec docs reflect the new syntax
- All existing tests pass, new argument parsing is covered
