---
id: "139"
title: "Document all non-documented CLI commands"
status: completed
priority: medium
effort: medium
tags:
  - docs
  - cli
dependencies: ["138"]
created: 2026-02-16
---

# Document all non-documented CLI commands

## Objective

Several CLI commands are registered but missing from the CLI guide (`apps/docs/guide/cli.md`). Add documentation for every undocumented command so the guide is complete. Each command should have a Quick Reference table entry, a dedicated section with description, flags table, and usage examples -- matching the style of existing documented commands.

## Commands to Document

The following commands exist in the codebase but are absent from the CLI guide:

| Command | Source file | Short description |
|---------|------------|-------------------|
| `commit-msg` | `commit_msg.go` | Generate a conventional commit message from task metadata |
| `verify` | `verify.go` | Run verification checks for a task |
| `add` | `add.go` | Create a new task |
| `search` | `search.go` | Full-text search across task titles and bodies |
| `context` | `context.go` | Show file context for a task |
| `status` | `status.go` | Get lightweight metadata for a task (no body, no resolved deps) |
| `worklog` | `worklog.go` | View or add worklog entries for a task |
| `import` | `importcmd.go` | Import tasks from external sources |
| `spec` | `spec.go` | Generate the taskmd specification file |
| `mcp` | `mcp.go` | Start MCP server over stdio |
| `man` | `man.go` | Generate man pages |

**Note:** `commit-msg` is covered separately by task 138. The `show` and `update` commands are aliases for `get` and `set` respectively -- mention them as aliases in those sections rather than documenting separately.

## Tasks

- [x] Add all missing commands to the Quick Reference table
- [x] Document `verify` -- flags, examples, verify types (bash, assert)
- [x] Document `add` -- flags, examples for creating tasks interactively or via flags
- [x] Document `search` -- flags, examples for full-text search
- [x] Document `context` -- flags, examples for showing file context from `touches`/`context` fields
- [x] Document `status` -- flags, examples for lightweight task metadata lookup
- [x] Document `worklog` -- flags, examples for viewing and adding worklog entries
- [x] Document `import` -- flags, examples for importing from external sources
- [x] Document `spec` -- flags, examples for generating the specification
- [x] Document `mcp` -- brief section explaining MCP server usage
- [x] Document `man` -- brief section explaining man page generation
- [x] Note `show`/`update` as aliases in the `get`/`set` sections

## Acceptance Criteria

- Every registered command appears in the Quick Reference table
- Each command has a dedicated section with description, flags table, and examples
- Alias commands (`show`, `update`) are mentioned in their parent command sections
- Documentation style is consistent with existing command sections
