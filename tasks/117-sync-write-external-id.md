---
id: "117"
title: "Write external_id in synced task files"
status: completed
priority: medium
effort: small
dependencies:
  - "116"
tags:
  - cli
  - sync
created: 2026-02-15
---

# Write external_id in Synced Task Files

## Objective

Update the sync writer to use the new `external_id` frontmatter field (from task 116) when creating and updating synced task files. Replace the current ad-hoc `sync_id` field with the spec-level `external_id`.

## Context

The sync writer (`internal/sync/writer.go`) currently writes `sync_source` and `sync_id` to synced task frontmatter. Now that `external_id` is a proper spec field, synced files should use `external_id` instead of `sync_id`. Both GitHub and Jira sources already populate `ExternalID` on `ExternalTask` -- this task ensures that value flows through the writer into the frontmatter.

## Tasks

- [x] In `internal/sync/writer.go`, replace `sync_id` with `external_id` in `renderTaskFile`
- [x] Ensure `UpdateSyncedTaskFile` preserves `external_id` on updates
- [x] Keep `sync_source` as-is (it indicates provenance, separate from the external ID)
- [x] Update existing tests in `internal/sync/` that assert on rendered frontmatter
- [x] Add a test that verifies `external_id` appears in the written file for both create and update paths

## Acceptance Criteria

- Newly synced task files contain `external_id: "PROJ-123"` (or equivalent) in frontmatter
- Updated synced task files retain their `external_id`
- `sync_id` is no longer written to new files
- All sync tests pass
