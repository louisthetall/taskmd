#!/usr/bin/env bash
# Tests for the task ID extraction logic used in the taskmd-complete action.
# Run: bash .github/actions/taskmd-complete/test-extract.sh
#
# Note: The action itself uses grep -oP (GNU grep, available on Ubuntu runners).
# This test uses bash regex for portability across macOS and Linux.

set -euo pipefail

PASS=0
FAIL=0
PATTERN='task[: ]+([^ ]+)'

extract_task_id() {
  local pr_body="${1:-}"
  local pr_branch="${2:-}"
  local task_id=""

  # Try PR body first (line by line)
  if [[ -n "$pr_body" ]]; then
    while IFS= read -r line; do
      if [[ "$line" =~ $PATTERN ]]; then
        task_id="${BASH_REMATCH[1]}"
        break
      fi
    done <<< "$pr_body"
  fi
  # Fall back to branch name
  if [[ -z "$task_id" && -n "$pr_branch" ]]; then
    if [[ "$pr_branch" =~ $PATTERN ]]; then
      task_id="${BASH_REMATCH[1]}"
    fi
  fi
  echo "$task_id"
}

assert_eq() {
  local test_name="$1"
  local expected="$2"
  local actual="$3"
  if [[ "$expected" == "$actual" ]]; then
    echo "  PASS: $test_name"
    PASS=$((PASS + 1))
  else
    echo "  FAIL: $test_name (expected '$expected', got '$actual')"
    FAIL=$((FAIL + 1))
  fi
}

echo "=== taskmd-complete extraction tests ==="
echo ""

# --- PR body extraction ---
echo "PR body extraction:"
assert_eq "task: 042 in body" "042" "$(extract_task_id "task: 042" "")"
assert_eq "task:042 no space" "042" "$(extract_task_id "task:042" "")"
assert_eq "task: cli-049" "cli-049" "$(extract_task_id "task: cli-049" "")"
assert_eq "multiline body" "015" "$(extract_task_id $'Some description\ntask: 015\nMore text' "")"
assert_eq "no match in body" "" "$(extract_task_id "just a regular PR description" "")"
echo ""

# --- Branch name extraction ---
echo "Branch name extraction:"
assert_eq "task:042 in branch" "042-some-feature" "$(extract_task_id "" "task:042-some-feature")"
assert_eq "no match in branch" "" "$(extract_task_id "" "feature/some-thing")"
echo ""

# --- Fallback from body to branch ---
echo "Fallback behavior:"
assert_eq "body match preferred" "001" "$(extract_task_id "task: 001" "task:002-branch")"
assert_eq "falls back to branch" "002-branch" "$(extract_task_id "just a description" "task:002-branch")"
assert_eq "neither matches" "" "$(extract_task_id "no match" "feature/some-branch")"
echo ""

# --- Action YAML validation ---
echo "Action YAML validation:"
ACTION_FILE="$(dirname "$0")/action.yml"
if [[ -f "$ACTION_FILE" ]]; then
  # Verify required structure
  if grep -q "using: 'composite'" "$ACTION_FILE"; then
    echo "  PASS: action is composite type"
    PASS=$((PASS + 1))
  else
    echo "  FAIL: action should be composite type"
    FAIL=$((FAIL + 1))
  fi

  # Verify dry-run input exists
  if grep -q "dry-run:" "$ACTION_FILE"; then
    echo "  PASS: dry-run input defined"
    PASS=$((PASS + 1))
  else
    echo "  FAIL: dry-run input missing"
    FAIL=$((FAIL + 1))
  fi

  # Verify all bash steps set shell explicitly
  BASH_STEPS=$(grep -c "shell: bash" "$ACTION_FILE" || true)
  if [[ "$BASH_STEPS" -ge 4 ]]; then
    echo "  PASS: all steps specify shell: bash ($BASH_STEPS steps)"
    PASS=$((PASS + 1))
  else
    echo "  FAIL: expected at least 4 shell: bash declarations, found $BASH_STEPS"
    FAIL=$((FAIL + 1))
  fi

  # Verify env vars are used instead of inline expressions in run blocks (security)
  INLINE_IN_RUN=$(grep -A1 'run: |' "$ACTION_FILE" | grep -c '\${{' || true)
  if [[ "$INLINE_IN_RUN" -eq 0 ]]; then
    echo "  PASS: no inline \${{ }} expressions in run blocks (uses env vars)"
    PASS=$((PASS + 1))
  else
    echo "  FAIL: found \${{ }} in run blocks — use env vars for security"
    FAIL=$((FAIL + 1))
  fi

  # Verify dry-run gates destructive steps
  DRY_RUN_GATES=$(grep -c "inputs.dry-run != 'true'" "$ACTION_FILE" || true)
  if [[ "$DRY_RUN_GATES" -ge 3 ]]; then
    echo "  PASS: dry-run gates destructive steps ($DRY_RUN_GATES gates)"
    PASS=$((PASS + 1))
  else
    echo "  FAIL: expected at least 3 dry-run gates, found $DRY_RUN_GATES"
    FAIL=$((FAIL + 1))
  fi
else
  echo "  SKIP: action.yml not found at $ACTION_FILE"
fi
echo ""

# --- Summary ---
echo "=== Results: $PASS passed, $FAIL failed ==="
if [[ "$FAIL" -gt 0 ]]; then
  exit 1
fi
