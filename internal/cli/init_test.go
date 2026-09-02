package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joaooncode/envguard/internal/cli"
)

func TestCLIInitHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"init", "--help"}, &stdout, &stderr)

	if code != cli.ExitCodeSuccess {
		t.Fatalf("expected exit code %d on init --help, got %d. stderr: %s", cli.ExitCodeSuccess, code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "Usage of init:") {
		t.Errorf("expected usage output in stderr, got: %s", stderr.String())
	}
}

func TestCLIInitDefault(t *testing.T) {
	tmpDir := t.TempDir()
	var stdout, stderr bytes.Buffer

	code := cli.Run([]string{"init", "--path", tmpDir}, &stdout, &stderr)
	if code != cli.ExitCodeSuccess {
		t.Fatalf("expected exit code %d, got %d. stderr: %s", cli.ExitCodeSuccess, code, stderr.String())
	}

	configPath := filepath.Join(tmpDir, ".envguard.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("expected .envguard.yaml to be created at %s", configPath)
	}

	templatePath := filepath.Join(tmpDir, ".env.example")
	if _, err := os.Stat(templatePath); !os.IsNotExist(err) {
		t.Fatalf("expected .env.example NOT to be created when --template is not passed")
	}

	if !strings.Contains(stdout.String(), "Created configuration file:") {
		t.Errorf("expected stdout to report created config, got: %s", stdout.String())
	}
}

func TestCLIInitWithTemplate(t *testing.T) {
	tmpDir := t.TempDir()

	// Create dummy .env
	envPath := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(envPath, []byte("API_SECRET=mysecretvalue\nPORT=4000\n"), 0644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"init", "-p", tmpDir, "--template"}, &stdout, &stderr)
	if code != cli.ExitCodeSuccess {
		t.Fatalf("expected exit code %d, got %d. stderr: %s", cli.ExitCodeSuccess, code, stderr.String())
	}

	templatePath := filepath.Join(tmpDir, ".env.example")
	data, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("failed to read .env.example: %v", err)
	}

	content := string(data)
	if strings.Contains(content, "mysecretvalue") {
		t.Errorf("expected sensitive value to be stripped, got: %s", content)
	}
	if !strings.Contains(content, "API_SECRET=") || !strings.Contains(content, "PORT=") {
		t.Errorf("expected keys to be preserved, got: %s", content)
	}
}

func TestCLIInitWithTemplateFrom(t *testing.T) {
	tmpDir := t.TempDir()
	sourceEnv := filepath.Join(tmpDir, ".env.production")
	if err := os.WriteFile(sourceEnv, []byte("PROD_DB=supersecret\n"), 0644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"init", "-p", tmpDir, "--template-from", sourceEnv}, &stdout, &stderr)
	if code != cli.ExitCodeSuccess {
		t.Fatalf("expected exit code %d, got %d. stderr: %s", cli.ExitCodeSuccess, code, stderr.String())
	}

	templatePath := filepath.Join(tmpDir, ".env.example")
	data, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("failed to read .env.example: %v", err)
	}

	if !strings.Contains(string(data), "PROD_DB=") || strings.Contains(string(data), "supersecret") {
		t.Errorf("unexpected template content: %s", string(data))
	}
}

func TestCLIInitCollisionWithoutForce(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".envguard.yaml")
	if err := os.WriteFile(configPath, []byte("existing config\n"), 0644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"init", "-p", tmpDir}, &stdout, &stderr)
	if code != cli.ExitCodeInternalError {
		t.Fatalf("expected exit code %d on file collision, got %d", cli.ExitCodeInternalError, code)
	}

	if !strings.Contains(stderr.String(), "already exists") {
		t.Errorf("expected stderr to mention already exists, got: %s", stderr.String())
	}
}

func TestCLIInitCollisionWithForce(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".envguard.yaml")
	if err := os.WriteFile(configPath, []byte("existing config\n"), 0644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"init", "-p", tmpDir, "--force"}, &stdout, &stderr)
	if code != cli.ExitCodeSuccess {
		t.Fatalf("expected exit code %d with --force, got %d. stderr: %s", cli.ExitCodeSuccess, code, stderr.String())
	}
}

func TestCLIInitInvalidFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{"init", "--invalid-flag"}, &stdout, &stderr)
	if code != cli.ExitCodeUsageError {
		t.Fatalf("expected exit code %d for invalid flag, got %d", cli.ExitCodeUsageError, code)
	}
}
