---
title: "Add rm command to delete tasks by ID"
id: "207"
status: pending
priority: medium
type: feature
tags: []
created: "2026-02-26"
---

# Add rm command to delete tasks by ID

## Objective

Add a `taskmd rm <id>` command that safely deletes a task file by ID. Before deleting, the command validates that removing the task would not break any references (e.g., dependencies from other tasks). If validation fails, the deletion is blocked unless `--force` is used.

## Tasks

- [ ] Create `internal/cli/rm.go` with the `rm` command registered under `rootCmd`
- [ ] Resolve task ID to file path using the scanner
- [ ] Run validation to check the task is not referenced by other tasks (dependencies, parent, etc.)
- [ ] If validation passes, delete the task file (and its worklog if present)
- [ ] If validation fails, print the referencing tasks and abort unless `--force` is provided
- [ ] With `--force`, skip reference checks and delete the file regardless
- [ ] Add comprehensive tests in `internal/cli/rm_test.go`
- [ ] Add e2e tests for the rm command

## Acceptance Criteria

- `taskmd rm <id>` deletes the task file when no other tasks reference it
- `taskmd rm <id>` prints an error listing referencing tasks and exits non-zero when references exist
- `taskmd rm <id> --force` deletes the task even when references exist
- Associated worklog file is also deleted if present
- `taskmd validate` passes after a successful (non-forced) deletion
- Command prints a confirmation message with the deleted file path
