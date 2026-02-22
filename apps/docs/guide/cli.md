# CLI Guide

Complete reference for using taskmd from the command line.

## Command Reference

### Quick Reference

| Command | Description |
|---------|-------------|
| `list` | List tasks in a quick textual format |
| `get` | Get detailed information about a specific task |
| `set` | Set a task's frontmatter fields |
| `next` | Recommend what task to work on next |
| `validate` | Lint and validate tasks |
| `graph` | Export task dependency graph |
| `board` | Display tasks grouped in a kanban-like board view |
| `stats` | Show computed metrics about tasks |
| `tags` | List all tags with task counts |
| `snapshot` | Produce a frozen, machine-readable representation of tasks |
| `report` | Generate a comprehensive project report |
| `tracks` | Show parallel work tracks based on scope overlap |
| `archive` | Archive or delete completed/cancelled tasks |
| `rm` | Delete a task file permanently |
| `deduplicate` | Detect and resolve duplicate task IDs |
| `next-id` | Show the next available task ID |
| `add` | Create a new task file with proper frontmatter |
| `search` | Full-text search across task titles and bodies |
| `templates` | List and manage task templates |
| `verify` | Run verification checks for a task |
| `status` | Get lightweight metadata for a task |
| `context` | Show file context for a task |
| `worklog` | View or add worklog entries for a task |
| `import` | Import tasks from external sources |
| `spec` | Generate the taskmd specification file |
| `sync` | Sync tasks from external sources |
| `web` | Web dashboard commands |
| `init` | Initialize a project with agent configuration and spec files |
| `commit-msg` | Generate conventional commit messages from task metadata |
| `mcp` | Start MCP server over stdio |
| `todos` | Find TODO/FIXME comments in source code |
| `completion` | Generate shell completion scripts |

---

### list - View and Filter Tasks

Display tasks in various formats with filtering and sorting.

**Basic usage:**
```bash
# List all tasks
taskmd list

# List tasks in specific directory
taskmd list ./tasks

# Different output formats
taskmd list --format table   # Default
taskmd list --format json
taskmd list --format yaml
```

**Filtering:**
```bash
# Filter by status
taskmd list --filter status=pending
taskmd list --filter status=in-progress

# Filter by priority
taskmd list --filter priority=high

# Filter by multiple criteria (AND logic)
taskmd list --filter status=pending --filter priority=high

# Filter by tag
taskmd list --filter tag=cli

# Filter by effort
taskmd list --filter effort=small
```

**Sorting:**
```bash
# Sort by priority
taskmd list --sort priority

# Sort by status
taskmd list --sort status

# Sort by created date
taskmd list --sort created
```

**Custom columns:**
```bash
# Show specific columns
taskmd list --columns id,title,status

# Show more columns
taskmd list --columns id,title,status,priority,effort,deps
```

**Limiting results:**
```bash
# Show only 5 tasks
taskmd list --sort priority --limit 5
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--filter` | | Filter tasks (repeatable, AND logic) |
| `--sort` | | Sort by field (`id`, `title`, `status`, `priority`, `effort`, `created`) |
| `--columns` | `id,title,status,priority,file` | Comma-separated list of columns to display |
| `--limit` | `0` | Maximum number of tasks to display (0 = unlimited) |
| `--format` | `table` | Output format (`table`, `json`, `yaml`) |

**Examples:**
```bash
# High-priority pending tasks
taskmd list --filter status=pending --filter priority=high

# Small tasks (quick wins)
taskmd list --filter effort=small --filter status=pending

# All CLI-related tasks
taskmd list --filter tag=cli --sort priority

# Top 5 by priority
taskmd list --sort priority --limit 5

# Export to JSON for scripting
taskmd list --format json > tasks.json
```

### validate - Check Task Files

Validate task files for errors and consistency issues.

```bash
# Validate all tasks
taskmd validate

# Validate specific directory
taskmd validate ./tasks

# Strict mode (enable warnings)
taskmd validate --strict

# JSON output
taskmd validate --format json
```

