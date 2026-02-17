---
id: "140"
title: "Fix critical path calculation counting completed dependencies"
status: completed
priority: medium
effort: small
tags:
  - cli
  - next
  - bugfix
touches:
  - cli/next
created: 2026-02-17
---

# Fix Critical Path Calculation Counting Completed Dependencies

## Objective

Fix a bug in `calculateDepthMap` where completed dependencies were counted toward chain depth, causing tasks with resolved dependencies to be falsely marked as "critical path." The critical path should only reflect remaining work.

## Tasks

- [x] Identify the bug: `calculateDepthMap` includes completed tasks in depth calculation
- [x] Write failing unit tests in `internal/next/next_test.go`:
  - [x] `TestCalculateCriticalPathTasks_IgnoresCompletedDependencies`
  - [x] `TestCalculateCriticalPathTasks_PendingChainIsCritical`
  - [x] `TestCalculateCriticalPathTasks_MixedCompletedPendingChain`
- [x] Fix `calculateDepthMap` to return depth 0 for resolved (completed/cancelled) tasks
- [x] Fix `markCriticalPathDependencies` to use existence check on depthMap to avoid marking resolved tasks
- [x] Add `Status.IsResolved()` helper method to model package
- [x] Add unit tests for `IsResolved()` in `internal/model/task_test.go`
- [x] Update existing test `TestNext_Critical_NoCriticalTasksAvailable` that relied on buggy behavior
- [x] Add integration test `TestNext_Critical_CompletedDepsIgnored`
- [x] Verify all tests pass

## Acceptance Criteria

- Tasks whose only dependencies are completed are NOT marked as critical path (unless they form the longest remaining chain)
- Completed and cancelled tasks return depth 0 in the depth map
- The critical path reflects only remaining (pending/in-progress) work
- All existing tests continue to pass
- `taskmd next` no longer shows tasks 128/129 as "on critical path"
