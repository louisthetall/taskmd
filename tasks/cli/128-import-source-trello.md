---
id: "128"
title: "Import source: Trello"
status: pending
priority: low
effort: small
dependencies: ["125"]
tags:
  - cli
  - import
  - trello
touches:
  - cli/import
created: 2026-02-16
---

# Import Source: Trello

## Objective

Implement the Trello source for the `taskmd import` command so users can import cards from a Trello board into taskmd task files.

## Tasks

- [ ] Create `internal/import/trello/trello.go` implementing the `Source` interface
- [ ] Interactive prompts:
  - [ ] Authentication (API key + token)
  - [ ] Board selection (list user's boards, let them pick)
  - [ ] Filter: all cards, by list (e.g., "To Do", "In Progress"), by label
- [ ] Non-interactive flags: `--board`, `--list`, `--labels`, `--api-key`, `--token`
- [ ] Use Trello REST API
- [ ] Map Trello fields to taskmd:
  - [ ] Card name → `title`
  - [ ] Card ID → `external_id`
  - [ ] List name → `status` (map common list names: To Do→pending, Doing→in-progress, Done→completed)
  - [ ] Labels → `tags` (use label name, lowercased)
  - [ ] Members → `owner`
  - [ ] Card description → markdown body
  - [ ] Checklists → `## Tasks` section with checkbox items
  - [ ] Due date → note in body
- [ ] Convert Trello checklists to taskmd subtask checkboxes
- [ ] Handle pagination for boards with many cards
- [ ] Add tests with mock Trello API responses

## Acceptance Criteria

- `taskmd import --source trello --board <id>` imports cards from a board
- Trello list names map to taskmd statuses with sensible defaults
- Trello checklists are converted to markdown checkbox lists in the task body
- Labels map to tags
- Each imported task includes a link back to the original Trello card
- Tests cover field mapping, checklist conversion, and error handling
