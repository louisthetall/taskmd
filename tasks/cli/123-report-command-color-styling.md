---
id: "123"
title: "Add color styling to the report command"
status: pending
priority: medium
effort: medium
tags:
  - cli
  - ux
created: 2026-02-16
---

# Add Color Styling to the Report Command

## Objective

Add color styling to the `report` command's markdown (`md`) terminal output to match the visual polish of other CLI commands (`list`, `board`, `graph`). The existing `colors.go` helpers (`formatStatus`, `formatPriority`, `formatTaskID`, `formatHeading`, etc.) should be used for consistency.

## Context

The `report` command currently outputs plain unformatted text via `fmt.Fprintf` in `report_md.go`. Other commands already use lipgloss-based color helpers from `colors.go` to colorize statuses, priorities, task IDs, headings, and effort levels. The report command should follow the same pattern.

Color should only apply to the **terminal (md) format** — HTML and JSON output should remain unchanged since they have their own styling mechanisms.

## Tasks

- [ ] Update `outputReportMarkdown` to accept/use a lipgloss renderer from `getRenderer()`
- [ ] Colorize status labels in the "By Status" breakdown (`writeMarkdownStatusBreakdown`)
- [ ] Colorize priority labels in the "By Priority" breakdown (`writeMarkdownPriorityBreakdown`)
- [ ] Colorize group headings using `formatHeading` in `writeMarkdownGroups`
- [ ] Colorize task IDs using `formatTaskID` in task listings
- [ ] Colorize status text in critical path entries (`writeMarkdownCriticalPath`)
- [ ] Colorize blocked task details in `writeMarkdownBlockedTasks`
- [ ] Use `formatLabel` for section headings ("Summary", "Critical Path", etc.)
- [ ] Skip color when output is directed to a file (`--out` flag) — ensure `colorsEnabled()` handles this correctly
- [ ] Add tests for colored vs plain output in `report_test.go`

## Acceptance Criteria

- Running `taskmd report tasks/` in a terminal shows colored output matching the style of `taskmd list` and `taskmd board`
- Status values are colored (green=completed, yellow=in-progress, red=blocked, gray=pending)
- Priority values are colored (red=critical, yellow=high, blue=medium, gray=low)
- Task IDs are cyan and bold
- Section headings are bold
- `--no-color` flag and `NO_COLOR` env var disable all colors
- Output to file (`--out report.md`) produces plain text without ANSI codes
- JSON and HTML formats are unaffected
- Existing report tests continue to pass
