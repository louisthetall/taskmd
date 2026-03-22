---
title: "Add default_project config for automatic project scoping"
id: "01kma473k"
status: pending
priority: low
type: feature
tags: ["global-registry", "config"]
dependencies: ["01kma46wt"]
created: "2026-03-22"
---

# Add default_project config for automatic project scoping

## Objective

Add a `default_project` key to `~/.taskmd.yaml` that automatically scopes commands to a project when the user is not inside any project directory. This avoids the "no .taskmd.yaml found" error for users who primarily work on one project.

## Tasks

- [ ] Read `default_project` from `~/.taskmd.yaml` in `LoadGlobalRegistry()` (or a sibling function)
- [ ] During config resolution, when no local `.taskmd.yaml` is found and no `--project` flag is set, check for `default_project`
- [ ] If set, resolve it against the global registry and use that project's config as if `--project` was passed
- [ ] If the default project id is not found in the registry, warn and fall back to normal behavior
- [ ] Add tests for default project resolution and fallback

## Acceptance Criteria

- With `default_project: foo` set, running `taskmd list` from `/tmp` lists tasks from the `foo` project
- `--project` flag overrides `default_project` when both are present
- If `default_project` references a non-existent registry entry, a clear warning is shown
- When inside a project directory, `default_project` is ignored (local config takes precedence)
