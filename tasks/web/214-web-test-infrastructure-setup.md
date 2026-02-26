---
id: "214"
title: "Set up Vitest test infrastructure for web app"
status: pending
priority: medium
type: chore
effort: small
tags: ["testing", "quality"]
parent: "211"
created: "2026-02-26"
---

# Set up Vitest test infrastructure for web app

## Objective

Configure Vitest as the test runner for the `apps/web` React/TypeScript app and write an initial set of tests to validate the setup works end-to-end.

## Tasks

- [ ] Install Vitest and required dependencies (`vitest`, `@testing-library/react`, `jsdom`, etc.)
- [ ] Add Vitest configuration (vitest.config.ts or within vite.config.ts)
- [ ] Add `test` script to `apps/web/package.json`
- [ ] Write at least one test each for: a component, a utility, and a page
- [ ] Verify `pnpm test` runs successfully from `apps/web`

## Acceptance Criteria

- `pnpm test` runs the test suite successfully from `apps/web`
- At least one test exists for a component, a utility, and a page
