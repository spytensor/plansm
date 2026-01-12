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

### From Source
```bash
git clone https://github.com/spytensor/plansm.git
cd plansm
go build -o plansm ./cmd/plansm
sudo mv plansm /usr/local/bin/
```

### From Releases (Recommended)
```bash
# Download latest binary for your platform
curl -L https://github.com/spytensor/plansm/releases/latest/download/plansm-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m) -o plansm
chmod +x plansm
sudo mv plansm /usr/local/bin/
```

## Quick Start

```bash
# 1. Initialize plan in your project
plansm init --claude

# 2. View current step (low token)
plansm current

# 3. Implement the step...

# 4. Verify (runs tests/checks)
plansm verify --current

# 5. Advance to next step
plansm advance
```

### Claude Code Integration

After `plansm init --claude`, use these commands:

- `/pwork` — Show current step
- `/pverify` — Verify with machine proofs
- `/pstatus` — Show all steps
- `/pnext` — Advance (only if verified)

**Recommended**: Set up Stop hook in Claude Code:
```bash
# In Claude Code, run: /hooks
# Add Stop hook: plansm verify --current (blocking)
```

This prevents Claude from stopping until verification passes.

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
