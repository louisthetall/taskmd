---
title: "Milestone support in web dashboard"
id: "01kka732z"
status: completed
priority: medium
type: feature
dependencies: ["01kka72zy"]
tags: ["milestone", "web"]
touches: ["web", "web/tasks", "web/board", "web/stats"]
created: "2026-03-09"
---

# Milestone support in web dashboard

## Objective

Display and filter by milestone in the web dashboard. Add milestone to the task list, board view, and stats view.

## Tasks

- [x] Display milestone badge/chip on task cards in list view
- [x] Add milestone filter dropdown to the task list sidebar
- [x] Add milestone column/grouping option to the board view
- [x] Show per-milestone progress in the stats view
- [x] Include milestone in task detail panel
- [x] Handle tasks with no milestone gracefully (show as ungrouped)

## Acceptance Criteria

- Task list shows milestone on each task when present
- Users can filter the task list by milestone
- Board view supports grouping by milestone
- Stats view shows a milestone-based progress breakdown
- Tasks without a milestone display without errors
