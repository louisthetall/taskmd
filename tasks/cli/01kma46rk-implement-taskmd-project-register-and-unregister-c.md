---
title: "Implement taskmd project register and unregister commands"
id: "01kma46rk"
status: pending
priority: high
type: feature
tags: ["global-registry", "cli-command"]
dependencies: ["01kma460m"]
created: "2026-03-22"
---

# Implement taskmd project register and unregister commands

## Objective

Add `taskmd project register` and `taskmd project unregister` commands that manage the global project registry in `~/.taskmd.yaml`. This lets users add and remove projects without manually editing the config file.

## Tasks

- [ ] Create `project.go` command file with `project` parent command and `register`/`unregister` subcommands
- [ ] `register`: resolve target path (cwd or `--path`), derive `id` from basename or `--id` flag, validate directory has `.taskmd.yaml`
- [ ] `register`: read existing `~/.taskmd.yaml`, append to `projects` list, write back preserving other config
- [ ] `register`: error if `id` already exists in the registry
- [ ] `unregister`: match by cwd path or `--id` flag, remove matching entry from `projects` list, write back
- [ ] `unregister`: error if no matching entry found
- [ ] Implement safe YAML read-modify-write for `~/.taskmd.yaml` (create file if it doesn't exist)
- [ ] Add tests for register, unregister, duplicate detection, missing `.taskmd.yaml` error

## Acceptance Criteria

- `taskmd project register` from a project directory adds it to `~/.taskmd.yaml`
- `taskmd project register --id foo --path /some/path` registers an explicit path
- `taskmd project register` errors if the directory has no `.taskmd.yaml`
- `taskmd project register` errors if the `id` is already registered
- `taskmd project unregister` removes the current directory from the registry
- `taskmd project unregister --id foo` removes by id
- Existing content in `~/.taskmd.yaml` (non-projects keys) is preserved
- `~/.taskmd.yaml` is created if it doesn't exist
