---
id: "192"
title: "Improve test coverage for Tier 1 (core) packages"
status: completed
priority: high
effort: large
type: chore
tags: [testing, coverage]
created: 2026-02-21
---

# Improve test coverage for Tier 1 (core) packages

## Objective

Add tests for the highest-importance uncovered areas: the core CLI commands and packages that guard data integrity, display critical project metrics, or underpin the data model. These are the areas where bugs have the highest blast radius.

## Background

Current coverage analysis (2026-02-21) identified these Tier 1 gaps:

| Area | Uncovered stmts | Total stmts | Coverage | Risk |
|------|-----------------|-------------|----------|------|
| `cli/validate.go` | 104 | 125 | 16.8% | Pre-commit hook gate; bugs block or silently pass bad files |
| `cli/stats.go` | 73 | 75 | 2.7% | Project metrics command; data aggregation entirely untested |
| `cli/snapshot*.go` | 137 | 184 | 25.5% | Snapshot command with topological sort; all outputs untested |
| `model/task.go` | 13 | 14 | 7.1% | Core data model used everywhere; `ValidateVerifySteps`, `IsValid`, `GetGroup` at 0% |
| `board` package | 45 | 75 | 40.0% | Sorting functions (`statusOrder`, `priorityOrder`, `effortOrder`) all at 0% |

## Tasks

- [x] Add tests for `validate.go`: `runValidate`, `outputValidationText`, `outputValidationJSON`, `printIssue`, `validateConfig`, `loadConfigForValidation`, `collectArchivedIDs`, `mergeValidationResults`
- [x] Add tests for `stats.go`: `runStats`, `outputStatsJSON`, `outputStatsTable`, `printStatsBreakdownByStatus`, `printStatsBreakdownByPriority`, `printStatsBreakdownByEffort`
- [x] Add tests for `snapshot.go`, `snapshot_analysis.go`, `snapshot_output.go`: `runSnapshot`, `taskToSnapshot`, `groupSnapshots`, `calculateTopologicalOrder`, all output formatters
- [x] Add tests for `model/task.go`: `ValidateVerifySteps`, `IsValid`, `GetGroup` (already at 100%)
- [x] Add tests for `board` package: `statusOrder`, `priorityOrder`, `effortOrder`, `sortedKeys`, additional `GroupTasks` scenarios (already at 96%)

## Acceptance Criteria

- All Tier 1 functions listed above have at least one test covering the happy path
- `validate.go` coverage reaches 80%+
- `stats.go` coverage reaches 80%+
- `snapshot*.go` combined coverage reaches 80%+
- `model/task.go` coverage reaches 90%+
- `board` package coverage reaches 80%+
- All tests pass: `cd apps/cli && make test`
- Linter passes: `cd apps/cli && make lint`
