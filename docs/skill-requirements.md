# Claude Code Skill Requirements

## skill structure and format

Claude Code skills are markdown files placed in `.claude-plugin/skills/` directory.

### File Location
```
.claude-plugin/
└── skills/
    └── plansm.md
```

### Skill File Format

```markdown
---
description: Short description of what the skill does
---

# Skill Instructions

Detailed instructions for how Claude should execute this skill.

## When to Use This Skill

Describe the trigger conditions.

## Execution Steps

1. Step 1
2. Step 2
...

## Output Format

What the skill should produce.
```

### Key Requirements

1. **Frontmatter**: Must include `description` field
2. **Autonomous Execution**: Skill runs until completion without user intervention
3. **Clear Instructions**: Should guide Claude through the entire workflow
4. **Verification**: Should include validation steps to ensure completion

### Skill Invocation

Users invoke skills with:
```
/skillname "arguments"
```

Example:
```
/plansm "Add dark mode toggle to the app"
```

### Skill vs Command

**Commands** (`.claude-plugin/commands/`):
- Short, single-action tools
- User needs to run multiple commands
- Example: `/pwork`, `/pverify`

**Skills** (`.claude-plugin/skills/`):
- Complete workflows
- Runs autonomously to completion
- Example: `/plansm` generates plan + executes all steps

### Best Practices for plansm Skill

1. **Accept user requirement as argument**
2. **Generate plan.json automatically**
3. **Execute loop**: implement → verify → advance
4. **Use internal tools**: pwork/pverify/pnext scripts
5. **Report completion with summary**
