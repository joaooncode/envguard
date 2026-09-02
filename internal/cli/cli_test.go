package cli_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joaooncode/envguard/internal/cli"
	"github.com/joaooncode/envguard/internal/reporter"
)

func TestCLIVersion(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"version subcommand", []string{"version"}},
		{"-v flag", []string{"-v"}},
		{"--version flag", []string{"--version"}},
		{"-version flag", []string{"-version"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := cli.Run(tt.args, &stdout, &stderr)

			if code != cli.ExitCodeSuccess {
				t.Fatalf("expected exit code %d, got %d. stderr: %s", cli.ExitCodeSuccess, code, stderr.String())
			}

			out := stdout.String()
			if !strings.Contains(out, "envguard v0.1.0") {
				t.Fatalf("expected stdout to contain 'envguard v0.1.0', got: %s", out)
			}
		})
	}
}

func TestCLIHelp(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"no args", []string{}},
		{"help subcommand", []string{"help"}},
		{"-h flag", []string{"-h"}},
		{"--help flag", []string{"--help"}},
		{"-help flag", []string{"-help"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := cli.Run(tt.args, &stdout, &stderr)

			if code != cli.ExitCodeSuccess {
				t.Fatalf("expected exit code %d, got %d", cli.ExitCodeSuccess, code)
			}

			out := stdout.String()
			if !strings.Contains(out, "Usage:") || !strings.Contains(out, "scan") || !strings.Contains(out, "check") {
				t.Fatalf("expected help text in stdout, got: %s", out)
			}
		})
	}
}

func TestCLIUnknownCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"unknowncmd"}, &stdout, &stderr)

	if code != cli.ExitCodeUsageError {
		t.Fatalf("expected exit code %d, got %d", cli.ExitCodeUsageError, code)
	}

	errOut := stderr.String()
	if !strings.Contains(errOut, "unknown command or flag") {
		t.Fatalf("expected error message in stderr, got: %s", errOut)
	}
}

func TestCLIScanFlagErrors(t *testing.T) {
	t.Run("invalid flag", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		code := cli.Run([]string{"scan", "--nonexistent-flag"}, &stdout, &stderr)
		if code != cli.ExitCodeUsageError {
			t.Fatalf("expected exit code %d, got %d", cli.ExitCodeUsageError, code)
		}
	})

	t.Run("invalid format", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		code := cli.Run([]string{"scan", "--format", "yaml"}, &stdout, &stderr)
		if code != cli.ExitCodeUsageError {
			t.Fatalf("expected exit code %d, got %d", cli.ExitCodeUsageError, code)
		}
		if !strings.Contains(stderr.String(), "invalid format") {
			t.Fatalf("expected format error in stderr, got: %s", stderr.String())
		}
	})

	t.Run("invalid severity", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		code := cli.Run([]string{"scan", "--severity", "extreme"}, &stdout, &stderr)
		if code != cli.ExitCodeUsageError {
			t.Fatalf("expected exit code %d, got %d", cli.ExitCodeUsageError, code)
		}
		if !strings.Contains(stderr.String(), "invalid severity") {
			t.Fatalf("expected severity error in stderr, got: %s", stderr.String())
		}
	})
}

func TestCLIScanNonExistentDirectory(t *testing.T) {
	var stdout, stderr bytes.Buffer
	nonExistentPath := filepath.Join(t.TempDir(), "non_existent_folder")
	code := cli.Run([]string{"scan", "--path", nonExistentPath}, &stdout, &stderr)

	if code != cli.ExitCodeInternalError {
		t.Fatalf("expected exit code %d, got %d", cli.ExitCodeInternalError, code)
	}

	if !strings.Contains(stderr.String(), "failed to scan directory") {
		t.Fatalf("expected scan failure message in stderr, got: %s", stderr.String())
	}
}

func TestCLIScanCleanDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	// Add an allowed template file
	err := os.WriteFile(filepath.Join(tmpDir, ".env.example"), []byte("PORT=8080"), 0644)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	t.Run("terminal format default", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		code := cli.Run([]string{"scan", "--path", tmpDir, "--no-color"}, &stdout, &stderr)

		if code != cli.ExitCodeSuccess {
			t.Fatalf("expected exit code %d, got %d. stderr: %s", cli.ExitCodeSuccess, code, stderr.String())
		}

		out := stdout.String()
		if !strings.Contains(out, "PASSED") {
			t.Fatalf("expected PASSED status in stdout, got: %s", out)
		}
	})

	t.Run("json format", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		code := cli.Run([]string{"scan", "--path", tmpDir, "--format", "json"}, &stdout, &stderr)

		if code != cli.ExitCodeSuccess {
			t.Fatalf("expected exit code %d, got %d. stderr: %s", cli.ExitCodeSuccess, code, stderr.String())
		}

		var report reporter.JSONReport
		if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
			t.Fatalf("failed to parse json output: %v. Output:\n%s", err, stdout.String())
		}

		if !report.Summary.Passed {
			t.Fatalf("expected report summary passed to be true, got false")
		}
		if report.Version != "0.1.0" {
			t.Fatalf("expected report version 0.1.0, got %s", report.Version)
		}
	})
}

func TestCLIScanDirectoryWithFindings(t *testing.T) {
	tmpDir := t.TempDir()

	// Create an unprotected .env file (not ignored by git)
	err := os.WriteFile(filepath.Join(tmpDir, ".env"), []byte("SECRET_KEY=12345"), 0644)
	if err != nil {
		t.Fatalf("failed to create temp .env: %v", err)
	}

	t.Run("terminal format failure", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		code := cli.Run([]string{"scan", "--path", tmpDir, "--no-color"}, &stdout, &stderr)

		if code != cli.ExitCodeFindingsFound {
			t.Fatalf("expected exit code %d, got %d", cli.ExitCodeFindingsFound, code)
		}

		out := stdout.String()
		if !strings.Contains(out, "FAILED") || !strings.Contains(out, ".env") {
			t.Fatalf("expected FAILED and .env in stdout, got: %s", out)
		}
	})

	t.Run("json format failure", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		code := cli.Run([]string{"scan", "--path", tmpDir, "--format", "json"}, &stdout, &stderr)

		if code != cli.ExitCodeFindingsFound {
			t.Fatalf("expected exit code %d, got %d", cli.ExitCodeFindingsFound, code)
		}

		var report reporter.JSONReport
		if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
			t.Fatalf("failed to parse json output: %v. Output:\n%s", err, stdout.String())
		}

		if report.Summary.Passed {
			t.Fatalf("expected report summary passed to be false")
		}
		if len(report.Findings) != 1 {
			t.Fatalf("expected 1 finding, got %d", len(report.Findings))
		}
	})

	t.Run("severity threshold filtering", func(t *testing.T) {
		// An unignored .env in a non-git dir yields WARNING severity.
		// Filtering for CRITICAL severity should filter it out and result in PASSED.
		var stdout, stderr bytes.Buffer
		code := cli.Run([]string{"scan", "--path", tmpDir, "--severity", "critical"}, &stdout, &stderr)

		if code != cli.ExitCodeSuccess {
			t.Fatalf("expected exit code %d when filtering for critical, got %d", cli.ExitCodeSuccess, code)
		}

		// Filtering for WARNING severity should retain it and result in failure.
		stdout.Reset()
		stderr.Reset()
		code = cli.Run([]string{"scan", "--path", tmpDir, "--severity", "warning"}, &stdout, &stderr)

		if code != cli.ExitCodeFindingsFound {
			t.Fatalf("expected exit code %d when filtering for warning, got %d", cli.ExitCodeFindingsFound, code)
		}
	})
}

