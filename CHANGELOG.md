# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2026-01-19

### Changed
- **Architecture Migration**: Converted from Go CLI to pure shell script Claude Code skill
- Simplified installation: No compilation needed, just copy `.claude/skills/plansm/`
- All commands now integrated as Claude Code skill (`/plansm`)
- Updated all documentation to reflect skill-based architecture
- Verification scripts remain shell-based (no dependency on Go)

### Added
- Claude Code skill integration with SKILL.md
- Skill-based workflow with automatic plan generation and execution
- Stop hook prevents session end until all steps verified
- Template system with `plan.json.example` as reference
- Comprehensive skill documentation and examples

### Removed
- Go CLI binary and compilation requirements
- Standalone CLI commands (replaced by skill commands)
- Go module dependencies

## [0.1.0] - 2026-01-12

Initial release. This is the first public version of plansm.
