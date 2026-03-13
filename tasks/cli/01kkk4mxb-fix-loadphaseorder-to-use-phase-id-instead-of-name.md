---
title: "Fix loadPhaseOrder to use phase id instead of name"
id: "01kkk4mxb"
status: completed
priority: high
type: bug
tags: ["bug", "cli"]
created: "2026-03-13"
---

# Fix loadPhaseOrder to use phase id instead of name

## Objective

Fix `loadPhaseOrder()` in `apps/cli/internal/cli/next.go:119-139` to use phase `id` (falling back to `name`) instead of always using `name`. This is inconsistent with the rest of the system where task `phase` values are matched against the phase `id`.

## Steps to Reproduce

1. Define phases in `.taskmd.yaml` with both `id` and `name` fields (e.g., `id: core-cli`, `name: "Core CLI"`)
2. Set a task's `phase` frontmatter to the phase `id` (e.g., `phase: core-cli`)
3. Run `taskmd next` — the phase-based scoring bonus is not applied because `loadPhaseOrder` returns `["Core CLI"]` but the task has `phase: core-cli`

## Expected Behavior

`loadPhaseOrder` should return phase `id` values (falling back to `name` when `id` is omitted), matching the convention used by validation (`parsePhasesConfig` in `validate.go`). Phase-based scoring in `scorePhase` should then correctly match tasks.

## Actual Behavior

`loadPhaseOrder` always reads the `name` field, so `scorePhase` cannot match tasks that reference phases by `id`.

## Tasks

- [x] Update `loadPhaseOrder` in `apps/cli/internal/cli/next.go` to use `id` when present, falling back to `name`
- [x] Add/update tests in `apps/cli/internal/cli/next_test.go` covering the id-vs-name fallback
- [x] Verify existing phase-related tests still pass

## Acceptance Criteria

- `loadPhaseOrder` returns `id` values when phases define `id`
- `loadPhaseOrder` falls back to `name` when `id` is omitted (backwards compatibility)
- All existing tests pass, new tests cover the fix
