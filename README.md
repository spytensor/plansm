# plansm — Plan as a Verifiable State Machine

[![CI](https://github.com/spytensor/plansm/workflows/CI/badge.svg)](https://github.com/spytensor/plansm/actions)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/spytensor/plansm)](https://goreportcard.com/report/github.com/spytensor/plansm)

**plansm** turns "planning" into a **machine-verifiable state machine**, preventing LLMs from claiming completion without proof.

## Why plansm?

Traditional LLM planning tools use Markdown checklists that can be marked complete without verification. **plansm is different**:

- Plan is **data** (`plan.json`), not documentation
- "Done" is **proof-based**: machine-checkable verification rules (tests, commands, file checks)
- `VERIFIED` status written **only by CLI**, never by LLM or manual edit
- **Low token cost**: LLM reads only current step, not entire plan
- **Anti-cheating**: Git hooks + CI gates prevent fake completion

> Read [Why Markdown Planning Fails](docs/concept.md) for the philosophy.

## Core Principles

| Markdown Planning | plansm |
|------------------|---------|
| `- [ ] Implement API` | `{"status": "PENDING", "verify": [...]}` |
| Self-reported completion | Machine-verified proof required |
| Easy to fake | Hard to fake (tests must pass) |
| High token cost | Minimal context (current step only) |

### State Machine

```
LOCKED    → Dependencies not met
PENDING   → Ready to work on
FAILED    → Verification failed
VERIFIED  → Tests/proofs passed (CLI only)
```

## Installation

**No Go, no compilation, just copy files!**

```bash
# 1. Clone the repository
git clone https://github.com/spytensor/plansm.git

# 2. Copy plugin to your project
cd your-project
cp -r ../plansm/.claude-plugin .

# 3. Install jq (only dependency)
brew install jq  # macOS
# or: sudo apt-get install jq  # Linux

# That's it! Ready to use in Claude Code
```

**What you get**:
- Shell-based verification engine (no binaries needed)
- Claude Code commands: `/pwork`, `/pverify`, `/pstatus`, `/pnext`
- Machine-verifiable proof system
- Zero maintenance (pure shell scripts)

## Quick Start

### Automated Workflow (Recommended)

**Just give Claude your requirement, and let plansm handle everything:**

In Claude Code, simply say:

> "I need to add user authentication to my app"

Claude with plansm will **automatically**:

1. **Analyze** your requirement and codebase
2. **Generate** `plan.json` with all steps and verification rules
3. **Execute** each step one by one:
   - Implement the code (using appropriate subagents)
   - Verify it passes tests/checks
   - Advance to next step
4. **Report** when all steps are VERIFIED

**You don't write plan.json manually. You don't run commands manually. Just state your need.**

### Available Commands

- `/plan` — Give requirements, auto-generate plan.json and execute
- `/pwork` — Show current step (low token)
- `/pverify` — Run verification (tests/checks)
- `/pstatus` — Show all steps
- `/pnext` — Advance (only if verified)

### Manual Planning (Advanced)

If you prefer to write `plan.json` yourself:

```json
{
  "version": 1,
  "current_step": "STEP_001",
  "steps": [
    {
      "id": "STEP_001",
      "title": "Implement login API",
      "status": "PENDING",
      "verification": {
        "rules": [
          {"type": "command", "cmd": "npm test login"}
        ]
      }
    }
  ]
}
```

Then tell Claude: "Follow plan.json using /pwork, /pverify, /pnext"

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

## CLI Commands

```bash
plansm init [--claude]              # Create plan.json and Claude integration
plansm status                       # Show all steps and statuses
plansm current [--json]             # Show current step (low token)
plansm verify --current|--all       # Run verification proofs
plansm advance                      # Move to next step (if verified)
plansm doctor                       # Health check
```

## Documentation

- [Quick Start Guide](docs/quickstart.md) — Get started in 5 minutes
- [Concept](docs/concept.md) — Why Markdown planning fails
- [Claude Code Integration](docs/claude-code.md) — Hooks, commands, workflows
- [Examples](examples/) — Node.js, Python projects

## Anti-Cheating Features

1. **State Machine Enforcement**: Only CLI can set VERIFIED status
2. **Git Pre-commit Hook**: Warns on manual status edits
3. **CI Verification**: `plansm verify --all` in CI/CD
4. **Audit Trail**: All state changes tracked in git history
5. **JSON Schema**: Structure validation

Install git hooks:
```bash
./scripts/install-hooks.sh
```

## Contributing

Contributions welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md).

### Development Setup
```bash
git clone https://github.com/spytensor/plansm.git
cd plansm
go mod download
go build -o plansm ./cmd/plansm
go test ./...
```

## Roadmap

- [ ] SQLite backend (in addition to JSON)
- [ ] Web UI for plan visualization
- [ ] More verify rule types (Docker, API contracts)
- [ ] Plan templates library
- [ ] Homebrew formula
- [ ] VS Code extension

## License

MIT License - see [LICENSE](LICENSE) for details.

## Credits

Inspired by the discussion on LLM planning failures and the need for verifiable, machine-enforced completion criteria.

---

**plansm** — Because "looks done" ≠ "is done"
