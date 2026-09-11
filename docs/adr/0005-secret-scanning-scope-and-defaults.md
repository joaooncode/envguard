# 0005: Secret Scanning Scope, Detection Strategy, and Defaults

To implement the v0.3.0 roadmap goal of content-based secret scanning, `envguard` introduces a **Secret Scanner** that inspects the text inside environment files for embedded secrets, producing **Secret Match** records distinct from the existing **Finding** concept.

## Context & Problem

Filename/path-level detection (the existing `Detector`) tells you a file *might* contain secrets, but not whether it actually does. The roadmap calls for two capabilities: known-provider signature matching (AWS, Stripe, GitHub, etc.) and entropy-based heuristics for unrecognized secret formats. These have very different false-positive profiles, and the feature needed to be scoped without breaking the `hook run` `<10ms` performance budget or the "never print secret values" design principle.

## Decision

1. **Scope**: The Secret Scanner only inspects files already classified as environment files (the `Detector`'s existing territory). It does not scan arbitrary repository files (source code, general config) in v0.3.0 — that's a larger surface with a different false-positive character, deferred to a later version.

2. **Detection strategy**: Pattern matching (known-provider signatures) and entropy calculation are independent checks, both able to produce `Secret Match` records on the same file. Pattern matching ships **on by default**; entropy-based detection is **opt-in** via `detector.entropy_scan`, since its higher false-positive rate risks alert fatigue and broken CI before the heuristic is tuned.

3. **Command boundaries**: The Secret Scanner runs in `scan` and `check` only. It is excluded from `hook run` to preserve its `<10ms` budget, and `fix` remains completely unaware of it — `fix` is scoped strictly to `.gitignore` edits and has no remediation path for a real secret (rotation isn't automatable).

4. **Secret Match is a distinct concept from Finding**, not a merge into it: a `Finding` (env file, severity from git tracked/staged status) can contain zero or more `Secret Match` records, each with its own severity floor (`HIGH` when untracked, `CRITICAL` when tracked/staged) — independent of the parent `Finding`'s severity. `check`'s exit-code aggregation takes the max severity across a file's `Finding` and its `Secret Match` records.

5. **Existing features extend to cover it**: files on the `Allowlist` (e.g. `.env.example`) are skipped entirely by the Secret Scanner. A file's `Severity Override` (e.g. `.env.test` → `warning`) acts as a ceiling on that file's `Secret Match` severities, so intentionally fake secrets in test fixtures don't trigger `CRITICAL`. A separate `detector.secret_ignore` list suppresses specific false-positive matches by key name/pattern regardless of file.

6. **Output stays zero-value-exposure**: a `Secret Match` is reported as key name (if identifiable) + file + line number only — never a value fragment, consistent with envguard's existing guarantee.

## Considered Options

- **Scanning all repository files, not just env files** — rejected for v0.3.0: much larger surface, higher false-positive risk from code that legitimately references key *names*, and a bigger traversal/performance question. Left as a future expansion.
- **Entropy detection on by default alongside pattern matching** — rejected: shipping a high-false-positive heuristic on-by-default in a security tool risks alert fatigue and broken CI runs before it's proven in the wild.
- **Folding Secret Match into the existing Finding object** — rejected: a `Finding`'s severity is driven by git exposure status, while a `Secret Match`'s severity is driven by the secret itself; conflating them would have made the severity model harder to reason about and prevented the `Severity Override` ceiling from applying cleanly.

## Consequences

- Widening scope later (e.g. scanning all repository files) will require revisiting the traversal and ignore-directory logic that today is scoped to already-detected env files.
- Enabling entropy detection by default in a future release is a config default change, not an architectural one — low cost to reverse.
