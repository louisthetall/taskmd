---
title: "Filter next-up tasks by folder path in web UI"
id: "225"
status: completed
priority: medium
type: feature
tags: ["ui", "filtering"]
created: "2026-02-26"
---

# Filter next-up tasks by folder path in web UI

## Objective

Add a folder/group filter to the "Next Up" page in the web UI so users can scope task recommendations to a specific subdirectory (group) relative to the task directory. For example, filtering by `cli` would only show recommended tasks from `tasks/cli/`.

The backend already supports filtering via `?filter=group=<value>` on the `/api/next` endpoint. This task focuses on exposing that capability in the frontend UI.

## Tasks

- [x] Add a folder/group filter control to the `NextPage` or `NextView` component (e.g. a dropdown or text input for the group path)
- [x] Pass the selected group as a `filter` query parameter to the `useNext` hook's API call (`/api/next?filter=group=<value>`)
- [x] Update the `useNext` hook to accept an optional group/filter parameter
- [x] Persist the selected filter in the URL query string so it survives page reloads
- [x] Add tests for the filtering behavior

## Acceptance Criteria

- The Next Up page displays a filter control that lets users select or type a folder path (group) relative to the task directory
- When a group filter is applied, only tasks from that group appear in the recommendations
- The filter value is reflected in the URL query string and persists across page reloads
- Clearing the filter shows all recommended tasks (default behavior)
- The filter integrates with the existing limit control without conflicts
