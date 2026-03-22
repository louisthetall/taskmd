---
title: "Implement taskmd projects command"
id: "01kma46mh"
status: pending
priority: high
type: feature
tags: ["global-registry", "cli-command"]
dependencies: ["01kma460m"]
created: "2026-03-22"
---

# Implement taskmd projects command

## Objective

Add a `taskmd projects` command that lists all globally registered projects with summary task stats. This gives users a dashboard view of all their projects from any directory.

## Tasks

- [ ] Create `projects.go` command file in `internal/cli/`
- [ ] Call `LoadGlobalRegistry()` to get registered projects
- [ ] For each project, resolve its config and scan its task directory to compute stats (total, pending, in-progress, completed)
- [ ] Render table output with columns: PROJECT, PATH, TASKS, PENDING, IN-PROGRESS, COMPLETED
- [ ] Support `--format json|yaml|table` flag
- [ ] Handle errors gracefully: warn for unreachable projects, skip and continue
- [ ] Add tests with mock registry and temp project directories

## Acceptance Criteria

- `taskmd projects` lists all registered projects with correct task counts
- JSON and YAML output formats include project metadata and stats
- Unreachable projects (deleted path) show a warning but don't fail the command
- Works from any directory (does not depend on cwd)
