package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// DefaultConfigFilenames lists the standard config filenames checked in target directories.
var DefaultConfigFilenames = []string{
	".envguard.yaml",
	".envguard.yml",
}

// Config represents the complete configuration structure for envguard.
type Config struct {
	Version  string         `yaml:"version,omitempty"`
	Scanner  ScannerConfig  `yaml:"scanner,omitempty"`
	Detector DetectorConfig `yaml:"detector,omitempty"`
}

// ScannerConfig holds scanner-specific options.
type ScannerConfig struct {
	IgnoreDirs []string `yaml:"ignore_dirs,omitempty"`
}

// DetectorConfig holds detector-specific options.
type DetectorConfig struct {
	CustomPatterns    []string           `yaml:"custom_patterns,omitempty"`
	Allowlist         []string           `yaml:"allowlist,omitempty"`
	SeverityOverrides []SeverityOverride `yaml:"severity_overrides,omitempty"`
}

// SeverityOverride overrides the calculated severity for files matching a pattern.
type SeverityOverride struct {
	Pattern  string `yaml:"pattern"`
	Severity string `yaml:"severity"`
}

// NewDefault returns a new Config with empty/default settings.
func NewDefault() *Config {
	return &Config{
		Version: "1",
		Scanner: ScannerConfig{
			IgnoreDirs: make([]string, 0),
		},
		Detector: DetectorConfig{
			CustomPatterns:    make([]string, 0),
			Allowlist:         make([]string, 0),
			SeverityOverrides: make([]SeverityOverride, 0),
		},
	}
}

// Parse decodes raw YAML bytes into a Config struct using strict field validation.
func Parse(data []byte) (*Config, error) {
	cfg := NewDefault()
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)

	if err := decoder.Decode(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration syntax: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Load reads and parses a YAML configuration file from the specified path.
func Load(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", filePath, err)
	}

	cfg, err := Parse(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", filePath, err)
	}

	return cfg, nil
}

// DiscoverAndLoad resolves configuration either from an explicit path or standard locations in dir.
func DiscoverAndLoad(dir string, explicitPath string) (*Config, bool, error) {
	if explicitPath != "" {
		cfg, err := Load(explicitPath)
		if err != nil {
			return nil, false, err
		}
		return cfg, true, nil
	}

	if dir == "" {
		dir = "."
	}

	for _, filename := range DefaultConfigFilenames {
		candidate := filepath.Join(dir, filename)
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			cfg, err := Load(candidate)
			if err != nil {
				return nil, false, err
			}
			return cfg, true, nil
		}
	}

	return NewDefault(), false, nil
}

// Validate ensures all configuration fields have valid and well-formed values.
func (c *Config) Validate() error {
	for i, override := range c.Detector.SeverityOverrides {
		if strings.TrimSpace(override.Pattern) == "" {
			return fmt.Errorf("severity_overrides[%d]: pattern cannot be empty", i)
		}
		sev := strings.ToLower(strings.TrimSpace(override.Severity))
		switch sev {
		case "info", "warning", "warn", "high", "critical":
			// valid
		default:
			return fmt.Errorf("severity_overrides[%d]: invalid severity %q (supported: info, warning, high, critical)", i, override.Severity)
		}
	}
	return nil
}
