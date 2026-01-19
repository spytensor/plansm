---
description: Advance to next step (only if current step is VERIFIED).
---
Run: `bash .claude-plugin/scripts/fsm.sh advance`

After advancing:
- Automatically unlocks any LOCKED steps whose dependencies are now met
- Updates current_step to the next PENDING step
- If all steps VERIFIED, reports completion

After running this, IMMEDIATELY continue with the next step:
1. Run `/pwork` to see the new current step
2. Implement it
3. Run `/pverify`
4. Run `/pnext` again
5. Repeat until done

**Keep the loop going automatically - don't wait for user input between steps!**
