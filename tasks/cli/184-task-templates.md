---
id: "184"
title: "Add task templates support"
status: completed
priority: medium
effort: medium
type: feature
tags:
  - cli
  - templates
  - productivity
created: 2026-02-20
---

# Add Task Templates Support

## Objective

Allow users to define reusable task templates for common task types (e.g., bug report, feature request, sprint retrospective). Templates reduce boilerplate and enforce consistency across tasks.

## Tasks

- [ ] Define a template format (markdown files in `.taskmd/templates/` or `tasks/.templates/`)
- [ ] Support template variables (e.g., `{{title}}`, `{{date}}`, `{{id}}`)
- [ ] Add `taskmd add --template <name>` flag to create a task from a template
- [ ] Add `taskmd templates list` subcommand to list available templates
- [ ] Ship 2-3 built-in templates (feature, bug, chore) as defaults
- [ ] Support project-level templates (in the project directory) and user-level templates (in `~/.taskmd/templates/`)
- [ ] Auto-fill `created` date and next available `id` in templates
- [ ] Add template selection to the `init` command for bootstrapping
- [ ] Add tests for template discovery, variable substitution, and task creation
- [ ] Document the template format and usage

## Acceptance Criteria

- `taskmd add --template bug` creates a new task pre-filled from the bug template
- `taskmd templates list` shows all available templates with descriptions
- Templates support variable substitution for dynamic fields
- Project-level templates override built-in templates of the same name
- Created tasks pass `taskmd validate`
