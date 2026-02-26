---
id: "216"
title: "Integrate web test coverage into CI"
status: pending
priority: medium
type: chore
effort: small
tags: ["testing", "quality", "ci"]
parent: "211"
dependencies: ["215"]
created: "2026-02-26"
---

# Integrate web test coverage into CI

## Objective

Add the web app test suite and coverage enforcement to the CI pipeline so threshold regressions block merges, and document the local testing workflow.

## Tasks

- [ ] Add web test step to CI workflow (run `pnpm test:coverage` in `apps/web`)
- [ ] Ensure CI fails when coverage drops below configured thresholds
- [ ] Document how to run tests and view coverage locally (in README or CONTRIBUTING)

## Acceptance Criteria

- CI runs `pnpm test:coverage` for the web app on every PR
- CI fails if coverage thresholds regress
- Documentation exists explaining how to run tests and view coverage locally
