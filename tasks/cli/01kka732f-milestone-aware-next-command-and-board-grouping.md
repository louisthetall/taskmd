---
title: "Milestone-aware next command and board grouping"
id: "01kka732f"
status: pending
priority: medium
type: feature
dependencies: ["01kka72zy"]
tags: ["milestone", "cli"]
touches: ["cli/next", "cli/commands"]
created: "2026-03-09"
---

# Milestone-aware next command and board grouping

## Objective

Make `taskmd next` prefer tasks from the earliest open milestone (based on ordering in `.taskmd.yaml`). Add `--group-by milestone` support to `taskmd board` and `taskmd stats`.

## Tasks

- [ ] Update `taskmd next` ranking to factor in milestone ordering (earliest milestone first)
- [ ] Add `--milestone` filter to `taskmd next`
- [ ] Add `milestone` as a `--group-by` option in `taskmd board`
- [ ] Add `milestone` as a `--group-by` option in `taskmd stats`
- [ ] Add tests for next command milestone preference
- [ ] Add tests for board/stats grouping by milestone

## Acceptance Criteria

- `taskmd next` prefers tasks from earlier milestones when milestones are configured in `.taskmd.yaml`
- `taskmd next --milestone v0.2` restricts suggestions to that milestone
- `taskmd board --group-by milestone` groups columns by milestone
- `taskmd stats --group-by milestone` shows per-milestone statistics
- Tasks with no milestone are grouped under a "(no milestone)" bucket
- Tests cover ranking, filtering, and grouping scenarios
