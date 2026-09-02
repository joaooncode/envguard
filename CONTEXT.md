# envguard

Security-focused CLI tool to detect and prevent committed or exposed environment files in Git repositories.

## Language

**Scanner**:
The recursive filesystem traversal and Git status coordinator that produces scan findings.
_Avoid_: Crawler, inspector, walker

**Detector**:
The rule evaluator that determines if a file path is an environment file and whether it matches safe allowlist templates.
_Avoid_: Matcher, filter, classifier

**Configuration**:
The project-level settings loaded from `.envguard.yaml`, `.envguard.yml`, or via `--config` to customize scanning and detection rules.
_Avoid_: Options, settings, preferences

**Allowlist**:
The collection of glob patterns for safe environment templates or sample files (e.g., `.env.example`) that should not raise security warnings.
_Avoid_: Whitelist, safe-list, permitted files

**Ignore Directory**:
A directory name skipped during recursive filesystem traversal (e.g., `node_modules`, `.git`).
_Avoid_: Excluded path, blacklisted folder

**Severity Override**:
An explicit configuration rule that replaces the calculated severity level for files matching a specific pattern.
_Avoid_: Custom rule, priority tweak

**Finding**:
The detected environment file with its assigned severity level, git status, and mitigation suggestions.
_Avoid_: Vulnerability, issue, report item

**Initializer**:
The CLI component responsible for bootstrapping repository configuration (`.envguard.yaml`) and safe environment templates (`.env.example`).
_Avoid_: Setup generator, config creator, scaffolder

**Sanitization**:
The process of stripping secret values from environment definitions while preserving comments, formatting, and key names to produce safe templates.
_Avoid_: Masking, redacting, cleaning

**Fixer**:
The remediation component responsible for automatically generating and applying non-destructive `.gitignore` rules for unignored environment files.
_Avoid_: Patcher, autofixer, corrector, remediator
