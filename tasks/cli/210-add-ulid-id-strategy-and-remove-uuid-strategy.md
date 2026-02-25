---
title: "Add ULID ID strategy and remove UUID strategy"
id: "210"
status: completed
priority: medium
type: feature
tags: ["id-strategy"]
created: "2026-02-25"
---

# Add ULID ID strategy and remove UUID strategy

## Objective

Add ULID (Universally Unique Lexicographically Sortable Identifier) as a new ID strategy and remove the existing UUID strategy. ULIDs combine a millisecond-precision timestamp with cryptographic randomness, producing IDs that are both globally unique and lexicographically sortable by creation time — a strict improvement over the current UUID strategy.

## Tasks

- [x] Add a Go ULID generation function in `internal/nextid/nextid.go` (Crockford Base32, 26 chars by default, with configurable length via truncation)
- [x] Remove `GenerateUUID` from `internal/nextid/nextid.go`
- [x] Update `internal/cli/nextid.go` to route `"ulid"` strategy and remove `"uuid"` case
- [x] Update `internal/validator/validator.go` ID config to accept `"ulid"` and reject `"uuid"`
- [x] Update `internal/sync/` engines that reference UUID strategy
- [x] Update `.taskmd.yaml` config schema/docs references from `"uuid"` to `"ulid"`
- [x] Update `docs/taskmd_specification.md` strategy table and config example
- [x] Run `make sync-spec` to propagate spec changes
- [x] Add/update unit tests in `internal/nextid/nextid_test.go` for ULID generation
- [x] Update CLI tests in `internal/cli/nextid_test.go`
- [x] Run `make check` and `make e2e` to verify nothing breaks

## Acceptance Criteria

- `taskmd next-id` with `strategy: ulid` produces a valid ULID (Crockford Base32, sortable by time)
- `strategy: uuid` is no longer accepted in config and produces a clear error
- All existing tests pass; new tests cover ULID generation, collision avoidance, and sortability
- Specification and embedded docs are updated and in sync (`TestSpecTemplate_MatchesCanonicalSpec` passes)
