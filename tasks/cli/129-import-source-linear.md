---
id: "129"
title: "Import source: Linear"
status: pending
priority: medium
effort: small
dependencies: ["125"]
tags:
  - cli
  - import
  - linear
touches:
  - cli/import
created: 2026-02-16
---

# Import Source: Linear

## Objective

Implement the Linear source for the `taskmd import` command so users can import issues from a Linear team/project into taskmd task files.

## Tasks

- [ ] Create `internal/import/linear/linear.go` implementing the `Source` interface
- [ ] Interactive prompts:
  - [ ] Authentication (API key)
  - [ ] Team selection (list user's teams)
  - [ ] Project selection (optional, within team)
  - [ ] Filter: all active issues, by status, by assignee, by cycle, by label
- [ ] Non-interactive flags: `--team`, `--project`, `--filter`, `--api-key`
- [ ] Use Linear GraphQL API
- [ ] Map Linear fields to taskmd:
  - [ ] Issue title → `title`
  - [ ] Issue identifier (e.g., `ENG-123`) → `external_id`
  - [ ] State → `status` (map Triage/Backlog/Todo→pending, In Progress→in-progress, Done→completed, Cancelled→cancelled)
  - [ ] Priority (Urgent/High/Medium/Low/No priority) → `priority` (critical/high/medium/low)
  - [ ] Labels → `tags`
  - [ ] Assignee → `owner`
  - [ ] Description (markdown) → markdown body (Linear already uses markdown)
  - [ ] Estimate → `effort` mapping (configurable thresholds)
  - [ ] Sub-issues → note parent/child relationships via `parent` field if both are imported
- [ ] Handle pagination via GraphQL cursors
- [ ] Add tests with mock GraphQL responses

## Acceptance Criteria

- `taskmd import --source linear --team <slug>` imports issues from a Linear team
- Linear states map correctly to taskmd statuses
- Linear priorities map correctly to taskmd priorities
- Linear markdown descriptions are preserved as-is (no conversion needed)
- Sub-issue relationships are preserved via `parent` field when both parent and child are imported
- Each imported task includes a link back to the original Linear issue
- Tests cover field mapping, pagination, sub-issue relationships, and error handling
