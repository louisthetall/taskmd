---
id: "215"
title: "Add coverage reporting and thresholds to web app"
status: pending
priority: medium
type: chore
effort: small
tags: ["testing", "quality"]
parent: "211"
dependencies: ["214"]
created: "2026-02-26"
---

# Add coverage reporting and thresholds to web app

## Objective

Configure code coverage reporting with `@vitest/coverage-v8` and set up baseline coverage thresholds so regressions are caught early.

## Tasks

- [ ] Install `@vitest/coverage-v8`
- [ ] Configure coverage reporting in Vitest config (line, branch, function metrics)
- [ ] Add `test:coverage` script to `apps/web/package.json`
- [ ] Set initial coverage thresholds (10-20% baseline)
- [ ] Verify HTML coverage report is generated and viewable locally

## Acceptance Criteria

- `pnpm test:coverage` generates a coverage report with line, branch, and function metrics
- Coverage thresholds are configured and enforced (Vitest fails if below threshold)
- A developer can view an HTML coverage report locally
