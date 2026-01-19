# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **glob_pattern_check verification rule**: New rule type that verifies ALL files matching a glob pattern contain a required pattern
  - Fixes Issue #1: Sampling verification vulnerability where checking only one file can miss violations in other files
  - Supports `min_count`, `max_count`, and `exact_count` validation
  - Reports all files that fail pattern matching for easy debugging
  - Example: Verify all TypeScript files in a directory export a function
- Comprehensive documentation for glob_pattern_check in verification-rules.md
- Sampling verification trap warnings in SKILL.md and verification-rules.md
- Example usage in plan.json.example

### Changed
- Updated documentation to reflect 5 verification rule types (was 4)
- Enhanced anti-patterns section with sampling verification examples

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
