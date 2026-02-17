---
id: "149"
title: "Group tracks by dependency connectivity in addition to scopes"
status: pending
priority: high
effort: medium
tags: [cli, tracks, dependencies]
touches: [cli/tracks]
created: 2026-02-17
---

# Group tracks by dependency connectivity in addition to scopes

## Objective

The `tracks` command currently groups tasks into parallel tracks based only on scope overlap (`touches`). It ignores dependency relationships entirely. Tasks connected by dependency chains should be placed in the same track, since they form a natural sequential pipeline for one contributor.

## Background

A "track" is a sequence of tasks for one human contributor. Two signals indicate tasks belong in the same track:

1. **Scopes (hard constraint):** Tasks sharing a scope would cause merge conflicts if worked on in parallel. They must be sequential.
2. **Dependencies (ordering constraint):** If A depends on B, they're inherently sequential. The contributor who finishes B has context to continue with A. No parallelism is lost by grouping them.

Currently only scopes are used. Dependencies are ignored, causing related tasks (e.g., the Windows installation chain: 147 -> 148 -> 144) to scatter into the "Flexible" bucket.

## Design

### Dependency-connected components

Build connected components from the **full** task list (including non-actionable tasks) using dependency edges treated as undirected. Two actionable tasks in the same connected component should be unioned in the same track.

This is necessary because among actionable-only tasks there are zero direct dependency edges (if A depends on B and B is pending, A is blocked). The full graph bridges actionable tasks through non-actionable intermediaries.

### Combined union-find

In `assignTracks`, after the existing scope-based union step, add a dependency-based union step:

1. Build dependency connected components from all tasks (passed in from the caller)
2. Map each actionable task to its component ID
3. For actionable tasks sharing the same component, union them in the union-find

### Flexible determination

Remove the `splitByTouches` pre-filtering. Instead, after the combined union-find:
- Groups with >1 member OR at least one scope → Track
- Singletons with no scopes → Flexible

## Tasks

- [ ] Add a `buildDependencyComponents` function that computes connected components from all tasks using dependency edges (undirected BFS/DFS)
- [ ] Pass full task list (including non-actionable) into the grouping logic so dependency components can be computed
- [ ] Modify `assignTracks` to accept all scored items (not just those with touches) and a component map
- [ ] Add dependency-based union step after the scope-based union step
- [ ] Replace `splitByTouches` with post-union-find flexible determination (singleton + no scopes = flexible)
- [ ] Update existing tests to account for the new grouping behavior
- [ ] Add test: two actionable tasks connected through a blocked intermediary are placed in the same track
- [ ] Add test: tasks connected only by dependencies (no scopes) form a track, not flexible
- [ ] Add test: singleton task with no scopes and no dependency connections remains flexible
- [ ] Run `make check` (test + lint + vet) to verify

## Acceptance Criteria

- Tasks connected by dependency chains (even through non-actionable intermediaries) are grouped into the same track
- Scope-based grouping continues to work as before
- Tasks with no scopes and no dependency connections remain in the Flexible bucket
- All existing tests pass (with adjustments for the new grouping behavior)
- New tests cover dependency-based grouping scenarios
