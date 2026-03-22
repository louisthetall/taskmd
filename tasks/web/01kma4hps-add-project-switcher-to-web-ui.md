---
title: "Add project switcher to web UI"
id: "01kma4hps"
status: pending
priority: medium
type: feature
tags: ["global-registry", "web", "frontend"]
dependencies: ["01kma4hkf"]
created: "2026-03-22"
---

# Add project switcher to web UI

## Objective

Add a project selector dropdown to the Shell header that lets users switch between registered projects. When a project is selected, all views (tasks, board, graph, stats, etc.) filter down to that project's data. Mirrors how the existing phase selector works.

## Tasks

- [ ] Add `useProjects()` SWR hook that fetches `GET /api/projects`
- [ ] Add `useProject()` hook that manages the selected project via URL search param (`?project=<id>`), similar to `usePhase()`
- [ ] Add project selector dropdown in Shell header (before or alongside the phase selector)
- [ ] Pass the selected project id to all data-fetching hooks (`useTasks`, `useBoard`, `useGraph`, `useStats`, `useNext`, `useTracks`, `useValidate`, `useSearch`) as a `?project=` query parameter
- [ ] When a project is selected, re-fetch `/api/config?project=<id>` to update available phases for that project
- [ ] Show "(local)" or the current directory name when no project is selected (default/current behavior)
- [ ] Persist selected project in URL so it survives page refreshes and is shareable
- [ ] Handle edge cases: registry is empty (hide selector), selected project becomes unavailable (clear selection with warning)

## Acceptance Criteria

- A project dropdown appears in the Shell header when projects are registered
- Selecting a project reloads all views with that project's tasks
- The phase selector updates to show the selected project's phases
- The URL reflects the selected project (`?project=foo`)
- Refreshing the page preserves the project selection
- When no projects are registered, the selector is hidden (no UI change from today)
- Switching projects clears any active phase filter (phases differ per project)
