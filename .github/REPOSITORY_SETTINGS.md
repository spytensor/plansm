# GitHub Repository Settings

This file contains recommended settings for the plansm GitHub repository.

## About Section

**Description**:
```
Machine-verified planning for Claude Code. Zero install: cp files + jq. LLMs can't fake completion - tests must pass.
```

**Website** (optional):
```
https://github.com/spytensor/plansm
```

**Topics** (tags):
```
claude-code
llm
planning
verification
automation
proof-based
shell-scripts
testing
```

## Repository Details

**Social Image**:
Consider creating a simple diagram showing:
```
Markdown Planning (❌ self-reported) → plansm (✅ machine-verified)
```

## README Badges

Already included in README.md:
- CI Status
- License
- Go Report Card

## Key Features to Highlight

1. **Zero Install** - No Go, no compilation, just `cp -r .claude-plugin`
2. **Machine Verification** - Tests must actually pass, not self-reported
3. **Claude Automation** - Claude follows plan.json automatically
4. **Low Token** - Only reads current step, not entire plan
5. **Anti-Cheating** - Status only updated by verification scripts

## Community Files

Already present:
- ✅ LICENSE (MIT)
- ✅ README.md
- ✅ CONTRIBUTING.md
- ✅ CODE_OF_CONDUCT.md (via CONTRIBUTING.md reference)

## How to Update GitHub About

1. Go to: https://github.com/spytensor/plansm
2. Click "⚙️" (Settings icon) next to "About" on the right sidebar
3. Add:
   - Description (see above)
   - Website (optional)
   - Topics (see above)
4. Save changes

The "About" section will then show:
```
🤖 plansm
Machine-verifiable planning for LLMs. No fake completion - tests must pass.
Zero install: just copy files.

Topics: llm claude planning verification testing automation...
```
