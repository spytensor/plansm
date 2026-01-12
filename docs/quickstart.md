# Quick Start Guide

Get started with plansm in 2 minutes (zero install!).

## Installation

### Option 1: Zero Install (Recommended)

**No Go, no binary - just copy files!**

```bash
# 1. Clone repo
git clone https://github.com/spytensor/plansm.git

# 2. Copy to your project
cd your-project
cp -r ../plansm/.claude-plugin .

# 3. Install jq (if not already installed)
brew install jq  # macOS
# or: sudo apt-get install jq  # Linux

# That's it! Claude Code commands work immediately.
```

### Option 2: Full CLI (Optional)

For additional features and better performance:

**From Releases:**
```bash
curl -L https://github.com/spytensor/plansm/releases/latest/download/plansm-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m) -o plansm
chmod +x plansm
sudo mv plansm /usr/local/bin/
```

**From Source:**
```bash
git clone https://github.com/spytensor/plansm.git
cd plansm
go build -o plansm ./cmd/plansm
sudo mv plansm /usr/local/bin/
```

## Basic Usage

### Zero Install Mode

```bash
# In Claude Code, use slash commands:
/pwork    # Show current step
/pverify  # Verify with tests
/pstatus  # Show all steps
/pnext    # Advance
```

Or run scripts directly:
```bash
bash .claude-plugin/scripts/fsm.sh current
bash .claude-plugin/scripts/verify.sh --current
bash .claude-plugin/scripts/fsm.sh advance
```

### CLI Mode

If you installed the full CLI:

```bash
plansm init --claude  # One-time setup
plansm current        # Show current step
plansm verify --current
plansm advance
```

Both modes create:
- `plan.json` - Your verifiable state machine
- `.claude/commands/` - Claude Code slash commands (auto-generated)
- `.claude-plugin/` - Plugin structure with shell scripts

### 2. Check Current Status

```bash
plansm status
```

Output:
```
current_step: STEP_001

STEP        STATUS     OBJECTIVE
--------------------------------------------------------------------------------
STEP_001    PENDING    Initialize project skeleton
STEP_002    LOCKED     Add first real proof rules (edit this step in plan.json)
```

### 3. View Current Step (Low Token)

```bash
plansm current
```

Output:
```
CURRENT_STEP: STEP_001
STATUS: PENDING
OBJECTIVE: Initialize project skeleton
VERIFY:
  - command: test -d .
```

### 4. Verify Current Step

```bash
plansm verify --current
```

Output:
```
STEP STEP_001: OK
  ✓ command: test -d .

OVERALL: OK
```

If verification passes, the step status automatically changes to `VERIFIED`.

### 5. Advance to Next Step

```bash
plansm advance
```

Output:
```
advanced to STEP_002
```

This only works if the current step is `VERIFIED`. It also unlocks dependent steps.

### 6. Health Check

```bash
plansm doctor
```

Output:
```
plansm doctor
------------
✅ plan load/validate ok (version=1)
✅ current step ok: STEP_001 (PENDING)
✅ .claude/commands present
✅ .claude-plugin present
```

## Claude Code Integration

### Using Project Commands

Once you've run `plansm init --claude`, you can use these commands in Claude Code:

- `/pwork` - Show current step (low token)
- `/pverify` - Verify current step proofs
- `/pstatus` - Show plan status table
- `/pnext` - Advance to next step

### Setting Up the Stop Hook

For maximum safety, configure a Stop hook:

1. In Claude Code, run `/hooks`
2. Select "Stop" hook
3. Set command to: `plansm verify --current`
4. Mark as "blocking"

This prevents Claude from stopping until verification passes.

## Creating Your First Real Plan

Edit `plan.json` to add real steps:

```json
{
  "version": 1,
  "current_step": "STEP_001",
  "invariants": [
    "Do not mark VERIFIED without running plansm verify.",
    "Only work on current_step unless explicitly unlocked."
  ],
  "steps": [
    {
      "id": "STEP_001",
      "objective": "Set up project structure",
      "status": "PENDING",
      "verify": [
        {"type": "file_exists", "file": "package.json"},
        {"type": "file_exists", "file": "src/index.js"}
      ]
    },
    {
      "id": "STEP_002",
      "objective": "Implement core API",
      "status": "LOCKED",
      "depends_on": ["STEP_001"],
      "allow_paths": ["src/**", "tests/**"],
      "verify": [
        {
          "type": "command",
          "cmd": "npm test",
          "expect": {"exit_code": 0}
        }
      ]
    },
    {
      "id": "STEP_003",
      "objective": "API health check works",
      "status": "LOCKED",
      "depends_on": ["STEP_002"],
      "verify": [
        {
          "type": "http",
          "url": "http://localhost:3000/health",
          "expect": {"http_status": 200}
        }
      ]
    }
  ]
}
```

## Verification Rule Types

### 1. Command
```json
{
  "type": "command",
  "cmd": "npm test",
  "expect": {
    "exit_code": 0,
    "stdout_regex": "passed"
  }
}
```

### 2. File Exists
```json
{
  "type": "file_exists",
  "file": "src/api/login.js"
}
```

### 3. File Contains
```json
{
  "type": "file_contains",
  "file": "src/api/login.js",
  "pattern": "export async function login"
}
```

### 4. HTTP
```json
{
  "type": "http",
  "url": "http://localhost:3000/api/status",
  "method": "GET",
  "expect": {
    "http_status": 200
  }
}
```

## Best Practices

1. **Small Steps**: Break work into verifiable chunks
2. **Clear Objectives**: Each step should have a single, testable goal
3. **Automatic Tests**: Use `command` type with test suites
4. **Dependencies**: Use `depends_on` to enforce order
5. **Path Restrictions**: Use `allow_paths` to limit scope
6. **Fail Fast**: Don't advance until verification passes

## Common Workflows

### Feature Development
```bash
# 1. Create plan with test-driven steps
vim plan.json

# 2. Work on current step
plansm current

# 3. Implement feature
# ... code ...

# 4. Verify (runs tests)
plansm verify --current

# 5. Advance when green
plansm advance
```

### CI/CD Integration
```bash
# In .github/workflows/ci.yml
- name: Verify all steps
  run: plansm verify --all --json
```

### Code Review
```bash
# Reviewer checks plan state
plansm status
plansm verify --all

# Review only if verification passes
git diff plan.json
```

## Next Steps

- Read [Concept](concept.md) for the philosophy
- See [Spec](spec.md) for complete reference
- Check [Examples](../examples/) for real projects
- Learn about [Claude Code integration](claude-code.md)
