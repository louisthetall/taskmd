---
id: "183"
title: "Add Docker image for taskmd web server"
status: completed
priority: low
effort: medium
type: feature
tags:
  - devops
  - docker
  - distribution
created: 2026-02-20
---

# Add Docker Image for taskmd Web Server

## Objective

Provide an official Docker image so users can run the taskmd web server in a container. This simplifies deployment for teams, CI environments, and users who prefer containerized workflows.

## Tasks

- [x] Create a multi-stage Dockerfile (build Go binary + embed web assets, then copy to minimal runtime image)
- [x] Use a minimal base image (e.g., `alpine` or `distroless`)
- [x] Configure the entrypoint to run `taskmd web start`
- [x] Expose port 8080 by default
- [x] Support mounting the tasks directory as a volume
- [x] Support passing config via environment variables or mounted `.taskmd.yaml`
- [x] Add a `.dockerignore` file to keep the image small
- [x] Add a `docker-compose.yml` example for quick setup
- [x] Set up GitHub Actions to build and push the image on release
- [x] Document usage in the docs site (installation page)
- [x] Add tests to verify the image builds and the server starts

## Acceptance Criteria

- `docker run -v ./tasks:/tasks ghcr.io/owner/taskmd` starts the web server
- The image is under 50MB
- Tasks directory is mounted as a volume (not baked into the image)
- The image is published automatically on each release
- Documentation includes docker-compose example for common setups
