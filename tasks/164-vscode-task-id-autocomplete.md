---
id: "164"
title: "VSCode extension: task ID autocomplete"
status: completed
priority: low
effort: medium
tags: []
touches:
  - vscode
created: 2026-02-20
---

# VSCode Extension: Task ID Autocomplete

## Objective

Suggest existing task IDs when editing `dependencies` and `parent` fields, so users don't have to remember or look up IDs manually.

## Tasks

- [ ] Scan the task directory to collect all task IDs and titles
- [ ] Implement a `CompletionItemProvider` that triggers inside array values for `dependencies` and after `parent:`
- [ ] Show task ID as the completion label and task title as the detail/description
- [ ] Cache the task list and refresh on file system changes
- [ ] Add tests for completion triggering and suggestion content

## Acceptance Criteria

- Typing inside `dependencies: ["` suggests existing task IDs with their titles
- Typing after `parent: ` suggests existing task IDs with their titles
- Completions update when task files are added or removed
- No completions appear for non-ID fields
