---
title: "Add --strict-phases flag to next command for strict phase ordering"
id: "01kkkj7jh"
status: completed
priority: medium
type: feature
tags: ["cli", "phases"]
created: "2026-03-13"
---

# Add --strict-phases flag to next command for strict phase ordering

## Objective

Add a `--strict-phases` flag to the `next` command that enforces strict phase sequentiality when ranking tasks. Currently, phase ordering is one scoring factor among many (priority, critical path, effort, etc.), so a high-priority task in a later phase can outrank a lower-priority task in an earlier phase. When `--strict-phases` is active, tasks in earlier phases (per `.taskmd.yaml` phase order) must always appear before tasks in later phases, regardless of other scoring factors. Within the same phase, the existing scoring logic applies as usual.

## Tasks

- [x] Add `--strict-phases` boolean flag to the `next` command in `apps/cli/internal/cli/next.go`
- [x] Pass the flag value through to the recommendation engine in `sdk/go/next/next.go`
- [x] Implement strict phase sorting: when enabled, group actionable tasks by phase index, then sort within each group by existing score
- [x] Handle edge cases: tasks with no phase assigned (sort after all phased tasks, or treat as last phase)
- [x] Add unit tests in `apps/cli/internal/cli/next_test.go` covering:
  - [x] Flag is off by default (existing behavior unchanged)
  - [x] With flag on, earlier-phase tasks always rank above later-phase tasks
  - [x] Within same phase, normal scoring applies
  - [x] Tasks with no phase are sorted after phased tasks
  - [x] Interaction with `--phase` filter flag
- [x] Update command help text to document the new flag

## Acceptance Criteria

- `taskmd next` without `--strict-phases` behaves exactly as before (no regression)
- `taskmd next --strict-phases` returns tasks strictly ordered by phase: all actionable tasks from phase A appear before any task from phase B, where A precedes B in `.taskmd.yaml`
- Within the same phase group, tasks are ranked by the existing scoring algorithm
- Tasks with no phase assigned appear after all phase-grouped tasks
- All existing tests pass; new tests cover the flag behavior
- `taskmd next --help` documents the `--strict-phases` flag
