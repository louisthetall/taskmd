---
id: "01kjn8050"
title: "Fix scope wildcard matching to cross path separators"
status: completed
priority: medium
effort: small
dependencies: []
tags: ["cli", "bugfix"]
created: 2026-03-01
---

# Fix scope wildcard matching to cross path separators

## Objective

The `--scope` flag wildcard matching used `filepath.Match`, where `*` does not cross path separators (`/`). This meant `cli*` would not match scopes like `cli/import` or `cli/tracks`, even though scopes are logical identifiers, not file paths. Replace with a simple `*`-only matcher that crosses `/`.

## Tasks

- [x] Replace `filepath.Match` in `MatchScope` with a custom `matchStar` function where `*` matches any character including `/`
- [x] Remove support for `?` and `[]` wildcards (not needed for scope matching)
- [x] Update `filter.go` to use `strings.Contains(value, "*")` instead of removed `containsWildcard` helper
- [x] Add tests for cross-separator matching (`cli*` → `cli/import`)
- [x] Add tests for mixed slash+wildcard patterns (`cli/imp*` → `cli/import`)
- [x] Add tests confirming `?` and `[]` are treated as literals
- [x] Verify all commands with `--scope` (next, list, graph, status, feed, tracks) use `filter.MatchScope`

## Acceptance Criteria

- `taskmd next --scope "cli*"` matches tasks with scopes like `cli/import`, `cli/tracks`
- Only `*` is treated as a wildcard; `?` and `[]` are literal characters
- All existing tests pass
- Fix applies uniformly to all commands that support `--scope`
