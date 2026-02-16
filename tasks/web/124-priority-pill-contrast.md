---
id: "124"
title: "Improve contrast of medium and low priority pills"
status: pending
priority: medium
effort: small
tags:
  - web
  - ux
  - accessibility
created: 2026-02-16
---

# Improve Contrast of Medium and Low Priority Pills

## Objective

Increase the visual contrast of the "medium" and "low" priority badge/pill colors so they are clearly distinguishable from the disabled (inactive) filter pill state. Currently both use gray tones that are too similar to the inactive state (`bg-gray-50 text-gray-400`), making it hard to tell active from inactive at a glance.

## Context

In `apps/web/src/components/tasks/TaskTable/constants.ts`, the current priority colors are:

- **medium**: `bg-gray-100 text-gray-600 ring-gray-300` — very close to disabled
- **low**: `bg-gray-50 text-gray-400 ring-gray-200` — nearly identical to disabled

The disabled filter pill state uses: `bg-gray-50 text-gray-400`

These need distinct, non-gray colors so users can immediately see which priority a task has, and whether a filter pill is active or not.

## Tasks

- [ ] Choose new color palettes for medium and low priorities that are visually distinct from gray disabled state (e.g., blue/indigo for medium, slate/zinc-blue for low)
- [ ] Update `PRIORITY_COLORS.medium` in `apps/web/src/components/tasks/TaskTable/constants.ts`
- [ ] Update `PRIORITY_COLORS.low` in `apps/web/src/components/tasks/TaskTable/constants.ts`
- [ ] Verify contrast in both light and dark modes
- [ ] Verify the updated pills look distinct across all usages: TaskTable badges, FilterBar pills, BoardFilterBar, TrackCard, and RecommendationCard

## Acceptance Criteria

- Medium and low priority pills are clearly distinguishable from the disabled/inactive filter pill state
- Colors are visually distinct from each other (medium vs low)
- Both light mode and dark mode maintain good contrast and readability
- The overall color hierarchy still communicates priority ranking (critical > high > medium > low)
- No regressions in other components that consume `PRIORITY_COLORS`
