---
id: "135"
title: "Windows installation support"
status: pending
priority: low
effort: medium
tags: [distribution, windows, packaging]
created: 2026-02-16
---

# Windows Installation Support

## Objective

Make taskmd easy to install on Windows, providing a native package manager experience comparable to the existing Homebrew installation on macOS. Windows users should be able to install taskmd with a single command.

## Sub-Tasks

This is a tracking task. The work is split into the following sub-tasks:

- [ ] [141 - Add Windows ARM64 binary to release workflow](./141-windows-arm64-release-binary.md)
- [ ] [142 - Create Scoop bucket for Windows package manager](./142-scoop-bucket-windows-package-manager.md)
- [ ] [143 - Add Windows CI testing](./143-windows-ci-testing.md)
- [ ] [144 - Update installation docs with Windows instructions](./144-windows-installation-docs.md)

## Already Done

- [x] Windows AMD64 binary is built in the release workflow (`.github/workflows/release.yml`)
- [x] Binary is compressed as `.zip` and published to GitHub releases
- [x] SHA256 checksums include the Windows binary

## Acceptance Criteria

- [ ] Windows users can install taskmd via a single package manager command (e.g., `scoop install taskmd`)
- [ ] Windows binaries are automatically built and published with each GitHub release (AMD64 + ARM64)
- [ ] Installation docs cover Windows alongside macOS/Linux
- [ ] Core CLI commands work correctly on Windows (path handling, file operations)
