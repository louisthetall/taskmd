---
id: "156"
title: "Extend todos list with rich structured output for agent consumption"
status: in-progress
priority: medium
effort: large
type: feature
dependencies:
  - "155"
tags:
  - cli
  - developer-experience
  - agent-tooling
created: 2026-02-19
---

# Extend todos list with Rich Structured Output for Agent Consumption

## Objective

Extend the `taskmd todos list` command so every TODO is a structured object that agents can consume directly without re-parsing. The current default output has only `file`, `line`, `marker`, and `text`. This task adds cheap, always-available fields to the **default** output, and introduces a `--rich` flag for expensive derived fields that require git or AST-like analysis.

### Field breakdown

**Default output** (always included — cheap to compute):

| Field | Description |
|-------|-------------|
| `id` | Stable fingerprint hash (file + line + marker + text) |
| `file` | Relative file path (already exists) |
| `line` | Line number (already exists) |
| `column` | Character offset of the marker within the line |
| `language` | Language name derived from extension (e.g. "go", "python") |
| `tag` | Marker type: TODO, FIXME, HACK, etc. (rename of existing `marker` field) |
| `text` | Cleaned comment text (already exists) |

**Optional flags:**

| Flag | Field | Description |
|------|-------|-------------|
| `--raw-text` | `raw_text` | Include unprocessed original comment line(s) as they appear in source |

**Rich output** (`--rich` flag — requires git calls / code analysis):

| Field | Description |
|-------|-------------|
| `scope` | Enclosing function/class/component name |
| `blame.author` | Author from `git blame` |
| `blame.commit` | Commit hash from `git blame` |
| `blame.date` | Commit date from `git blame` |
| `age` | Days since the blame commit date |

## Tasks

### Default fields (always included)

- [ ] Generate a stable `id` fingerprint for each TODO (e.g. SHA-256 of file + marker + text, truncated)
- [ ] Add `column` field (character offset of the marker within the line)
- [ ] Add `language` field (derived from file extension via the language registry)
- [ ] Add `--raw-text` flag that includes `raw_text` field (unprocessed original comment lines)
- [ ] Rename `marker` JSON/YAML key to `tag` (keep backward compat alias or just rename)
- [ ] Update default table format to include `id` (short form) and `language` columns
- [ ] Write tests for fingerprint generation (stability across runs, uniqueness)

### Rich fields (`--rich` flag)

- [ ] Add `--rich` flag to `todos list`
- [ ] Implement scope detection: extract enclosing function, class, or component name
  - [ ] Go: detect enclosing `func` declaration
  - [ ] JavaScript/TypeScript: detect enclosing function, method, or component
  - [ ] Python: detect enclosing `def` or `class`
  - [ ] Fallback to empty string for unsupported languages
- [ ] Implement git blame integration for each TODO location
  - [ ] Extract author, commit hash, and date via `git blame -L <line>,<line> <file>`
  - [ ] Calculate `age` in days from blame date to current date
  - [ ] Handle non-git repos gracefully (omit blame/age fields)
- [ ] Update table format under `--rich` to show scope and age columns
- [ ] Write tests for scope detection across Go, JS/TS, and Python
- [ ] Write tests for blame info extraction (with and without git)
- [ ] Write tests for rich JSON/YAML output round-trip

## Acceptance Criteria

- `taskmd todos list --format json` default output includes: id, file, line, column, language, tag, text
- `taskmd todos list --raw-text --format json` additionally includes raw_text
- `taskmd todos list --rich --format json` adds: scope, blame (author, date, commit), age
- The `id` fingerprint is stable across runs for the same TODO
- Scope detection correctly identifies the enclosing function/class for Go, JS/TS, and Python
- Git blame info is populated in git repos, gracefully omitted otherwise
- `age` is calculated correctly as days since the blame commit date
- Table output without `--rich` is not cluttered (shows id short hash, file, line, tag, text)
- Tests cover fingerprint stability, scope detection, blame parsing, and edge cases
