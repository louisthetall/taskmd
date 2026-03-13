---
title: "Sort phase warnings deterministically in taskmd phases command"
id: "01kkk3bkw"
status: pending
priority: high
type: bug
tags: []
created: "2026-03-13"
---

# Sort phase warnings deterministically in taskmd phases command

## Objective

When running `taskmd phases` and some tasks reference phases not defined in `.taskmd.yaml`, warnings are printed listing those undefined phases. The order of these warnings is currently non-deterministic (likely due to Go map iteration order), causing inconsistent output between runs.

## Steps to Reproduce

1. Have multiple tasks referencing phases not defined in `.taskmd.yaml` config
2. Run `taskmd phases` multiple times
3. Observe that the order of "undefined phase" warnings changes between runs

## Expected Behavior

Warnings about undefined phases should always appear in a consistent, sorted order (e.g., alphabetical by phase name).

## Actual Behavior

Warning order varies between runs due to non-deterministic map iteration in Go.

## Tasks

- [ ] Locate the code in the `phases` command that collects and prints undefined phase warnings
- [ ] Sort the warnings (or the phase names) before printing
- [ ] Add a test verifying deterministic warning order

## Acceptance Criteria

- Running `taskmd phases` with undefined phases produces warnings in a stable, alphabetically sorted order
- A test confirms the ordering is deterministic
