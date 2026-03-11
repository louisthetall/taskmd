---
title: "Add milestone field to spec and parser"
id: "01kka72zy"
status: in-progress
priority: high
type: feature
tags: ["milestone", "spec", "parser"]
touches: ["cli/scanner"]
created: "2026-03-09"
---

# Add milestone field to spec and parser

## Objective

Add `milestone` as an optional single-string frontmatter field to the taskmd specification and Go parser. A milestone represents a time-based grouping (sprint, phase, iteration, release) that a task belongs to. The value is free-form (e.g., `"v0.2"`, `"2026-Q1"`, `"beta-launch"`).

Optionally, milestones can be defined in `.taskmd.yaml` with metadata (description, due date, ordering):

```yaml
milestones:
  - name: "v0.2"
    description: "Core CLI features"
    due: 2026-04-01
  - name: "v0.3"
    description: "Web dashboard"
    due: 2026-06-01
```

## Tasks

- [x] Add `milestone` field to the spec Field Summary table and Optional Fields section
- [x] Add `milestone` to the task model struct in the Go parser
- [x] Parse `milestone` from frontmatter (string type)
- [x] Add `milestones` config section to `.taskmd.yaml` schema
- [x] Parse milestones config (name, description, due date)
- [x] Add validation: warn if task references a milestone not defined in config (soft warning, not error)
- [x] Add unit tests for milestone parsing
- [x] Add unit tests for milestones config parsing
- [x] Add validation tests for undefined milestone warnings
- [x] Run `make sync-spec` to sync spec copies

## Acceptance Criteria

- `milestone` appears in the spec as an optional string field
- The parser reads `milestone` from frontmatter and populates the task model
- `.taskmd.yaml` supports a `milestones` list with `name`, `description`, and `due` fields
- Validation warns (does not error) when a task references an undefined milestone
- Tasks without a `milestone` field are unaffected
- All new code has tests; `make test` and `make lint` pass
