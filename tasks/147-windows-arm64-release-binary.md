---
id: "147"
title: "Add Windows ARM64 binary to release workflow"
status: completed
priority: low
effort: small
tags: [distribution, windows, ci]
created: 2026-02-17
---

# Add Windows ARM64 binary to release workflow

## Objective

The release workflow already builds a Windows AMD64 binary. Add a Windows ARM64 build so ARM-based Windows devices (e.g., Surface Pro X, Snapdragon laptops) are also supported.

## Tasks

- [x] Add `GOOS=windows GOARCH=arm64` build step in `.github/workflows/release.yml`
- [x] Compress the ARM64 binary as `.zip` and include it in release artifacts
- [x] Include the ARM64 zip in SHA256 checksum generation
- [x] Update `build-all` target in `apps/cli/Makefile` to include `windows-arm64`

## Acceptance Criteria

- Release workflow produces a `taskmd-windows-arm64.zip` artifact alongside the existing AMD64 zip
- SHA256 checksums file includes the Windows ARM64 binary
- `make build-all` builds the Windows ARM64 binary
