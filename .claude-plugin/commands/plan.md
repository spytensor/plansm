---
description: Generate plan.json from requirements and auto-execute
---

# Plan Generation and Auto-Execution

When the user provides a development requirement, you should:

## Step 1: Analyze Requirements

Deeply understand what the user wants to build:
- What is the core functionality?
- What files need to be created/modified?
- What are the dependencies between tasks?
- What verification is needed for each step?

## Step 2: Generate plan.json

Create a NEW plan.json with:

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
      "title": "Descriptive title",
      "status": "PENDING",
      "description": "What needs to be done",
      "dependencies": [],
      "verification": {
        "rules": [
          {
            "type": "command|file_exists|file_contains|http",
            "description": "What this verifies",
            "...": "rule-specific fields"
          }
        ]
      }
    }
  ]
}
```

### Verification Rule Types:

**command**: Run a command, check exit code
```json
{
  "type": "command",
  "description": "Tests pass",
  "cmd": "npm test",
  "expect": {"exit_code": 0}
}
```

**file_exists**: Check if file exists
```json
{
  "type": "file_exists",
  "description": "Component file created",
  "path": "src/components/NewFeature.tsx"
}
```

**file_contains**: Check file content
```json
{
  "type": "file_contains",
  "description": "Export added to index",
  "path": "src/index.ts",
  "expect": {
    "contains": "export { NewFeature }"
  }
}
```

**http**: Check HTTP response
```json
{
  "type": "http",
  "description": "API endpoint responds",
  "url": "http://localhost:3000/api/health",
  "expect": {
    "status": 200,
    "body_contains": "ok"
  }
}
```

## Step 3: Task Breakdown Principles

- **Atomic steps**: Each step should be independently verifiable
- **Clear dependencies**: Use `dependencies: ["STEP_XXX"]` for ordering
- **Start LOCKED**: Steps with dependencies start as LOCKED
- **Reasonable granularity**: Not too fine (10+ steps for simple task), not too coarse (1 step for complex feature)
- **Verification first**: Always think "how do I PROVE this is done?" before writing the step

## Step 4: Auto-Execution Loop

After generating plan.json, IMMEDIATELY start executing:

```bash
# 1. Check current step
bash .claude-plugin/scripts/fsm.sh current

# 2. Implement the step
# Use appropriate subagents:
# - Task tool with subagent_type="general-purpose" for implementation
# - Task tool with subagent_type="Bash" for git/command operations
# - Edit/Write tools for file changes

# 3. Verify the step
bash .claude-plugin/scripts/verify.sh --current

# 4. If verification passes, advance
bash .claude-plugin/scripts/fsm.sh advance

# 5. Repeat until all steps VERIFIED
```

## Important Rules

1. **NEVER manually edit status fields** in plan.json
2. **ALWAYS verify** before marking complete
3. **Use subagents** for complex tasks - don't do everything yourself
4. **Keep going** until ALL steps are VERIFIED
5. **Each task gets a NEW plan.json** - don't append to existing plans

## Example Workflow

User: "Add dark mode toggle to the app"

You should:
1. Analyze: Need theme context, toggle component, CSS updates, tests
2. Generate plan.json with 5 steps:
   - STEP_001: Create theme context (verify: file exists, exports ThemeContext)
   - STEP_002: Create toggle component (verify: file exists, tests pass)
   - STEP_003: Add theme CSS variables (verify: file contains CSS vars)
   - STEP_004: Integrate toggle in header (verify: file contains import)
   - STEP_005: Run all tests (verify: command exit 0)
3. Auto-execute: Implement → Verify → Advance, repeat
4. Report completion when all steps VERIFIED

## Integration with Subagents

For each step implementation, choose the right tool:
- **Simple edits**: Use Edit/Write directly
- **Complex logic**: Use Task tool with general-purpose agent
- **Git operations**: Use Bash tool or Task with Bash agent
- **Research needed**: Use Task with Explore agent first

Start now - generate the plan.json and begin execution!
