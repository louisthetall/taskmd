## 2026-03-01T12:00:00Z

Completed implementation of replacing `current` command with `status` no-args mode.

**Changes:**
- Deleted `current.go` and `current_test.go`
- Updated `status` command to accept optional query arg (`MaximumNArgs(1)`)
- Added `--statusline` flag for compact shell integration output (`#ID title (+N more)`)
- Added `--scope` flag to filter by group/directory in no-args mode
- Added 10 new tests covering all no-args scenarios (zero/one/multiple tasks, statusline, scope, JSON, YAML)
- Updated `apps/docs/guide/cli.md`: removed `current` section, updated `status` section with new flags and statusline examples

**Verification:** All tests pass, lint clean, build succeeds.
