---
id: "204"
title: "Support uuid ID strategy"
status: in-progress
priority: medium
effort: medium
type: feature
tags: [id, cli]
parent: "200"
dependencies: ["201", "202"]
created: 2026-02-22
---

# Support uuid ID strategy

## Objective

Add `uuid` as a fourth ID strategy option alongside `sequential`, `prefixed`, and `random`. When configured, `taskmd add` should generate UUIDs (v4) as task IDs, producing filenames like `f47ac10b-slug.md`. This is useful for teams that want globally unique, non-sequential identifiers without risk of collision across repositories.

## Tasks

- [x] Add `"uuid"` to `validIDStrategies` in `validator/validator.go`
- [x] Add `GenerateUUID()` function to `nextid` package using `crypto/rand` (RFC 4122 v4 format, or a shortened variant configurable via `length`)
- [x] Update `resolveNextID()` in `cli/add.go` to dispatch to UUID generation when strategy is `uuid`
- [x] Update `runNextID()` in `cli/nextid.go` to support `uuid` strategy
- [x] Update `splitFilenameID()` in `parser/markdown.go` to recognize UUID-formatted ID segments (8-4-4-4-12 hex pattern or truncated UUID)
- [x] Update `docs/taskmd_specification.md` ID Generation section to document `uuid` strategy and run `make sync-spec`
- [x] Add tests for UUID generation (format, uniqueness) and parser recognition

## Acceptance Criteria

- `strategy: uuid` is accepted as a valid config value (no validation error)
- `taskmd add "Fix bug"` with `strategy: uuid` produces a file like `f47ac10b-fix-bug.md`
- `taskmd next-id` with `strategy: uuid` outputs a valid UUID (or truncated UUID based on `length`)
- `deriveFieldsFromFilename()` correctly parses UUID-prefixed filenames
- All existing tests continue to pass (backward compatibility)
- Invalid strategies (e.g. `strategy: guid`) still produce validation errors
