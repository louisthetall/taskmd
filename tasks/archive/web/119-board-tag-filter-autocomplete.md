---
id: "119"
title: "Replace board tag filter pills with autocomplete search"
status: completed
priority: medium
effort: medium
tags: [board, filtering, ux]
created: 2026-02-15
---

# Replace board tag filter pills with autocomplete search

## Objective

Simplify the board page tag filtering UI. Projects can have a large number of tags, and displaying them all as pill buttons creates visual clutter and makes the filter bar unwieldy. Replace the pill list with a compact autocomplete search input that lets users find and add tags one at a time.

## Tasks

- [x] Remove the tag pill list from the board filter bar
- [x] Add an autocomplete/combobox search input for tags
  - Should search/filter available tags as the user types
  - Show matching tags in a dropdown
  - Allow selecting a tag to add it as an active filter
- [x] Display currently active tag filters as removable chips below or beside the search input
- [x] Ensure keyboard navigation works (arrow keys, enter to select, escape to close)
- [x] Handle edge cases: no matching tags, all tags already selected, empty tag list

## Acceptance Criteria

- Tag filters on the board page use an autocomplete search input instead of a pill list
- Users can type to search and filter available tags
- Selected tags appear as removable chips
- Removing a chip deactivates that tag filter
- The UI remains clean and compact even with many tags (50+)
- Existing filter behavior (status, priority) is unchanged
