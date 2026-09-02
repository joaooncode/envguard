# 0003: Subcommand `fix` and Automatic Remediation

To assist developers in preventing accidental leaks before commits occur, `envguard` provides the `fix` subcommand. By default, `envguard fix` scans the target directory, identifies unprotected local environment files (severity `WARNING`), and automatically updates the root `.gitignore` with the corresponding relative patterns.

## Context & Problem

Unignored `.env` files in local working trees pose an imminent risk of accidental staging and committing. While `envguard scan` identifies these files, manual updating of `.gitignore` across deep subdirectories is prone to omissions and syntax mistakes. Furthermore, tracked files (`CRITICAL`) cannot be remediated solely via `.gitignore` and require active Git cache removal (`git rm --cached`) and secret rotation.

## Decision

1. **Centralized Root `.gitignore` Remediation**:
   `envguard fix` writes remediation patterns to the root repository `.gitignore`. Filepaths located in subdirectories are converted to root-relative ignore paths (e.g., `/packages/backend/.env.local`).

2. **Identified & Non-Destructive Formatting**:
   New ignore patterns are inserted under an identified section (`# Added by envguard`) while preserving all existing comments, whitespace, and formatting. Existing rules are checked to avoid duplicate entries.

3. **Dry-Run Mode**:
   When `--dry-run` is supplied, `envguard fix` previews proposed changes to `.gitignore` without altering the filesystem.

4. **Handling of Critical Findings & Exit Codes**:
   When tracked environment files (`CRITICAL`) are encountered, `envguard fix` outputs actionable remediation guidance (instructing developers to run `git rm --cached <path>` and rotate credentials) and returns an exit code of `1` to signal remaining unresolved security risks in the repository. If all findings are resolved or no findings exist, it returns `0`.
