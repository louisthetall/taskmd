---
id: "187"
title: "Format markdown content in get command output"
status: completed
priority: medium
effort: medium
type: improvement
tags: [cli, ux, formatting]
created: 2026-02-21
---

# Format markdown content in get command output

## Objective

Improve the `taskmd get` command to render markdown content with terminal-friendly formatting instead of displaying raw markdown. Headings should appear **bold**, backticked content should be visually distinct (e.g., dimmed or highlighted), and common markdown constructs like `**bold**`, `*italic*`, lists, and links should be rendered using ANSI styling — similar to how Claude Code presents markdown in the terminal.

A `--raw-markdown` flag should be added to bypass formatting and print the original markdown as-is.

## Tasks

- [x] Implement a markdown-to-terminal renderer that handles:
  - `# Headings` → bold (and optionally colored)
  - `**bold**` → bold ANSI
  - `*italic*` → italic or dimmed ANSI
  - `` `inline code` `` → distinct styling (e.g., dimmed, background, or colored)
  - ```` ``` code blocks ``` ```` → indented or boxed with distinct styling
  - `- [ ] / - [x]` checkboxes → rendered with visual indicators
  - `- list items` → preserved with proper indentation
  - `[links](url)` → show text with URL in parentheses or dimmed
- [x] Integrate the renderer into the `get` command's default output path
- [x] Add `--raw-markdown` flag to `get` command that prints unformatted markdown
- [x] Ensure formatted output degrades gracefully when piped (no ANSI codes) or when `--no-color` is set
- [x] Add tests for the markdown renderer covering each supported construct
- [x] Add tests for the `--raw-markdown` flag

## Acceptance Criteria

- `taskmd get <id>` displays markdown body with ANSI-formatted headings, bold, italic, inline code, and code blocks
- `taskmd get <id> --raw-markdown` prints the raw markdown without any ANSI formatting
- Piped output (non-TTY) omits ANSI codes automatically
- All existing `get` command tests continue to pass
- New unit tests cover each markdown construct handled by the renderer
