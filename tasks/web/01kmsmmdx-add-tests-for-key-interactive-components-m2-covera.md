---
title: "Add tests for key interactive components (M2 coverage target)"
id: "01kmsmmdx"
status: completed
priority: low
type: chore
tags: ["testing", "quality"]
created: "2026-03-28"
phase: web-ui
depends-on: ["01kmsmma7"]
---

# Add tests for key interactive components (M2 coverage target)

## Objective

Add tests for the highest-value interactive components to reach the M2 milestone (75% statements, 85% branches, 80% functions). These components have meaningful UI logic and user interaction flows.

## Tasks

- [x] Add tests for `Shell.tsx` (0% → target 70%+) — routing, layout, global state
- [x] Add tests for `SearchDialog.tsx` (0% → target 70%+) — keyboard shortcuts, search logic, 202 lines
- [x] Improve tests for `BoardFilterBar.tsx` (36% → target 80%+) — filter interactions
- [x] Improve tests for `BoardView.tsx` (45% → target 80%+) — board orchestration
- [x] Add tests for `KeyboardList.tsx` (45% → target 80%+) — keyboard navigation
- [x] Add tests for `WorklogSection.tsx` (14% → target 80%+) — worklog rendering
- [x] Verify overall coverage reaches M2 thresholds

## Acceptance Criteria

- All listed components have tests covering primary user interactions
- `Shell.tsx` tests verify routing and layout rendering
- `SearchDialog.tsx` tests cover keyboard shortcut activation and search filtering
- Overall statement coverage is at or above 75%
- Overall function coverage is at or above 80%
