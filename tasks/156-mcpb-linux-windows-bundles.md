---
id: "156"
title: "Add MCPB bundles for Linux and Windows platforms"
status: pending
priority: low
effort: small
tags:
  - release
  - mcp
  - distribution
created: 2026-02-19
---

# Add MCPB Bundles for Linux and Windows Platforms

## Objective

Extend the MCPB bundle generation (added in task 122 for macOS only) to include Linux and Windows platforms once those platforms support MCPB desktop extensions.

## Background

Task 122 added `.mcpb` artifact generation for macOS (darwin-amd64, darwin-arm64). The build script (`scripts/build-mcpb.sh`) already handles all platforms generically, but the CI workflow and Makefile only produce macOS bundles. This task extends coverage to all 6 platform/arch combinations.

## Tasks

- [ ] Update `.github/workflows/release.yml` to call `build-mcpb.sh` for linux-amd64, linux-arm64, windows-amd64, windows-arm64
- [ ] Update release.yml release file list to include the 4 new `.mcpb` artifacts
- [ ] Update `apps/cli/Makefile` `mcpb-all` target to build Linux and Windows bundles
- [ ] Update `scripts/release.sh` success output to list all 6 `.mcpb` artifacts
- [ ] Verify `scripts/build-mcpb.sh` correctly handles Windows `.exe` suffix and `win32` platform name
- [ ] Test bundles install correctly on target platforms

## Acceptance Criteria

- Each release produces 6 `.mcpb` files (one per platform/arch)
- Windows bundles contain `taskmd.exe` and manifest references `.exe` suffix
- Checksums file includes all `.mcpb` artifacts
- `make mcpb-all` builds all 6 bundles locally
