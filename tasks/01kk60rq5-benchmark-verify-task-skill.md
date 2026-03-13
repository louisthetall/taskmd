---
id: "01kk60rq5"
title: "Benchmark verify-task skill"
status: pending
priority: medium
dependencies: []
tags: ["benchmark", "skill-eval"]
created: 2026-03-08
phase: skill-benchmarks
---

# Benchmark verify-task skill

## Objective

Benchmark the verify-task skill by running it with and without the skill loaded, comparing quality, timing, token usage, and cost.

## Prerequisites

- Read `benchmark/CLAUDE.md` for methodology, control case setup, and known pitfalls
- Use `benchmark/run_eval.sh` for all eval runs (handles stream-json, timing extraction)
- Reference `benchmark/evals.json` for the eval prompts and assertions

## Tasks

- [ ] Set up isolated projects using `benchmark/fixtures/setup.sh`
  - **with_skill**: full `taskmd init` project (CLAUDE.md, .taskmd.yaml, TASKMD_SPEC.md present)
  - **without_skill**: bare project — remove CLAUDE.md, TASKMD_SPEC.md, .taskmd.yaml, .taskmd/; block `taskmd` from PATH using shadow dir
- [ ] Run with_skill eval using `benchmark/run_eval.sh`:
  ```
  bash benchmark/run_eval.sh <project-dir> "verify the acceptance criteria for task 003" benchmark/iteration-1/eval-10-verify-task/with_skill/outputs --allowedTools "Bash,Skill,Read,Glob,Grep,taskmd:verify-task"
  ```
- [ ] Run without_skill baseline using `benchmark/run_eval.sh` with taskmd blocked:
  ```
  PATH="$SHADOW_DIR:$PATH" bash benchmark/run_eval.sh <project-dir> "verify the acceptance criteria for task 003" benchmark/iteration-1/eval-10-verify-task/without_skill/outputs --allowedTools "Bash,Read,Glob,Grep"
  ```
- [ ] Write `eval_metadata.json` with assertions from `evals.json`
- [ ] Grade both outputs against assertions
- [ ] Write `benchmark.json` with pass rates, timing deltas, token/cost comparison
- [ ] Write `report.md` summarizing results (see `benchmark/iteration-1/report.md` for format)
- [ ] Write improvement suggestions to `benchmark/suggestions/verify-task.md`

## Acceptance Criteria

- Both with_skill and without_skill runs are executed and saved to `benchmark/iteration-1/`
- `timing.json` exists for both runs with duration_ms, output_tokens, total_cost_usd
- `benchmark.json` contains pass rate deltas AND timing/cost comparison
- `report.md` exists with results table, timing table, analysis, and recommendations
- `benchmark/suggestions/verify-task.md` written with improvement ideas
