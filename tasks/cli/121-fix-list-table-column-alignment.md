---
id: "121"
title: "Fix list command table column alignment"
status: pending
priority: low
effort: small
tags: [cli, output, ux]
created: 2026-02-16
---

# Fix list command table column alignment

## Objective

Fix the misaligned separator line in the `list` command's table output. Currently, `outputTable()` in `list.go` uses hardcoded 10-dash separators (`----------`) for all columns regardless of actual content width. This causes visual misalignment between headers, separators, and data rows — especially visible in the "title" column where content is much wider than 10 characters.

## Tasks

- [ ] Update `outputTable()` in `internal/cli/list.go` to generate dynamic-width separator dashes that match the column header lengths (at minimum)
- [ ] Verify alignment looks correct for all default columns (`id`, `title`, `status`, `priority`, `file`)
- [ ] Add a test case validating separator alignment matches header widths
- [ ] Check that colored output (ANSI codes) doesn't break alignment

## Acceptance Criteria

- Separator dashes align visually with column headers and content
- All existing list command tests continue to pass
- Table output remains readable with and without color enabled
