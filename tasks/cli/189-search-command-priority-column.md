---
id: "189"
title: "Add priority column to search command table output"
status: completed
priority: low
effort: small
tags: [cli, search]
created: 2026-02-21
---

# Add priority column to search command table output

## Objective

Add a `PRIORITY` column to the table output of the `taskmd search` command so users can see task priority alongside existing columns (ID, TITLE, STATUS, MATCH, SNIPPET).

## Tasks

- [x] Add `Priority` field to `search.Result` struct (if not already present)
- [x] Populate priority in search results from task metadata
- [x] Add `PRIORITY` column header to the table output in `outputSearchTable`
- [x] Format priority with color styling consistent with other commands (e.g., `list`)
- [x] Include priority in JSON/YAML output if not already present
- [x] Add/update tests for the new column

## Acceptance Criteria

- Running `taskmd search <query>` in table format displays a `PRIORITY` column between `STATUS` and `MATCH`
- Priority values are color-formatted consistently with the `list` command
- JSON and YAML output includes the priority field
- Existing tests pass and new tests cover the priority column
