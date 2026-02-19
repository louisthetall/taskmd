---
id: "160"
title: "Add a Claude Code skill to convert TODOs into tasks"
status: completed
priority: medium
effort: medium
dependencies:
  - "155"
tags:
  - plugin
  - developer-experience
  - agent-tooling
created: 2026-02-19
---

# Add a Claude Code Skill to Convert TODOs into Tasks

## Objective

Create a new skill for the taskmd Claude Code plugin that allows a coding agent to discover TODO/FIXME comments in the codebase via `taskmd todos list`, present them to the user as a numbered list, let the user select which ones to convert into task files, and then use the existing `/add-task` skill to create each selected task.

## Tasks

- [x] Create a new skill definition (e.g. `todos-to-tasks` or `import-todos`) in the plugin skill directory
- [x] Skill instructions should tell the agent to:
  - [x] Run `taskmd todos list --format json` to discover all TODOs
  - [x] Present results as a numbered list showing: number, file:line, marker, text
  - [x] Ask the user which TODOs to convert (by number, range, or "all")
  - [x] For each selected TODO, invoke `/add-task` with context from the TODO (marker as type hint, text as title/description, file path as context)
- [x] Register the skill in the plugin manifest so it appears as a slash command (e.g. `/import-todos`)
- [x] Ensure the skill handles edge cases:
  - [x] No TODOs found — inform the user
  - [x] User selects none — exit gracefully
  - [x] Duplicate detection — warn if a TODO's text closely matches an existing task title
- [x] Write documentation for the skill in the plugin README or help text

## Acceptance Criteria

- Agent can invoke `/import-todos` (or equivalent) and see a numbered list of TODOs from the codebase
- User can select specific TODOs by number (e.g. "1, 3, 5"), range (e.g. "1-5"), or "all"
- Selected TODOs are converted into properly formatted task files via `/add-task`
- The created tasks reference the source file and line in their context field
- When no TODOs are found, the agent reports this clearly
- The skill is discoverable in the plugin's skill list
