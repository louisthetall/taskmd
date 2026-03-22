---
title: "Add project-aware API endpoints to web server"
id: "01kma4hkf"
status: pending
priority: medium
type: feature
tags: ["global-registry", "web", "api"]
dependencies: ["01kma460m"]
created: "2026-03-22"
---

# Add project-aware API endpoints to web server

## Objective

Extend the web server so the frontend can discover registered projects and switch between them. Add a `GET /api/projects` endpoint and a `?project=<id>` query parameter on all existing endpoints that changes which project's task directory is scanned.

## Tasks

- [ ] Add `GET /api/projects` endpoint that calls `LoadGlobalRegistry()` and returns the list of registered projects (id, name, path)
- [ ] Accept an optional `?project=<id>` query parameter on all task-related endpoints (`/api/tasks`, `/api/board`, `/api/graph`, `/api/stats`, `/api/next`, `/api/tracks`, `/api/validate`, `/api/search`)
- [ ] When `?project=<id>` is set, resolve the project's path from the registry, load its `.taskmd.yaml`, and create a scanner rooted at that project's task directory
- [ ] When `?project=` is omitted, use the current scan directory (existing behavior, backwards-compatible)
- [ ] Include the project's phases in `/api/config?project=<id>` response so the frontend can show project-specific phases
- [ ] Return appropriate errors: 404 if project id not found, 500 if project path is unreachable
- [ ] Add handler tests with mock registry entries pointing to temp project directories

## Acceptance Criteria

- `GET /api/projects` returns the global registry entries
- `GET /api/tasks?project=foo` returns tasks from the `foo` project's task directory
- All existing endpoints accept the `?project=` parameter and scan the correct directory
- Omitting `?project=` returns the same results as today (no regression)
- `/api/config?project=foo` returns that project's phases
- Unknown project id returns a 404 with a clear error message
