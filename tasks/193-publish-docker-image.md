---
id: "193"
title: "Publish Docker image to GitHub Container Registry"
status: pending
priority: medium
effort: medium
type: chore
tags: [docker, distribution]
created: 2026-02-22
---

# Publish Docker image to GitHub Container Registry

## Objective

Publish the taskmd Docker image to GitHub Container Registry (ghcr.io) so that users can pull it from a public registry instead of building locally.

## Tasks

- [ ] Review and update the existing `Dockerfile` as needed
- [ ] Set up GitHub Actions workflow to build and push the Docker image to `ghcr.io`
- [ ] Configure image tagging strategy (latest, semver, git SHA)
- [ ] Ensure the package visibility is set to public on GitHub
- [ ] Add labels and metadata to the Docker image (version, description, source URL)
- [ ] Test pulling and running the image from the public registry
- [ ] Document Docker usage in the README or docs

## Acceptance Criteria

- Docker image is publicly available at `ghcr.io/<org>/taskmd`
- Users can pull with `docker pull ghcr.io/<org>/taskmd:latest`
- A CI workflow builds and pushes new images on release or version tag
- Image includes proper labels and metadata
