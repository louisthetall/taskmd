# Benchmarking Guide

Lessons learned from running skill benchmarks. Reference this before running future evals.

## Control Case Setup

The without-skill (baseline) run must be a clean environment with **no access to taskmd at all** — no CLI, no config, no docs. The goal is to measure what Claude does with only raw markdown files.

### Remove all taskmd context from the project

- **Remove** `tasks/CLAUDE.md` — contains taskmd command references
- **Remove** `tasks/TASKMD_SPEC.md` — documents task format and CLI usage
- **Remove** `.taskmd.yaml` — config file that signals a taskmd project
- **Remove** `.taskmd/` directory — templates that hint at taskmd structure
- **Keep** only the raw task `.md` files in `tasks/`

### Hide taskmd from PATH

Claude will discover `taskmd` if it's on PATH, which defeats the purpose of the control. Use a shadow directory that overrides `taskmd` with a "not found" stub:

```bash
# Create a shadow dir with a stub that blocks taskmd
SHADOW_DIR=$(mktemp -d)
echo '#!/bin/sh
echo "taskmd: command not found" >&2; exit 127' > "$SHADOW_DIR/taskmd"
chmod +x "$SHADOW_DIR/taskmd"

# Prepend shadow dir so the stub takes priority over the real binary
unset CLAUDECODE && PATH="$SHADOW_DIR:$PATH" claude -p "<prompt>" --allowedTools "Bash" --permission-mode acceptEdits
```

**Why not strip the directory from PATH?** Because `claude` and `taskmd` may live in the same directory (e.g. `/opt/homebrew/bin`). Removing that directory would also hide `claude` itself (exit code 127).

Without this, Claude finds `taskmd` on its own and the baseline measures nothing.

## Running Evals

Use `benchmark/run_eval.sh` to run each eval. It handles `claude -p` invocation with `--output-format stream-json --verbose` and extracts timing/token/cost metrics automatically.

```bash
# With-skill run
bash benchmark/run_eval.sh <project-dir> "<prompt>" <output-dir> --allowedTools "Bash,taskmd:<skill>"

# Without-skill run (with taskmd blocked)
PATH="$SHADOW_DIR:$PATH" bash benchmark/run_eval.sh <project-dir> "<prompt>" <output-dir> --allowedTools "Bash"
```

The script outputs:
- `result.md` — the model's text response
- `timing.json` — duration_ms, output_tokens, total_cost_usd, num_turns
- `raw_output.jsonl` — full stream-json for debugging

## Running `claude -p` from Claude Code

- **Nested session block**: `claude -p` fails inside Claude Code with "cannot be launched inside another Claude Code session." The `run_eval.sh` script handles this by unsetting `CLAUDECODE`.
- **Timing**: `date +%s%3N` doesn't work in zsh on macOS. Use `--output-format stream-json --verbose` instead — the final `result` event includes `duration_ms`, `total_cost_usd`, and full `usage` with token breakdowns.
- **Background runs**: `claude -p` with variable capture can silently fail in background tasks. The `run_eval.sh` script works in both foreground and background.
- **Subagents can't run claude -p**: Subagents spawned via the Agent tool don't have Bash access by default. Run evals from the main session.

## Assertions

- **Focus on output quality, not implementation**: Don't assert "uses --filter flag" — assert "output contains only bug-type tasks." Claude may achieve the same result through different means.
- **Non-discriminating assertions are noise**: If both with-skill and without-skill pass an assertion equally, it doesn't measure skill value. Flag these in the analysis.
- **Test what the skill uniquely enables**: Assertions should target behaviors the skill teaches that Claude wouldn't do on its own.

## Iteration Report

Every iteration must produce a `report.md` at `iteration-N/report.md`. Generate it after all evals are graded. It should include:

1. **Header** — skill name, iteration number, commit (from snapshot.json), date
2. **Test conditions table** — what with_skill and without_skill setups include
3. **Results table** — one row per eval with prompt, pass rates for both configs, and delta
4. **Assertion detail** — per-eval breakdown showing each assertion's pass/fail for both configs
5. **Analysis** — headline finding, non-discriminating assertions, failed assertions with explanation
6. **Recommendations** — actionable next steps for improving the skill or the benchmark
7. **Files** — links to related files (benchmark.json, snapshot.json, eval dirs)

See `iteration-1/report.md` for the reference format.

## Iteration Snapshots

Every iteration must include a `snapshot.json` capturing the git commit and metadata at the time of the run. Create it at the start of each iteration:

```json
{
  "commit": "<full sha>",
  "short": "<short sha>",
  "subject": "<commit message>",
  "date": "<ISO 8601>",
  "author": "<author>",
  "branch": "<branch>",
  "skill_path": "<path to skills being tested>",
  "notes": "<brief summary of what changed since last iteration>"
}
```

Generate with: `git log -1 --format='{"commit": "%H", "short": "%h", "subject": "%s", "date": "%ci", "author": "%an"}'`

## Confounding Factors

- **Globally installed CLI tools**: If `taskmd` is on PATH, Claude discovers it without any skill. Always strip `taskmd` from PATH in baseline runs (see Control Case Setup above). Iteration 1 did NOT do this, which is why it showed 0% delta.
- **Model knowledge**: Claude may know about common tools from training data. A truly novel tool would show more skill lift than a well-known one.
- **Single-run variance**: LLM outputs are non-deterministic. A single run per eval is directional but not statistically significant. For rigorous benchmarks, run 3+ times per eval and report mean +/- stddev.

## Directory Structure

```
benchmark/
├── CLAUDE.md                # this file
├── README.md                # overview
├── evals.json               # eval definitions
├── suggestions/             # per-skill improvement notes
│   └── <skill-name>.md
├── fixtures/
│   ├── setup.sh
│   ├── tasks/
│   └── src/
└── iteration-N/
    ├── snapshot.json         # git commit & metadata for this iteration
    ├── report.md             # human-readable summary of results
    ├── benchmark.json        # machine-readable results
    └── eval-{id}-{skill}/
        ├── eval_metadata.json
        ├── with_skill/outputs/
        │   └── result.md
        └── without_skill/outputs/
            └── result.md
```
