---
title: "Add --depends-on flag to set command"
id: "208"
status: completed
priority: medium
type: feature
tags: []
created: "2026-02-26"
---

# Add --depends-on flag to set command

## Objective

Add a `--depends-on` flag to the `taskmd set` command so users can add or update task dependencies from the CLI without manually editing frontmatter. For example: `taskmd set 042 --depends-on 010,015`.

## Tasks

- [x] Add `--depends-on` string flag to the `set` command in `internal/cli/set.go`
- [x] Parse comma-separated task IDs from the flag value
- [x] Update the task's `dependencies` frontmatter field (replace semantics)
- [x] Validate that referenced dependency IDs exist
- [x] Validate no circular dependencies are introduced
- [x] Add tests in `internal/cli/set_test.go` for the new flag
- [x] Add e2e tests covering `--depends-on` usage

## Acceptance Criteria

- `taskmd set <id> --depends-on 010,015` updates the task's dependencies to include tasks 010 and 015
- Invalid dependency IDs (non-existent tasks) produce a clear error
- Circular dependencies are detected and rejected
- `taskmd validate` passes after setting dependencies
- Flag works in combination with other `set` flags (e.g., `--status`, `--priority`)
