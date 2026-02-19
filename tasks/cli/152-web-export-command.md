---
id: "152"
title: "Add `taskmd web export` subcommand for static site export"
status: in-progress
priority: medium
effort: medium
type: feature
tags: [cli, web, export]
created: 2026-02-17
---

# Add `taskmd web export` subcommand for static site export

## Objective

Add a `taskmd web export` subcommand that exports a fully static bundle of the taskmd web UI, ready for standalone deployment. This allows users to generate a self-contained site (HTML, CSS, JS, and pre-rendered data) that can be hosted on any static file server (e.g., GitHub Pages, Netlify, S3) without requiring the `taskmd web` server to be running.

## Tasks

- [ ] Add `export` subcommand under the existing `web` command
- [ ] Implement static data generation (scan tasks and write JSON data files)
- [ ] Bundle the embedded web frontend assets into the export output directory
- [ ] Configure the frontend to load from static JSON files instead of API routes
- [ ] Add `--output` / `-o` flag to specify the export directory (default: `./taskmd-export`)
- [ ] Add `--base-path` flag for deployments under a subpath (e.g., `/projects/myapp/`)
- [ ] Ensure all assets use relative paths so the bundle works from any hosting root
- [ ] Write tests for the export command (happy path, flags, error handling)
- [ ] Verify the exported bundle opens correctly in a browser with no server

## Acceptance Criteria

- Running `taskmd web export` produces a self-contained directory with all HTML, CSS, JS, and data files
- The exported site displays the same task data as `taskmd web` (read-only)
- The `--output` flag controls the destination directory
- The `--base-path` flag correctly adjusts asset and routing paths
- The exported bundle works when served by any static file server (e.g., `python -m http.server`)
- Tests cover command flags, output structure, and error cases
