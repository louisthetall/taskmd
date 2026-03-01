---
title: "Show blocked/unblocked state in status command for tasks with dependencies"
id: "01kjmggj3"
status: in-progress
priority: medium
type: feature
tags: []
created: "2026-03-01"
---

# Show blocked/unblocked state in status command for tasks with dependencies

## Objective

Update the `status` command so that when a task has dependencies, the output clearly indicates whether the task is **blocked** (has unmet dependencies) or **unblocked** (all dependencies are completed). This gives users immediate visibility into whether a task is ready to be worked on without needing to manually check each dependency.

## Tasks

- [ ] Inspect the current `status` command implementation to understand its output format
- [ ] Add dependency resolution logic: for each dependency, check if it is completed
- [ ] Add a "Blocked" / "Unblocked" indicator to the status output when a task has dependencies
- [ ] List the specific blocking dependency IDs (those not yet completed) when blocked
- [ ] Add tests for blocked state (task with incomplete dependencies)
- [ ] Add tests for unblocked state (task with all dependencies completed)
- [ ] Add tests for tasks with no dependencies (should not show blocked/unblocked indicator)

## Acceptance Criteria

- When running `taskmd status <id>` on a task with dependencies where some are not completed, the output includes a clear "Blocked" indicator along with the IDs of the blocking tasks
- When running `taskmd status <id>` on a task with dependencies where all are completed, the output includes an "Unblocked" indicator
- When running `taskmd status <id>` on a task with no dependencies, no blocked/unblocked indicator is shown
- All new behavior is covered by unit tests
