package initializer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joaooncode/envguard/internal/config"
)

func TestSanitizeEnv(t *testing.T) {
	input := `# Project configuration
PORT=8080
DATABASE_URL="postgres://user:pass@localhost:5432/db"

# Secret credentials
SECRET_KEY='super-secret-123'
JWT_TOKEN=xyz.abc.123
export API_ENDPOINT=https://api.example.com
export API_KEY=secret_key

# Empty variable
EMPTY_VAR=

# Comments and whitespaces

ANOTHER_VAR=123 # inline comments
`

	expected := `# Project configuration
PORT=
DATABASE_URL=

# Secret credentials
SECRET_KEY=
JWT_TOKEN=
export API_ENDPOINT=
export API_KEY=

# Empty variable
EMPTY_VAR=

# Comments and whitespaces

ANOTHER_VAR=
`

	sanitized := SanitizeEnv(strings.NewReader(input))
	if sanitized != expected {
		t.Errorf("SanitizeEnv mismatch.\nExpected:\n%s\nGot:\n%s", expected, sanitized)
	}
}

func TestGenerateConfig_Success(t *testing.T) {
	tmpDir := t.TempDir()

	err := GenerateConfig(tmpDir, false)
	if err != nil {
		t.Fatalf("unexpected error generating config: %v", err)
	}

	configPath := filepath.Join(tmpDir, ".envguard.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read generated config: %v", err)
	}

	// Verify that generated config is valid YAML accepted by our config parser
	cfg, err := config.Parse(data)
	if err != nil {
		t.Fatalf("generated config failed strict YAML validation: %v", err)
	}

	if cfg.Version != "1" {
		t.Errorf("expected version '1', got %q", cfg.Version)
	}
	if len(cfg.Scanner.IgnoreDirs) == 0 {
		t.Errorf("expected default ignore dirs in generated config")
	}
	if len(cfg.Detector.Allowlist) == 0 {
		t.Errorf("expected default allowlist in generated config")
	}
}

func TestGenerateConfig_CollisionWithoutForce(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".envguard.yaml")
	if err := os.WriteFile(configPath, []byte("version: '1'\n"), 0644); err != nil {
		t.Fatalf("failed to write initial file: %v", err)
	}

	err := GenerateConfig(tmpDir, false)
	if err == nil {
		t.Fatal("expected collision error when generating existing config without force, got nil")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("expected error message to mention 'already exists', got: %v", err)
	}
}

func TestGenerateConfig_CollisionWithForce(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".envguard.yaml")
	if err := os.WriteFile(configPath, []byte("dummy content\n"), 0644); err != nil {
		t.Fatalf("failed to write initial file: %v", err)
	}

	err := GenerateConfig(tmpDir, true)
	if err != nil {
		t.Fatalf("unexpected error overwriting config with force: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read overwritten config: %v", err)
	}

	if string(data) == "dummy content\n" {
		t.Error("expected config content to be updated with template")
	}
}

func TestGenerateTemplate_FromExistingEnv(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	envContent := "# App\nAPP_NAME=my-app\nSECRET=supersecret\n"
	if err := os.WriteFile(envPath, []byte(envContent), 0644); err != nil {
		t.Fatalf("failed to write .env: %v", err)
	}

	err := GenerateTemplate(tmpDir, "", false)
	if err != nil {
		t.Fatalf("unexpected error generating template from .env: %v", err)
	}

	templatePath := filepath.Join(tmpDir, ".env.example")
	data, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("failed to read .env.example: %v", err)
	}

	expected := "# App\nAPP_NAME=\nSECRET=\n"
	if string(data) != expected {
		t.Errorf("template content mismatch.\nExpected:\n%s\nGot:\n%s", expected, string(data))
	}
}

func TestGenerateTemplate_FromExplicitSource(t *testing.T) {
	tmpDir := t.TempDir()
	customEnvPath := filepath.Join(tmpDir, ".env.staging")
	envContent := "DB_HOST=localhost\nDB_PASS=123456\n"
	if err := os.WriteFile(customEnvPath, []byte(envContent), 0644); err != nil {
		t.Fatalf("failed to write custom env: %v", err)
	}

	err := GenerateTemplate(tmpDir, customEnvPath, false)
	if err != nil {
		t.Fatalf("unexpected error generating template from custom source: %v", err)
	}

	templatePath := filepath.Join(tmpDir, ".env.example")
	data, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("failed to read .env.example: %v", err)
	}

	expected := "DB_HOST=\nDB_PASS=\n"
	if string(data) != expected {
		t.Errorf("template content mismatch.\nExpected:\n%s\nGot:\n%s", expected, string(data))
	}
}

func TestGenerateTemplate_DefaultBoilerplateWhenNoEnv(t *testing.T) {
	tmpDir := t.TempDir()

	err := GenerateTemplate(tmpDir, "", false)
	if err != nil {
		t.Fatalf("unexpected error generating boilerplate template: %v", err)
	}

	templatePath := filepath.Join(tmpDir, ".env.example")
	data, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("failed to read .env.example: %v", err)
	}

	if !strings.Contains(string(data), "PORT=") {
		t.Errorf("expected boilerplate to contain PORT=, got: %s", string(data))
	}
}

func TestGenerateTemplate_NonExistentExplicitSource(t *testing.T) {
	tmpDir := t.TempDir()
	nonExistent := filepath.Join(tmpDir, "missing.env")

	err := GenerateTemplate(tmpDir, nonExistent, false)
	if err == nil {
		t.Fatal("expected error when explicit source file does not exist, got nil")
	}
}

func TestGenerateTemplate_CollisionWithoutForce(t *testing.T) {
	tmpDir := t.TempDir()
	templatePath := filepath.Join(tmpDir, ".env.example")
	if err := os.WriteFile(templatePath, []byte("EXISTING=true\n"), 0644); err != nil {
		t.Fatalf("failed to write existing template: %v", err)
	}

	err := GenerateTemplate(tmpDir, "", false)
	if err == nil {
		t.Fatal("expected collision error when .env.example already exists, got nil")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("expected error message to mention 'already exists', got: %v", err)
	}
}
