---
id: "120"
title: "Add web interface screenshots to documentation"
status: pending
priority: medium
effort: medium
tags: [documentation, web]
created: 2026-02-15
---

# Add Web Interface Screenshots to Documentation

## Objective

Add screenshots of the web interface to `apps/docs/guide/web.md` (and optionally `getting-started/index.md`) so users can see what each view looks like before launching the app. The docs currently have zero images.

Screenshots should be placed in `apps/docs/public/images/web/` and referenced from the markdown files.

## Process

Work through the screenshots **one at a time**. For each screenshot:

1. Ask the user to take a screenshot of the specified view
2. Confirm the file path once provided
3. Add the image reference to the correct location in the docs
4. Move on to the next screenshot

## Tasks

- [ ] Create `apps/docs/public/images/web/` directory
- [ ] **Tasks View** — Filtered table showing search bar, status/priority filters, and sortable columns. Place under `## Views > ### Tasks View` in `guide/web.md`
- [ ] **Board View (by status)** — Kanban board with Pending/In-Progress/Completed columns and task cards. Place under `### Board View` in `guide/web.md`
- [ ] **Graph View** — Interactive dependency graph with colored nodes (yellow=pending, blue=in-progress, green=completed). Place under `### Graph View` in `guide/web.md`
- [ ] **Stats View** — Dashboard showing metric cards and breakdown sections (by status, priority, effort, tags). Place under `### Stats View` in `guide/web.md`
- [ ] **Task Detail Page** — Single task showing metadata (status badge, priority, effort, dependencies, tags) and rendered markdown body. Place under `### Tasks View` after the table screenshot in `guide/web.md`
- [ ] **Next View** — Recommendation cards with scores and reasons. Place in a new `### Next View` section or under an existing relevant heading in `guide/web.md`
- [ ] **Tracks View** — Parallel work tracks with task cards grouped by scope. Place in a new `### Tracks View` section in `guide/web.md`
- [ ] **Validate View** — Validation results with errors/warnings grouped by file. Place in a new `### Validate View` section in `guide/web.md`
- [ ] **Quick Start hero screenshot** — The Tasks view on first launch. Place at step 9 of `getting-started/index.md` where `taskmd web` is introduced

## Acceptance Criteria

- At least the 4 main view screenshots (Tasks, Board, Graph, Stats) are added
- Images are stored in `apps/docs/public/images/web/`
- Each image is referenced with an alt-text description in the markdown
- Screenshots are reasonably sized (aim for ~1200px wide, compressed PNG or WebP)
