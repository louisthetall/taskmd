---
id: "144"
title: "Update installation docs with Windows instructions"
status: pending
priority: medium
effort: small
tags: [docs, windows]
dependencies: ["148"]
created: 2026-02-17
---

# Update installation docs with Windows instructions

## Objective

Add Windows installation instructions to all relevant documentation so Windows users know how to install taskmd via Scoop or direct binary download.

## Tasks

- [ ] Add Scoop installation section to `README.md` (after the Homebrew section)
- [ ] Add Scoop section to `apps/docs/getting-started/installation.md` (after the Homebrew section)
- [ ] Add Windows-specific binary download instructions (PowerShell `Expand-Archive` + add to PATH)
- [ ] Mention both AMD64 and ARM64 binary availability

## Acceptance Criteria

- README.md includes Scoop install command alongside Homebrew
- Installation docs cover Windows with both Scoop and manual binary methods
- Instructions are clear and tested
