---
title: "Create Python library for taskmd"
id: "229"
status: pending
priority: low
effort: large
type: feature
tags: ["library", "python", "api"]
dependencies:
  - "227"
  - "230"
created: "2026-02-28"
---

# Create Python library for taskmd

## Objective

Create a native Python library that implements the taskmd specification, enabling Python users to parse, scan, filter, validate, and query taskmd files programmatically. The library should be idiomatic Python (not a Go wrapper) and published to PyPI.

## Tasks

- [ ] Set up a new package (e.g. `apps/lib-py` or a separate repo) with Python tooling (pyproject.toml, pytest, mypy)
- [ ] Implement task file parser (YAML frontmatter + markdown body)
- [ ] Implement directory scanner (recursive scan with group inference)
- [ ] Implement task filtering (by status, priority, tags, etc.)
- [ ] Implement validator (required fields, enum values, duplicate IDs, circular dependencies)
- [ ] Implement dependency graph construction
- [ ] Implement next-task ranking logic
- [ ] Add comprehensive tests — reuse conformance test cases from the Go library (task 227)
- [ ] Write API documentation and README with usage examples
- [ ] Publish to PyPI

## Acceptance Criteria

- Library is published on PyPI and installable via `pip install`
- Passes the shared conformance test suite from task 227
- Supports parsing, scanning, filtering, validation, graph, and next-task ranking
- Has full type annotations (mypy-compatible)
- Supports Python 3.10+
- README includes quick-start examples for common operations
