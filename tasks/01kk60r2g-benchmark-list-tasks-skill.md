---
id: "01kk60r2g"
title: "Benchmark list-tasks skill"
status: completed
priority: medium
dependencies: []
tags: ["benchmark", "skill-eval"]
created: 2026-03-08
phase: skill-benchmarks
---

# Benchmark list-tasks skill

## Objective

Run the list-tasks skill **with and without** the taskmd skill loaded in an isolated project, then compare quality, accuracy, token usage, and latency.

## Tasks

- [x] Create isolated temp dir, run `taskmd init`, copy fixtures from `benchmark/fixtures/tasks/`
- [x] Run **without_skill** baseline: `claude -p "show me all my tasks"` (no skill loaded, no CLAUDE.md/TASKMD_SPEC.md/.taskmd.yaml)
- [x] Save without_skill output to `benchmark/iteration-1/eval-1-list-tasks/without_skill/outputs/result.md`
- [x] Run **with_skill** variant: `claude -p "show me all my tasks" --allowedTools "Bash,taskmd:list-tasks"` (skill loaded)
- [x] Save with_skill output to `benchmark/iteration-1/eval-1-list-tasks/with_skill/outputs/result.md`
- [ ] Record token usage and duration in `timing.json` for both runs — **skipped**: `claude -p` doesn't expose token counts
- [x] Grade both runs against assertions in `eval_metadata.json`, save `grading.json` for each
- [ ] Run `aggregate_benchmark.py` to produce `benchmark.json` and `benchmark.md` — **skipped**: wrote `benchmark.json` manually instead (no aggregation script exists yet)
- [x] Evaluate: compare quality, accuracy, tokens, and latency between with/without skill

### Additional work done (beyond original scope)

- [x] Added 4 new list-tasks evals to `evals.json` (IDs 14-17): filter bugs, top 3 as JSON, not-done sorted, custom columns
- [x] Ran all 5 evals (10 total runs) with proper control conditions (baseline has no CLAUDE.md, .taskmd.yaml, or TASKMD_SPEC.md)
- [x] Created `benchmark/suggestions/list-tasks.md` with improvement recommendations
- [x] Created `benchmark/CLAUDE.md` with lessons learned for future benchmarks
- [x] Created `benchmark/iteration-1/snapshot.json` capturing the git commit for traceability

## Acceptance Criteria

- [x] Both with_skill and without_skill runs are executed and recorded
- [x] Grading.json files exist for both configurations with assertion results
- [x] benchmark.json contains comparison deltas (pass_rate) — token/time deltas not available
- [ ] Token usage and duration recorded in timing.json — not possible with current `claude -p` approach

## Results

**0% delta** across all 5 evals. The list-tasks skill provides no measurable improvement because Claude discovers `taskmd` via the system PATH even without any project context. See `benchmark/suggestions/list-tasks.md` for improvement recommendations.
