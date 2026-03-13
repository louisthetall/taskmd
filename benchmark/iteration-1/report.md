# list-tasks Benchmark — Iteration 1

**Commit:** `6878c4c` ("prepare benchmarks")
**Date:** 2026-03-13

## Test Conditions

| Config | Setup |
|--------|-------|
| **with_skill** | `taskmd init` project + CLAUDE.md + .taskmd.yaml + TASKMD_SPEC.md + `taskmd:list-tasks` skill + `taskmd` on PATH |
| **without_skill** | Bare project — only raw task `.md` files. No CLAUDE.md, no config, no spec, no skills, **`taskmd` blocked from PATH** |

## Results

| Eval | Prompt | with_skill | without_skill | Delta |
|------|--------|:----------:|:-------------:|:-----:|
| 1 | "show me all my tasks" | 100% | 100% | 0% |
| 14 | "which bugs still need fixing?" | 67% | 67% | 0% |
| 15 | "top 3 highest priority tasks as json" | 75% | 75% | 0% |
| 16 | "stuff not done yet, sorted by priority" | 100% | 100% | 0% |
| 17 | "just the id and title columns" | 100% | 100% | 0% |
| **Mean** | | **88%** | **88%** | **0%** |

## Timing & Cost

| Eval | with_skill | without_skill | Delta |
|------|-----------|---------------|-------|
| 1 | 24.0s / 816 tok / $0.139 | 17.4s / 843 tok / $0.121 | +6.6s / +$0.018 |
| 14 | 22.9s / 886 tok / $0.127 | 13.5s / 537 tok / $0.095 | +9.4s / +$0.032 |
| 15 | 19.8s / 964 tok / $0.128 | 15.2s / 828 tok / $0.109 | +4.5s / +$0.019 |
| 16 | 17.1s / 769 tok / $0.117 | 23.3s / 770 tok / $0.109 | **-6.2s** / +$0.008 |
| 17 | 14.1s / 446 tok / $0.100 | 15.6s / 428 tok / $0.098 | **-1.4s** / +$0.002 |
| **Mean** | **19.6s / 776 tok / $0.122** | **17.0s / 681 tok / $0.106** | **+2.6s / +$0.016** |

## Analysis

**The list-tasks skill provides zero quality improvement and is slightly slower and more expensive than baseline.**

### Quality
Pass rates are identical (88%) across all 5 evals. Both configs produce correct, well-formatted output. Even without `taskmd` on PATH, Claude reads the raw markdown files, parses YAML frontmatter, and produces equivalent results.

### Performance
- **Duration**: with-skill averages 19.6s vs 17.0s for baseline (+15%)
- **Turns**: with-skill uses more turns (8.2 vs 7.2) — it reads CLAUDE.md and TASKMD_SPEC.md as extra context
- **Cost**: with-skill costs ~$0.016 more per query ($0.122 vs $0.106)
- **Exception**: Evals 16 and 17 were slightly faster with the skill, but this is likely variance

### Key Insight
The skill's overhead (loading CLAUDE.md, TASKMD_SPEC.md, .taskmd.yaml into context) costs more tokens than the value it adds. For 5 task files, Claude's native file-reading ability is sufficient — and cheaper.

### Failed Assertions (both configs)
- `uses-filter` (eval 14) — neither used `--filter type=bug`
- `uses-sort` (eval 15) — unclear if `--sort priority` was used
- `uses-limit` (eval 15) — neither used `--limit 3`

These are implementation-specific assertions that should be replaced with output-focused checks.

## Recommendations

1. **Drop implementation-specific assertions** — focus on output correctness
2. **Test at scale** (50-100 tasks) to find where `taskmd list` outperforms file-by-file reading
3. **Reduce skill overhead** — the skill loads unnecessary context that slows things down
4. **Add unique value** — the skill should teach flag combinations, suggest follow-ups, or handle queries Claude can't do natively

## Files

- `benchmark.json` — machine-readable results with timing
- `snapshot.json` — git commit metadata
- `eval-{1,14,15,16,17}-list-tasks/` — per-eval outputs, timing, and raw stream-json logs
