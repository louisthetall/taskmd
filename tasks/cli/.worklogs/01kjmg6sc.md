## 2026-03-01T13:00:00Z

Implemented `taskmd feed` command.

**Completed:**
- [x] Created `internal/cli/feed.go` with cobra command, git log parsing, text/JSON output
- [x] Created `internal/cli/feed_test.go` with 15 tests covering all flags, formats, parsing, and error cases
- [x] Flags: `--format text|json`, `--limit N`, `--since 2d/1w/date`, `--scope subdir`
- [x] Added `normalizeSince()` to convert shorthand durations (2d, 1w) to git-compatible format
- [x] Extracted `parseEntryLine()` to keep cognitive complexity under lint threshold
- [x] Added legend line (`[A] added  [M] modified  [R] renamed`) to text output
- [x] All tests pass, lint clean, manual verification done
