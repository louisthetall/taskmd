# taskmd Specification

> Trimmed for agent use. Full spec: `taskmd spec --stdout`

## Quick Reference

Each task is a `.md` file with YAML frontmatter and a markdown body.

```yaml
---
id: "01h5a3mpk"
title: "Task title"
status: pending
---

# Task Title

Description and subtasks go here.
```

### Field Summary

| Field | Type | Required | Values / Format |
|-------|------|----------|-----------------|
| `id` | string | **Yes** | Unique identifier (e.g., `"001"`, `"42"`, `"cli-049"`) |
| `title` | string | **Yes** | Brief, descriptive text |
| `status` | enum | Recommended | `pending`, `in-progress`, `completed`, `in-review`, `blocked`, `cancelled` |
| `priority` | enum | No | `low`, `medium`, `high`, `critical` |
| `effort` | enum | No | `small`, `medium`, `large` |
| `type` | enum | No | `feature`, `bug`, `improvement`, `chore`, `docs` |
| `dependencies` | array | No | List of task ID strings (e.g., `["001", "015"]`) |
| `tags` | array | No | Lowercase, hyphen-separated strings |
| `group` | string | No | Logical grouping (derived from directory if omitted) |
| `owner` | string | No | Free-form assignee name or identifier |
| `phase` | string | No | Free-form phase name (e.g., `"v0.2"`) |
| `touches` | array | No | Abstract scope identifiers (e.g., `["cli/graph"]`) |
| `context` | array | No | Explicit file paths relevant to the task |
| `parent` | string | No | Single task ID (e.g., `"045"`) |
| `created` | date | No | `YYYY-MM-DD` |
| `verify` | array | No | List of typed verification checks |
| `pr` | array | No | List of pull request URLs |

## Frontmatter Schema

### Required Fields

**`id`** -- Unique identifier. Any non-empty string. Must be unique across all tasks.

**`title`** -- Brief, action-oriented description.

### Optional Fields

**`status`** -- Current state (recommended for all tasks):

```
pending -> in-progress -> in-review -> completed
   |            |              |
   v            v              v
blocked         +-> cancelled <+
```

| Status | Meaning |
|--------|---------|
| `pending` | Not started (initial state) |
| `in-progress` | Currently being worked on |
| `in-review` | Submitted for review (PR open) |
| `completed` | Finished and verified |
| `blocked` | Cannot proceed due to a blocker |
| `cancelled` | Will not be completed |

**`priority`** -- `low`, `medium` (default), `high`, `critical`

**`effort`** -- `small` (< 2h), `medium` (2-8h), `large` (> 8h)

**`type`** -- `feature`, `bug`, `improvement`, `chore`, `docs`

**`dependencies`** -- Task IDs that must complete before this task can start:

```yaml
dependencies: ["001", "015"]
```

**`tags`** -- Lowercase, hyphen-separated labels:

```yaml
tags: [core, api]
```

**`group`** -- Logical grouping. Derived from parent directory if omitted.

**`owner`** -- Free-form assignee string.

**`phase`** -- Sprint, phase, or release identifier.

**`touches`** -- Scope identifiers declaring which code areas a task modifies. Scopes can be mapped to paths in `.taskmd.yaml`. Two tasks sharing a scope risk merge conflicts.

```yaml
touches: [cli/graph, cli/output]
```

**`context`** -- Explicit file paths (relative to project root) relevant to the task:

```yaml
context: ["docs/api-design.md"]
```

**`parent`** -- Task ID of parent task for hierarchical grouping. Organizational only (no blocking or status cascading).

**`created`** -- Date in `YYYY-MM-DD` format.

**`verify`** -- Acceptance checks run with `taskmd verify <ID>`:

| Type | Fields | Behavior |
|------|--------|----------|
| `bash` | `run` (required), `dir` (optional) | Runs shell command; pass if exit code 0 |
| `assert` | `check` (required) | Text for agent evaluation (not executed) |

```yaml
verify:
  - type: bash
    run: "go test ./..."
    dir: "apps/cli"
  - type: assert
    check: "API returns paginated results"
```

**`pr`** -- Pull request URLs. Managed via `taskmd set --add-pr <url>`.

## File Organization

### File Naming

```
ID-descriptive-title.md
```

Examples: `001-project-scaffolding.md`, `042-implement-user-auth.md`

### Directory Structure

```
tasks/
  001-spec.md                   # No group
  cli/                          # Group: "cli"
    015-go-cli-scaffolding.md
    016-task-model-parsing.md
  web/                          # Group: "web"
    001-project-scaffolding.md
```

Group resolution: explicit `group` in frontmatter > parent directory name > no group.

## Validation

A valid task **must**:

1. Have YAML frontmatter enclosed in `---` delimiters
2. Include required fields: `id`, `title`
3. Use valid enum values for `status`, `priority`, `effort`, `type`
4. Have unique IDs across the project
5. Reference only existing tasks in `dependencies`
6. Have no circular dependency chains
7. Reference an existing task in `parent` (if set), with no self-reference or cycles

## Examples

### Minimal Task

```markdown
---
id: "01h5a3mpk"
title: "Fix login button alignment"
status: pending
---

# Fix Login Button Alignment

The login button on the homepage is misaligned. Update the CSS to center it.
```

### Full Task

```markdown
---
id: "015"
title: "Implement user authentication"
status: in-progress
priority: high
effort: large
type: feature
phase: "v1.0"
dependencies: ["012", "013"]
parent: "012"
tags: [auth, security, api]
created: 2026-02-08
verify:
  - type: bash
    run: "go test ./internal/auth/..."
  - type: assert
    check: "JWT tokens expire after 24 hours"
---

# Implement User Authentication

## Objective

Add JWT-based authentication to the API.

## Tasks

- [x] Design authentication flow
- [ ] Create login endpoint
- [ ] Add authentication middleware
- [ ] Write integration tests

## Acceptance Criteria

- Users can log in with email and password
- Protected routes require valid JWT
- All endpoints have > 90% test coverage
```


## ID Generation

This project uses **ULID** (Universally Unique Lexicographically Sortable Identifier) IDs.

- Format: lowercase ULID strings (e.g., `01h5a3mpk`)
- Generated automatically by `taskmd add`
- ULIDs are time-ordered, so tasks sort chronologically by creation
- Default length: 9 characters
