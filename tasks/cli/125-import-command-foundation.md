---
id: "125"
title: "Add interactive 'import' command foundation"
status: pending
priority: high
effort: medium
tags:
  - cli
  - onboarding
  - import
touches:
  - cli/commands
  - cli/import
created: 2026-02-16
---

# Add Interactive "import" Command Foundation

## Objective

Implement a `taskmd import` command that provides an interactive, guided experience for importing tasks from external sources into taskmd. This solves the cold-start problem: developers with existing tasks in GitHub Issues, Jira, etc. need a frictionless path to populate their `tasks/` directory without manually writing files or configuring YAML sync sources.

Unlike the existing `sync` command (designed for ongoing bidirectional sync with field mappings, conflict strategies, and `.taskmd.yaml` config), `import` is a one-time onboarding tool that prioritizes ease of use over configurability.

## Design

### Interactive Flow

```
$ taskmd import

? Where are your tasks?
  > GitHub Issues
    Jira
    Trello
    Linear

? GitHub repository (owner/repo): driangle/taskmd

? Which issues do you want to import?
  > All open issues (23 found)
    Issues with specific labels...
    Issues assigned to me...

? Import directory: ./tasks

Importing 23 issues...
  ✓ #12 Fix login redirect        → tasks/001-fix-login-redirect.md
  ✓ #15 Add dark mode              → tasks/002-add-dark-mode.md
  ✓ #18 API rate limiting          → tasks/003-api-rate-limiting.md
  ...
  ✓ 23/23 imported

Done! Run `taskmd list` to see your imported tasks.
```

### Non-Interactive Mode

Support flags for scripting and AI usage:

```bash
taskmd import --source github --repo driangle/taskmd --filter "state:open"
taskmd import --source jira --project PROJ --filter "status!=Done"
```

## Tasks

- [ ] Create `internal/cli/import.go` with the `importCmd` cobra command
- [ ] Create `internal/import/` package with the core import engine:
  - [ ] `Source` interface: `Name() string`, `Prompt() (Config, error)`, `Fetch(Config) ([]ExternalTask, error)`
  - [ ] `ExternalTask` struct: generic representation of an external task (title, body, status, priority, labels, assignee, URL)
  - [ ] `Mapper`: converts `ExternalTask` to taskmd frontmatter + body
  - [ ] `Writer`: generates task files with sequential IDs and slugs
- [ ] Implement interactive prompts using a lightweight prompt library (e.g., `survey` or `bubbletea` prompter)
- [ ] Support `--source` flag to skip source selection prompt
- [ ] Support `--dir` flag for target directory (default from `.taskmd.yaml` or `./tasks`)
- [ ] Support `--dry-run` to preview what would be imported without writing files
- [ ] Support `--format json` to output import results as JSON
- [ ] Map external statuses to taskmd statuses (open→pending, closed→completed, etc.)
- [ ] Map external labels/tags to taskmd tags
- [ ] Map external priority/severity to taskmd priority where available
- [ ] Preserve the external source URL in the task body as a reference link
- [ ] Store `external_id` in frontmatter for traceability
- [ ] Detect and skip duplicates if `external_id` already exists in task files
- [ ] Print summary on completion (count imported, skipped, errors)
- [ ] Add tests for the core engine (mapper, writer, duplicate detection)

## Acceptance Criteria

- `taskmd import` launches an interactive wizard that guides the user through source selection and configuration
- `taskmd import --source <name> [flags]` works non-interactively for scripting
- `--dry-run` shows what would be imported without writing files
- Imported tasks pass `taskmd validate`
- Duplicate imports are detected and skipped via `external_id`
- The `Source` interface is simple enough that adding a new provider requires only one file
- Tests cover the engine, mapper, writer, and duplicate detection (source-specific fetching tested in source tasks)
