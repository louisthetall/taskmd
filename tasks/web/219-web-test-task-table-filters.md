---
id: "219"
title: "Test TaskTable filtering and URL sync"
status: completed
priority: medium
type: chore
effort: medium
tags: ["testing", "quality"]
dependencies: ["214"]
created: "2026-02-26"
---

# Test TaskTable filtering and URL sync

## Objective

Add tests for the task list's filter logic, which is the most-used interactive feature in the web app. The filter state, URL synchronization, and multi-criteria intersection logic should all be covered.

## Tasks

- [x] Refactor: extract filter application logic from `TaskTable.tsx` into a pure function (e.g. `applyFilters(tasks, filters)`) that can be tested without rendering
- [x] Test filter intersection (status AND priority AND type AND effort AND tags)
- [x] Test `hasActiveFilters` calculation
- [x] Test URL synchronization (filters read from and written to search params)
- [x] Test `clearFilters()` resets all filter state
- [x] Test `SearchDialog`'s `Highlight` subcomponent (match found, no match, case insensitivity)
- [x] Refactor: extract `Highlight` from `SearchDialog.tsx` into its own file if it's useful elsewhere

## Acceptance Criteria

- Filter logic has tests for single and combined filter criteria
- URL sync round-trips correctly (set filters → URL updates → page reload → filters restored)
- Highlight component is tested independently
- All new tests pass via `pnpm test`
