# list-tasks: Benchmark Findings & Improvement Suggestions

## Benchmark Results (iteration-1)

5 evals, 0% delta between with-skill and without-skill across all prompts.

The skill currently adds no measurable value. Claude discovers `taskmd list` on its own via the system PATH and produces equivalent output — even without CLAUDE.md or any project context in the baseline.

## Why It Doesn't Help

The skill's instructions are essentially "run `taskmd list $ARGUMENTS`". Claude already does this without the skill. The skill doesn't teach Claude anything it can't figure out by running `taskmd list --help`.

## What Would Make It Valuable

### 1. Translate natural language to flags

The skill should map common user intents to specific flag combinations that Claude might not discover on its own:

- "what's urgent" → `--filter priority=critical --filter status!=completed`
- "my backlog" → `--filter status=pending --sort priority`
- "what did I finish this week" → `--filter status=completed` (+ date filtering if available)
- "quick wins" → `--filter effort=small --filter status=pending --sort priority`

### 2. Enrich output beyond raw CLI

After running the command, the skill could instruct Claude to:

- Highlight actionable items (e.g., "003 is critical and unstarted")
- Show progress stats (X of Y completed, N blocked)
- Suggest what to pick up next based on priority + effort

### 3. Chain with other commands

For filter-heavy queries, the skill could suggest combining `taskmd list` with `taskmd next` or `taskmd graph` to give richer answers.

### 4. Handle ambiguous queries

The skill could provide guidance for mapping vague prompts like "which bugs still need fixing?" to the right filter combination, rather than leaving Claude to guess whether to use `--filter type=bug` or read files manually.

## Assertion Improvements

- Drop implementation-specific assertions like "uses --filter flag" — what matters is the output, not how it got there
- Add assertions for output enrichment: "includes actionable recommendation", "highlights priority items"
- Add assertions for flag efficiency: "uses a single taskmd command rather than reading files individually"
