---
title: "Exclude cancelled tasks from phase completion percentage"
id: "01kkk461x"
status: completed
priority: high
type: bug
tags: ["phases", "metrics"]
created: "2026-03-13"
---

# Exclude cancelled tasks from phase completion percentage

## Steps to Reproduce

1. Create tasks assigned to a phase (e.g., `phase: mvp`)
2. Set some tasks to `status: cancelled`
3. Run `taskmd phases`
4. Observe the completion percentage includes cancelled tasks in the total count

## Expected Behavior

Cancelled tasks should be excluded from the phase completion percentage calculation. For example, if a phase has 10 tasks, 3 completed, and 2 cancelled, the progress should be 3/8 = 37% (not 3/10 = 30%).

## Actual Behavior

Cancelled tasks are counted in `summary.Tasks` (the denominator) but not in `summary.Done`, which deflates the completion percentage.

## Tasks

- [x] Update `computePhaseSummaries` in `apps/cli/internal/cli/phases.go` to skip cancelled tasks (do not count them in `summary.Tasks` or `summary.ByStatus`)
- [x] Decide whether cancelled tasks should still appear in `ByStatus` breakdown (likely yes for visibility, but excluded from Tasks/Done/Progress)
- [x] Add tests for phase completion with cancelled tasks

## Acceptance Criteria

- Cancelled tasks are not counted in the `Tasks` total or `Progress` percentage
- The `ByStatus` map still includes cancelled task counts for visibility
- `taskmd phases` output reflects accurate completion percentages when cancelled tasks exist
- Unit tests cover the cancelled-task exclusion scenario