func TestCLICheckCommand(t *testing.T) {
	t.Run("flag error", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		code := cli.Run([]string{"check", "--bad-flag"}, &stdout, &stderr)
		if code != cli.ExitCodeUsageError {
			t.Fatalf("expected exit code %d, got %d", cli.ExitCodeUsageError, code)
		}
	})

	t.Run("non existent directory", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		code := cli.Run([]string{"check", "--path", "/path/that/does/not/exist"}, &stdout, &stderr)
		if code != cli.ExitCodeInternalError {
			t.Fatalf("expected exit code %d, got %d", cli.ExitCodeInternalError, code)
		}
	})

	t.Run("clean directory check", func(t *testing.T) {
		tmpDir := t.TempDir()
		var stdout, stderr bytes.Buffer
		code := cli.Run([]string{"check", "-p", tmpDir, "--no-color"}, &stdout, &stderr)

		if code != cli.ExitCodeSuccess {
			t.Fatalf("expected exit code %d, got %d", cli.ExitCodeSuccess, code)
		}
		if !strings.Contains(stdout.String(), "PASSED") {
			t.Fatalf("expected PASSED status in output, got: %s", stdout.String())
		}
	})

	t.Run("violation directory check", func(t *testing.T) {
		tmpDir := t.TempDir()
		err := os.WriteFile(filepath.Join(tmpDir, ".env.production"), []byte("PROD=true"), 0644)
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}

		var stdout, stderr bytes.Buffer
		code := cli.Run([]string{"check", "-p", tmpDir, "-f", "json"}, &stdout, &stderr)

		if code != cli.ExitCodeFindingsFound {
			t.Fatalf("expected exit code %d, got %d", cli.ExitCodeFindingsFound, code)
		}

		var report reporter.JSONReport
		if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
			t.Fatalf("failed to parse json output: %v", err)
		}
		if report.Summary.Passed {
			t.Fatalf("expected summary passed to be false")
		}
	})
}

func TestCLIConfigFileIntegration(t *testing.T) {
	t.Run("auto-discovered .envguard.yaml with allowlist", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Write .envguard.yaml allowing .env.custom
		configContent := `
detector:
  allowlist:
    - ".env.custom"
`
		if err := os.WriteFile(filepath.Join(tmpDir, ".envguard.yaml"), []byte(configContent), 0644); err != nil {
			t.Fatal(err)
		}

		// Write .env.custom (which would normally fail without config)
		if err := os.WriteFile(filepath.Join(tmpDir, ".env.custom"), []byte("CUSTOM=1"), 0644); err != nil {
			t.Fatal(err)
		}

		var stdout, stderr bytes.Buffer
		code := cli.Run([]string{"scan", "--path", tmpDir, "--no-color"}, &stdout, &stderr)

		if code != cli.ExitCodeSuccess {
			t.Fatalf("expected exit code %d (PASSED due to allowlist in .envguard.yaml), got %d. stderr: %s", cli.ExitCodeSuccess, code, stderr.String())
		}
		if !strings.Contains(stdout.String(), "PASSED") {
			t.Fatalf("expected PASSED in output, got: %s", stdout.String())
		}
	})

	t.Run("explicit --config flag", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "custom-rules.yaml")

		configContent := `
detector:
  custom_patterns:
    - "*.env.vault"
`
		if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(filepath.Join(tmpDir, "app.env.vault"), []byte("SECRET=1"), 0644); err != nil {
			t.Fatal(err)
		}

		var stdout, stderr bytes.Buffer
		code := cli.Run([]string{"check", "--path", tmpDir, "--config", configFile, "--no-color"}, &stdout, &stderr)

		if code != cli.ExitCodeFindingsFound {
			t.Fatalf("expected exit code %d (FINDINGS due to custom pattern), got %d", cli.ExitCodeFindingsFound, code)
		}
		if !strings.Contains(stdout.String(), "app.env.vault") {
			t.Fatalf("expected app.env.vault in output, got: %s", stdout.String())
		}
	})

	t.Run("invalid config file returns UsageError", func(t *testing.T) {
		tmpDir := t.TempDir()
		configFile := filepath.Join(tmpDir, "bad.yaml")
		_ = os.WriteFile(configFile, []byte("scanner: [invalid"), 0644)

		var stdout, stderr bytes.Buffer
		code := cli.Run([]string{"scan", "--path", tmpDir, "--config", configFile}, &stdout, &stderr)

		if code != cli.ExitCodeUsageError {
			t.Fatalf("expected exit code %d for invalid config file, got %d", cli.ExitCodeUsageError, code)
		}
		if !strings.Contains(stderr.String(), "Error:") {
			t.Fatalf("expected error message in stderr, got: %s", stderr.String())
		}
	})
}

