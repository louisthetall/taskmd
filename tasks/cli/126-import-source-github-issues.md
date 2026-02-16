---
id: "126"
title: "Import source: GitHub Issues"
status: pending
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

- [ ] Create `internal/import/github/github.go` implementing the `Source` interface
- [ ] Interactive prompts:
  - [ ] Repository (`owner/repo`) — auto-detect from git remote if possible
  - [ ] Filter: all open, by label, by assignee, by milestone
- [ ] Non-interactive flags: `--repo`, `--filter`, `--labels`, `--milestone`, `--assignee`
- [ ] Use GitHub API via `gh` CLI or `go-github` library (prefer `gh api` for zero-config auth)
- [ ] Map GitHub fields to taskmd:
  - [ ] Issue title → `title`
  - [ ] Issue number → `external_id`
  - [ ] Issue state (open/closed) → `status` (pending/completed)
  - [ ] Issue labels → `tags`
  - [ ] Issue assignee → `owner`
  - [ ] Issue body → markdown body with original URL as reference link
  - [ ] Issue milestone → tag or group (configurable)
- [ ] Handle pagination for repos with many issues
- [ ] Add tests with mock GitHub API responses

## Acceptance Criteria

- `taskmd import --source github --repo owner/repo` imports all open issues
- Auto-detects the repo from the current git remote in interactive mode
- Labels map cleanly to tags (lowercased, spaces replaced with hyphens)
- Issue body is preserved in the task markdown body
- Each imported task includes a link back to the original GitHub issue
- Pagination works correctly for repos with 100+ issues
- Tests cover field mapping, pagination, and error handling
