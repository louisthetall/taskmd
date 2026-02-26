---
title: "Clickable status, priority, and effort on Stats page to filter List view"
id: "223"
status: completed
priority: medium
type: feature
tags: ["ui", "navigation"]
created: "2026-02-25"
---

# Clickable status, priority, and effort on Stats page to filter List view

## Objective

Make status, priority, and effort values on the Stats page clickable so that clicking one navigates the user to the Tasks (List) page with that value already filtered. Currently the List page filters status/priority/type via local state only — URL param support needs to be added (similar to how `?tag=` already works) so that the Stats page can link directly to a filtered view.

## Tasks

- [x] Add URL param support for `status`, `priority`, and `effort` filters in `TasksPage.tsx` and `TaskTable.tsx` (mirroring the existing `?tag=` pattern)
- [x] Sync `selectedStatuses`, `selectedPriorities`, and `selectedTypes` state to/from URL params
- [x] In `StatsView.tsx`, make each value in the `BreakdownCard` components ("By Status", "By Priority", "By Effort") clickable
- [x] On click, navigate to `/tasks?status=<value>`, `/tasks?priority=<value>`, or `/tasks?effort=<value>`
- [x] Style the breakdown values to indicate they are clickable (cursor pointer, hover effect)

## Acceptance Criteria

- Clicking a status value on the Stats page navigates to the Tasks page showing only tasks with that status
- Clicking a priority value on the Stats page navigates to the Tasks page showing only tasks with that priority
- Clicking an effort value on the Stats page navigates to the Tasks page showing only tasks with that effort level
- The URL updates to reflect the active filters so the filtered view is bookmarkable
- Breakdown values have a visible hover state indicating they are interactive
