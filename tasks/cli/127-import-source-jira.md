---
id: "127"
title: "Import source: Jira"
status: in-progress
priority: medium
effort: small
dependencies: ["125"]
tags:
  - cli
  - import
  - jira
touches:
  - cli/import
created: 2026-02-16
---

# Import Source: Jira

## Objective

Implement the Jira source for the `taskmd import` command so users can import issues from a Jira project into taskmd task files.

## Tasks

- [x] Create `internal/import/jira/jira.go` implementing the `Source` interface
- [x] Interactive prompts:
  - [x] Jira instance URL (e.g., `https://company.atlassian.net`)
  - [x] Authentication (API token + email, or personal access token)
  - [x] Project key (e.g., `PROJ`)
  - [ ] Filter: all non-done issues, by status, by assignee, by sprint, or JQL query
- [x] Non-interactive flags: `--url`, `--project`, `--filter`, `--jql`
- [x] Use Jira REST API v3
- [x] Map Jira fields to taskmd:
  - [x] Summary → `title`
  - [x] Issue key (e.g., `PROJ-123`) → `external_id`
  - [x] Status → `status` (map To Do→pending, In Progress→in-progress, Done→completed)
  - [x] Priority (Highest/High/Medium/Low/Lowest) → `priority` (critical/high/medium/low)
  - [x] Labels → `tags`
  - [x] Assignee → `owner`
  - [x] Description (Atlassian Document Format) → markdown body
  - [ ] Story points → `effort` mapping (configurable thresholds)
- [x] Convert Atlassian Document Format (ADF) to markdown for issue descriptions
- [x] Handle pagination for large projects
- [x] Add tests with mock Jira API responses

## Acceptance Criteria

- `taskmd import --source jira --url <url> --project PROJ` imports issues
- Jira statuses map correctly to taskmd statuses
- Jira priorities map correctly to taskmd priorities
- ADF descriptions are converted to readable markdown
- Each imported task includes a link back to the original Jira issue
- Authentication credentials are prompted securely (not echoed) in interactive mode
- Tests cover field mapping, ADF conversion, pagination, and error handling
