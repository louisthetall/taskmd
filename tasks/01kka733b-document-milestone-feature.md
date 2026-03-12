---
title: "Document milestone feature"
id: "01kka733b"
status: completed
priority: medium
type: docs
dependencies: ["01kka72zy", "01kka730t", "01kka731b"]
tags: ["milestone", "docs"]
created: "2026-03-09"
---

# Document milestone feature

## Objective

Add documentation for the milestone feature across the docs site, including the field reference, configuration guide, and CLI command updates.

## Tasks

- [x] Add `milestone` to the frontmatter field reference on the docs site
- [x] Document `milestones` configuration in `.taskmd.yaml` reference
- [x] Update CLI command reference for `list --milestone`, `set --milestone`, `add --milestone`
- [x] Update `next` command docs to mention milestone-aware ranking
- [x] Update `board` and `stats` docs for `--group-by milestone`
- [x] Add a "Milestones" section to the user guide / best practices with usage examples
- [x] Update the spec (`docs/taskmd_specification.md`) examples to include milestone in the "Full Task" example

## Acceptance Criteria

- Milestone appears in the field reference with type, description, and example
- `.taskmd.yaml` milestones config is documented with all fields (name, description, due)
- All CLI flags related to milestone are documented
- A usage guide section explains when and how to use milestones
- `docs/taskmd_specification.md` includes milestone in examples
