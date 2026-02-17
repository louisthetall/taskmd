---
id: "146"
title: "Add Homebrew formula audit checks"
status: pending
priority: low
effort: small
tags: [homebrew, distribution]
created: 2026-02-16
---

# Add Homebrew formula audit checks

## Objective

Ensure the Homebrew formula passes `brew audit --strict` so it meets the style guidelines required for submission to the official `homebrew-core` tap.

## Tasks

- [ ] Run `brew audit --strict taskmd` against the current formula
- [ ] Fix any style guideline violations reported by the audit
- [ ] Add a CI step to run `brew audit` on formula changes
- [ ] Verify the formula meets homebrew-core submission requirements

## Acceptance Criteria

- `brew audit --strict taskmd` passes with no errors or warnings
- Formula follows all Homebrew style guidelines for homebrew-core
