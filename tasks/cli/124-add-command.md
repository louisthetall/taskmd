---
id: "124"
title: "Add CLI 'add' command for quick task creation"
status: pending
priority: high
effort: medium
tags:
  - cli
  - ux
  - onboarding
touches:
  - cli/commands
created: 2026-02-16
---

# Add CLI "add" Command for Quick Task Creation

## Objective

Implement a `taskmd add` command that creates a new task file from the command line without requiring users to manually write YAML frontmatter or know the file naming convention. This is the single biggest onboarding friction point — every other task tool has a quick-add mechanism, but taskmd currently requires manual file creation.

## Tasks

- [ ] Create `internal/cli/add.go` with a new `addCmd` cobra command
- [ ] Accept a positional argument for the task title: `taskmd add "Fix the login bug"`
- [ ] Auto-assign the next sequential ID using the existing `next-id` logic
- [ ] Support flags for common fields:
  - `--priority` (`low`, `medium`, `high`, `critical`) — default: `medium`
  - `--effort` (`small`, `medium`, `large`)
  - `--tags` (comma-separated list)
  - `--status` — default: `pending`
  - `--owner`
  - `--depends-on` (comma-separated list of task IDs)
  - `--parent` (task ID)
  - `--group` (subdirectory name, e.g., `cli`, `web`)
- [ ] Generate a slug from the title for the filename (lowercase, hyphenated)
- [ ] Write the task file to the correct directory (`tasks/` or `tasks/<group>/`)
- [ ] Include a markdown body template with `# Title`, `## Objective`, `## Tasks`, and `## Acceptance Criteria` sections
- [ ] Print the created file path and task ID to stdout on success
- [ ] Support `--format json` to output the created task as JSON (for scripting/AI integration)
- [ ] Open the file in `$EDITOR` if `--edit` flag is passed
- [ ] Add comprehensive tests in `internal/cli/add_test.go`
- [ ] Add shell completion for flag values (status, priority, effort enums)

## Acceptance Criteria

- `taskmd add "My task"` creates a valid task file with correct frontmatter and body
- The generated file passes `taskmd validate`
- The ID is always unique and sequential
- The filename follows the `NNN-slug.md` convention
- All flags correctly map to frontmatter fields
- The command works with both `--dir` flag and `.taskmd.yaml` config for the tasks directory
- Error messages are clear when the tasks directory doesn't exist
- Tests cover happy path, all flags, edge cases (special characters in title, duplicate handling), and error conditions
