# taskmd Conformance Test Suite

This directory contains fixture files and expected outputs for verifying that taskmd library implementations conform to the [Operations Specification](../../docs/taskmd_operations.md).

## Directory Structure

```
tests/conformance/
  fixtures/
    tasks/              Task files used as input for all tests
      001-setup-project.md
      002-add-authentication.md
      003-create-api-endpoints.md
      004-design-database-schema.md
      005-write-unit-tests.md
      006-setup-ci-pipeline.md
      007-add-logging.md
      008-frontend-dashboard.md
      009-api-rate-limiting.md
      010-write-documentation.md
      011-parent-task.md
      012-child-of-epic.md
    config/             Configuration file variations
      default.yaml
      with-ignore.yaml
  expected/
    scan/               Expected scan results
      default.json
    filter/             Expected filter results
      status-pending.json
      tag-api.json
      no-matches.json
    validate/           Expected validation results
      valid.json
    next/               Expected next-task ranking results
      default.json
    search/             Expected search results
      query-api.json
    graph/              Expected graph output
      default.json
```

## How to Use

### For implementers

1. Load all task files from `fixtures/tasks/` using your scanner.
2. Run the operation being tested (scan, filter, validate, next, search, graph).
3. Compare your output against the corresponding file in `expected/`.

### Comparison rules

- **Task ordering**: Tasks in scan results are ordered by file path (lexicographic). Next-task results are ordered by score descending, then ID ascending.
- **Null vs empty arrays**: `null` and `[]` for dependency lists are treated as equivalent (no dependencies).
- **File paths**: Expected `file_path` values are relative to the scan root. Your implementation may return absolute paths — normalize before comparison.
- **Timestamps**: Created dates are serialized as RFC 3339 (`2026-01-01T00:00:00Z`). Implementations may use date-only format.

### Running with the Go CLI

```bash
cd tests/conformance/fixtures

# Scan
taskmd list --dir tasks --format json

# Filter
taskmd list --dir tasks --format json --filter "status=pending"
taskmd list --dir tasks --format json --filter "tag=api"
taskmd list --dir tasks --format json --filter "status=nonexistent"

# Validate
taskmd validate --dir tasks --format json

# Next
taskmd next --dir tasks --format json --limit 10

# Search
taskmd search --dir tasks --format json "API"

# Graph
taskmd graph --dir tasks --format json
```

## Fixture Design

The fixture set covers:

| Scenario | Tasks |
|----------|-------|
| Completed task (dependency satisfied) | 001 |
| Critical priority with dependencies | 002 |
| Multiple dependencies (001 + 002) | 003 |
| In-progress status | 004 |
| Blocked by multiple tasks | 005 |
| Low priority, quick win | 006 |
| Cancelled status | 007 |
| Deep dependency chain (001→002→003→008) | 008 |
| Quick win candidate behind blocker | 009 |
| Low priority behind blocker | 010 |
| Parent task (epic) | 011 |
| Child task with parent reference | 012 |

## Adding New Test Cases

1. Add fixture files to `fixtures/tasks/`.
2. Run the CLI against the fixtures to generate actual output.
3. Verify the output matches expected behavior per the operations spec.
4. Save the output as a new expected file.
5. Document the scenario in this README.
