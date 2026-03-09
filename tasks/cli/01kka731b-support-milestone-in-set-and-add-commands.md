---
title: "Support milestone in set and add commands"
id: "01kka731b"
status: pending
priority: high
type: feature
dependencies: ["01kka72zy"]
tags: ["milestone", "cli"]
touches: ["cli/set", "cli/commands"]
created: "2026-03-09"
---

# Support milestone in set and add commands

## Objective

Allow users to set and modify the `milestone` field on tasks via `taskmd set` and specify it when creating tasks with `taskmd add`.

## Tasks

- [ ] Add `--milestone` flag to `taskmd set` command
- [ ] Implement writing milestone to task frontmatter via set
- [ ] Add `--milestone` flag to `taskmd add` command
- [ ] Include milestone in task creation when specified
- [ ] Add tests for set milestone (add, change, clear)
- [ ] Add tests for add with milestone

## Acceptance Criteria

- `taskmd set 042 --milestone v0.2` sets the milestone field
- `taskmd set 042 --milestone ""` clears the milestone field
- `taskmd add "New task" --milestone v0.2` creates a task with the milestone set
- Tests cover set (add/change/clear) and add scenarios
