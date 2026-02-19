---
id: "122"
title: "Add MCPB artifact to release process"
status: completed
priority: high
effort: medium
tags:
  - release
  - mcp
  - distribution
created: 2026-02-16
---

# Add MCPB Artifact to Release Process

## Objective

Update the release process to produce an `.mcpb` bundle (MCP Bundle) alongside existing release artifacts. MCPB is the [Desktop Extensions](https://github.com/modelcontextprotocol/mcpb) format for one-click local MCP server installation in desktop apps like Claude for macOS/Windows. This enables users to install taskmd's MCP server directly from a downloadable bundle.

## Background

MCPB bundles are ZIP archives containing a complete MCP server and a `manifest.json` describing the server's capabilities. Desktop apps that support MCPB (e.g., Claude) allow users to open the `.mcpb` file to trigger installation with no manual configuration.

Key references:
- MCPB spec and CLI: https://github.com/modelcontextprotocol/mcpb
- Bundle format: ZIP archive with `manifest.json` + server files
- Supported runtimes: Node.js, Python, Binary

## Tasks

- [ ] Review the MCPB specification and manifest format (`MANIFEST.md` in the mcpb repo)
- [ ] Determine the appropriate runtime type for the taskmd MCP server (likely binary)
- [ ] Create a `manifest.json` template for the taskmd MCP server bundle
- [ ] Add a build step that packages the MCP server binary + manifest into an `.mcpb` archive
- [ ] Integrate the MCPB build step into the existing release workflow
- [ ] Upload the `.mcpb` artifact alongside other release assets (GitHub release, Homebrew, etc.)
- [ ] Test installation of the `.mcpb` bundle in a supporting desktop app
- [ ] Document the MCPB distribution option for users

## Acceptance Criteria

- Each release produces a `.mcpb` file as a release artifact
- The `.mcpb` bundle contains a valid `manifest.json` and the taskmd MCP server
- The bundle can be installed via one-click in a desktop app that supports MCPB
- The release workflow (CI/scripts) is updated to build and upload the MCPB artifact automatically
