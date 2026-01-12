# Concept: Why Markdown Planning Fails

## The Problem with Markdown-Based Planning

Most LLM planning tools use Markdown files with checklists like:

```markdown
- [ ] Implement login API
- [ ] Add authentication middleware
- [x] Write tests
```

This creates a fundamental problem: **The LLM can mark tasks as complete without any verification**.

### Why This Is Problematic

1. **No Verification Mechanism**: Markdown checkboxes are just text. The LLM can mark `- [x]` without actually completing the work.

2. **Self-Reporting Bias**: LLMs are optimized to appear helpful and complete tasks. They will often claim completion based on "looking complete" rather than "being verified."

3. **Hallucinated Completeness**: When context windows get compressed or reset, LLMs lose track of what was actually done versus what was discussed.

4. **Token-Expensive**: Constantly re-reading large plan files consumes tokens without adding verification value.

5. **No Enforcement**: There's nothing preventing manual edits that claim completion without doing the work.

## The plansm Solution: Plan as State Machine

plansm treats planning as a **verifiable state machine** rather than a document:

### Core Principles

1. **Plan is Data, Not Documentation**
   - `plan.json` (structured, machine-readable)
   - NOT `plan.md` (narrative, human-readable)

2. **Completion Requires Proof**
   - Every step has `verify` rules (commands, file checks, API tests)
   - Status changes only through verification: `plansm verify`

3. **Separation of Authority**
   - LLM: proposes changes, implements features
   - CLI: verifies completion, updates state machine
   - Human: reviews and approves plans

4. **States Are Explicit**
   ```
   LOCKED    → Dependencies not met
   PENDING   → Ready to work on
   FAILED    → Verification failed
   VERIFIED  → Proof passed (only CLI can set)
   ```

5. **Low Token Cost**
   - LLM only reads current step context
   - Full plan stays in file system, not conversation

## Comparison

| Aspect | Markdown Planning | plansm |
|--------|------------------|---------|
| Format | `- [ ] Task` | `{"status": "PENDING", "verify": [...]}` |
| Verification | None (self-reported) | Machine-executable proofs |
| Status Update | LLM/Human can edit | Only CLI after verification |
| Token Cost | Re-read entire plan | Read current step only |
| Falsifiability | Easy to fake | Hard to fake (requires proof) |
| Auditability | No history | JSON diffs in git |
| Enforcement | None | Git hooks, CI gates |

## The Anti-Cheating Stack

plansm provides multiple layers to prevent "fake completion":

1. **JSON Schema Validation**: Structure must be correct
2. **Proof Execution**: Commands/tests must pass
3. **State Machine Rules**: Can't mark VERIFIED without passing verify
4. **Git Hooks** (optional): Reject manual status edits
5. **CI Gates** (optional): PR must pass `plansm verify --all`

## Real-World Impact

**Without plansm:**
- LLM claims "API implemented"
- Human reviews code, finds bugs
- Manual testing reveals missing cases
- Back-and-forth review cycles

**With plansm:**
- LLM implements API
- Runs `plansm verify --current`
- Tests fail → LLM fixes
- Tests pass → Automatic VERIFIED status
- Human reviews clean, working code

## Philosophy

> "If you can't measure it, you can't manage it."

Planning should be:
- **Executable**, not narrative
- **Verifiable**, not self-reported
- **State-based**, not document-based
- **Machine-first**, then human-readable

plansm brings software engineering rigor to LLM planning.
