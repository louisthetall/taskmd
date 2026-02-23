---
id: "195"
title: "Publish VS Code extension to marketplace"
status: pending
priority: low
effort: medium
type: chore
tags: [vscode, distribution]
created: 2026-02-22
---

# Publish VS Code extension to marketplace

## Objective

Publish the taskmd VS Code extension to the Visual Studio Code Marketplace so that users can discover and install it like any other public extension, instead of manually installing a `.vsix` file.

## Tasks

- [ ] Create a Visual Studio Marketplace publisher account (or verify existing `taskmd` publisher)
- [ ] Generate a Personal Access Token (PAT) with Marketplace publish scope
- [ ] Remove `"private": true` from `apps/vscode/package.json` (or set to `false`)
- [ ] Add required marketplace metadata to `package.json` (icon, repository, license, homepage, etc.)
- [ ] Add a `CHANGELOG.md` for the extension
- [ ] Verify `vsce package` produces a valid `.vsix`
- [ ] Publish with `vsce publish` using the PAT
- [ ] Add a CI/CD workflow (GitHub Actions) to automate publishing on version tag or release
- [ ] Verify the extension is publicly visible and installable from the marketplace

## Acceptance Criteria

- Extension is listed on the VS Code Marketplace under the `taskmd` publisher
- Users can install it by searching "taskmd" in the VS Code extensions panel
- A CI workflow exists to publish new versions automatically on release
