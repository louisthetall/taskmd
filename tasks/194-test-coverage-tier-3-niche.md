---
id: "194"
title: "Improve test coverage for Tier 3 (niche/low-risk) packages"
status: completed
priority: low
effort: medium
type: chore
tags: [testing, coverage]
dependencies: ["193"]
created: 2026-02-21
---

# Improve test coverage for Tier 3 (niche/low-risk) packages

## Objective

Add tests for lower-importance uncovered areas: niche commands, thin wrappers, and small utilities. These are low-risk individually but contribute to overall coverage and catch edge-case regressions.

## Background

Current coverage analysis (2026-02-21) identified these Tier 3 gaps:

| Area | Uncovered stmts | Total stmts | Coverage | Risk |
|------|-----------------|-------------|----------|------|
| `cli/web.go` + `web_export.go` | 33 | 80 | 58.8% | CLI wrappers for web server/export |
| `cli/completion.go` | 6 | 7 | 14.3% | Shell completion; delegates to Cobra |
| `cli/mcp.go` | 2 | 3 | 33.3% | MCP server launcher; thin wrapper |
| `slug` package | 7 | 7 | 0.0% | `Slugify` function; trivial but untested |
| `cli/project_init.go` | 37 | 166 | 77.7% | Interactive agent selection prompt |

## Tasks

- [ ] Add tests for `slug` package: `Slugify` with various inputs (spaces, special chars, unicode)
- [ ] Add tests for `cli/web.go`: test flag parsing and validation (not actual server startup)
- [ ] Add tests for `cli/web_export.go`: `runWebExport` with mock scanner
- [ ] Add tests for `cli/project_init.go`: non-interactive logic paths
- [ ] Add a smoke test for `cli/completion.go` if feasible
- [ ] Add a smoke test for `cli/mcp.go` if feasible

## Acceptance Criteria

- `slug` package reaches 100% coverage
- `cli/web.go` and `cli/web_export.go` testable paths are covered
- `cli/project_init.go` coverage reaches 85%+
- Overall CLI tool coverage reaches 82%+ (up from 79.1%)
- All tests pass: `cd apps/cli && make test`
- Linter passes: `cd apps/cli && make lint`
