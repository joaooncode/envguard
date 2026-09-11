# envguard

Security-focused CLI tool to detect and prevent committed or exposed environment files in Git repositories.

## Language

**Scanner**:
The recursive filesystem traversal and Git status coordinator that produces scan findings.
_Avoid_: Crawler, inspector, walker

**Detector**:
The rule evaluator that determines if a file path is an environment file and whether it matches safe allowlist templates.
_Avoid_: Matcher, filter, classifier. Not to be confused with **Secret Scanner**, which inspects file *content* rather than file *paths*.

**Configuration**:
The project-level settings loaded from `.envguard.yaml`, `.envguard.yml`, or via `--config` to customize scanning and detection rules.
_Avoid_: Options, settings, preferences

**Allowlist**:
The collection of glob patterns for safe environment templates or sample files (e.g., `.env.example`) that should not raise security warnings. Files matching the Allowlist are also skipped entirely by the **Secret Scanner**.
_Avoid_: Whitelist, safe-list, permitted files

**Ignore Directory**:
A directory name skipped during recursive filesystem traversal (e.g., `node_modules`, `.git`).
_Avoid_: Excluded path, blacklisted folder

**Severity Override**:
An explicit configuration rule that replaces the calculated severity level for files matching a specific pattern. When set on a file, it also acts as a ceiling on that file's Secret Match severities (a Secret Match's severity is `min(computed severity, file's overridden severity)`).
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
The remediation component responsible for automatically generating and applying non-destructive `.gitignore` rules for unignored environment files. Scoped strictly to `.gitignore` edits; it never invokes the Secret Scanner and has no awareness of Secret Matches.
_Avoid_: Patcher, autofixer, corrector, remediator

**Hook Manager**:
The component responsible for installing, verifying, and uninstalling local Git pre-commit hook scripts (`.git/hooks/pre-commit`).
_Avoid_: Hook installer, hook script, hook handler

**Hook Runner**:
The fast execution mode invoked during Git pre-commit lifecycle to inspect staged repository files and prevent accidental commits of environment secrets.
_Avoid_: Commit watcher, stage scanner, commit blocker

**Secret Scanner**:
The content-inspection component that examines the text inside environment files for embedded secrets, using known-provider signature patterns (AWS, Stripe, GitHub, etc.) and entropy-based heuristics as independent checks. Skips any file covered by the Allowlist and any file that fails a binary/size check.
_Avoid_: Detector, content matcher, secret detector

**Secret Match**:
A single detected secret instance found by the Secret Scanner within a Finding's file: line number, key name (if identifiable), detection method (pattern or entropy), provider (if pattern-matched), and its own severity, independent of the parent Finding's severity.
_Avoid_: Secret finding, leak, hit
