---
title: "Create CLI-free Claude Code plugin with embedded spec"
id: "01kkjy95y"
status: pending
priority: high
type: feature
effort: large
tags:
  - claude-code
  - plugin
created: "2026-03-13"
---

# Create CLI-free Claude Code plugin with embedded spec

## Objective

Create a standalone Claude Code plugin (`claude-code-plugin-lite`) that provides all taskmd skills without requiring the `taskmd` CLI binary. Instead of shelling out to `taskmd`, each skill embeds the taskmd specification and uses Claude's native tools (Read, Edit, Write, Glob, Grep) to directly read/write/find task markdown files. This makes the plugin zero-dependency — users only need the plugin installed, no CLI.

The plugin should be registered in the `.claude-plugin` marketplace alongside the existing `taskmd` and `taskmd-mcp` plugins.

## Background

The existing `claude-code-plugin` relies heavily on `taskmd` CLI commands (e.g., `taskmd list`, `taskmd get`, `taskmd set`, `taskmd validate`, `taskmd add`, `taskmd next`). This requires users to install the Go CLI binary first. A CLI-free version would:

- Lower the barrier to entry (no binary installation needed)
- Work in environments where the CLI can't be installed
- Serve as a lightweight alternative for users who only need Claude Code integration

## Tasks

- [ ] Create `claude-code-plugin-lite/` directory with `.claude-plugin/plugin.json` manifest
- [ ] Write a shared spec reference document (embedded in plugin or referenced by skills) containing the taskmd specification essentials: frontmatter schema, field values, file naming, directory structure, validation rules
- [ ] Create `list-tasks` skill — uses Glob to find `tasks/**/*.md`, reads frontmatter via Read, filters/sorts/displays results
- [ ] Create `get-task` skill — finds task by ID pattern matching filenames with Glob, reads full file content
- [ ] Create `get-task-status` skill — lightweight version of get-task, reads only frontmatter
- [ ] Create `next-task` skill — scans tasks, filters pending/non-blocked, ranks by priority/dependencies/effort
- [ ] Create `add-task` skill — generates ID (per `.taskmd.yaml` strategy or default sequential), creates file with Write tool following spec format
- [ ] Create `update-task` skill — finds task file, uses Edit to modify frontmatter fields or body content
- [ ] Create `complete-task` skill — finds task file, updates status to completed via Edit, checks workflow mode in `.taskmd.yaml`
- [ ] Create `validate-tasks` skill — scans all task files, checks required fields, valid enum values, unique IDs, dependency references, no circular deps
- [ ] Create `verify-task` skill — reads task's `verify` field, runs bash checks via Bash, evaluates assert checks by inspecting code
- [ ] Create `do-task` skill — orchestrates: find task → read → mark in-progress → execute → complete
- [ ] Create `split-task` skill — reads task, assesses complexity, creates sub-task files as siblings
- [ ] Create `divide-and-conquer` skill — reads task, plans workstreams, launches parallel subagents
- [ ] Create `import-todos` skill — uses Grep to find TODO/FIXME comments, presents list, invokes add-task for selected items
- [ ] Register plugin in `.claude-plugin/marketplace.json` as a third plugin entry
- [ ] Test all skills end-to-end without the CLI installed

## Design Decisions

### Embedded Spec vs. Inline

Each skill should reference a shared spec document (e.g., `SPEC_REFERENCE.md` at the plugin root) rather than duplicating the specification in every SKILL.md. Skills can include a note like "See SPEC_REFERENCE.md for field definitions and valid values."

### ID Generation (add-task)

The add-task skill needs to handle ID generation without the CLI:
1. Read `.taskmd.yaml` for `id.strategy` (default: `sequential`)
2. Scan existing files to determine next ID
3. Generate ID according to the configured strategy

### Validation (validate-tasks)

The validate skill must implement core validation logic in instructions:
- Required fields check (`id`, `title`)
- Enum validation (`status`, `priority`, `effort`, `type`)
- Unique ID check across all scanned files
- Dependency reference validation
- Circular dependency detection

### Ranking (next-task)

The next-task skill needs ranking logic embedded in instructions:
1. Filter to `pending` status only
2. Exclude tasks with unmet dependencies (dependencies whose status != `completed`)
3. Sort by: priority (critical > high > medium > low), then effort (small first), then created date

## Acceptance Criteria

- Plugin directory exists at `claude-code-plugin-lite/` with valid `plugin.json`
- All 13 skills from the original plugin have CLI-free equivalents
- No skill uses `Bash` to invoke `taskmd` — only standard Claude tools (Read, Write, Edit, Glob, Grep) plus Bash for non-taskmd commands (e.g., running verify bash checks)
- Each skill includes the taskmd spec essentials or references the shared spec document
- Plugin is registered in `.claude-plugin/marketplace.json`
- Skills correctly handle `.taskmd.yaml` configuration (task dir, ID strategy, workflow mode, worklogs)
- A user with no `taskmd` binary installed can use all skills successfully
