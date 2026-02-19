---
id: "161"
title: "Next command: priority-aware downstream scoring"
status: pending
priority: medium
effort: small
tags: [cli, next, scoring]
created: 2026-02-19
---

# Next Command: Priority-Aware Downstream Scoring

## Objective

Adjust the `next` command's scoring algorithm so that downstream/critical-path bonuses are weighted by the priority of the downstream chain. Currently, a low-priority task that unblocks other low-priority tasks can outscore an unblocked medium-priority task because the downstream bonus (+15) and critical-path bonus (+15) are applied regardless of the priority of the tasks being unblocked.

For example, a low-priority task (base 10) with 5 low-priority downstream tasks scores 10 + 15 (downstream) + 15 (critical path) = 40, beating a standalone medium-priority task at 20. If the entire chain is low priority, the medium task should rank higher.

## Tasks

- [ ] In `ScoreTask()` (`apps/cli/internal/next/next.go`), scale the downstream bonus by the max priority in the downstream chain (e.g., full bonus if downstream contains high/critical tasks, reduced bonus if all downstream are low)
- [ ] Similarly scale the critical-path bonus — if the critical path consists entirely of low-priority tasks, reduce its weight relative to higher-priority standalone tasks
- [ ] Update or add tests in `apps/cli/internal/next/next_test.go` covering:
  - Low chain (all low-priority deps) does not outscore unblocked medium task
  - Mixed chain (low task unblocking high task) still gets full downstream bonus
  - Existing scoring behavior for high/critical priority chains is preserved

## Acceptance Criteria

- An unblocked medium-priority task ranks higher than a low-priority task whose entire downstream chain is also low priority
- A low-priority task that unblocks a high or critical priority task still gets a meaningful downstream bonus
- Existing tests continue to pass
- `make lint` passes
