---
title: "Walk up directory tree to find project root for task scanning"
id: "01kka7e28"
status: completed
priority: medium
type: feature
tags: []
created: "2026-03-09"
---

# Walk up directory tree to find project root for task scanning

## Objective

When `taskmd` is run from a subdirectory (e.g., `apps/cli/`), it cannot find `.taskmd.yaml` or task files because config discovery only searches cwd and `$HOME`. Update the config/scanning logic to walk up the directory tree — like `git` does — until it finds a `.taskmd.yaml` or `.git` root, then resolve the task directory relative to that project root.

## Tasks

- [x] Update `initConfig()` in `root.go` to walk up from cwd, adding ancestor directories as viper config paths until `.git` or filesystem root is reached
- [x] Add `resolveRelativeToConfig()` helper so that relative `dir` values in `.taskmd.yaml` resolve against the config file's location, not cwd
- [x] Update `resolveTaskDir()` to use the new helper when the dir comes from config
- [x] Simplify `resolveProjectRoot()` in `verify.go` to leverage the improved config discovery
- [x] Add unit tests for walk-up discovery from subdirectories
- [x] Add e2e test: run `taskmd list` from a subdirectory and verify it finds tasks

## Acceptance Criteria

- Running `taskmd list` from any subdirectory within a project finds tasks correctly
- Running `taskmd set <id> --done` from a subdirectory works
- Relative `dir` config values (e.g., `dir: ./tasks`) resolve relative to the `.taskmd.yaml` location, not cwd
- Existing behavior when running from the project root is unchanged
- Walk-up stops at `.git` boundary or filesystem root
