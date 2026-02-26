---
title: "Add tests and coverage reporting to web app"
id: "211"
status: completed
priority: medium
type: feature
tags: ["testing", "quality"]
created: "2026-02-25"
---

# Add tests and coverage reporting to web app

## Objective

Set up a testing framework and coverage reporting for the web app (`apps/web`). The web app currently has no test infrastructure, so this task establishes the foundation: test runner configuration, a coverage reporting pipeline, and an initial set of tests to validate the setup works end-to-end.

## Tasks

- [ ] Set up Vitest (or similar) as the test runner for the React/TypeScript web app
- [ ] Configure coverage reporting (e.g. `vitest --coverage` with `@vitest/coverage-v8`)
- [ ] Add test scripts to `apps/web/package.json` (`test`, `test:coverage`)
- [ ] Write initial tests for a representative set of components/utilities to validate the setup
- [ ] Add coverage thresholds configuration (start low, e.g. 10-20%, as a baseline)
- [ ] Integrate coverage report into CI (fail on threshold regression)
- [ ] Document how to run tests and view coverage locally

## Sub-tasks

- **214** — Set up Vitest test infrastructure for web app
- **215** — Add coverage reporting and thresholds to web app
- **216** — Integrate web test coverage into CI

## Acceptance Criteria

- `pnpm test` runs the test suite successfully from `apps/web`
- `pnpm test:coverage` generates a coverage report with line, branch, and function metrics
- At least one test exists for a component, a utility, and a page
- Coverage thresholds are configured and enforced in CI
- A developer can view an HTML coverage report locally
