---
id: "157"
title: "Handle mixed task change types in commit-msg command"
status: completed
priority: medium
effort: medium
type: improvement
tags: [cli, git, dx]
created: 2026-02-19
---

# Handle mixed task change types in commit-msg command

## Objective

Update the `commit-msg` command to gracefully handle commits that contain multiple types of task changes simultaneously. Currently the command uses a strict priority hierarchy: it looks for completed tasks first, and only if none are found does it look for new pending tasks. If a commit contains both completed tasks and new pending tasks (or tasks marked in-progress), only the completed tasks are reflected in the generated message.

The command should detect all categories of task changes in the staged diff and produce a commit message that accurately describes the full set of changes, for example:

```
chore: complete tasks 042, 043; add tasks 044, 045; start task 046
```

## Tasks

- [x] Refactor diff parsing to detect all status change types in a single pass (completed, in-progress, pending/new, blocked, cancelled)
- [x] Create a unified struct to hold categorized task changes (e.g., `DiffResult` with fields for completed, added, started, etc.)
- [x] Update `runCommitMsg` to collect all change categories instead of short-circuiting on the first match
- [x] Build a combined subject line that lists each change type present (e.g., "complete tasks X; add tasks Y; start task Z")
- [x] Support `--body` flag with mixed changes: group subtask bullets by change category
- [x] Support `--short` flag with mixed changes
- [x] Handle edge case: task status changed to a non-completed status (e.g., `in-progress`, `blocked`)
- [x] Handle edge case: task status changed between non-new statuses (e.g., `pending` to `in-progress`)
- [x] Preserve existing behavior for single-type commits (no regression)
- [x] Add tests for mixed completed + new pending tasks in same commit
- [x] Add tests for mixed completed + in-progress tasks in same commit
- [x] Add tests for three or more change types in same commit
- [x] Add tests for `--task-id` flag (should remain unchanged)

## Acceptance Criteria

- When a commit contains both completed and newly added tasks, the message reflects both (e.g., `chore: complete tasks 042, 043; add tasks 044, 045`)
- When a commit contains tasks marked in-progress alongside other changes, those are included in the message
- Single-type commits produce the same output as before (no regression)
- `--task-id` flag behavior is unchanged (generates message for that specific task only)
- `--type` flag still overrides the commit prefix
- `--body` flag works correctly with mixed change types, grouping bullets by category
- `--short` flag produces a single subject line covering all change types
- All new code paths have test coverage