**What it checks:**
- Required fields present (id, title, status)
- Valid field values
- Duplicate task IDs
- Missing dependencies (references to non-existent tasks)
- Circular dependencies
- YAML syntax errors

**Exit codes:**
- `0` - Valid (no errors)
- `1` - Invalid (errors found)
- `2` - Valid with warnings (strict mode only)

### next - Find What to Work On

Analyze tasks and recommend the best ones to work on next.

taskmd scores tasks based on:
- **Priority**: High priority scores higher
- **Critical path**: Tasks on the critical path score higher
- **Downstream impact**: Tasks blocking many others score higher
- **Effort**: Smaller tasks get a boost (quick wins)
- **Actionability**: Only tasks with satisfied dependencies

```bash
# Get top 5 recommendations
taskmd next

# Get top 3 recommendations
taskmd next --limit 3

# Next high-priority task
taskmd next --filter priority=high

# Next small task (quick win)
taskmd next --filter effort=small --limit 1

# Show only quick wins (effort: small)
taskmd next --quick-wins

# Show only critical path tasks
taskmd next --critical --limit 1

# JSON for automation
taskmd next --format json
```

### graph - Visualize Dependencies

Export task dependency graphs in various formats.

```bash
# ASCII art (terminal-friendly)
taskmd graph --format ascii

# Mermaid diagram
taskmd graph --format mermaid

# Graphviz DOT
taskmd graph --format dot

# JSON structure
taskmd graph --format json
```

**Filtering:**
```bash
# Exclude completed tasks (default)
taskmd graph

# Include all tasks
taskmd graph --all

# Exclude specific statuses
taskmd graph --exclude-status completed --exclude-status blocked
```

**Focus on specific tasks:**
```bash
# Show task and its dependencies (upstream)
taskmd graph --root 022 --upstream

# Show task and what depends on it (downstream)
taskmd graph --root 022 --downstream

# Show full subgraph
taskmd graph --root 022
```

**Output to file:**
```bash
taskmd graph --format mermaid --out deps.mmd
taskmd graph --format dot --out deps.dot

# Generate PNG with Graphviz
taskmd graph --format dot | dot -Tpng > graph.png
```

### stats - Project Metrics

Display computed statistics about your task set.

```bash
# Show all statistics
taskmd stats

# Specific directory
taskmd stats ./tasks

# JSON output
taskmd stats --format json
```

**Metrics provided:**
- Total tasks and count by status
- Priority breakdown
- Effort breakdown
- Blocked tasks count
- Completion rate
- Critical path length
- Max dependency depth
- Average dependencies per task

### board - Kanban View

Display tasks grouped by a field in a board layout.

```bash
# Group by status (default)
taskmd board

# Group by priority
taskmd board --group-by priority

# Group by effort
taskmd board --group-by effort

# Group by tag
taskmd board --group-by tag

# Output formats
taskmd board --format md    # Markdown (default)
taskmd board --format txt   # Plain text
taskmd board --format json  # JSON
```

### snapshot - Machine-Readable Export

Produce a static, machine-readable representation for automation.

```bash
# Full snapshot (JSON)
taskmd snapshot

# Core fields only
taskmd snapshot --core

# Include derived analysis
taskmd snapshot --derived

# Output formats
taskmd snapshot --format json
taskmd snapshot --format yaml
taskmd snapshot --format md

# Grouping
taskmd snapshot --group-by status
taskmd snapshot --group-by priority

# Output to file
taskmd snapshot --out snapshot.json
```

### get - View Task Details

> **Alias:** `show` is a deprecated alias for `get`. Use `get` instead.

Display detailed information about a specific task, identified by ID, title, or file path.

**Matching priority:**
1. Exact match by task ID (case-sensitive)
2. Exact match by task title (case-insensitive)
3. Match by file path or filename
4. Fuzzy match across IDs and titles (unless `--exact` is set)

