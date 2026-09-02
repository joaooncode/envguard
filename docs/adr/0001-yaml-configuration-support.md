# 0001: YAML Configuration File Support (.envguard.yaml)

To allow repositories to customize detection rules, allowlists, ignored directories, and severity levels, `envguard` will support project configuration files (`.envguard.yaml` and `.envguard.yml`) and an explicit `--config` CLI flag. Configuration parsing uses strict decoding via `gopkg.in/yaml.v3`, failing fast on invalid syntax or unknown fields to prevent silent security misconfigurations. User-provided allowlists and ignore directories append to the built-in safe defaults by default.
