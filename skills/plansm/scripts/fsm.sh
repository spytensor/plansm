#!/bin/bash
# fsm.sh - State machine operations for plansm
# Usage: fsm.sh <command> [plan.json] [args...]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

COMMAND="$1"
PLAN_FILE="${2:-plan.json}"

if [ ! -f "$PLAN_FILE" ]; then
    echo "Error: Plan file not found: $PLAN_FILE" >&2
    exit 1
fi

# Check dependencies
if ! command -v jq &> /dev/null; then
    echo "Error: jq is required. Install: brew install jq" >&2
    exit 1
fi

# Unlock steps whose dependencies are VERIFIED
unlock_ready() {
    local tmp_file=$(mktemp)

    # Read current plan
    local plan=$(cat "$PLAN_FILE")

    # For each LOCKED step, check if dependencies are met
    local steps=$(echo "$plan" | jq -c '.steps[] | select(.status == "LOCKED")')

    while IFS= read -r step; do
        local step_id=$(echo "$step" | jq -r '.id')
        local depends_on=$(echo "$step" | jq -r '.depends_on // [] | .[]')

        local all_deps_met=1
        for dep_id in $depends_on; do
            local dep_status=$(echo "$plan" | jq -r ".steps[] | select(.id == \"$dep_id\") | .status")
            if [ "$dep_status" != "VERIFIED" ]; then
                all_deps_met=0
                break
            fi
        done

        if [ "$all_deps_met" = "1" ]; then
            # Unlock this step
            plan=$(echo "$plan" | jq "(.steps[] | select(.id == \"$step_id\") | .status) = \"PENDING\"")
        fi
    done <<< "$steps"

    echo "$plan" > "$tmp_file"
    mv "$tmp_file" "$PLAN_FILE"
}

# Advance to next PENDING/FAILED step
advance_step() {
    local current_step=$(bash "$SCRIPT_DIR/parse-plan.sh" "$PLAN_FILE" current_step)

    # Check that current step is VERIFIED
    local current_status=$(jq -r ".steps[] | select(.id == \"$current_step\") | .status" "$PLAN_FILE")
    if [ "$current_status" != "VERIFIED" ]; then
        echo "Error: Current step $current_step is not VERIFIED (status=$current_status)" >&2
        echo "Run verify.sh --current first" >&2
        exit 1
    fi

    # Unlock ready steps
    unlock_ready

    # Find current step index
    local current_idx=$(jq "[.steps[] | .id] | index(\"$current_step\")" "$PLAN_FILE")

    # Find next PENDING/FAILED step after current
    local next_step=$(jq -r ".steps[$((current_idx + 1)):] | .[] | select(.status == \"PENDING\" or .status == \"FAILED\") | .id" "$PLAN_FILE" | head -1)

    # If no next step after current, wrap around to beginning
    if [ -z "$next_step" ]; then
        next_step=$(jq -r '.steps[] | select(.status == "PENDING" or .status == "FAILED") | .id' "$PLAN_FILE" | head -1)
    fi

    if [ -z "$next_step" ]; then
        echo "No more pending steps. Plan complete."
        exit 0
    fi

    # Update current_step
    local tmp_file=$(mktemp)
    jq ".current_step = \"$next_step\"" "$PLAN_FILE" > "$tmp_file"
    mv "$tmp_file" "$PLAN_FILE"

    echo "advanced to $next_step"
}

# Show status table
show_status() {
    unlock_ready

    local current_step=$(bash "$SCRIPT_DIR/parse-plan.sh" "$PLAN_FILE" current_step)

    echo "current_step: $current_step"
    echo ""
    printf "%-12s %-10s %s\n" "STEP" "STATUS" "OBJECTIVE"
    printf "%s\n" "--------------------------------------------------------------------------------"

    jq -r '.steps[] | "\(.id)\t\(.status)\t\(.objective)"' "$PLAN_FILE" | while IFS=$'\t' read -r id status obj; do
        # Truncate objective if too long
        if [ ${#obj} -gt 55 ]; then
            obj="${obj:0:55}…"
        fi
        printf "%-12s %-10s %s\n" "$id" "$status" "$obj"
    done
}

# Show current step details
show_current() {
    unlock_ready

    local current_step=$(bash "$SCRIPT_DIR/parse-plan.sh" "$PLAN_FILE" current_step)
    local step=$(jq ".steps[] | select(.id == \"$current_step\")" "$PLAN_FILE")

    echo "CURRENT_STEP: $current_step"
    echo "STATUS: $(echo "$step" | jq -r '.status')"
    echo "OBJECTIVE: $(echo "$step" | jq -r '.objective')"

    local allow_paths=$(echo "$step" | jq -r '.allow_paths // [] | .[]')
    if [ -n "$allow_paths" ]; then
        echo "ALLOW_PATHS:"
        echo "$allow_paths" | while read -r path; do
            echo "  - $path"
        done
    fi

    echo "VERIFY:"
    echo "$step" | jq -r '.verify[] | "  - \(.type): \(.cmd // .file // .url // "")"'
}

# Main command dispatch
case "$COMMAND" in
    status)
        show_status
        ;;
    current)
        show_current
        ;;
    advance)
        advance_step
        ;;
    unlock)
        unlock_ready
        echo "Unlocked ready steps"
        ;;
    *)
        echo "Usage: fsm.sh <command> [plan.json]" >&2
        echo "" >&2
        echo "Commands:" >&2
        echo "  status    Show status table" >&2
        echo "  current   Show current step details" >&2
        echo "  advance   Advance to next step" >&2
        echo "  unlock    Unlock steps whose dependencies are met" >&2
        exit 1
        ;;
esac
