---
id: "167"
title: "VSCode extension: unknown field warnings"
status: completed
priority: low
effort: small
tags: []
touches:
  - vscode
created: 2026-02-20
---

# VSCode Extension: Unknown Field Warnings

## Objective

Warn users when frontmatter contains field names not in the taskmd schema, catching typos like `stauts` instead of `status` or `dependecies` instead of `dependencies`.

## Tasks

- [x] Add validation rule that checks each frontmatter key against the known field set in `schema.ts`
- [x] Report unknown fields as warnings (not errors, since the CLI silently ignores them)
- [x] Highlight the key range for unknown fields
- [x] Add tests for unknown field detection

## Acceptance Criteria

- A typo like `stauts: pending` produces a warning diagnostic
- Valid fields produce no warnings
- The `verify` field and all its sub-fields are not flagged
- Warning message includes the unknown field name