```bash
# Look up by task ID
taskmd get cli-037

# Look up by title
taskmd get "Add show command"

# Look up by file path
taskmd get tasks/cli/037-task.md

# Fuzzy search
taskmd get sho

# Strict lookup — fail if no exact match
taskmd get sho --exact
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--format` | `text` | Output format (`text`, `json`, `yaml`) |
| `--exact` | `false` | Disable fuzzy matching |
| `--threshold` | `0.6` | Fuzzy match sensitivity (0.0–1.0) |

### set - Update Task Fields

> **Alias:** `update` is a deprecated alias for `set`. Use `set` instead.

Modify a task's frontmatter fields by ID.

```bash
# Change status
taskmd set 042 --status in-progress

# Change priority and effort
taskmd set 042 --priority high --effort large

# Mark as completed (shortcut)
taskmd set 042 --done

# Preview changes without writing
taskmd set 042 --priority critical --dry-run
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `[task-id]` | | Task ID as positional argument |
| `--task-id` | | Task ID to update (alternative to positional) |
| `--status` | | New status (`pending`, `in-progress`, `in-review`, `completed`, `blocked`, `cancelled`) |
| `--priority` | | New priority (`low`, `medium`, `high`, `critical`) |
| `--effort` | | New effort (`small`, `medium`, `large`) |
| `--owner` | | Owner/assignee |
| `--parent` | | Parent task ID (empty string to clear) |
| `--done` | `false` | Alias for `--status completed` |
| `--dry-run` | `false` | Preview changes without writing to disk |
| `--add-tag` | | Add a tag (repeatable) |
| `--remove-tag` | | Remove a tag (repeatable) |
| `--add-pr` | | Add a PR URL (repeatable) |
| `--remove-pr` | | Remove a PR URL (repeatable) |
| `--type` | | Work type (`feature`, `bug`, `improvement`, `chore`, `docs`) |
| `--verify` | `false` | Run verification checks before completing a task |

**Tag management:**
```bash
# Add tags
taskmd set 042 --add-tag backend --add-tag api

# Remove a tag
taskmd set 042 --remove-tag deprecated

# Add and remove in one command
taskmd set 042 --add-tag v2 --remove-tag v1
```

### tags - List Tags

Display all tags used across task files with usage counts.

```bash
# List all tags
taskmd tags

# List tags in specific directory
taskmd tags ./tasks

# Tags used by pending tasks only
taskmd tags --filter status=pending

# JSON output
taskmd tags --format json
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--format` | `table` | Output format (`table`, `json`, `yaml`) |
| `--filter` | | Filter tasks before aggregating (repeatable) |

### archive - Archive Completed Tasks

Move completed or cancelled task files into an `archive/` subdirectory, or permanently delete them.

```bash
# Archive all completed tasks
taskmd archive --all-completed -y

# Archive all cancelled tasks
taskmd archive --all-cancelled -y

# Archive specific tasks by ID
taskmd archive --id 042 --id 043 -y

# Preview what would be archived
taskmd archive --all-completed --dry-run

# Permanently delete cancelled tasks
taskmd archive --all-cancelled --delete -f
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--id` | | Archive task(s) by ID (repeatable) |
| `--status` | | Archive tasks matching this status |
| `--all-completed` | `false` | Archive all completed tasks |
| `--all-cancelled` | `false` | Archive all cancelled tasks |
| `--tag` | | Archive tasks with this tag |
| `--dry-run` | `false` | Preview changes without making them |
| `--yes`, `-y` | `false` | Skip confirmation prompt |
| `--delete` | `false` | Permanently delete instead of archive |
| `--force`, `-f` | `false` | Skip confirmation for delete |

### rm - Delete a Task

Permanently delete a task file by ID. Displays the task details and asks for confirmation before deleting.

```bash
# Delete a task (with confirmation prompt)
taskmd rm 042

# Skip confirmation
taskmd rm 042 --force

