# Contributing to plansm

Thank you for your interest in contributing to plansm!

## Core Design Principles

1. **Machine Verification Only**: Plans must be machine-verifiable, no narrative requirements
2. **Script-Verified Status**: `VERIFIED` can only be written by verification scripts, never by manual edits or LLM self-reporting
3. **Pure Shell Implementation**: All verification logic must be in bash + jq (no compilation required)
4. **Self-Contained Skill**: Everything needed must be in `skills/plansm/` directory

## Development Setup

```bash
git clone https://github.com/spytensor/plansm.git
cd plansm

# Install jq (only dependency)
brew install jq  # macOS
# or: sudo apt-get install jq  # Linux

# Test the skill locally
cp -r skills/plansm ~/.claude/skills/
# Then use /plansm in Claude Code
```

## Testing Changes

```bash
# Test verification scripts
bash skills/plansm/scripts/verify.sh --help
bash skills/plansm/scripts/fsm.sh status

# Test with example plan
cp plan.json.example plan.json
bash skills/plansm/scripts/verify.sh --current

# Validate JSON schema
jq empty plan.json.example
```

## Submitting Changes

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/my-feature`
3. **Make your changes**
4. **Test thoroughly**: Verify scripts work on both macOS and Linux
5. **Update documentation**: If you change behavior, update SKILL.md
6. **Commit with clear messages**: Use descriptive commit messages
7. **Submit a Pull Request**

## Code Style

- **Shell Scripts**: Follow shellcheck recommendations
- **JSON**: Validate with `jq empty filename.json`
- **Documentation**: Use clear, concise language

## Adding New Verification Rule Types

If you want to add a new verification rule type (beyond command, file_exists, file_contains, http):

1. Update `skills/plansm/scripts/verify.sh` to handle the new type
2. Add examples to `skills/plansm/references/verification-rules.md`
3. Update `plan.json.example` with an example
4. Test on both macOS and Linux

## Questions?

Open an issue for discussion before major changes.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
