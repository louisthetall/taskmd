---
id: "179"
title: "Add worklog editing in web UI"
status: pending
priority: low
effort: medium
type: feature
tags:
  - web
  - worklogs
created: 2026-02-20
---

# Add Worklog Editing in Web UI

## Objective

Allow users to add worklog entries from the web interface. Currently the web UI can display worklogs but adding entries requires the CLI. This creates unnecessary context-switching for users who prefer the web dashboard.

## Tasks

- [ ] Add an "Add Entry" button/form to the task detail worklog section
- [ ] Implement a `POST /api/tasks/{id}/worklog` endpoint that appends to the worklog file
- [ ] Auto-generate the timestamp in the correct format
- [ ] Support the standard worklog entry format (timestamped markdown)
- [ ] Create the worklog file if it doesn't exist yet
- [ ] Validate that worklogs are enabled in the project config before allowing edits
- [ ] Update the worklog display in real-time after adding an entry
- [ ] Add tests for the new API endpoint

## Acceptance Criteria

- Users can add timestamped worklog entries from the task detail page
- Entries are appended to the correct worklog file on disk
- The worklog display updates immediately after submission
- The feature is hidden when worklogs are disabled in `.taskmd.yaml`