# Preview what would be deleted
taskmd rm 042 --dry-run
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--force`, `-f` | `false` | Skip confirmation prompt |
| `--dry-run` | `false` | Preview what would be deleted without acting |

### deduplicate - Resolve Duplicate IDs

Detect and resolve duplicate task IDs that can occur when multiple contributors create tasks on separate branches.

For each collision, the oldest task (by created date) keeps its original ID. Newer tasks get reassigned a fresh ID, with file renames and cross-reference updates applied automatically.

```bash
# Detect and fix duplicates
taskmd deduplicate

# Scan a specific directory
taskmd deduplicate ./tasks

# Preview changes without modifying files
taskmd deduplicate --dry-run

# Skip interactive prompts for ambiguous references
taskmd deduplicate --no-interactive

# JSON output
taskmd deduplicate --format json
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--dry-run` | `false` | Preview changes without modifying files |
| `--format` | `text` | Output format (`text`, `json`) |
| `--no-interactive` | `false` | Skip interactive prompts for ambiguous references |

### next-id - Get Next Available ID

Scan task files and output the next available sequential ID. Finds the highest numeric ID and returns max + 1, preserving any common prefix and zero-padding.

```bash
# Get next ID
taskmd next-id

# Scan specific directory
taskmd next-id ./tasks/cli

# JSON output with metadata
taskmd next-id --format json
```

**Scripting example:**
```bash
# Create a new task file with the next ID
ID=$(taskmd next-id)
echo "---
id: \"$ID\"
title: \"My new task\"
status: pending
---" > "tasks/${ID}-my-new-task.md"
```

### commit-msg - Generate Commit Messages

Generate a conventional commit message derived from task metadata.

When `--task-id` is provided, the message is generated from that task. When no `--task-id` is provided, the command inspects staged changes (`git diff --cached`) to find task files whose status changed to `completed` and generates a message from those tasks automatically.

The subject line format is `type(scope): lowercase title (task ID)`, where the scope is the task's group directory (if any).

```bash
# Generate message for a specific task
taskmd commit-msg --task-id 042

# Use a custom commit type
taskmd commit-msg --task-id 042 --type feat

# Include completed subtasks as bullet points in the body
taskmd commit-msg --task-id 042 --body

# Subject line only (no body)
taskmd commit-msg --task-id 042 --short

# Auto-detect completed tasks from staged changes
taskmd commit-msg

# Use with git commit
git commit -m "$(taskmd commit-msg --task-id 042)"
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--task-id` | | Task ID to generate the message for (omit to auto-detect from staged changes) |
| `--type` | `chore` | Commit type prefix (`feat`, `fix`, `chore`, `docs`, `test`, `refactor`) |
| `--body` | `false` | Include completed subtasks (`- [x]`) as bullet points in the commit body |
| `--short` | `false` | Output the subject line only (no body) |

**Auto-detection:**

When `--task-id` is omitted, the command runs `git diff --cached` and looks for task files where `+status: completed` appears in the diff. It then generates a commit message from all matched tasks. If multiple tasks are found, the subject line lists all task IDs (e.g., `chore: complete tasks 042, 043`).

### add - Create a New Task

Create a new task markdown file with proper frontmatter. The title is used to generate both the task title and the filename slug. A sequential ID is automatically assigned based on existing tasks.

```bash
# Create a task with just a title
taskmd add "Fix the login bug"

# Set priority and tags
taskmd add "Implement OAuth" --priority high --tags backend,auth

# Create in a subdirectory group
taskmd add "Design mockups" --group design --effort large

# Open in $EDITOR after creation
taskmd add "Quick fix" --edit

# With dependencies
taskmd add "Deploy to staging" --depends-on 041,042

