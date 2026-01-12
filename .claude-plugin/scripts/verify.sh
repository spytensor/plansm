#!/bin/bash
# verify.sh - Shell-based verifier for plansm (zero-install mode)
# Usage: verify.sh [--current|--all] [--json] [--plan plan.json]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PLAN_FILE="plan.json"
MODE="current"
OUTPUT_JSON=0

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --help|-h)
            echo "Usage: verify.sh [--current|--all] [--plan PATH] [--json]"
            echo ""
            echo "Verify steps in plan.json using machine-checkable proofs."
            echo ""
            echo "Options:"
            echo "  --current    Verify only current step (default)"
            echo "  --all        Verify all PENDING/FAILED steps"
            echo "  --plan PATH  Path to plan.json (default: plan.json)"
            echo "  --json       JSON output"
            echo "  --help       Show this help"
            exit 0
            ;;
        --current)
            MODE="current"
            shift
            ;;
        --all)
            MODE="all"
            shift
            ;;
        --plan)
            PLAN_FILE="$2"
            shift 2
            ;;
        --json)
            OUTPUT_JSON=1
            shift
            ;;
        *)
            echo "Unknown option: $1" >&2
            exit 1
            ;;
    esac
done

# Check dependencies
if ! command -v jq &> /dev/null; then
    echo "Error: jq is required. Install: brew install jq" >&2
    exit 1
fi

if [ ! -f "$PLAN_FILE" ]; then
    echo "Error: Plan file not found: $PLAN_FILE" >&2
    exit 1
fi

# Get current step
CURRENT_STEP=$(bash "$SCRIPT_DIR/parse-plan.sh" "$PLAN_FILE" current_step)

# Verify a single rule
verify_rule() {
    local rule_json="$1"
    local rule_type=$(echo "$rule_json" | jq -r '.type')

    case "$rule_type" in
        command)
            verify_command "$rule_json"
            ;;
        file_exists)
            verify_file_exists "$rule_json"
            ;;
        file_contains)
            verify_file_contains "$rule_json"
            ;;
        http)
            verify_http "$rule_json"
            ;;
        *)
            echo "❌ Unknown rule type: $rule_type"
            return 1
            ;;
    esac
}

# Verify command rule
verify_command() {
    local rule_json="$1"
    local cmd=$(echo "$rule_json" | jq -r '.cmd')
    local expect_exit=$(echo "$rule_json" | jq -r '.expect.exit_code // "0"')

    # Run command
    set +e
    output=$(eval "$cmd" 2>&1)
    exit_code=$?
    set -e

    if [ "$exit_code" = "$expect_exit" ]; then
        echo "✓ command: $cmd"
        return 0
    else
        echo "✗ command: $cmd (exit=$exit_code, expected=$expect_exit)"
        return 1
    fi
}

# Verify file_exists rule
verify_file_exists() {
    local rule_json="$1"
    local file=$(echo "$rule_json" | jq -r '.file')

    if [ -f "$file" ] || [ -d "$file" ]; then
        echo "✓ file_exists: $file"
        return 0
    else
        echo "✗ file_exists: $file (not found)"
        return 1
    fi
}

# Verify file_contains rule
verify_file_contains() {
    local rule_json="$1"
    local file=$(echo "$rule_json" | jq -r '.file')
    local pattern=$(echo "$rule_json" | jq -r '.pattern')

    if [ ! -f "$file" ]; then
        echo "✗ file_contains: $file (file not found)"
        return 1
    fi

    if grep -E -q "$pattern" "$file"; then
        echo "✓ file_contains: $file (pattern: $pattern)"
        return 0
    else
        echo "✗ file_contains: $file (pattern not found: $pattern)"
        return 1
    fi
}

# Verify http rule
verify_http() {
    local rule_json="$1"
    local url=$(echo "$rule_json" | jq -r '.url')
    local method=$(echo "$rule_json" | jq -r '.method // "GET"')
    local expect_status=$(echo "$rule_json" | jq -r '.expect.http_status // "200"')

    if ! command -v curl &> /dev/null; then
        echo "✗ http: curl not available"
        return 1
    fi

    set +e
    status=$(curl -s -o /dev/null -w "%{http_code}" -X "$method" "$url")
    set -e

    if [ "$status" = "$expect_status" ]; then
        echo "✓ http: $url (status=$status)"
        return 0
    else
        echo "✗ http: $url (status=$status, expected=$expect_status)"
        return 1
    fi
}

# Verify a step
verify_step() {
    local step_id="$1"
    local step_json=$(jq ".steps[] | select(.id == \"$step_id\")" "$PLAN_FILE")

    if [ -z "$step_json" ]; then
        echo "Error: Step not found: $step_id" >&2
        return 1
    fi

    echo "STEP $step_id:"

    local all_ok=1
    local verify_rules=$(echo "$step_json" | jq -c '.verify[]')

    while IFS= read -r rule; do
        if ! verify_rule "$rule"; then
            all_ok=0
        fi
    done <<< "$verify_rules"

    if [ "$all_ok" = "1" ]; then
        # Update status to VERIFIED
        local tmp_file=$(mktemp)
        jq "(.steps[] | select(.id == \"$step_id\") | .status) = \"VERIFIED\"" "$PLAN_FILE" > "$tmp_file"
        mv "$tmp_file" "$PLAN_FILE"
        echo "→ Status: VERIFIED"
        return 0
    else
        # Update status to FAILED
        local tmp_file=$(mktemp)
        jq "(.steps[] | select(.id == \"$step_id\") | .status) = \"FAILED\"" "$PLAN_FILE" > "$tmp_file"
        mv "$tmp_file" "$PLAN_FILE"
        echo "→ Status: FAILED"
        return 1
    fi
}

# Main logic
if [ "$MODE" = "current" ]; then
    if verify_step "$CURRENT_STEP"; then
        echo ""
        echo "OVERALL: OK"
        exit 0
    else
        echo ""
        echo "OVERALL: FAILED"
        exit 1
    fi
elif [ "$MODE" = "all" ]; then
    # Get all PENDING or FAILED steps
    steps=$(jq -r '.steps[] | select(.status == "PENDING" or .status == "FAILED") | .id' "$PLAN_FILE")

    overall_ok=1
    for step_id in $steps; do
        echo ""
        if ! verify_step "$step_id"; then
            overall_ok=0
            break
        fi
    done

    echo ""
    if [ "$overall_ok" = "1" ]; then
        echo "OVERALL: OK"
        exit 0
    else
        echo "OVERALL: FAILED"
        exit 1
    fi
fi
