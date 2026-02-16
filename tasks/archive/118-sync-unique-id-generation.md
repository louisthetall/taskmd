---
id: "118"
title: "Generate unique IDs for synced tasks across all local tasks"
status: completed
priority: high
effort: small
dependencies:
  - "116"
tags:
  - cli
  - sync
  - bug
created: 2026-02-15
---

# Generate Unique IDs for Synced Tasks Across All Local Tasks

## Objective

Fix the sync engine so that newly synced tasks receive IDs that are unique across the entire project, not just the sync output directory.

## Problem

Currently `Engine.scanExistingIDs` only scans the `outputDir` (e.g., `./tasks/jira/`) to determine the next available ID. If the project already has tasks `001`-`115` in `./tasks/` and `./tasks/cli/`, synced tasks will start at `001` again, colliding with existing tasks.

The relevant code in `engine.go`:

```go
func (e *Engine) scanExistingIDs(outputDir string) ([]string, error) {
    // BUG: only scans outputDir, not the full project
    taskScanner := scanner.NewScanner(outputDir, e.Verbose, nil)
    ...
}
```

## Tasks

- [x] Change `scanExistingIDs` to scan the project root (`ConfigDir`) instead of just `outputDir`, so it discovers all task IDs across all groups/directories
- [x] Ensure the scanner result includes tasks from all subdirectories (the scanner already does recursive scanning)
- [x] Add a test that creates tasks in a separate directory (simulating existing project tasks) and verifies that synced tasks get non-colliding IDs
- [x] Add a test that syncs from two different sources into different output directories and verifies IDs don't collide between them

## Acceptance Criteria

- Synced tasks receive IDs that don't collide with any existing task in the project
- If the project has tasks up to ID `115`, the first synced task gets `116` (or higher)
- Multiple sync sources writing to different directories produce unique IDs across all of them
- All existing sync tests continue to pass
