---
title: "Update status command to include children task IDs and their statuses"
id: "01kjm6k2f"
status: completed
priority: medium
type: feature
tags: []
created: "2026-03-01"
---

# Update status command to include children task IDs and their statuses

## Objective

Enhance the `taskmd status` command to recursively display the full children tree (tasks whose `parent` field matches the current task, and their children, etc.) along with each child's current status. This gives users a quick overview of a parent task's full subtask hierarchy without needing to run separate queries.

Add a `--minimal` flag that outputs only the task's own metadata with no children or extra info.

## Tasks

- [x] Add a `--minimal` flag to the status command that skips children and extra info
- [x] Add a `Children` field to `statusOutput` struct (recursive: each child has its own children)
- [x] Build a parent→children index from all scanned tasks
- [x] Recursively collect children tree, tracking visited IDs to detect and break circular loops
- [x] Include children tree in text output with indentation showing depth
- [x] Include children tree in JSON/YAML output as nested structured arrays
- [x] Add tests for recursive children display (parent → child → grandchild)
- [x] Add tests for circular parent loop detection (exits gracefully)
- [x] Add tests for `--minimal` flag (no children in output)
- [x] Add tests for tasks with no children (no regression)

## Acceptance Criteria

- Running `taskmd status <id>` on a parent task shows a recursive "Children" tree listing all descendant task IDs and their statuses
- Children are displayed with indentation to reflect nesting depth in text output
- JSON/YAML output includes a nested `children` array where each child can have its own `children`
- Circular parent references are detected and the loop is broken gracefully (no infinite recursion)
- `--minimal` flag outputs only the task's own metadata with no children or extra info
- When a task has no children, the children field is omitted (consistent with other optional fields)
- Existing status command behavior is unchanged for non-parent tasks when not using `--minimal`
