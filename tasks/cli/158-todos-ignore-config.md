---
id: "158"
title: "Add configurable ignore patterns for todos command in .taskmd.yaml"
status: pending
priority: medium
effort: medium
type: improvement
tags:
  - cli
  - configuration
  - todos
created: 2026-02-19
---

# Add Configurable Ignore Patterns for Todos Command in .taskmd.yaml

## Objective

Allow users to configure persistent file and directory exclude patterns for the `taskmd todos list` command via `.taskmd.yaml`. Currently, users must pass `--exclude` flags on every invocation. A config-based approach lets users define project-specific exclude patterns once (e.g., test files, generated code, vendor directories) and have them applied automatically.

Patterns should support glob syntax so users can exclude both directories and file patterns, for example:

```yaml
# .taskmd.yaml
todos:
  exclude:
    - "*_test.go"
    - "apps/cli/internal/todos/*_test.go"
    - "dist"
    - "*.generated.*"
```

## Context

- The todos scanner (`internal/todos/scanner.go`) already supports `--include` and `--exclude` glob flags via `ScanOptions.IncludeGlobs` and `ScanOptions.ExcludeGlobs`
- The hardcoded `skipDirs` map provides baseline directory skipping (node_modules, .git, vendor, etc.)
- The existing `ignore` config field (task 097) is for the **task scanner**, not the todos scanner
- Config is loaded via Viper in `internal/cli/root.go`

## Tasks

- [ ] Add `todos.exclude` field to `.taskmd.yaml` config schema
- [ ] Read `todos.exclude` from Viper config in the todos command (`internal/cli/todos.go`)
- [ ] Merge config-based exclude patterns with CLI `--exclude` flag patterns (additive)
- [ ] Add tests for config-based exclude patterns (single pattern, multiple patterns, glob patterns)
- [ ] Test that CLI `--exclude` flags combine with config patterns (not replace)
- [ ] Test with path-based glob patterns (e.g., `apps/cli/internal/todos/*_test.go`)
- [ ] Update help text for the todos command to mention `.taskmd.yaml` configuration

## Acceptance Criteria

- Users can define `todos.exclude` patterns in `.taskmd.yaml` that are applied on every `taskmd todos list` invocation
- Glob patterns work for both file names (e.g., `*_test.go`) and paths (e.g., `apps/cli/internal/todos/*_test.go`)
- CLI `--exclude` flags are additive with config patterns (both are applied)
- When no config patterns are present, behavior is unchanged
- Tests cover config loading, pattern merging, and glob matching
