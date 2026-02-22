---
id: "200"
title: "Configurable task ID strategy"
status: in-progress
priority: high
effort: large
type: feature
tags: [id, config, multi-user]
created: 2026-02-22
---

# Configurable task ID strategy

## Objective

When multiple contributors create tasks independently on separate branches, sequential numeric IDs inevitably collide on merge. Make the ID generation strategy configurable via `.taskmd.yaml` so teams can choose the approach that fits their workflow.

Supported strategies:
- **sequential** (current default) — `001`, `002`, `003`
- **prefixed** — user-configured prefix + sequential number, e.g. `dr-001`, `jk-002`
- **random** — short random alphanumeric string, e.g. `a3f9x2`

## Tasks

- [ ] Add `id` config section to `.taskmd.yaml` (task 201)
- [ ] Implement strategy-aware ID generation (task 202)
- [ ] Add `taskmd deduplicate` command for collision resolution (task 203)

## Acceptance Criteria

- Users can set `id.strategy` in `.taskmd.yaml` to `sequential`, `prefixed`, or `random`
- `taskmd add` and `taskmd next-id` respect the configured strategy
- `taskmd deduplicate` can resolve ID collisions after merge
- Default behavior is unchanged when no `id` config is present
- Specification and documentation are updated
