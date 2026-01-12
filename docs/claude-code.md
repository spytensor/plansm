# Claude Code Integration

This guide shows how to use plansm with Claude Code for verifiable, low-token planning.

## Why Claude Code + plansm?

Traditional Claude Code workflows can suffer from:
- **Attention drift**: LLM forgets original goals after many tool calls
- **Token bloat**: Repeatedly reading plan files
- **Fake completion**: LLM claims done without verification
- **Context loss**: Plan details lost after summarization

plansm solves this by:
- **Explicit state machine**: Current step always clear
- **Proof-based completion**: Can't claim done without passing verify
- **Minimal context**: Only current step in conversation
- **Machine verification**: Tests/commands gate advancement

## Installation

### In Your Project

```bash
# Install plansm CLI
go install github.com/spytensor/plansm/cmd/plansm@latest

# Initialize in project
cd your-project
plansm init --claude
```

This creates:
```
your-project/
├── plan.json                    # State machine
├── .claude/
│   └── commands/               # Project slash commands
│       ├── pwork.md
│       ├── pverify.md
│       ├── pstatus.md
│       └── pnext.md
└── .claude-plugin/             # Plugin structure
    ├── plugin.json
    ├── commands/
    └── hooks/
```

### Global Installation (Optional)

For commands available across all projects:

```bash
# Create user commands directory
mkdir -p ~/.claude/commands

# Copy commands
cp .claude/commands/*.md ~/.claude/commands/
```

## Project Commands

Once installed, these slash commands are available in Claude Code:

### `/pwork` - Show Current Step

Shows only the current step with minimal token usage.

```
CURRENT_STEP: STEP_002
STATUS: PENDING
OBJECTIVE: Implement login API
ALLOW_PATHS:
  - src/api/**
  - tests/**
VERIFY:
  - command: npm test -- login
```

**When to use**: At the start of work, or when Claude seems confused about what to do.

### `/pverify` - Verify Current Step

Runs machine verification on the current step.

```bash
# Runs: plansm verify --current
```

Output:
```
STEP STEP_002: OK
  ✓ command: npm test -- login

OVERALL: OK
```

**When to use**: After implementing a feature, before considering it done.

### `/pstatus` - Show Plan Overview

Shows all steps and their states.

```
current_step: STEP_002

STEP        STATUS     OBJECTIVE
--------------------------------------------------------------------------------
STEP_001    VERIFIED   Initialize project skeleton
STEP_002    PENDING    Implement login API
STEP_003    LOCKED     Add authentication middleware
STEP_004    LOCKED     Integration tests
```

**When to use**: To understand overall progress or check what's blocked.

### `/pnext` - Advance to Next Step

Advances to the next step (only if current step is VERIFIED).

```bash
# Runs: plansm advance
```

Output:
```
advanced to STEP_003
```

**When to use**: After verification passes, to move to the next task.

## Hooks: Automatic Verification

Hooks ensure verification happens at critical points.

### Stop Hook (Recommended)

Prevents Claude from stopping until the current step is verified.

**Setup in Claude Code:**

1. Run `/hooks` in Claude Code
2. Select **Stop** hook
3. Set command: `plansm verify --current`
4. Set as **blocking** (if option available)

**Effect:**
- Claude tries to stop → Hook runs `plansm verify --current`
- If verification fails → Stop is blocked, Claude sees error
- If verification passes → Stop succeeds

This prevents Claude from claiming completion without proof.

### Manual Hook Configuration

If `/hooks` UI is not available, edit `.claude/settings.json`:

```json
{
  "hooks": {
    "Stop": [
      {
        "matcher": "",
        "command": "plansm verify --current",
        "blocking": true
      }
    ]
  }
}
```

Or use project-local: `.claude/settings.local.json`

### PreToolUse Hook (Optional)

Reminds Claude of current step before major actions.

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "^(Write|Edit|Bash)",
        "command": "echo \"Current step: $(plansm current --json | jq -r .objective)\""
      }
    ]
  }
}
```

**Note**: This adds token cost. Use only if Claude frequently loses focus.

## Recommended Workflow

### 1. Plan Phase

Create your plan in `plan.json`:

```json
{
  "version": 1,
  "current_step": "STEP_001",
  "steps": [
    {
      "id": "STEP_001",
      "objective": "Setup Express server",
      "status": "PENDING",
      "verify": [
        {"type": "command", "cmd": "node src/server.js & sleep 2 && curl http://localhost:3000"},
        {"type": "file_contains", "file": "src/server.js", "pattern": "express"}
      ]
    },
    ...
  ]
}
```

### 2. Work Phase

Tell Claude:

> "Follow plan.json. Use /pwork to see current step. Use /pverify when done."

Claude will:
1. Run `/pwork` to see current objective
2. Implement the feature
3. Run `/pverify` to check proof
4. If failed: fix and re-verify
5. If passed: run `/pnext`

### 3. Review Phase

As a human:

```bash
# Check overall state
plansm status

