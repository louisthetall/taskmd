---
title: "Support milestone filtering in list command"
id: "01kka730t"
status: pending
priority: high
type: feature
dependencies: ["01kka72zy"]
tags: ["milestone", "cli"]
touches: ["cli/commands"]
created: "2026-03-09"
---

# Support milestone filtering in list command

## Objective

Add a `--milestone` flag to `taskmd list` so users can filter tasks by milestone. Also include the milestone column in table/JSON/YAML output.

## Tasks

- [ ] Add `--milestone` flag to the list command
- [ ] Filter tasks by milestone value (exact match)
- [ ] Include `milestone` in table output as a column
- [ ] Include `milestone` in JSON and YAML output
- [ ] Add tests for milestone filtering (matching, non-matching, empty)
- [ ] Add tests for milestone in output formats

## Acceptance Criteria

- `taskmd list --milestone v0.2` shows only tasks with `milestone: v0.2`
- `taskmd list` (without flag) shows all tasks as before, with milestone column visible when any task has a milestone
- JSON/YAML output includes the `milestone` field
- Tests cover filtering and output formats
