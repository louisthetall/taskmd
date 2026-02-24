---
id: "206"
title: "Show children task status in `taskmd get` output"
status: completed
priority: medium
effort: small
tags: [cli, get]
created: 2026-02-24
---

# Show children task status in `taskmd get` output

## Objective

Update the `taskmd get` command to display the status of children tasks when a task has children (i.e., when other tasks reference it via the `parent` field). This gives users a quick overview of sub-task progress without needing to run a separate `list` or `graph` command.

## Tasks

- [x]Identify how `taskmd get` currently renders task output (inspect `internal/cli/get.go` or equivalent)
- [x]After loading the target task, scan for tasks whose `parent` field matches the target task's ID
- [x]Add a "Children" section to the `get` output that lists each child task's ID, title, and status
- [x]Ensure the children section is omitted when the task has no children
- [x]Support all output formats (table/text, JSON, YAML) with the children data included
- [x]Add tests for the new children status display
- [x]Update CLI help text if needed

## Acceptance Criteria

- When running `taskmd get <ID>` on a task that has children, the output includes a section showing each child task's ID, title, and status
- When the task has no children, the output is unchanged from current behavior
- JSON and YAML output formats include a `children` array with id, title, and status for each child
- Tests cover both parent-with-children and leaf-task scenarios
