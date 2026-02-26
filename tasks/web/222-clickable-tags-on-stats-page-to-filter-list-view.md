---
title: "Clickable tags on Stats page to filter List view"
id: "222"
status: completed
priority: medium
type: feature
tags: ["ui", "navigation"]
created: "2026-02-25"
---

# Clickable tags on Stats page to filter List view

## Objective

Make tags on the Stats page clickable so that clicking a tag navigates the user to the Tasks (List) page with that tag already applied as a filter. The List page already supports tag filtering via `?tag=<name>` URL parameters, so this is primarily a navigation change on the Stats page.

## Tasks

- [x] In `StatsView.tsx`, wrap each tag name in the "Tags" section with a clickable element (link or button)
- [x] On click, navigate to `/tasks?tag=<tagName>` using React Router's `useNavigate` or a `<Link>` component
- [x] Style the tag to indicate it is clickable (cursor pointer, hover effect)

## Acceptance Criteria

- Clicking a tag on the Stats page navigates to the Tasks page with that tag pre-selected in the filter
- The URL updates to include `?tag=<tagName>` so the filtered view is bookmarkable
- Tags have a visible hover state indicating they are interactive
