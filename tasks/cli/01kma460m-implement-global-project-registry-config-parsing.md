---
title: "Implement global project registry config parsing"
id: "01kma460m"
status: completed
priority: critical
type: feature
tags: ["global-registry", "config"]
created: "2026-03-22"
---

# Implement global project registry config parsing

## Objective

Add the ability to read a `projects` list from `~/.taskmd.yaml` independently of the existing viper config walk-up. This is the foundation for all global project registry features — every subsequent task depends on this config layer existing.

## Tasks

- [x] Define `GlobalProjectEntry` struct (`ID`, `Name`, `Path` fields)
- [x] Implement `LoadGlobalRegistry()` that reads `~/.taskmd.yaml` (or `$TASKMD_HOME_CONFIG`), parses only the `projects` key, and returns `[]GlobalProjectEntry`
- [x] Handle `~` expansion in paths
- [x] Validate entries: `path` required, must be a directory, should contain `.taskmd.yaml` (warn if missing), unique `id` values
- [x] Derive `id` from directory basename when omitted
- [x] Add unit tests with temp home config files

## Acceptance Criteria

- `LoadGlobalRegistry()` correctly parses a `~/.taskmd.yaml` with `projects` entries
- Tilde paths are expanded to absolute paths
- Missing `id` falls back to directory basename
- Validation errors are returned for missing `path`, duplicate `id`, non-existent directories
- Loading does not interfere with the existing viper config walk-up
