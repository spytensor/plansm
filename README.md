# plansm — Plan as a Verifiable State Machine (for Claude Code)

**plansm** turns “planning” into a **machine-verifiable state machine**.

- Plan is **data**: `plan.json`
- “Done” is **proof-based**: each step has machine-checkable verification rules
- `VERIFIED` is written only by the verifier (CLI), not by Claude and not by humans

> This repo intentionally avoids Markdown-based planning as a source of truth.
> Markdown can still exist as docs, but **never** as the plan authority.

## Quickstart

```bash
go build -o plansm ./cmd/plansm

# initialize in your project
./plansm init --claude

# show current step (small output → low token)
./plansm current

# run verification for current step (updates plan.json status)
./plansm verify --current

# if VERIFIED, advance to next step and unlock dependents
./plansm advance
```

Claude Code project commands created under `.claude/commands/`:
- `/pwork`  → `plansm current`
- `/pverify`→ `plansm verify --current`
- `/pstatus`→ `plansm status`
- `/pnext`  → `plansm advance`

## Plan file (plan.json)

```json
{
  "version": 1,
  "current_step": "STEP_001",
  "steps": [
    {
      "id": "STEP_001",
      "objective": "Initialize project skeleton",
      "status": "PENDING",
      "verify": [
        { "type": "command", "cmd": "test -d ." }
      ]
    },
    {
      "id": "STEP_002",
      "objective": "Implement login API",
      "status": "LOCKED",
      "depends_on": ["STEP_001"],
      "verify": [
        { "type": "command", "cmd": "npm test -- login", "expect": { "exit_code": 0 } }
      ]
    }
  ]
}
```

Statuses: `LOCKED | PENDING | FAILED | VERIFIED`

Verify rule types:
- `command`: run bash command, check exit code and optional stdout regex
- `file_exists`
- `file_contains` (regex)
- `http` (GET by default; checks status code)

## CLI
- `plansm init [--claude]`
- `plansm status`
- `plansm current [--json]`
- `plansm verify --current|--all [--json]`
- `plansm advance`
- `plansm doctor`

## License
MIT (see LICENSE).
