---
id: "176"
title: "E2e tests for error handling and edge cases"
status: pending
priority: medium
effort: medium
type: improvement
tags:
  - testing
  - cli
parent: "173"
dependencies:
  - "174"
created: 2026-02-20
---

# E2e tests for error handling and edge cases

## Objective

Test that the CLI handles error conditions gracefully — returning non-zero exit codes, printing useful stderr messages, and not crashing on malformed input.

## Tasks

- [ ] Test unknown command: `taskmd nonexistent` exits non-zero with helpful message
- [ ] Test missing required args: e.g. `taskmd set` with no ID
- [ ] Test invalid flag values: e.g. `taskmd list --status bogus`
- [ ] Test validate on malformed task files: missing frontmatter, invalid YAML, missing required fields
- [ ] Test operations on empty directory: list, next, graph with no task files
- [ ] Test stdin/pipe behavior: pipe content to `taskmd validate --stdin`
- [ ] Test invalid `--task-dir` path: non-existent directory
- [ ] Verify exit codes: 0 for success, non-zero for all error cases
- [ ] Verify stderr contains actionable error messages (not stack traces)

## Acceptance Criteria

- Every error scenario returns a non-zero exit code
- stderr output is user-friendly and describes what went wrong
- No panics or stack traces in error output
- stdin validation works via pipe
- Empty directory cases are handled gracefully (not treated as errors where appropriate)
