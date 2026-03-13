#!/bin/bash
# Run a single benchmark eval and extract metrics from stream-json output.
# Usage: ./run_eval.sh <project-dir> <prompt> <output-dir> [extra-claude-args...]
#
# Outputs:
#   <output-dir>/result.md    — the model's text response
#   <output-dir>/timing.json  — duration, tokens, cost extracted from stream-json

set -e

PROJECT_DIR="$1"
PROMPT="$2"
OUTPUT_DIR="$3"
shift 3
EXTRA_ARGS=("$@")

mkdir -p "$OUTPUT_DIR"

# Run claude -p with stream-json to capture metrics
unset CLAUDECODE
RAW=$(cd "$PROJECT_DIR" && claude -p "$PROMPT" \
  --permission-mode acceptEdits \
  --output-format stream-json \
  --verbose \
  "${EXTRA_ARGS[@]}" 2>&1)

# Extract the final result line (type=result)
RESULT_LINE=$(echo "$RAW" | grep '"type":"result"' | tail -1)

if [ -z "$RESULT_LINE" ]; then
  echo "ERROR: No result line found in output" >&2
  echo "$RAW" > "$OUTPUT_DIR/raw_output.jsonl"
  exit 1
fi

# Extract text result
echo "$RESULT_LINE" | python3 -c "
import sys, json
data = json.load(sys.stdin)
print(data.get('result', ''))
" > "$OUTPUT_DIR/result.md"

# Extract timing/usage metrics
echo "$RESULT_LINE" | python3 -c "
import sys, json
data = json.load(sys.stdin)
usage = data.get('usage', {})
timing = {
    'duration_ms': data.get('duration_ms'),
    'duration_api_ms': data.get('duration_api_ms'),
    'num_turns': data.get('num_turns'),
    'total_cost_usd': data.get('total_cost_usd'),
    'input_tokens': usage.get('input_tokens'),
    'output_tokens': usage.get('output_tokens'),
    'cache_read_input_tokens': usage.get('cache_read_input_tokens'),
    'cache_creation_input_tokens': usage.get('cache_creation_input_tokens'),
}
json.dump(timing, sys.stdout, indent=2)
print()
" > "$OUTPUT_DIR/timing.json"

# Save full stream for debugging
echo "$RAW" > "$OUTPUT_DIR/raw_output.jsonl"

echo "Done: $(cat "$OUTPUT_DIR/timing.json" | python3 -c "import sys,json; d=json.load(sys.stdin); print(f\"{d['duration_ms']}ms, {d['output_tokens']} output tokens, \${d['total_cost_usd']:.4f}\")")"
