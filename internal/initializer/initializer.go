package initializer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// DefaultConfigContent contains the well-commented default .envguard.yaml configuration.
const DefaultConfigContent = `# envguard configuration file
# For more information, visit https://github.com/joaooncode/envguard

version: "1"

scanner:
  # Directories to ignore during recursive filesystem traversal
  ignore_dirs:
    - node_modules
    - vendor
    - .git

detector:
  # Additional regex or glob patterns to detect as environment files
  custom_patterns: []

  # Safe template patterns or example files that should not raise security warnings
  allowlist:
    - "*.example"
    - "*.sample"
    - "*.template"

  # Explicit severity overrides for specific patterns
  # Supported severities: info, warning, high, critical
  severity_overrides: []
`

// DefaultTemplateContent contains the boilerplate .env.example template.
const DefaultTemplateContent = `# Environment variables example template
# Copy this file to .env and fill in your actual values

# Application
PORT=3000
APP_ENV=development

# Database
DATABASE_URL=

# Authentication & Secrets
API_KEY=
JWT_SECRET=
`

// SanitizeEnv reads an environment file and strips all sensitive values,
// keeping variable names, empty lines, and comments intact.
func SanitizeEnv(r io.Reader) string {
	scanner := bufio.NewScanner(r)
	var sb strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			sb.WriteString(line)
			sb.WriteString("\n")
			continue
		}

		if idx := strings.Index(line, "="); idx != -1 {
			keyPart := strings.TrimRight(line[:idx], " \t")
			sb.WriteString(keyPart)
			sb.WriteString("=\n")
		} else {
			sb.WriteString(line)
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// GenerateConfig creates the default .envguard.yaml in targetDir.
// If the file already exists and force is false, it returns an error.
func GenerateConfig(targetDir string, force bool) error {
	if targetDir == "" {
		targetDir = "."
	}
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory %s: %w", targetDir, err)
	}

	configPath := filepath.Join(targetDir, ".envguard.yaml")
	if !force {
		if _, err := os.Stat(configPath); err == nil {
			return fmt.Errorf("configuration file already exists at %s (use --force to overwrite)", configPath)
		}
	}

	if err := os.WriteFile(configPath, []byte(DefaultConfigContent), 0644); err != nil {
		return fmt.Errorf("failed to write configuration file %s: %w", configPath, err)
	}

	return nil
}

// GenerateTemplate creates a safe .env.example file in targetDir.
// If sourceEnvPath is provided, it sanitizes that file. If empty, it checks
// for an existing .env in targetDir, or falls back to DefaultTemplateContent.
// If the file already exists and force is false, it returns an error.
func GenerateTemplate(targetDir, sourceEnvPath string, force bool) error {
	if targetDir == "" {
		targetDir = "."
	}
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory %s: %w", targetDir, err)
	}

	templatePath := filepath.Join(targetDir, ".env.example")
	if !force {
		if _, err := os.Stat(templatePath); err == nil {
			return fmt.Errorf("template file already exists at %s (use --force to overwrite)", templatePath)
		}
	}

	var content string
	if sourceEnvPath != "" {
		file, err := os.Open(sourceEnvPath)
		if err != nil {
			return fmt.Errorf("failed to open source env file %s: %w", sourceEnvPath, err)
		}
		defer file.Close()
		content = SanitizeEnv(file)
	} else {
		defaultEnv := filepath.Join(targetDir, ".env")
		if info, err := os.Stat(defaultEnv); err == nil && !info.IsDir() {
			file, err := os.Open(defaultEnv)
			if err != nil {
				return fmt.Errorf("failed to open .env file %s: %w", defaultEnv, err)
			}
			defer file.Close()
			content = SanitizeEnv(file)
		} else {
			content = DefaultTemplateContent
		}
	}

	if err := os.WriteFile(templatePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write template file %s: %w", templatePath, err)
	}

	return nil
}
