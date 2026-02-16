---
id: "121"
title: "Fix archived dependencies treated as unmet"
status: pending
priority: medium
effort: medium
tags:
  - cli
  - go
  - bug
touches:
  - cli/next
  - cli/scanner
created: 2026-02-16
---

# Fix archived dependencies treated as unmet

## Objective

When a task depends on a completed task that has been archived (moved to `tasks/archive/`), the dependency is treated as unmet because the scanner doesn't include archived tasks in the task map. This causes dependent tasks to be excluded from `next` and `tracks` even though their dependencies are satisfied.

Fix by including archived tasks in the task map for dependency resolution.

## Tasks

- [ ] Scan `tasks/archive/` directory and include archived tasks in the task map used by `HasUnmetDependencies`
- [ ] Ensure archived tasks are only used for dependency resolution, not surfaced in `next`/`tracks`/`list` output
- [ ] Add tests: task with archived completed dependency should be actionable
- [ ] Add tests: task with archived non-completed dependency should remain non-actionable

## Acceptance Criteria

- A task depending on a completed+archived task is correctly identified as actionable
- `taskmd next` and `taskmd tracks` include tasks whose dependencies are all completed (even if archived)
- Archived tasks themselves do not appear in `next`, `tracks`, or `list` output
- Existing tests continue to pass
