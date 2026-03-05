---
name: divide-and-conquer
description: Pick up a task and execute it using subagents to parallelize independent workstreams. Use when the user wants to work on a task with maximum concurrency.
allowed-tools: Bash, Read, Glob, Grep, Write, Edit, Task, Agent, EnterPlanMode
---

# Divide and Conquer

Pick up a task and execute it by splitting the work into independent workstreams that run in parallel via subagents.

## Instructions

The user's query is in `$ARGUMENTS` (a task ID like `077` or a task name/keyword).

1. **Look up the task**: Run `taskmd get $ARGUMENTS` to find the task
   - If not found, run `taskmd list` to show available tasks and ask the user which one they meant
2. **Read the task file** with the `Read` tool to get the full description, subtasks, and acceptance criteria
3. **Mark the task as in-progress**: Run `taskmd set <ID> --status in-progress`
4. **Start a worklog entry** (if worklogs are enabled):
   - Check `.taskmd.yaml` for `worklogs: false` -- if set, skip worklog steps
   - Otherwise, find or create the worklog file at `tasks/<group>/.worklogs/<ID>.md` (or `tasks/.worklogs/<ID>.md` for root tasks)
   - Append a timestamped entry noting your approach and initial findings
5. **Plan and identify workstreams**:
   - Use `EnterPlanMode` to design the overall approach
   - In the plan, include a reference to the original task ID and task file path
   - Analyze the task and break it into **independent workstreams** — pieces of work that can proceed in parallel without depending on each other's output
   - Examples of independent workstreams:
     - Implementation code vs. tests vs. documentation
     - Changes to separate packages or modules
     - Backend changes vs. frontend changes
   - If the task is simple enough that parallelization adds no benefit, just do it directly (skip to step 7)
6. **Launch subagents in parallel**:
   - Use the `Agent` tool to launch one subagent per independent workstream
   - Give each subagent a clear, self-contained prompt describing exactly what to do, including relevant file paths and context
   - Launch all independent subagents in a **single message** so they run concurrently
   - Use `isolation: "worktree"` for subagents that modify files, to avoid conflicts
   - Wait for all subagents to complete
7. **Coordinate and integrate**:
   - Review all subagent results for correctness
   - If subagents ran in worktrees, merge their changes (review diffs, resolve any conflicts)
   - If any subagent failed, handle the failure directly rather than re-launching
   - Run tests and linting to verify the integrated result
   - Check off subtasks (`- [x]`) in the task file as they are completed
   - Append worklog entries when you make key decisions, hit blockers, or complete significant subtasks
8. **Write a final worklog entry** summarizing what was done, which workstreams ran in parallel, decisions made, and any open items
9. **Mark the task as done**:
   - Check `.taskmd.yaml` for `workflow: pr-review` -- if set, use the PR-review workflow below
   - **Solo workflow** (default): Run `taskmd set <ID> --status completed --verify`
     - The `--verify` flag will run any verification checks defined in the task before applying the status change
     - If verification fails, fix the issues and try again
   - **PR-review workflow**: Open a PR, then run `taskmd set <ID> --status in-review --add-pr <PR-URL>` and stop

## Worklog Format

Each worklog entry uses a timestamp heading followed by free-form notes:

```markdown
## 2026-02-15T10:30:00Z

Started divide-and-conquer execution of the search feature task.

**Workstreams identified:**

1. Core search implementation (subagent — worktree)
2. Test suite (subagent — worktree)
3. Documentation updates (subagent)

**Completed:**

- [x] All subagents finished successfully
- [x] Merged worktree changes
- [x] Tests passing after integration

**Decisions:** Used full-text search with SQLite rather than Elasticsearch.
```
