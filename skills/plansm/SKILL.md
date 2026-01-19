---
name: plansm
version: "1.0.0"
description: Machine-verified planning with JSON state machines and proof-based completion. Generates plan.json, auto-executes with verification, prevents fake completion through tests.
user-invocable: true
allowed-tools:
  - Read
  - Write
  - Edit
  - Bash
  - Glob
  - Grep
  - Task
  - WebFetch
  - WebSearch
hooks:
  Stop:
    - hooks:
        - type: command
          command: |
            if [ -f plan.json ]; then
              bash "${CLAUDE_PLUGIN_ROOT}/scripts/verify.sh" --all || exit 1
            fi
---

# plansm: Machine-Verified Planning Skill

This skill implements end-to-end machine-verified planning and execution.

## When to Use This Skill

User provides a development requirement like:
- "Add dark mode toggle to my app"
- "Implement user authentication with JWT"
- "Add real-time notifications feature"

You run: `/plansm` and it handles everything automatically.

## Core Workflow

### Phase 1: Analyze Requirements

Deeply understand what the user wants to build:
- What is the core functionality?
- What files need to be created/modified?
- What are the dependencies between tasks?
- What verification is needed for each step?

### Phase 2: Generate plan.json

Create a NEW plan.json with this structure:

```json
{
  "version": 1,
  "current_step": "STEP_001",
  "invariants": [
    "Do not mark VERIFIED without running verification.",
    "Only work on current_step unless explicitly unlocked.",
    "Test each step before advancing."
  ],
  "steps": [
    {
      "id": "STEP_001",
      "objective": "Descriptive objective of what to accomplish",
      "status": "PENDING",
      "depends_on": [],
      "verify": [
        {
          "type": "command|file_exists|file_contains|http",
          "...": "rule-specific fields"
        }
      ]
    }
  ]
}
```

### Phase 3: Auto-Execution Loop

After generating plan.json, IMMEDIATELY start executing:

```bash
# Loop until all steps VERIFIED:
while true; do
  # 1. Check current step
  bash ${CLAUDE_PLUGIN_ROOT}/scripts/fsm.sh current

  # 2. Implement the step
  # Use appropriate subagents:
  # - Task tool with subagent_type="general-purpose" for implementation
  # - Task tool with subagent_type="Bash" for git/command operations
  # - Edit/Write tools for file changes

  # 3. Verify the step
  bash ${CLAUDE_PLUGIN_ROOT}/scripts/verify.sh --current

  # 4. If verification passes, advance
  if verification_passed; then
    bash ${CLAUDE_PLUGIN_ROOT}/scripts/fsm.sh advance
  else
    # Fix the issue and retry
    fix_and_retry
  fi

  # 5. Check if done
  if all_steps_verified; then
    break
  fi
done
```

### Phase 4: Report Completion

When all steps VERIFIED:
1. Show summary of completed steps
2. Confirm all verifications passed
3. Report success to user

## verification rules Types

### command: Run a command, check exit code
```json
{
  "type": "command",
  "cmd": "npm test",
  "expect": {"exit_code": 0}
}
```

### file_exists: Check if file exists
```json
{
  "type": "file_exists",
  "file": "src/components/NewFeature.tsx"
}
```

### file_contains: Check file content
```json
{
  "type": "file_contains",
  "file": "src/index.ts",
  "pattern": "export.*NewFeature"
}
```

### http: Check HTTP response
```json
{
  "type": "http",
  "url": "http://localhost:3000/api/health",
  "expect_status": 200,
  "expect_body": "ok"
}
```

## Task Breakdown Principles

- **Atomic steps**: Each step should be independently verifiable
- **Clear dependencies**: Use `depends_on: ["STEP_XXX"]` for ordering
- **Start LOCKED**: Steps with dependencies start as LOCKED status
- **Reasonable granularity**: Not too fine (10+ steps for simple task), not too coarse (1 step for complex feature)
- **Verification first**: Always think "how do I PROVE this is done?" before writing the step

## Important Rules

1. **NEVER manually edit status fields** in plan.json - only verification scripts can do this
2. **ALWAYS verify** before marking complete
3. **Use subagents** for complex tasks - don't do everything yourself
4. **Keep going** until ALL steps are VERIFIED - no stopping midway
5. **Each task gets a NEW plan.json** - don't append to existing plans

## Example Execution

User runs: `/plansm`

You respond:
```
I'll implement this feature using machine-verified planning.

Step 1: Analyzing requirements...
- Need: Theme context provider
- Need: Dark mode toggle component
- Need: CSS variable system
- Need: Integration tests

Step 2: Generating plan.json with 5 steps...
✓ Plan generated

Step 3: Auto-executing...

STEP_001: Create theme context provider
[Implements ThemeContext.tsx]
[Runs verification]
✓ Verified - file exists and exports ThemeContext

STEP_002: Create dark mode toggle component
[Implements DarkModeToggle.tsx]
[Runs verification]
✓ Verified - file exists and tests pass

STEP_003: Add CSS variables for theming
[Updates global.css]
[Runs verification]
✓ Verified - CSS contains theme variables

STEP_004: Integrate toggle in app header
[Updates Header.tsx]
[Runs verification]
✓ Verified - import statement found

STEP_005: Run full test suite
[Executes npm test]
✓ Verified - all tests pass

All steps completed and verified!

Summary:
- 5 steps executed
- 5 verifications passed
- Feature fully implemented
```

## Integration with Subagents

For each step implementation, choose the right tool:
- **Simple edits**: Use Edit/Write directly
- **Complex logic**: Use Task tool with general-purpose agent
- **Git operations**: Use Bash tool or Task with Bash agent
- **Research needed**: Use Task with Explore agent first

## Key Difference from Manual Planning

**Without plansm**: LLM says "done" → code might be broken → user discovers bugs later

**With plansm**: LLM must pass verification → tests prove completion → guaranteed working code

Start now - analyze the requirement, generate the plan.json, and begin automatic execution!
