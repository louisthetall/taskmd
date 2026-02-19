---
id: "162"
title: "Add code coverage metric to project dashboard/README on GitHub"
status: pending
priority: medium
effort: medium
type: improvement
tags: [ci, testing, documentation]
created: 2026-02-19
---

# Add code coverage metric to project dashboard/README on GitHub

## Objective

Display a code coverage badge and/or metric in the project README on GitHub so that contributors and users can quickly see the current test coverage level. This improves project transparency and encourages maintaining high test quality.

## Tasks

- [ ] Choose a code coverage service (e.g., Codecov, Coveralls, or GitHub Actions native)
- [ ] Add a CI step to generate Go coverage reports (`go test -coverprofile=coverage.out ./...`)
- [ ] Upload coverage data to the chosen service (or publish as a GitHub Actions artifact)
- [ ] Generate a coverage badge and add it to the project README
- [ ] Verify the badge renders correctly on GitHub and updates on new pushes

## Acceptance Criteria

- A code coverage badge is visible in the project README on GitHub
- Coverage data is automatically updated on each push/PR via CI
- The coverage metric accurately reflects the Go test suite results
- Setup is documented so contributors understand how coverage is tracked
