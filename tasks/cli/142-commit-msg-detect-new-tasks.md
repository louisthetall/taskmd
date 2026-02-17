---
id: "142"
title: "Detect new pending task files in commit-msg and adjust message"
status: pending
priority: medium
effort: small
tags: [cli, git, dx]
created: 2026-02-17
---

# Detect new pending task files in commit-msg and adjust message

## Objective

Update the `commit-msg` command so that when auto-inferring from staged changes (no `--task-id`), it detects when **only** new task files with `status: pending` were added and generates an appropriate commit message like `chore: added task 142` (or `chore: added tasks 142, 143` for multiple).

Currently the command only looks for tasks changing to `completed`. When a user stages only new pending task files, it errors with "no completed tasks found in staged changes." Instead, it should recognize this common workflow and produce a sensible commit message automatically.

## Tasks

- [ ] Add a `parseNewPendingFilesFromDiff` function (or extend existing diff parsing) to detect newly added files where the diff contains `+status: pending`
- [ ] Extract task IDs from the matched file paths or parsed frontmatter
- [ ] When no completed tasks are found but new pending tasks are detected, generate a message like `chore: added task <ID>` (single) or `chore: added tasks <ID1>, <ID2>` (multiple)
- [ ] Respect existing `--type` flag override for the commit prefix
- [ ] Ensure completed tasks still take priority (if both completed and new pending are staged, use the existing completed-task message)
- [ ] Add tests for the new pending-task detection path
- [ ] Add tests for mixed scenarios (completed + new pending in same commit)

## Acceptance Criteria

- `taskmd commit-msg` with only new pending task files staged outputs `chore: added task <ID>`
- Multiple new pending task files produce `chore: added tasks <ID1>, <ID2>`
- `--type` flag overrides the prefix (e.g., `--type feat` produces `feat: added task <ID>`)
- When both completed and new pending tasks are staged, the completed-task message is used (existing behavior preserved)
- No regression in existing commit-msg functionality
- All new code paths have test coverage
