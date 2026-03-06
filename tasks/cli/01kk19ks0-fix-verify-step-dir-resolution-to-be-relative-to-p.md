---
title: "Fix verify step dir resolution to be relative to project root"
id: "01kk19ks0"
status: completed
priority: high
type: bug
tags: ["bug", "cli"]
created: "2026-03-06"
verify:
  - type: bash
    run: "go test ./internal/cli/ -run TestResolveProjectRoot -count=1"
    dir: "apps/cli"
  - type: bash
    run: "make test"
    dir: "apps/cli"
  - type: bash
    run: "make lint"
    dir: "apps/cli"
---

# Fix verify step dir resolution to be relative to project root

## Steps to Reproduce

1. Create a task with `verify` steps that use `dir: "apps/cli"`
2. Run `taskmd verify <id>` from a subdirectory (e.g. `apps/cli/`) with `--task-dir ../../tasks`
3. The verify step fails because `dir` resolves relative to cwd, not project root

## Expected Behavior

`dir: "apps/cli"` in verify steps should resolve relative to the project root (where `.taskmd.yaml` lives), regardless of the user's current working directory.

## Actual Behavior

`resolveProjectRoot()` returned a relative path (`"."`) derived from `filepath.Dir(".taskmd.yaml")`. When running from a subdirectory, `dir: "apps/cli"` resolved to `apps/cli/apps/cli/` (doubled path), causing the command to fail.

## Tasks

- [x] Make `resolveProjectRoot()` return an absolute path via `filepath.Abs()`
- [x] Add upward directory walk to find `.taskmd.yaml` when viper hasn't loaded a config
- [x] Add tests for `resolveProjectRoot()` absolute path and config dir matching
- [x] Fix template tests that relied on the old buggy relative path behavior
- [x] Verify all tests pass (unit, e2e, lint)

## Acceptance Criteria

- `resolveProjectRoot()` always returns an absolute path
- Verify steps with `dir` resolve correctly from any working directory
- All existing tests pass
