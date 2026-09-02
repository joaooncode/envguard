package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joaooncode/envguard/internal/cli"
)

func TestFixCommand_BasicRemediation(t *testing.T) {
	tempDir := t.TempDir()

	// Create an unprotected .env file
	envPath := filepath.Join(tempDir, ".env")
	if err := os.WriteFile(envPath, []byte("SECRET=123\n"), 0644); err != nil {
		t.Fatalf("failed to create env file: %v", err)
	}

	var stdout, stderr bytes.Buffer
	app := cli.New(&stdout, &stderr)

	code := app.Run([]string{"fix", "--path", tempDir, "--no-color"})
	if code != cli.ExitCodeSuccess {
		t.Fatalf("expected exit code %d, got %d. stderr: %s", cli.ExitCodeSuccess, code, stderr.String())
	}

	gitignorePath := filepath.Join(tempDir, ".gitignore")
	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("failed to read .gitignore: %v", err)
	}

	if !strings.Contains(string(content), ".env") {
		t.Errorf("expected .gitignore to contain .env, got:\n%s", string(content))
	}

	outStr := stdout.String()
	if !strings.Contains(outStr, ".env") {
		t.Errorf("expected stdout to mention .env, got:\n%s", outStr)
	}
}

func TestFixCommand_DryRun(t *testing.T) {
	tempDir := t.TempDir()

	envPath := filepath.Join(tempDir, ".env")
	if err := os.WriteFile(envPath, []byte("SECRET=123\n"), 0644); err != nil {
		t.Fatalf("failed to create env file: %v", err)
	}

	var stdout, stderr bytes.Buffer
	app := cli.New(&stdout, &stderr)

	code := app.Run([]string{"fix", "--path", tempDir, "--dry-run", "--no-color"})
	if code != cli.ExitCodeSuccess {
		t.Fatalf("expected exit code %d, got %d. stderr: %s", cli.ExitCodeSuccess, code, stderr.String())
	}

	gitignorePath := filepath.Join(tempDir, ".gitignore")
	if _, err := os.Stat(gitignorePath); !os.IsNotExist(err) {
		t.Errorf(".gitignore should not exist in dry-run mode")
	}

	outStr := stdout.String()
	if !strings.Contains(outStr, "dry run") && !strings.Contains(outStr, "Dry run") {
		t.Errorf("expected stdout to mention dry run, got:\n%s", outStr)
	}
	if !strings.Contains(outStr, ".env") {
		t.Errorf("expected stdout to mention proposed rule .env, got:\n%s", outStr)
	}
}

func TestFixCommand_NoFindings(t *testing.T) {
	tempDir := t.TempDir()

	var stdout, stderr bytes.Buffer
	app := cli.New(&stdout, &stderr)

	code := app.Run([]string{"fix", "--path", tempDir, "--no-color"})
	if code != cli.ExitCodeSuccess {
		t.Fatalf("expected exit code %d, got %d. stderr: %s", cli.ExitCodeSuccess, code, stderr.String())
	}

	outStr := stdout.String()
	if !strings.Contains(outStr, "No unprotected environment files found") {
		t.Errorf("expected stdout to indicate no unprotected files, got:\n%s", outStr)
	}
}

func TestFixCommand_InvalidFlags(t *testing.T) {
	var stdout, stderr bytes.Buffer
	app := cli.New(&stdout, &stderr)

	code := app.Run([]string{"fix", "--unknown-flag"})
	if code != cli.ExitCodeUsageError {
		t.Fatalf("expected exit code %d for invalid flags, got %d", cli.ExitCodeUsageError, code)
	}
}

func TestFixCommand_SkippedExistingRules(t *testing.T) {
	tempDir := t.TempDir()

	envPath := filepath.Join(tempDir, ".env")
	if err := os.WriteFile(envPath, []byte("SECRET=123\n"), 0644); err != nil {
		t.Fatalf("failed to create env file: %v", err)
	}

	gitignorePath := filepath.Join(tempDir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte("# Added by envguard\n.env\n"), 0644); err != nil {
		t.Fatalf("failed to create gitignore: %v", err)
	}

	var stdout, stderr bytes.Buffer
	app := cli.New(&stdout, &stderr)

	code := app.Run([]string{"fix", "--path", tempDir, "--no-color"})
	if code != cli.ExitCodeSuccess {
		t.Fatalf("expected exit code %d, got %d", cli.ExitCodeSuccess, code)
	}

	outStr := stdout.String()
	if !strings.Contains(outStr, "already present") {
		t.Errorf("expected stdout to mention already present rules, got:\n%s", outStr)
	}
}
