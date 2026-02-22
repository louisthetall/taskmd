---
id: "193"
title: "Improve test coverage for Tier 2 (user-facing) packages"
status: completed
priority: medium
effort: large
type: chore
tags: [testing, coverage]
dependencies: ["192"]
created: 2026-02-21
---

# Improve test coverage for Tier 2 (user-facing) packages

## Objective

Add tests for medium-importance uncovered areas: user-facing commands and web handlers that are less critical than Tier 1 but still affect real users. Some areas (import wizard, web server) may need refactoring to become testable.

## Background

Current coverage analysis (2026-02-21) identified these Tier 2 gaps:

| Area | Uncovered stmts | Total stmts | Coverage | Risk |
|------|-----------------|-------------|----------|------|
| `cli/tags.go` | 22 | 49 | 55.1% | User-facing tag listing command |
| `web/server.go` | 86 | 86 | 0.0% | Web server startup, CORS, static file mounting |
| `web/handlers.go` | 74 | 208 | 64.4% | File update error handling + worklog handler at 0% |
| `web/export.go` | 48 | 142 | 66.2% | Static HTML export; failures produce broken exports |
| `cli/importcmd.go` (wizard) | 110 | 215 | 48.8% | Interactive import wizard; 9 functions at 0% |
| `cli/get.go` | 35 | 255 | 86.3% | Core `get` command with output formatting gaps |

## Tasks

- [ ] Add tests for `tags.go`: `runTags` with various task sets and output formats
- [ ] Add tests for `web/handlers.go`: `handleFileUpdateError`, `handleWorklog`
- [ ] Add tests for `web/export.go`: `Export` function with mock data
- [ ] Add tests for `web/server.go`: `NewServer`, `corsMiddleware` (may need refactoring for testability)
- [ ] Add tests for `cli/importcmd.go`: extract pure logic from wizard functions and test separately; test `projectHint` and config-building logic
- [ ] Improve `cli/get.go` coverage for uncovered output formatting paths

## Acceptance Criteria

- `tags.go` `runTags` has tests covering text and JSON output
- `web/handlers.go` coverage reaches 75%+
- `web/export.go` coverage reaches 75%+
- `web` package overall coverage reaches 70%+
- Import wizard pure logic (non-interactive parts) is tested
- All tests pass: `cd apps/cli && make test`
- Linter passes: `cd apps/cli && make lint`
