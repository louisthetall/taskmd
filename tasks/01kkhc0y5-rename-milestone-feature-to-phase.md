---
title: "Rename milestone feature to phase"
id: "01kkhc0y5"
status: completed
priority: high
type: chore
tags: ["refactor"]
created: "2026-03-12"
---

# Rename milestone feature to phase

## Objective

Rename the "milestone" concept to "phase" throughout the entire codebase. This is a terminology change — the feature behavior remains the same, but all references to "milestone" (in the context of task grouping/phases) should become "phase".

## Tasks

- [ ] Update the taskmd specification (`docs/taskmd_specification.md`) — rename the `milestone` frontmatter field to `phase`, update all descriptions and examples
- [ ] Run `make sync-spec` to propagate spec changes to embedded CLI template and docs site
- [ ] Update the Go parser/scanner code — rename `milestone` field handling to `phase` in struct definitions, parsing logic, and field mappings
- [ ] Update CLI commands — rename `--milestone` flags to `--phase`, update help text and usage examples
- [ ] Update CLI output formatting — table headers, JSON/YAML keys, filter labels
- [ ] Update `.taskmd.yaml` config schema if milestone is referenced there
- [ ] Update all test files — rename test cases, assertions, and test data referencing milestone
- [ ] Update documentation (`docs/`, `apps/docs/`) — replace milestone references with phase
- [ ] Update `CLAUDE.md` and `PLAN.md` if they reference milestones
- [ ] Update any existing task files that use `milestone` in their frontmatter
- [ ] Run `make test`, `make e2e`, `make lint` and fix any issues
- [ ] Run `taskmd validate` to ensure all task files remain valid

## Acceptance Criteria

- No references to "milestone" remain in source code, docs, or config (except git history)
- The `phase` frontmatter field works identically to how `milestone` worked before
- All CLI flags/filters use `--phase` instead of `--milestone`
- All tests pass (`make test`, `make e2e`)
- Linter passes (`make lint`)
- `taskmd validate` passes
- Spec is in sync across all three locations
