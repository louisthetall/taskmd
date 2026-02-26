---
title: "Add update-task skill to plugin"
id: "224"
status: completed
priority: medium
type: feature
tags: []
created: "2026-02-26"
---

# Add update-task skill to plugin

## Objective

Add a new "update-task" skill to the Claude Code plugin (`claude-code-plugin/skills/`) that allows users to update an existing task's fields (status, priority, title, tags, dependencies, etc.) via a natural language command. This complements the existing `get-task`, `add-task`, `complete-task`, and `set` CLI functionality.

## Tasks

- [x] Create `claude-code-plugin/skills/update-task/` directory with skill definition
- [x] Define the skill prompt that parses user intent and determines which fields to update
- [x] For CLI-supported fields (status, priority, effort, tags, owner, parent, type, PRs), instruct the agent to use `taskmd set` with appropriate flags
- [x] For fields not supported by `taskmd set` (title, dependencies, custom frontmatter), instruct the agent to edit the task file directly
- [x] Document which fields use CLI vs direct file editing in the skill prompt
- [x] Register the skill in the plugin manifest
- [x] Handle edge cases (invalid IDs, invalid field values) with clear error messages
- [x] Test the skill with various update scenarios

## Acceptance Criteria

- `/taskmd:update-task` skill is available and listed in the plugin
- Users can update task fields via natural language (e.g., "set task 042 to high priority and in-progress")
- Fields supported by `taskmd set` (status, priority, effort, tags, owner, parent, type, PRs) are updated via CLI
- Fields not supported by `taskmd set` (title, dependencies, custom frontmatter) are updated by editing the task file directly
- The agent runs `taskmd validate` after direct file edits to ensure correctness
- Invalid task IDs or field values produce helpful error messages
- Skill follows the same conventions as existing skills (e.g., `add-task`, `complete-task`)
