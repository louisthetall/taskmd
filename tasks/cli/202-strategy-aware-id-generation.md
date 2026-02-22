---
id: "202"
title: "Implement strategy-aware ID generation"
status: pending
priority: high
effort: medium
type: feature
tags: [id, cli]
parent: "200"
dependencies: ["201"]
created: 2026-02-22
---

# Implement strategy-aware ID generation

## Objective

Update `taskmd add` and `taskmd next-id` to generate IDs according to the configured strategy in `.taskmd.yaml`. The `nextid` package currently only supports sequential IDs; extend it to support prefixed and random strategies.

## Tasks

- [ ] Add `GenerateRandom(existingIDs []string, length int)` function to `nextid` package using `crypto/rand` (base-36 alphanumeric, lowercase)
- [ ] Add `GeneratePrefixed(existingIDs []string, prefix string, padding int)` function to `nextid` package (prefix + sequential number)
- [ ] Update `resolveNextID()` in `cli/add.go` to read ID config and dispatch to the correct generation function
- [ ] Update `runNextID()` in `cli/nextid.go` to respect configured strategy
- [ ] Update `add-task` skill in `claude-code-plugin/skills/add-task/SKILL.md` to rely on `taskmd next-id` without hardcoded fallback patterns
- [ ] Add tests for random and prefixed generation (uniqueness, format, collision avoidance)

## Acceptance Criteria

- `taskmd add "Fix bug"` with `strategy: sequential` produces `042-fix-bug.md` (current behavior)
- `taskmd add "Fix bug"` with `strategy: prefixed, prefix: "dr-"` produces `dr-042-fix-bug.md`
- `taskmd add "Fix bug"` with `strategy: random, length: 6` produces `a3f9x2-fix-bug.md`
- `taskmd next-id` output matches the configured strategy
- Random IDs do not collide with existing IDs
- All existing tests continue to pass (backward compatibility)
