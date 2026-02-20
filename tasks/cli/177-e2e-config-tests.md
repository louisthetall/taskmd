---
id: "177"
title: "E2e tests for config loading and precedence"
status: pending
priority: medium
effort: small
type: improvement
tags:
  - testing
  - cli
parent: "173"
dependencies:
  - "174"
created: 2026-02-20
---

# E2e tests for config loading and precedence

## Objective

Test that `.taskmd.yaml` configuration files are loaded correctly and that the precedence rules (CLI flags > project config > home config) work as expected end-to-end.

## Tasks

- [ ] Test project-level `.taskmd.yaml`: create config in task dir, verify it affects command behavior
- [ ] Test home-level config: create `.taskmd.yaml` in overridden HOME, verify it's picked up
- [ ] Test CLI flag overrides config: set a value in config, override with flag, verify flag wins
- [ ] Test default behavior with no config file present
- [ ] Test config options that affect output: e.g. `worklogs: false`, default format settings

## Acceptance Criteria

- Project config is loaded and affects command output
- Home config is loaded as fallback when no project config exists
- CLI flags override config values in all tested cases
- Missing config file does not cause errors