# JSON output for scripting
taskmd add "Automated task" --format json
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--priority` | `medium` | Task priority (`low`, `medium`, `high`, `critical`) |
| `--effort` | | Task effort (`small`, `medium`, `large`) |
| `--tags` | | Comma-separated tags |
| `--status` | `pending` | Task status (`pending`, `in-progress`, `completed`, `blocked`, `cancelled`) |
| `--owner` | | Task owner/assignee |
| `--depends-on` | | Comma-separated dependency task IDs |
| `--parent` | | Parent task ID |
| `--group` | | Subdirectory to create the task in |
| `--format` | `plain` | Output format (`plain`, `json`) |
| `--edit` | `false` | Open the new task in `$EDITOR` |
| `--template` | | Use a task template (e.g., `bug`, `feature`, `chore`) |

**Templates:**
```bash
# Create from a template
taskmd add "Login fails on Safari" --template bug
taskmd add "Dark mode support" --template feature --priority high
```

### templates - Manage Task Templates

List and inspect task templates used by the `add` command. Templates are discovered from three sources in precedence order: project (`.taskmd/templates/`), user (`~/.taskmd/templates/`), and built-in.

```bash
# List available templates
taskmd templates list

# JSON output
taskmd templates list --format json

# YAML output
taskmd templates list --format yaml
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--format` | `table` | Output format (`table`, `json`, `yaml`) |

### search - Full-Text Search

Perform case-insensitive full-text search across all task titles and markdown body content. Results show where the match was found and a context snippet.

```bash
# Search for a keyword
taskmd search "authentication"

# JSON output
taskmd search deploy --format json

# Filter and sort results
taskmd search "auth" --filter priority=high
taskmd search "deploy" --filter status=pending --sort priority --limit 5

# YAML output
taskmd search "bug fix" --format yaml
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--format` | `table` | Output format (`table`, `json`, `yaml`) |
| `--filter` | | Filter tasks (repeatable, AND logic, e.g., `--filter status=pending --filter priority=high`) |
| `--sort` | | Sort by field (`id`, `title`, `status`, `priority`, `effort`, `created`) |
| `--limit` | `0` | Maximum number of results (0 = unlimited) |

### verify - Run Verification Checks

Run the acceptance checks defined in a task's `verify` field. Each verify step has a type:

- **bash** -- runs a shell command, reports pass/fail based on exit code
- **assert** -- displays a check for the agent to evaluate (not executed)

```bash
# Verify a task
taskmd verify --task-id 042

# JSON output
taskmd verify --task-id 042 --format json

# Preview checks without executing
taskmd verify --task-id 042 --dry-run

# Custom timeout (seconds) per command
taskmd verify --task-id 042 --timeout 120
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--task-id` | *(required)* | Task ID to verify |
| `--format` | `table` | Output format (`table`, `json`) |
| `--dry-run` | `false` | List checks without executing |
| `--timeout` | `60` | Per-command timeout in seconds |

**Exit codes:**
- `0` - All executable checks passed
- `1` - One or more executable checks failed

### status - Lightweight Task Metadata

Display only the frontmatter metadata of a task, without body content, resolved dependency info, context files, or worklog data. Use this when you just need to quickly check a task's status, priority, or other metadata.

Matching uses the same logic as `get` (ID, title, file path, fuzzy).

```bash
# Look up by task ID
taskmd status 042

# Look up by title
taskmd status "Setup project"

# JSON output
taskmd status 042 --format json

# YAML output
taskmd status 042 --format yaml

# Strict lookup (no fuzzy matching)
taskmd status sho --exact
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--format` | `text` | Output format (`text`, `json`, `yaml`) |
| `--exact` | `false` | Disable fuzzy matching, exact only |
| `--threshold` | `0.6` | Fuzzy match sensitivity (0.0-1.0) |

### context - Show File Context

Resolve all relevant files for a task into a structured output. Files come from two sources:

1. **Scope files** -- resolved from the task's `touches` field via scope definitions in `.taskmd.yaml`
2. **Explicit files** -- listed directly in the task's `context` field

```bash
# Show context for a task
taskmd context --task-id 042

# JSON output
taskmd context --task-id 042 --format json

# Include file contents and task body
taskmd context --task-id 042 --include-content --resolve

# Include files from dependency tasks
taskmd context --task-id 042 --include-deps

