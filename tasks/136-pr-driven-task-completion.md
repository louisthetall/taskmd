---
id: "136"
title: "PR-driven task completion workflow"
status: pending
priority: high
effort: large
tags:
  - workflow
  - github-actions
  - cli
  - agent
created: 2026-02-16
---

# PR-Driven Task Completion Workflow

## Objective

Enable a workflow where an AI agent works on a task, opens a PR (setting the task to `in-review`), a human reviews and gives feedback, and only after the PR is merged does the task get marked `completed`. This requires a new `in-review` status, a `pr` frontmatter field, and a reusable GitHub Action.

## Background

Currently, agents mark tasks as `completed` immediately after finishing work. This skips the human review step entirely. The desired workflow is:

1. Agent picks up a task, sets status to `in-progress`
2. Agent does the work, opens a PR
3. Agent sets status to `in-review` and **stops** (does NOT mark completed)
4. Human reviews the PR, requests changes if needed
5. PR is merged
6. GitHub Action automatically sets status to `completed`

This ensures humans stay in the loop and tasks accurately reflect their review state.

## Tasks

### Part 1: New `in-review` status

- [ ] Add `StatusInReview` constant to `apps/cli/internal/model/task.go`
- [ ] Add `in-review` to valid statuses in `apps/cli/internal/taskfile/taskfile.go`
- [ ] Update status lifecycle diagram in `docs/taskmd_specification.md`
- [ ] Run `make sync-spec` to sync specification copies
- [ ] Update agent templates (`docs/templates/CLAUDE.md`, `CODEX.md`, `GEMINI.md`) to instruct agents: after opening a PR, set status to `in-review` and stop
- [ ] Update CLI templates in `apps/cli/internal/cli/templates/` to match
- [ ] Add tests for the new status (parsing, validation, transitions)

### Part 2: `pr` frontmatter field

- [ ] Add `PRs []string` field to task model in `apps/cli/internal/model/task.go`
- [ ] Support reading/writing the `pr` field in `apps/cli/internal/taskfile/` parser
- [ ] Add `--add-pr` flag to the `set` command (similar to `--add-tag`)
- [ ] Add `--remove-pr` flag to the `set` command (similar to `--remove-tag`)
- [ ] Document the `pr` field in `docs/taskmd_specification.md` and sync
- [ ] Add tests for PR field parsing, `--add-pr`, and `--remove-pr`

### Part 3: GitHub Action

- [ ] Create a reusable GitHub Action at `.github/actions/taskmd-complete/action.yml`
- [ ] Action should extract task ID from PR body, branch name, or labels on merge
- [ ] Action runs `taskmd set --task-id X --status completed` and commits the change
- [ ] Create an example workflow YAML that users can copy into their repos
- [ ] Document the GitHub Action setup in `docs/` or the specification
- [ ] Add tests or validation for the action (e.g., shellcheck, dry-run mode)

## Acceptance Criteria

- [ ] `in-review` is a valid task status accepted by the CLI and specification
- [ ] Tasks can store one or more PR URLs in a `pr` frontmatter field
- [ ] `taskmd set --task-id X --add-pr <url>` adds a PR URL to the task
- [ ] `taskmd set --task-id X --remove-pr <url>` removes a PR URL from the task
- [ ] Agent templates instruct agents to set `in-review` after opening a PR
- [ ] A reusable GitHub Action marks tasks `completed` when their PR is merged
- [ ] Example workflow YAML is provided for users to adopt
- [ ] All new features have comprehensive tests
- [ ] `taskmd-dev validate` passes with no errors
