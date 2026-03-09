---
title: "Support full keyboard navigation on taskmd web"
id: "01kka70cf"
status: pending
priority: high
type: feature
tags: ["accessibility", "ux"]
created: "2026-03-09"
---

# Support full keyboard navigation on taskmd web

## Objective

Enable complete keyboard-driven navigation throughout the taskmd web interface. Users should be able to browse, select, and interact with all UI elements using only Tab, Arrow keys, and Enter — without ever needing a mouse. This improves accessibility (WCAG 2.1 compliance) and power-user efficiency.

## Tasks

- [ ] Audit all interactive elements for proper focus management and tab order
- [ ] Add visible focus indicators (outline/ring styles) to all focusable elements
- [ ] Implement arrow-key navigation within task lists (up/down to move between tasks)
- [ ] Implement Enter key to open/select the currently focused task
- [ ] Ensure Tab moves focus logically between major UI sections (sidebar, task list, detail panel)
- [ ] Add Escape key support to close modals, panels, or return to the previous focus context
- [ ] Add skip-navigation links for quickly jumping between sections
- [ ] Handle focus trapping in modals and dialogs
- [ ] Test keyboard navigation across all major views (board, list, detail, graph)
- [ ] Ensure no keyboard traps exist (user can always Tab out of any component)

## Acceptance Criteria

- All interactive elements (buttons, links, inputs, task rows) are reachable via Tab key
- Arrow keys navigate between items within lists and grids
- Enter activates/opens the currently focused element
- Escape closes open panels, modals, or menus
- A visible focus indicator is shown on the currently focused element at all times
- The entire app can be operated without a mouse
- No keyboard traps — focus can always be moved away from any component
- Focus order follows a logical reading/layout sequence
