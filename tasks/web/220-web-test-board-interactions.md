---
id: "220"
title: "Test board drag-and-drop and page containers"
status: completed
priority: low
type: chore
effort: medium
tags: ["testing", "quality"]
dependencies: ["214"]
created: "2026-02-26"
phase: web-ui
---

# Test board drag-and-drop and page containers

## Objective

Add tests for the board view's drag-and-drop interactions and the remaining page containers that have non-trivial logic (BoardPage filtering, TaskDetailPage edit flow).

## Tasks

- [x] Test `BoardColumn.tsx` drag handlers (dragOver, dragLeave, drop with source !== target validation)
- [x] Test `TaskCard.tsx` drag guard system (dragStart data setup, dragEnd cleanup)
- [x] Test `BoardPage.tsx` filter logic (`availableTags` extraction, `filteredGroups` memoization)
- [x] Refactor: extract board drag handler logic into a utility if the handlers are complex enough to warrant it
- [x] Test `TaskDetailPage.tsx` edit flow (save handler, error display with `ApiRequestError`)

## Acceptance Criteria

- Board drag-and-drop has tests for valid and invalid drop targets
- BoardPage filtering is tested for at least 2 groupBy modes
- TaskDetailPage edit and error handling are tested
- All new tests pass via `pnpm test`