# Limit number of files
taskmd context --task-id 042 --max-files 20
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--task-id` | *(required)* | Task ID to build context for |
| `--format` | `text` | Output format (`text`, `json`, `yaml`) |
| `--resolve` | `false` | Expand directory paths to individual files |
| `--include-content` | `false` | Inline file contents and task body |
| `--include-deps` | `false` | Include files from direct dependency tasks |
| `--max-files` | `0` | Cap number of files (0 = unlimited) |

### worklog - View or Add Worklog Entries

View or add timestamped worklog entries for a task. Worklog files are stored at `tasks/<group>/.worklogs/<ID>.md`.

```bash
# View worklog entries
taskmd worklog --task-id 015

# Add a new entry
taskmd worklog --task-id 015 --add "Started implementation"

# JSON output
taskmd worklog --task-id 015 --format json

# YAML output
taskmd worklog --task-id 015 --format yaml
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--task-id` | *(required)* | Task ID |
| `--add` | | Append a new worklog entry with the given text |
| `--format` | `text` | Output format (`text`, `json`, `yaml`) |

### import - Import Tasks from External Sources

Fetch tasks from an external source (GitHub Issues, Jira, etc.) and create local markdown task files. This is a one-time onboarding tool for populating your `tasks/` directory.

When run without `--source`, an interactive wizard guides you through setup.

```bash
# Interactive wizard
taskmd import

# GitHub: import from a repository
taskmd import --source github --repo owner/repo

# GitHub: with auth token and filters
taskmd import --source github --project owner/repo --token-env GITHUB_TOKEN

# GitHub: filter by labels and assignee
taskmd import --source github --repo owner/repo --labels bug,critical --assignee alice

# GitHub: filter by milestone
taskmd import --source github --repo owner/repo --milestone "v2.0"

# Jira: import from a project
taskmd import --source jira --project PROJ --url https://company.atlassian.net

# Jira: with JQL filter
taskmd import --source jira --project PROJ --url https://company.atlassian.net --jql "assignee = currentUser()"

# Preview without writing files
taskmd import --source github --repo owner/repo --dry-run

# JSON output for scripting
taskmd import --source github --repo owner/repo --format json
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--source` | | Source name (`github`, `jira`, etc.). Omit for interactive wizard |
| `--project` | | Project identifier (owner/repo for GitHub, project key for Jira) |
| `--token-env` | | Environment variable name for auth token |
| `--user-env` | | Environment variable name for username (Jira) |
| `--base-url` | | API base URL (for Jira or GitHub Enterprise) |
| `--output-dir` | `./tasks` | Target directory for imported task files |
| `--filter` | | Source-specific filters as key:value pairs (e.g. `"state:open labels:bug"`) |
| `--dry-run` | `false` | Preview import without writing files |
| `--format` | `table` | Output format (`table`, `json`, `yaml`) |

**GitHub-specific flags:**

| Flag | Description |
|------|-------------|
| `--repo` | Alias for `--project` (owner/repo) |
| `--labels` | Filter by labels (comma-separated) |
| `--milestone` | Filter by milestone |
| `--assignee` | Filter by assignee |

**Jira-specific flags:**

| Flag | Description |
|------|-------------|
| `--url` | Alias for `--base-url` (Jira instance URL) |
| `--jql` | Jira Query Language filter |

### spec - Generate Specification File

Generate the taskmd specification document in your project directory. The specification describes the task file format, including frontmatter fields, valid values, file naming conventions, and directory structure.

```bash
# Write TASKMD_SPEC.md to current directory
taskmd spec

# Print spec to stdout
taskmd spec --stdout

# Write to a specific directory
taskmd spec --dir ./docs

# Overwrite existing file
taskmd spec --force
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--force` | `false` | Overwrite existing `TASKMD_SPEC.md` |
| `--stdout` | `false` | Print spec to stdout instead of writing a file |

### mcp - Start MCP Server

Start a Model Context Protocol (MCP) server that communicates over stdin/stdout. This allows LLM-based tools (Cursor, Windsurf, Copilot agents, Claude Code, etc.) to interact with your taskmd project using the standard MCP protocol.

```bash
# Start MCP server
taskmd mcp
```

**Configuration example for Claude Code (`.mcp.json`):**
```json
{
  "mcpServers": {
    "taskmd": {
      "command": "taskmd",
      "args": ["mcp"]
    }
  }
}
```

### report - Generate Reports

Generate a comprehensive project report combining summary statistics, task groupings, critical-path analysis, and blocked tasks.

```bash
# Markdown report to stdout
taskmd report

