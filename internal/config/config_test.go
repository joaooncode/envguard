package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := NewDefault()
	if cfg == nil {
		t.Fatal("expected non-nil default config")
	}
	if len(cfg.Scanner.IgnoreDirs) != 0 {
		t.Errorf("expected 0 default custom ignore dirs, got %d", len(cfg.Scanner.IgnoreDirs))
	}
	if len(cfg.Detector.Allowlist) != 0 {
		t.Errorf("expected 0 default custom allowlist patterns, got %d", len(cfg.Detector.Allowlist))
	}
}

func TestParseYAML_Valid(t *testing.T) {
	yamlContent := `
version: "1"
scanner:
  ignore_dirs:
    - .custom_vendor
    - build_artifacts
detector:
  custom_patterns:
    - "*.env.vault"
    - ".env.secret"
  allowlist:
    - ".env.custom.example"
  severity_overrides:
    - pattern: ".env.production"
      severity: "critical"
    - pattern: ".env.local"
      severity: "info"
`
	cfg, err := Parse([]byte(yamlContent))
	if err != nil {
		t.Fatalf("unexpected error parsing valid yaml: %v", err)
	}

	if cfg.Version != "1" {
		t.Errorf("expected version '1', got %q", cfg.Version)
	}
	if len(cfg.Scanner.IgnoreDirs) != 2 {
		t.Errorf("expected 2 ignore dirs, got %d", len(cfg.Scanner.IgnoreDirs))
	}
	if len(cfg.Detector.CustomPatterns) != 2 {
		t.Errorf("expected 2 custom patterns, got %d", len(cfg.Detector.CustomPatterns))
	}
	if len(cfg.Detector.Allowlist) != 1 {
		t.Errorf("expected 1 allowlist pattern, got %d", len(cfg.Detector.Allowlist))
	}
	if len(cfg.Detector.SeverityOverrides) != 2 {
		t.Errorf("expected 2 severity overrides, got %d", len(cfg.Detector.SeverityOverrides))
	}
}

func TestParseYAML_InvalidSyntax(t *testing.T) {
	invalidYAML := `
scanner:
  ignore_dirs: [
`
	_, err := Parse([]byte(invalidYAML))
	if err == nil {
		t.Fatal("expected error parsing invalid YAML syntax, got nil")
	}
}

func TestParseYAML_UnknownFields(t *testing.T) {
	unknownFieldYAML := `
scanner:
  unknown_key: "value"
`
	_, err := Parse([]byte(unknownFieldYAML))
	if err == nil {
		t.Fatal("expected error on unknown fields with strict decoding, got nil")
	}
}

func TestParseYAML_InvalidSeverity(t *testing.T) {
	invalidSevYAML := `
detector:
  severity_overrides:
    - pattern: ".env.production"
      severity: "super-critical"
`
	_, err := Parse([]byte(invalidSevYAML))
	if err == nil {
		t.Fatal("expected error on invalid severity value, got nil")
	}
}

func TestDiscoverAndLoad(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. Without any config file, returns default and found=false, no error
	cfg, found, err := DiscoverAndLoad(tmpDir, "")
	if err != nil {
		t.Fatalf("unexpected error when no config exists: %v", err)
	}
	if found {
		t.Errorf("expected found=false when no config file exists")
	}
	if cfg == nil {
		t.Fatal("expected non-nil default config")
	}

	// 2. Discover .envguard.yaml
	yamlPath := filepath.Join(tmpDir, ".envguard.yaml")
	content := `
scanner:
  ignore_dirs:
    - test_dir
`
	if err := os.WriteFile(yamlPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test .envguard.yaml: %v", err)
	}

	cfg, found, err = DiscoverAndLoad(tmpDir, "")
	if err != nil {
		t.Fatalf("unexpected error discovering .envguard.yaml: %v", err)
	}
	if !found {
		t.Fatal("expected found=true when .envguard.yaml exists")
	}
	if len(cfg.Scanner.IgnoreDirs) != 1 || cfg.Scanner.IgnoreDirs[0] != "test_dir" {
		t.Errorf("expected ignore_dirs to contain 'test_dir', got %v", cfg.Scanner.IgnoreDirs)
	}

	// 3. Explicit config path override
	customPath := filepath.Join(tmpDir, "custom.yaml")
	customContent := `
detector:
  custom_patterns:
    - "*.env.custom"
`
	if err := os.WriteFile(customPath, []byte(customContent), 0644); err != nil {
		t.Fatalf("failed to write custom.yaml: %v", err)
	}

	cfg, found, err = DiscoverAndLoad(tmpDir, customPath)
	if err != nil {
		t.Fatalf("unexpected error loading explicit config: %v", err)
	}
	if !found {
		t.Fatal("expected found=true for explicit config")
	}
	if len(cfg.Detector.CustomPatterns) != 1 || cfg.Detector.CustomPatterns[0] != "*.env.custom" {
		t.Errorf("expected custom_patterns to contain '*.env.custom', got %v", cfg.Detector.CustomPatterns)
	}

	// 4. Non-existent explicit config path returns error
	_, _, err = DiscoverAndLoad(tmpDir, filepath.Join(tmpDir, "non_existent.yaml"))
	if err == nil {
		t.Fatal("expected error when explicit config does not exist, got nil")
	}
}
