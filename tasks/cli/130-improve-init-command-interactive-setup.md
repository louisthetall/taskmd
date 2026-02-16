---
id: "130"
title: "Improve init command with interactive project setup"
status: pending
priority: high
effort: medium
tags:
  - cli
  - onboarding
  - init
touches:
  - cli/commands
created: 2026-02-16
---

# Improve init Command with Interactive Project Setup

## Objective

Enhance the existing `taskmd init` command to guide the user through a full project setup interactively. Currently `init` only writes agent config files (CLAUDE.md, etc.) and a TASKMD_SPEC.md — it doesn't create the tasks directory or a `.taskmd.yaml` config file. After this change, `taskmd init` becomes the single entry point for bootstrapping a taskmd project, creating everything needed to start using the tool immediately.

## Current Behavior

```
$ taskmd init
Created /path/to/CLAUDE.md
Created /path/to/TASKMD_SPEC.md
```

No tasks directory. No config file. The user still has to manually set things up.

## Desired Behavior

### Prompt only for what's missing

Interactive prompts fill in gaps — if the user provides all values via flags, no prompts appear. Each piece of information has a flag, a sensible default, and a prompt that only fires when the flag wasn't set and stdin is a TTY.

```bash
# Fully explicit — no prompts, just runs
taskmd init --task-dir ./tasks --claude

# Nothing provided — prompts for task dir and agent selection
taskmd init

# Partial — prompts only for what's missing (agent selection)
taskmd init --task-dir ./work/tasks
```

### Interactive example (no flags provided)

```
$ taskmd init

? Task directory: (./tasks)
? Which AI assistants do you use?
  > [x] Claude Code
    [ ] Gemini
    [ ] Codex

Created:
  ./tasks/                  (task directory)
  ./.taskmd.yaml            (project config)
  ./CLAUDE.md               (agent config)
  ./tasks/TASKMD_SPEC.md    (task specification)

You're ready! Try:
  taskmd add "My first task"
  taskmd list
  taskmd web start --open
```

### Fully non-interactive

```bash
# Flags supply everything — no prompts at all
taskmd init --task-dir ./tasks --port 3000 --claude --gemini

# Current behavior still works (no breaking changes)
taskmd init --claude --gemini
```

## Tasks

- [ ] Add "prompt only when missing" logic to `runProjectInit` in `project_init.go`:
  - [ ] Task directory path — prompt if `--task-dir` not provided (default: `./tasks`)
  - [ ] Agent config selection — prompt if none of `--claude`/`--gemini`/`--codex` provided (multi-select, default: Claude)
  - [ ] Skip all prompts when stdin is not a TTY (piped/CI) — use defaults silently
- [ ] Create the tasks directory if it doesn't exist
- [ ] Generate a minimal `.taskmd.yaml` in the project root:
  ```yaml
  dir: ./tasks
  ```
- [ ] Add `--task-dir` flag to set the task directory non-interactively
- [ ] Skip `.taskmd.yaml` creation if one already exists (warn the user, respect `--force`)
- [ ] Skip tasks directory creation if it already exists (not an error, just note it)
- [ ] Move TASKMD_SPEC.md into the tasks directory instead of the project root
- [ ] Print a helpful "next steps" summary after setup (as shown above)
- [ ] Update command help text and examples
- [ ] Add tests covering:
  - [ ] All flags provided — no prompts fired
  - [ ] No flags — prompts for task dir and agent selection
  - [ ] Partial flags — prompts only for missing values
  - [ ] Non-TTY stdin — skips prompts, uses defaults
  - [ ] Existing `.taskmd.yaml` is not overwritten without `--force`
  - [ ] Existing tasks directory is handled gracefully

## Acceptance Criteria

- `taskmd init` with no flags prompts only for task directory and agent selection, then creates everything
- `taskmd init --task-dir ./tasks --claude` provides all info via flags — no prompts
- `taskmd init --task-dir ./tasks` prompts only for agent selection (the missing piece)
- Non-TTY stdin (piped/CI) skips prompts and uses defaults
- The generated `.taskmd.yaml` is minimal and valid (only `dir`)
- Existing files/directories are never overwritten without `--force`
- The "next steps" output gives the user a clear path forward
- Tests cover all-flags, no-flags, partial-flags, and non-TTY scenarios
