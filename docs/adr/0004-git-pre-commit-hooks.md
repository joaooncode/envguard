# 0004: Git Pre-Commit Hooks and Staged Inspection

To safeguard developers from inadvertently committing sensitive environment files, `envguard` provides native Git pre-commit hook management via `envguard hook install`, `envguard hook run`, and `envguard hook uninstall`, alongside official integration with the Python `pre-commit` framework (`.pre-commit-hooks.yaml`).

## Context & Problem

Scanning the entire filesystem during every `git commit` invocation can be slow and disruptive in large repositories. Furthermore, developers need a zero-friction mechanism to install local hooks without relying on third-party package managers if they prefer a standalone Go binary. When hooks are installed, they must be idempotent, non-destructive, and execute in milliseconds by inspecting only files currently staged for commit.

## Decision

1. **Subcommands Architecture (`envguard hook <subcommand>`)**:
   - `envguard hook install`: Resolves the `.git/hooks` directory (respecting repository root and standard Git hooks structure), writes an executable POSIX shell script to `.git/hooks/pre-commit` (with `0755` permissions on Unix platforms), and marks it with the signature `# Installed by envguard`. If a pre-existing hook without the envguard signature exists, it rejects installation unless `--force` is supplied.
   - `envguard hook run`: Queries only staged files via `git diff --name-only --cached --diff-filter=ACM` and evaluates each staged path through the `Detector` and `Allowlist`. If any non-allowlisted environment file is staged, it outputs findings and exits with `1`, blocking the commit. If no staged environment files are detected, it exits immediately with `0`.
   - `envguard hook uninstall`: Safely removes `.git/hooks/pre-commit` if it was installed by envguard (verified by header signature) or if `--force` is provided.

2. **Ultra-Fast Staged Inspection**:
   - Staged inspection bypasses full recursive filesystem crawling and only analyzes paths staged in the Git index.
   - Supports `--path` (repository root) and optional `--config` for custom rule loading.

3. **Official `pre-commit` Framework Integration**:
   - Root repository `.pre-commit-hooks.yaml` configures the `envguard` hook with `language: system`, `entry: envguard hook run`, `pass_filenames: false`, `always_run: true`, and `stages: [pre-commit]`.