# Verify all steps
plansm verify --all

# Review code changes
git diff
```

## Advanced Usage

### JSON Output

For programmatic use:

```bash
plansm current --json
plansm verify --current --json
```

Example output:
```json
{
  "ok": true,
  "results": [
    {
      "step_id": "STEP_002",
      "ok": true,
      "results": [
        {
          "rule": {"type": "command", "cmd": "npm test"},
          "ok": true
        }
      ]
    }
  ]
}
```

### Multi-Step Verification

Verify all pending/failed steps in order:

```bash
plansm verify --all
```

This is useful for:
- CI/CD pipelines
- Pre-commit hooks
- PR checks

### Custom Plan Path

If `plan.json` is not in current directory:

```bash
plansm status --plan ./plans/feature-x.json
plansm verify --current --plan ./plans/feature-x.json
```

## Plugin Distribution

The `.claude-plugin/` directory can be distributed as a plugin:

```
.claude-plugin/
├── plugin.json          # Plugin metadata
├── commands/            # Command definitions
│   ├── pwork.md
│   ├── pverify.md
│   ├── pstatus.md
│   └── pnext.md
└── hooks/
    └── hooks.json       # Hook definitions
```

Users can copy this directory to their projects or install it as a Claude Code plugin.

## Troubleshooting

### Commands Not Showing Up

1. **Restart Claude Code session**: Commands are loaded at session start
2. **Check path**: Commands should be in `.claude/commands/`
3. **Check format**: Files must have frontmatter:
   ```markdown
   ---
   description: Command description
   ---
   ```

### Hooks Not Running

1. **Check hooks config**: Must be in `.claude/settings.json` or `.claude/settings.local.json`
2. **Restart session**: Hooks load at session start
3. **Check command path**: `plansm` must be in PATH
4. **Check blocking**: Some Claude versions don't support blocking hooks

### Verification Fails in Hook

If verification fails in a hook, Claude will see the output:

```
STEP STEP_002: FAILED
  ✗ command: npm test — exit_code=1
```

Claude should then:
1. Understand what failed
2. Fix the issue
3. Re-run `/pverify`

If Claude ignores this, explicitly tell it:
> "Verification failed. Do not stop until /pverify succeeds."

## Token Optimization Tips

1. **Use `/pwork`, not full plan**: Shows only current step
2. **Avoid PreToolUse hooks**: Adds tokens on every tool call
3. **Stop hook only**: Minimal tokens, maximum safety
4. **Clear objectives**: Short, specific step objectives reduce clarification rounds
5. **Proof-first**: Write verification rules before implementation to guide Claude

## Best Practices

1. **Small, Verifiable Steps**: Each step should take 5-15 minutes of work
2. **Test-Driven**: Include tests in verify rules
3. **Fail Fast**: Don't advance until verification passes
4. **Explicit Dependencies**: Use `depends_on` to enforce order
5. **Path Restrictions**: Use `allow_paths` to scope changes
6. **Clear Invariants**: Remind Claude what NOT to do

## Example Session

```
User: Follow plan.json to implement the login feature

Claude: Let me check the current step.
[Runs /pwork]

CURRENT_STEP: STEP_003
STATUS: PENDING
OBJECTIVE: Implement login API endpoint
VERIFY:
  - command: npm test -- login

I'll implement the login API endpoint at src/api/login.js.
[... implements code ...]

Let me verify this works.
[Runs /pverify]

STEP STEP_003: OK
  ✓ command: npm test -- login

Great! The tests pass. Moving to next step.
[Runs /pnext]

advanced to STEP_004

Now working on: Add JWT authentication
```

## CI/CD Integration

Add to `.github/workflows/ci.yml`:

```yaml
- name: Verify Plan
  run: |
    if [ -f plan.json ]; then
      plansm verify --all || exit 1
    fi
```

This ensures:
- All steps are verifiable
- No fake completion in PRs
- Plan and code stay synchronized

## Summary

plansm + Claude Code = **Low-token, verifiable, machine-enforced planning**.

- Commands: `/pwork`, `/pverify`, `/pstatus`, `/pnext`
- Hook: Stop hook with `plansm verify --current`
- Workflow: Plan → Work → Verify → Advance
- Result: No fake completion, minimal tokens, clear progress