# HTML report to file
taskmd report --format html --out report.html

# JSON report grouped by priority
taskmd report tasks/ --group-by priority --format json

# Include dependency graph
taskmd report tasks/ --format html --include-graph --out report.html
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--format` | `md` | Output format (`md`, `html`, `json`) |
| `--group-by` | `status` | Field to group by (`status`, `priority`, `effort`, `group`, `tag`) |
| `--out`, `-o` | | Write output to file |
| `--include-graph` | `false` | Embed dependency graph in report |

### tracks - Parallel Work Tracks

Assign actionable tasks to parallel work tracks based on the `touches` frontmatter field. Tasks that share a scope are placed in separate tracks so they can be worked on without merge conflicts.

```bash
# Show work tracks
taskmd tracks

# Filter to CLI-related tasks
taskmd tracks --filter tag=cli

# Limit to top 3 tracks
taskmd tracks --limit 3

# Export track assignments
taskmd tracks --format json > tracks.json
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--format` | `table` | Output format (`table`, `json`, `yaml`) |
| `--filter` | | Filter tasks (repeatable) |
| `--limit` | `0` | Maximum number of tracks (0 = unlimited) |

### sync - Sync External Sources

Commands for syncing tasks with external sources (GitHub Issues, Jira, etc.). Running `taskmd sync` alone displays usage and available subcommands.

#### sync down

Fetch tasks from configured external sources and create or update local markdown task files. Configuration is read from `.taskmd.yaml`.

```bash
# Sync all configured sources
taskmd sync down

# Preview without writing files
taskmd sync down --dry-run

# Sync a specific source
taskmd sync down --source github
taskmd sync down --source jira

# Overwrite local changes with remote data
taskmd sync down --conflict remote
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--dry-run` | `false` | Preview changes without writing files |
| `--source` | | Sync only the named source |
| `--conflict` | `skip` | Conflict strategy: `skip`, `remote`, or `local` |

See [Configuration](/reference/configuration#sync-configuration) for how to set up sync sources in `.taskmd.yaml`.

### web - Web Dashboard

Commands for the taskmd web dashboard.

#### web start

Start the web interface server.

```bash
# Start server
taskmd web start

# Start and open browser
taskmd web start --open

# Custom port
taskmd web start --port 3000

# Read-only mode (disables editing)
taskmd web start --readonly

# Specific tasks directory
taskmd web start --task-dir ./my-tasks --open
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--port` | `8080` | Server port |
| `--open` | `false` | Open browser on start |
| `--dev` | `false` | Enable dev mode (CORS for Vite dev server) |
| `--readonly` | `false` | Start in read-only mode (disables editing) |

#### web export

Export the dashboard as a self-contained static site. The exported site can be deployed to GitHub Pages, Netlify, S3, or any static file host.

```bash
# Export to default directory (./taskmd-export)
taskmd web export

# Export to a specific directory
taskmd web export -o ./public

# Set base path for URLs (e.g., for GitHub Pages subpath)
taskmd web export --base-path /demo/

# Export with custom task directory
taskmd web export --task-dir ./tasks -o ./site
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `-o`, `--output` | `./taskmd-export` | Output directory |
| `--base-path` | `/` | Base path for URLs (e.g., `/demo/`) |

See the [Web Interface Guide](./web) for detailed web UI documentation.

### todos - Find TODO/FIXME Comments

Scan source code files recursively for marker comments (TODO, FIXME, HACK, XXX, NOTE, BUG, OPTIMIZE) and display them with file path, line number, marker type, and comment text.

