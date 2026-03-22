---
title: "Add --all-projects flag for aggregate task views"
id: "01kma4701"
status: pending
priority: medium
type: feature
tags: ["global-registry", "cli-command"]
dependencies: ["01kma46wt"]
created: "2026-03-22"
---

# Add --all-projects flag for aggregate task views

## Objective

Add an `--all-projects` flag to `list` and `next` that scans all registered projects and presents aggregated results. Task IDs are qualified as `<project-id>:<task-id>` to avoid ambiguity.

## Tasks

- [ ] Add `--all-projects` boolean flag to `list` and `next` commands
- [ ] When set, iterate over all entries from `LoadGlobalRegistry()`, scan each project independently
- [ ] Add a `PROJECT` column to table output showing which project each task belongs to
- [ ] Qualify task IDs as `<project-id>:<task-id>` in display output
- [ ] In JSON/YAML output, include a `project` field on each task object
- [ ] For `next`, rank across all projects using the same ranking algorithm, include project context in output
- [ ] Handle unreachable projects gracefully (warn and skip)
- [ ] Ensure `--all-projects` and `--project` are mutually exclusive (error if both set)
- [ ] Add tests with multiple temp project directories

## Acceptance Criteria

- `taskmd list --all-projects` shows tasks from all registered projects with a project column
- Task IDs are displayed as `projectid:taskid` in aggregate mode
- `taskmd next --all-projects` returns the highest-priority task across all projects
- Unreachable projects are skipped with a warning
- `--all-projects` and `--project` cannot be used together
