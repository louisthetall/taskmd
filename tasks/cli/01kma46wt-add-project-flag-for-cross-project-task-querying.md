---
title: "Add --project flag for cross-project task querying"
id: "01kma46wt"
status: pending
priority: high
type: feature
tags: ["global-registry", "cli-command"]
dependencies: ["01kma460m"]
created: "2026-03-22"
---

# Add --project flag for cross-project task querying

## Objective

Add a `--project <id>` flag to task-scanning commands (`list`, `next`, `graph`, `metrics`) that resolves a project from the global registry and scans its tasks. This lets users query any registered project from any directory.

## Tasks

- [ ] Add `--project` as a persistent flag on the root command (available to all subcommands)
- [ ] When `--project` is set and cwd is already inside that project's path, treat as no-op (use local config)
- [ ] When `--project` is set and cwd is outside the project, look up the project in the global registry, resolve its path, and load its `.taskmd.yaml` to determine task dir
- [ ] Override the scan directory and config context so the command operates on the resolved project
- [ ] Ensure `list`, `next`, `graph`, and `metrics` all respect the flag
- [ ] Add e2e tests: run commands with `--project` from a different directory, verify correct tasks are returned

## Acceptance Criteria

- `taskmd list --project foo` from any directory lists tasks from the project registered as `foo`
- `taskmd next --project foo` returns the next task from that project
- `--project` errors with a clear message if the id is not found in the registry
- When already inside the project directory, `--project` is a no-op (same results as without the flag)
- The project's own `.taskmd.yaml` config (phases, scopes, ID strategy) is used, not the cwd's config
