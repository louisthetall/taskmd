---
title: "Add --scope flag to next command"
id: "233"
status: completed
priority: medium
type: feature
tags: ["cli", "next"]
created: "2026-03-01"
---

# Add --scope flag to next command

## Objective

Add a `--scope` flag to the `taskmd next` command so users can get task recommendations filtered to a specific scope. When provided, `next` should only recommend actionable tasks whose `touches` field includes the given scope, similar to how `--scope` works in the `tracks` command.

## Reference

- `tracks` command already implements `--scope` via `tracks.Assign()` with `tracks.Options{Scope: tracksScope}` — see `apps/cli/internal/cli/tracks.go`
- Scope filtering logic in `sdk/go/tracks/tracks.go` (`assignScope()`) finds tasks where `touches` contains the scope, then expands via dependency components
- `next` command lives at `apps/cli/internal/cli/next.go`, scoring logic at `sdk/go/next/next.go`

## Tasks

- [x] Add `nextScope string` flag variable and register `--scope` flag in `next.go` `init()`
- [x] Pass scope to `next.Recommend()` options (add `Scope` field to `next.Options` if needed)
- [x] In the SDK (`sdk/go/next/next.go`), filter actionable tasks by scope (tasks whose `touches` contains the scope) before scoring
- [x] Add tests for `--scope` filtering in `next_test.go`
- [x] Add e2e test covering `taskmd next --scope <scope>`

## Acceptance Criteria

- `taskmd next --scope auth` only returns recommendations for tasks whose `touches` includes `auth`
- Without `--scope`, behavior is unchanged (no regression)
- `--scope` combines correctly with existing filters (`--filter`, `--quick-wins`, `--critical`)
- Table output adapts messaging to indicate scope filtering is active
- Tests cover: scope filtering, scope with no matches, scope combined with other flags
