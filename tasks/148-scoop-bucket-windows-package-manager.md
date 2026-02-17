---
id: "148"
title: "Create Scoop bucket for Windows package manager"
status: pending
priority: medium
effort: medium
tags: [distribution, windows, packaging]
dependencies: ["147"]
created: 2026-02-17
---

# Create Scoop bucket for Windows package manager

## Objective

Provide a native Windows package manager experience using Scoop, allowing Windows users to install taskmd with a single command: `scoop install taskmd`.

## Tasks

- [ ] Create a `driangle/scoop-bucket` GitHub repository (similar to `driangle/homebrew-tap`)
- [ ] Create a Scoop manifest JSON (`bucket/taskmd.json`) with architecture/URL/hash mappings for AMD64 and ARM64
- [ ] Add an auto-update step to `.github/workflows/release.yml` that updates the Scoop manifest on new releases (after the Homebrew formula update step)
- [ ] Add a `SCOOP_BUCKET_TOKEN` repository secret for pushing to the bucket repo
- [ ] Test installation via `scoop bucket add driangle https://github.com/driangle/scoop-bucket && scoop install taskmd`

## Acceptance Criteria

- Scoop manifest is auto-updated on each GitHub release with correct URLs and hashes
- `scoop install taskmd` installs the correct binary for the user's architecture
- `scoop update taskmd` upgrades to the latest version
