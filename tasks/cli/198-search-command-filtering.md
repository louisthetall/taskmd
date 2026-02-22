---
id: "198"
title: "Add filtering, sorting, and limit flags to search command"
status: completed
priority: medium
effort: small
tags: [cli, search]
created: 2026-02-22
---

# Add filtering, sorting, and limit flags to search command

## Objective

The `search` command currently only supports `--format` but lacks the filtering, sorting, and limit capabilities that `list` already provides. Add `--filter`, `--sort`, and `--limit` flags to the `search` command so users can narrow down search results (e.g., search for "authentication" but only in high-priority or pending tasks).

## Tasks

- [x] Add `--filter` flag (reuse `applyFilters` from list command)
- [x] Add `--sort` flag (reuse `sortTasks` from list command)
- [x] Add `--limit` flag
- [x] Update command help text and examples
- [x] Add tests for new flags
- [x] Add e2e tests for search with filters

## Acceptance Criteria

- `taskmd search "query" --filter priority=high` returns only high-priority matching tasks
- `taskmd search "query" --filter status=pending --filter priority=high` combines filters with AND logic
- `taskmd search "query" --sort priority` sorts search results by priority
- `taskmd search "query" --limit 5` limits output to 5 results
- All flags work together: `taskmd search "query" --filter status=pending --sort priority --limit 10`
- Flags work across all output formats (table, json, yaml)
- Existing search behavior is unchanged when no new flags are provided
