---
id: "116"
title: "Add external_id field to taskmd specification"
status: completed
priority: medium
effort: small
tags:
  - spec
  - sync
created: 2026-02-15
---

# Add external_id Field to taskmd Specification

## Objective

Add `external_id` as a new optional frontmatter field to the taskmd specification. This field stores the original identifier from an external system (e.g., GitHub issue number, Jira issue key) so synced tasks can be traced back to their source.

## Context

The sync engine currently writes `sync_source` and `sync_id` fields to synced task files (in `internal/sync/writer.go`), but these are not part of the official taskmd specification or validated by the parser/validator. This task formalizes `external_id` as a spec-level field.

## Tasks

- [x] Add `external_id` to the frontmatter schema in `docs/taskmd_specification.md`
  - Type: string, optional
  - Description: Identifier from an external system (e.g., `"PROJ-123"`, `"42"`)
  - Add to the Field Summary table and the Optional Fields section
- [x] Add `external_id` to the reference specification on the docs site (`apps/docs/reference/specification.md`)
- [x] Add `ExternalID` field to the task model in `internal/model/`
- [x] Update the parser (`internal/parser/`) to read `external_id` from frontmatter
- [x] Update the validator (`internal/validator/`) to accept `external_id` as a known field
- [x] Add tests for parsing and validating `external_id`

## Acceptance Criteria

- `external_id` appears in the specification as an optional string field
- The parser reads `external_id` from frontmatter into the task model
- The validator does not warn on `external_id`
- Existing tasks without `external_id` continue to work unchanged
