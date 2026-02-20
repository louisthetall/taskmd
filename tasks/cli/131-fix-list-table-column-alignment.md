---
id: "131"
title: "Fix table column alignment with color across all commands"
status: completed
priority: low
effort: small
tags: [cli, output, ux]
created: 2026-02-16
---

# Fix table column alignment with color across all commands

## Objective

Fix misaligned column headers/data in table output when ANSI color codes are present. The root cause is that `text/tabwriter` counts ANSI escape bytes as visible characters, inflating perceived column widths and causing misalignment.

This was previously fixed for the `list` command by replacing `tabwriter` with manual column padding (see worklog). However, the same bug still exists in all other commands that use `tabwriter` with colored output:

- `todos.go` — `tabwriter` at line ~199
- `tracks.go` — `tabwriter` at lines ~139, ~155
- `tags.go` — `tabwriter` at line ~132
- `search.go` — `tabwriter` at line ~79
- `stats.go` — `tabwriter` at line ~100
- `next.go` — `tabwriter` at line ~132

The bug does **not** occur with `--no-color` because without ANSI codes, `tabwriter` measures widths correctly.

### Root Cause

`tabwriter` pads columns based on byte length. ANSI color codes (e.g. `\033[32m...\033[0m`) add ~10+ invisible bytes per colored value. This inflates the column width `tabwriter` calculates, pushing subsequent columns further right than the uncolored headers expect.

### Previous Fix (list.go only)

The `list` command was fixed by replacing `tabwriter` with manual column padding: compute max visible width per column (stripping ANSI codes), then pad each cell explicitly with spaces. This approach should be extracted into a shared utility and applied to all affected commands.

## Tasks

- [x] Update `outputTable()` in `internal/cli/list.go` to use manual padding instead of `tabwriter`
- [x] Add alignment tests for `list` command
- [ ] Extract the manual-padding table approach from `list.go` into a shared utility (e.g. `internal/cli/tablewriter.go`)
- [ ] Replace `tabwriter` usage in `todos.go` with the shared utility
- [ ] Replace `tabwriter` usage in `tracks.go` with the shared utility
- [ ] Replace `tabwriter` usage in `tags.go` with the shared utility
- [ ] Replace `tabwriter` usage in `search.go` with the shared utility
- [ ] Replace `tabwriter` usage in `stats.go` with the shared utility
- [ ] Replace `tabwriter` usage in `next.go` with the shared utility
- [ ] Add tests verifying color alignment matches plain alignment across commands
- [ ] All existing tests pass

## Acceptance Criteria

- Table columns align correctly with and without color enabled across **all** commands that output tables
- ANSI color codes do not affect column width calculations
- A shared table-writing utility is used consistently (no more direct `tabwriter` for colored output)
- All existing tests continue to pass
