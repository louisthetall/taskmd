---
title: "Create TypeScript library for taskmd"
id: "228"
status: pending
priority: low
effort: large
type: feature
tags: ["library", "typescript", "api"]
dependencies:
  - "227"
  - "230"
created: "2026-02-28"
---

# Create TypeScript library for taskmd

## Objective

Create a native TypeScript library that implements the taskmd specification, enabling Node.js/Deno/Bun users to parse, scan, filter, validate, and query taskmd files programmatically. The library should be idiomatic TypeScript (not a Go wrapper) and published to npm.

## Tasks

- [ ] Set up a new package (e.g. `apps/lib-ts` or a separate repo) with TypeScript tooling (tsconfig, build, test)
- [ ] Implement task file parser (YAML frontmatter + markdown body)
- [ ] Implement directory scanner (recursive scan with group inference)
- [ ] Implement task filtering (by status, priority, tags, etc.)
- [ ] Implement validator (required fields, enum values, duplicate IDs, circular dependencies)
- [ ] Implement dependency graph construction
- [ ] Implement next-task ranking logic
- [ ] Add comprehensive tests — reuse conformance test cases from the Go library (task 227)
- [ ] Write API documentation and README with usage examples
- [ ] Publish to npm

## Acceptance Criteria

- Library is published on npm and installable via `npm install`
- Passes the shared conformance test suite from task 227
- Supports parsing, scanning, filtering, validation, graph, and next-task ranking
- Has full TypeScript type definitions
- Works in Node.js, Deno, and Bun runtimes
- README includes quick-start examples for common operations
