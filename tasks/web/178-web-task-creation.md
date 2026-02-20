---
id: "178"
title: "Add task creation dialog to web UI"
status: pending
priority: low
effort: medium
type: feature
tags:
  - web
  - ux
created: 2026-02-20
---

# Add Task Creation Dialog to Web UI

## Objective

Allow users to create new tasks directly from the web interface without dropping to the CLI. This is the most significant workflow friction point — users must currently switch to a terminal to add tasks.

## Tasks

- [ ] Design a "New Task" dialog/modal with fields for title, status, priority, effort, type, tags, and body
- [ ] Add a "New Task" button to the tasks page header and board view
- [ ] Implement a `POST /api/tasks` endpoint that creates the task file on disk
- [ ] Auto-assign the next available ID (using next-id logic)
- [ ] Auto-generate the filename slug from the title
- [ ] Support setting the group/subdirectory for the new task
- [ ] Add client-side validation for required fields (id, title)
- [ ] Trigger SSE reload so other views update after creation
- [ ] Add tests for the new API endpoint

## Acceptance Criteria

- Users can create a new task entirely from the web UI
- The created task file is valid and passes `taskmd validate`
- The dialog pre-fills sensible defaults (status: pending, priority: medium)
- The new task appears in all views (table, board, graph) without a page refresh
