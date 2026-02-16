---
id: "126"
title: "Import source: GitHub Issues"
status: completed
priority: high
effort: small
dependencies: ["125"]
tags:
  - cli
  - import
  - github
touches:
  - cli/import
created: 2026-02-16
---

# Import Source: GitHub Issues

## Objective

Implement the GitHub Issues source for the `taskmd import` command so users can import their open (or filtered) GitHub issues into taskmd task files with a single command.

## Tasks

- [x] Create `internal/import/github/github.go` implementing the `Source` interface
- [x] Interactive prompts:
  - [x] Repository (`owner/repo`) — auto-detect from git remote if possible
  - [x] Filter: all open, by label, by assignee, by milestone
- [x] Non-interactive flags: `--repo`, `--filter`, `--labels`, `--milestone`, `--assignee`
- [x] Use GitHub API via `gh` CLI or `go-github` library (prefer `gh api` for zero-config auth)
- [x] Map GitHub fields to taskmd:
  - [x] Issue title → `title`
  - [x] Issue number → `external_id`
  - [x] Issue state (open/closed) → `status` (pending/completed)
  - [x] Issue labels → `tags`
  - [x] Issue assignee → `owner`
  - [x] Issue body → markdown body with original URL as reference link
  - [x] Issue milestone → tag or group (configurable)
- [x] Handle pagination for repos with many issues
- [x] Add tests with mock GitHub API responses

## Acceptance Criteria

- `taskmd import --source github --repo owner/repo` imports all open issues
- Auto-detects the repo from the current git remote in interactive mode
- Labels map cleanly to tags (lowercased, spaces replaced with hyphens)
- Issue body is preserved in the task markdown body
- Each imported task includes a link back to the original GitHub issue
- Pagination works correctly for repos with 100+ issues
- Tests cover field mapping, pagination, and error handling
