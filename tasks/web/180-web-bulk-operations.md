---
id: "180"
title: "Add bulk operations to web UI"
status: pending
priority: low
effort: medium
type: feature
tags:
  - web
  - ux
created: 2026-02-20
---

# Add Bulk Operations to Web UI

## Objective

Allow users to select multiple tasks and perform batch actions from the web interface. Managing tasks one-at-a-time is tedious when archiving completed work or updating priorities across many tasks.

## Tasks

- [ ] Add row-level checkboxes to the tasks table view
- [ ] Add a "select all" checkbox in the table header
- [ ] Show a floating action bar when tasks are selected (with count)
- [ ] Implement bulk status change (e.g., mark all selected as completed)
- [ ] Implement bulk archive (move selected tasks to archive)
- [ ] Implement bulk priority/effort update
- [ ] Add a `PATCH /api/tasks/bulk` endpoint for batch updates
- [ ] Add a `POST /api/tasks/bulk/archive` endpoint for batch archiving
- [ ] Add confirmation dialog before destructive bulk actions
- [ ] Add keyboard shortcut for select all / deselect all
- [ ] Add tests for bulk API endpoints

## Acceptance Criteria

- Users can select multiple tasks via checkboxes in the table view
- Bulk status change updates all selected tasks in a single action
- Bulk archive moves all selected tasks to the archive directory
- A confirmation dialog appears before destructive operations
- The UI updates all affected rows after a bulk operation completes
