---
id: "203"
title: "Add deduplicate command for ID collision resolution"
status: pending
priority: medium
effort: medium
type: feature
tags: [id, cli, multi-user]
parent: "200"
dependencies: ["201"]
created: 2026-02-22
---

# Add deduplicate command for ID collision resolution

## Objective

After merging branches that independently created tasks, duplicate IDs may exist. Add a `taskmd deduplicate` command that detects and resolves ID collisions by reassigning new IDs and updating all cross-references.

## Tasks

- [ ] Create `cli/deduplicate.go` with a new `deduplicate` cobra command
- [ ] Scan for duplicate IDs (reuse logic from `validator.checkDuplicateIDs`)
- [ ] For each collision, assign a new ID to the newer file (by created date or file modification time) using the configured ID strategy
- [ ] Update all cross-references: `dependencies` and `parent` fields in other tasks that reference the old ID
- [ ] Rename the task file to match the new ID (`<new-id>-<slug>.md`)
- [ ] Support `--dry-run` flag to preview changes without applying them
- [ ] Support `--format json` for structured output
- [ ] Add tests covering: no duplicates (no-op), single collision, multiple collisions, cross-reference updates, dry-run mode

## Acceptance Criteria

- `taskmd deduplicate` with no duplicates reports "no duplicates found" and exits 0
- `taskmd deduplicate` with duplicates reassigns IDs, renames files, and updates references
- `taskmd deduplicate --dry-run` shows what would change without modifying files
- Cross-references (dependencies, parent) are updated to the new ID
- The command uses the configured ID strategy for generating replacement IDs
- `taskmd validate` passes after deduplication