Respects `.gitignore` and skips common non-source directories (`node_modules`, `.git`, `vendor`, etc.). Supports language-aware comment parsing for Go, JavaScript, TypeScript, Python, Ruby, Shell, CSS, HTML, Rust, YAML, and TOML.

```bash
# List all TODO/FIXME comments
taskmd todos list

# Scan a specific directory
taskmd todos list --dir ./src

# Filter by marker type
taskmd todos list --marker TODO --marker FIXME

# Include only specific file patterns
taskmd todos list --include "*.go"

# Exclude specific file patterns
taskmd todos list --exclude "*.test.go"

# JSON output
taskmd todos list --format json

# Rich output with scope and git blame info
taskmd todos list --rich
```

**Flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--dir` | `.` | Directory to scan for source code |
| `--marker` | *(all)* | Filter by marker type (repeatable) |
| `--include` | | Include only files matching glob pattern (repeatable) |
| `--exclude` | | Exclude files matching glob pattern (repeatable) |
| `--format` | `table` | Output format (`table`, `json`, `yaml`) |
| `--rich` | `false` | Include scope and git blame information (slower) |
| `--raw-text` | `false` | Include original source line text in output |

Exclude patterns can also be configured in `.taskmd.yaml` under `todos.exclude`. CLI `--exclude` flags are additive with config patterns.

## Global Flags

Available for all commands:

```bash
--config string       # Config file path
-d, --task-dir string # Task directory to scan (default ".")
--format string       # Output format (table, json, yaml)
--verbose             # Verbose logging
--quiet               # Suppress non-essential output
--stdin               # Read from stdin instead of files
--debug               # Enable debug output (prints to stderr)
--no-color            # Disable colored output
```

## Common Workflows

### Daily Task Management

```bash
# Morning: What should I work on?
taskmd next --limit 5

# Check project status
taskmd stats

# During work: Validate changes
taskmd validate

# End of day: What got done?
taskmd list --filter status=completed --sort created
```

### Weekly Planning

```bash
# Visual overview
taskmd board --group-by priority

# Identify bottlenecks
taskmd graph --exclude-status completed --format ascii

# Focus on priorities
taskmd list --filter status=pending --sort priority
```

### CI/CD Integration

```bash
# Validate in CI pipeline
if ! taskmd validate tasks/ --strict; then
    echo "Task validation failed"
    exit 1
fi

# Generate snapshot artifact
taskmd snapshot tasks/ --derived --out task-snapshot.json
```

### Scripting and Piping

```bash
# Pipe between commands
taskmd list --format json | jq '.[] | select(.priority == "high")'

# Validate from stdin
echo '---
id: "999"
title: "Test task"
status: pending
---
# Test' | taskmd validate --stdin

# Find quick wins
taskmd list tasks/ \
  --filter status=pending \
  --filter priority=high \
  --filter effort=small \
  --format json | jq -r '.[] | "\(.id): \(.title)"'
```

## Environment Variables

taskmd supports environment variables with the `TASKMD_` prefix:

```bash
export TASKMD_DIR=./tasks
export TASKMD_VERBOSE=true
```

Environment variables have lower precedence than config files and CLI flags.

## Troubleshooting

### "No tasks found"

1. Check that your tasks directory exists: `ls -la tasks/`
2. Ensure files have `.md` extension
3. Verify YAML frontmatter format
4. Run `taskmd validate tasks/` for specific errors
5. Try verbose output: `taskmd list tasks/ --verbose`

### "Invalid task format"

- Check YAML frontmatter is properly formatted
- Ensure required fields are present: `id`, `title`
- Verify status is valid: `pending`, `in-progress`, `completed`, `blocked`, `cancelled`
- Run `taskmd validate tasks/` for line-level error messages

### "Circular dependency detected"

Dependencies form a cycle (A depends on B, B depends on A). Use `taskmd graph --format ascii` to visualize the cycle and remove one dependency to break it.

### Command not found

```bash
# Check installation
which taskmd
taskmd --version

# If not found, add to PATH
export PATH=$PATH:$(go env GOPATH)/bin
```
