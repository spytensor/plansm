#!/bin/bash
# Install plansm git hooks

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
HOOKS_DIR="$PROJECT_ROOT/hooks"
GIT_HOOKS_DIR="$PROJECT_ROOT/.git/hooks"

echo "Installing plansm git hooks..."

# Check if we're in a git repository
if [ ! -d "$GIT_HOOKS_DIR" ]; then
    echo "Error: Not in a git repository"
    echo "Run this script from the project root"
    exit 1
fi

# Install pre-commit hook
if [ -f "$HOOKS_DIR/pre-commit" ]; then
    cp "$HOOKS_DIR/pre-commit" "$GIT_HOOKS_DIR/pre-commit"
    chmod +x "$GIT_HOOKS_DIR/pre-commit"
    echo "✅ Installed pre-commit hook"
else
    echo "⚠️  pre-commit hook not found in $HOOKS_DIR"
fi

echo ""
echo "Git hooks installed successfully!"
echo ""
echo "These hooks will:"
echo "  - Warn when plan.json status fields are manually edited"
echo "  - Encourage using 'plansm verify' and 'plansm advance' instead"
echo ""
echo "To bypass hooks (not recommended), use: git commit --no-verify"
