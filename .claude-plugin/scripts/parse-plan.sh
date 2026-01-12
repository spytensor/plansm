#!/bin/bash
# parse-plan.sh - Parse plan.json using jq or pure bash
# Usage: parse-plan.sh <plan.json> <field>
# Fields: current_step, version, steps, step:<id>:<field>

set -e

PLAN_FILE="${1:-plan.json}"
FIELD="${2:-current_step}"

if [ ! -f "$PLAN_FILE" ]; then
    echo "Error: Plan file not found: $PLAN_FILE" >&2
    exit 1
fi

# Check if jq is available
if command -v jq &> /dev/null; then
    USE_JQ=1
else
    USE_JQ=0
fi

# Parse using jq (preferred)
if [ "$USE_JQ" = "1" ]; then
    case "$FIELD" in
        current_step)
            jq -r '.current_step' "$PLAN_FILE"
            ;;
        version)
            jq -r '.version' "$PLAN_FILE"
            ;;
        steps)
            jq -c '.steps[]' "$PLAN_FILE"
            ;;
        step:*)
            # Format: step:<id>:<field>
            STEP_ID=$(echo "$FIELD" | cut -d: -f2)
            STEP_FIELD=$(echo "$FIELD" | cut -d: -f3)
            jq -r ".steps[] | select(.id == \"$STEP_ID\") | .$STEP_FIELD" "$PLAN_FILE"
            ;;
        *)
            jq -r ".$FIELD" "$PLAN_FILE"
            ;;
    esac
else
    # Fallback: pure bash parsing (basic, limited)
    case "$FIELD" in
        current_step)
            grep '"current_step"' "$PLAN_FILE" | head -1 | sed 's/.*"current_step"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/'
            ;;
        version)
            grep '"version"' "$PLAN_FILE" | head -1 | sed 's/.*"version"[[:space:]]*:[[:space:]]*\([0-9]*\).*/\1/'
            ;;
        *)
            echo "Error: Field '$FIELD' requires jq. Install jq: brew install jq" >&2
            exit 1
            ;;
    esac
fi
