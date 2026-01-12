# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial implementation of plansm CLI
- `plan.json` state machine format with JSON schema validation
- Four verification rule types: `command`, `file_exists`, `file_contains`, `http`
- CLI commands: `init`, `status`, `current`, `verify`, `advance`, `doctor`
- Claude Code integration with project commands
- Slash commands: `/pwork`, `/pverify`, `/pstatus`, `/pnext`
- Claude Code plugin structure (`.claude-plugin/`)
- Stop hook support for blocking on failed verification
- Comprehensive documentation:
  - Concept guide explaining why Markdown planning fails
  - Quick start guide
  - Claude Code integration guide
- Example project: Node.js REST API with test-driven steps
- Git pre-commit hook to prevent manual status edits
- GitHub Actions CI/CD workflows:
  - CI: Build, test, lint across multiple Go versions and OS
  - Release: Automated binary builds for multiple platforms
- Anti-cheating features:
  - State machine enforcement
  - Git hooks
  - CI verification gates
  - Audit trail through git history

### Fixed
- String escaping issues in Claude Code command generation
- Go module path updated to correct GitHub repository

## [0.1.0] - 2026-01-12

Initial release. This is the first public version of plansm.
