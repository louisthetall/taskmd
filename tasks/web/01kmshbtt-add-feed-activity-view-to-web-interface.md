---
title: "Add feed/activity view to web interface"
id: "01kmshbtt"
status: completed
priority: medium
type: feature
tags: ["feed", "activity"]
created: "2026-03-28"
---

# Add feed/activity view to web interface

## Objective

Add a Feed/Activity page to the web interface that mirrors the CLI `taskmd feed` command. The feed should show a chronological activity stream of task changes (from git history) and worklog entries, presented in a polished web UX.

The CLI feed shows:
- **Git-sourced entries**: commit hash, author, timestamp, message, changed task files with field changes (status transitions, priority changes) and subtask completions
- **Worklog-sourced entries**: timestamped worklog notes per task
- Filtering by time range (`--since`), scope (task subdirectory), source (git/worklog/all), and entry limit

## Tasks

- [x] Add `GET /api/feed` endpoint to the Go web server (`apps/cli/internal/web/server.go`) that reuses the existing feed logic from `feed.go` (git log parsing, worklog scanning, diff analysis) and returns JSON
- [x] Support query params: `limit`, `since`, `scope`, `source` matching the CLI flags
- [x] Create `use-feed.ts` hook in `apps/web/src/hooks/` to fetch from the feed API
- [x] Create `FeedView.tsx` component showing a timeline/activity stream with entries grouped or ordered chronologically
- [x] Render git entries with: timestamp, author, commit message, list of changed files with field change badges (e.g. "status: pending -> in-progress") and subtask change indicators
- [x] Render worklog entries with: timestamp, task ID link, worklog message content
- [x] Add filter controls for source (git/worklog/all), time range, and scope
- [x] Add `FeedPage.tsx` page and wire up routing + nav tab
- [x] Add tests for the API endpoint, hook, and component

## Acceptance Criteria

- Feed page is accessible from the web nav and shows the same data as `taskmd feed --format json`
- Entries display timestamps, sources, and relevant details (author, commit message, file changes, field diffs, worklog content)
- Users can filter by source, time range, and scope
- Empty state is handled gracefully ("No recent activity")
- Feed loads performantly for typical repositories
