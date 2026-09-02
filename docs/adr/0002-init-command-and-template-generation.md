# 0002: Subcommand `init` and Safe Template Generation

To streamline repository onboarding and enforce security best practices from day one, `envguard` provides the `init` subcommand. By default, `envguard init` generates a documented `.envguard.yaml` configuration file containing recommended security defaults and commented configuration blocks.

When invoked with `--template`, `envguard init` produces a safe `.env.example` file. If an existing `.env` file is present (or specified via `--template-from`), the initializer performs key sanitization, preserving comments, empty lines, and environment variable names while stripping all sensitive values to empty assignments (`KEY=`). If no source `.env` file exists, a curated default `.env.example` boilerplate is created.

To protect existing configurations, `envguard init` refuses to overwrite existing files unless the `--force` flag is explicitly provided.
