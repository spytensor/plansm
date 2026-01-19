---
description: Show the current step (low token).
---
Run: `bash .claude-plugin/scripts/fsm.sh current`

This shows you what to work on next. After viewing:
1. Implement the step (use Task tool with appropriate subagent if complex)
2. Run `/pverify` to verify
3. If passes, run `/pnext` to advance
4. Repeat until all steps VERIFIED

Rules:
- Do NOT edit plan.json status fields manually.
- Only work on CURRENT_STEP.
- Use subagents for complex tasks.
