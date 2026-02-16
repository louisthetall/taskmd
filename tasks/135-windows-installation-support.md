---
id: "135"
title: "Windows installation support"
status: pending
priority: medium
effort: medium
tags: [distribution, windows, packaging]
created: 2026-02-16
---

# Windows Installation Support

## Objective

Make taskmd easy to install on Windows, providing a native package manager experience comparable to the existing Homebrew installation on macOS. Windows users should be able to install taskmd with a single command.

## Tasks

- [ ] Research Windows package manager options (Scoop, Chocolatey, winget)
- [ ] Choose a primary package manager to support (Scoop recommended for CLI tools)
- [ ] Update the GitHub Release workflow to produce Windows binaries (`.exe`) for amd64 and arm64
- [ ] Create a Scoop bucket/manifest (or Chocolatey package/winget manifest)
- [ ] Add Windows CI/CD testing to ensure binaries work correctly
- [ ] Update installation documentation with Windows instructions
- [ ] Test end-to-end installation on a Windows machine

## Acceptance Criteria

- [ ] Windows users can install taskmd via a single package manager command (e.g., `scoop install taskmd`)
- [ ] Windows binaries are automatically built and published with each GitHub release
- [ ] Installation docs cover Windows alongside macOS/Linux
- [ ] Core CLI commands work correctly on Windows (path handling, file operations)
