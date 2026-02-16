---
id: "127"
title: "Import source: Jira"
status: pending
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

- [ ] Create `internal/import/jira/jira.go` implementing the `Source` interface
- [ ] Interactive prompts:
  - [ ] Jira instance URL (e.g., `https://company.atlassian.net`)
  - [ ] Authentication (API token + email, or personal access token)
  - [ ] Project key (e.g., `PROJ`)
  - [ ] Filter: all non-done issues, by status, by assignee, by sprint, or JQL query
- [ ] Non-interactive flags: `--url`, `--project`, `--filter`, `--jql`
- [ ] Use Jira REST API v3
- [ ] Map Jira fields to taskmd:
  - [ ] Summary → `title`
  - [ ] Issue key (e.g., `PROJ-123`) → `external_id`
  - [ ] Status → `status` (map To Do→pending, In Progress→in-progress, Done→completed)
  - [ ] Priority (Highest/High/Medium/Low/Lowest) → `priority` (critical/high/medium/low)
  - [ ] Labels → `tags`
  - [ ] Assignee → `owner`
  - [ ] Description (Atlassian Document Format) → markdown body
  - [ ] Story points → `effort` mapping (configurable thresholds)
- [ ] Convert Atlassian Document Format (ADF) to markdown for issue descriptions
- [ ] Handle pagination for large projects
- [ ] Add tests with mock Jira API responses

## Acceptance Criteria

- `taskmd import --source jira --url <url> --project PROJ` imports issues
- Jira statuses map correctly to taskmd statuses
- Jira priorities map correctly to taskmd priorities
- ADF descriptions are converted to readable markdown
- Each imported task includes a link back to the original Jira issue
- Authentication credentials are prompted securely (not echoed) in interactive mode
- Tests cover field mapping, ADF conversion, pagination, and error handling
