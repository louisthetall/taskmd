---
id: "01kjmx8ym"
title: "Show informational message when status has no in-progress tasks"
status: completed
priority: low
dependencies: []
tags: ["cli", "ux"]
created: 2026-03-01
---

# Show informational message when status has no in-progress tasks

## Objective

When `taskmd status` is run with no arguments and no tasks are in progress, print a minimal informational message instead of producing no output. Structured formats (JSON/YAML) should return an empty array for machine consumption.

## Tasks

- [x] Print "No tasks currently in progress." to stderr in text mode
- [x] Return empty array `[]` for --format json and --format yaml
- [x] Keep --statusline silent (no change) to avoid breaking shell integrations
- [x] Add tests for all three empty-result paths (text, json, yaml)

## Acceptance Criteria

- `taskmd status` with no in-progress tasks prints "No tasks currently in progress." to stderr
- `taskmd status --format json` with no in-progress tasks outputs `[]`
- `taskmd status --format yaml` with no in-progress tasks outputs `[]`
- `taskmd status --statusline` with no in-progress tasks remains silent
- All existing status tests continue to pass
