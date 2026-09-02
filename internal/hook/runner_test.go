package hook

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/joaooncode/envguard/internal/config"
	"github.com/joaooncode/envguard/internal/detector"
	"github.com/joaooncode/envguard/internal/git"
	"github.com/joaooncode/envguard/internal/scanner"
)

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_CONFIG_NOSYSTEM=1")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v, out: %s", args, err, string(out))
	}
}

func TestRunnerStagedCheck(t *testing.T) {
	repoDir := setupTestGitRepo(t)
	runGit(t, repoDir, "config", "user.email", "test@envguard.dev")
	runGit(t, repoDir, "config", "user.name", "Envguard Test")

	r := NewRunner(git.NewClient(), nil, nil)

	// 1. Empty stage - no findings
	findings, err := r.RunStagedCheck(repoDir)
	if err != nil {
		t.Fatalf("unexpected error on empty stage: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected 0 findings on empty stage, got %d", len(findings))
	}

	// 2. Stage safe file (.env.example) and regular file (main.go)
	examplePath := filepath.Join(repoDir, ".env.example")
	mainPath := filepath.Join(repoDir, "main.go")
	_ = os.WriteFile(examplePath, []byte("API_KEY=\n"), 0644)
	_ = os.WriteFile(mainPath, []byte("package main\n"), 0644)
	runGit(t, repoDir, "add", ".env.example", "main.go")

	findings, err = r.RunStagedCheck(repoDir)
	if err != nil {
		t.Fatalf("unexpected error checking allowed staged files: %v", err)
	}
	if len(findings) != 0 {
		t.Errorf("expected 0 findings when only safe files are staged, got %d", len(findings))
	}

	// 3. Stage sensitive .env file and nested .env.production
	secretPath := filepath.Join(repoDir, ".env")
	nestedDir := filepath.Join(repoDir, "services", "api")
	_ = os.MkdirAll(nestedDir, 0755)
	nestedSecret := filepath.Join(nestedDir, ".env.production")

	_ = os.WriteFile(secretPath, []byte("SECRET=123\n"), 0644)
	_ = os.WriteFile(nestedSecret, []byte("PROD_KEY=xyz\n"), 0644)
	runGit(t, repoDir, "add", ".env", "services/api/.env.production")

	findings, err = r.RunStagedCheck(repoDir)
	if err != nil {
		t.Fatalf("failed to run staged check: %v", err)
	}
	if len(findings) != 2 {
		t.Fatalf("expected 2 findings for staged secrets, got %d", len(findings))
	}

	for _, f := range findings {
		if f.Severity != scanner.SeverityHigh {
			t.Errorf("expected SeverityHigh, got %s for %s", f.Severity, f.Path)
		}
		if !f.GitStatus.IsStaged {
			t.Errorf("expected IsStaged=true for %s", f.Path)
		}
	}

	// 4. Test with custom severity override
	cfg := config.NewDefault()
	cfg.Detector.SeverityOverrides = []config.SeverityOverride{
		{
			Pattern:  ".env",
			Severity: "critical",
		},
	}
	det := detector.NewWithPatterns(cfg.Detector.CustomPatterns, cfg.Detector.Allowlist)
	rWithConfig := NewRunner(git.NewClient(), det, cfg)

	findingsOverride, err := rWithConfig.RunStagedCheck(repoDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	foundCritical := false
	for _, f := range findingsOverride {
		if f.Path == ".env" && f.Severity == scanner.SeverityCritical {
			foundCritical = true
		}
	}
	if !foundCritical {
		t.Errorf("expected .env to be overridden to critical severity")
	}

	// 5. Test non-git directory
	tempNonGit := t.TempDir()
	_, err = r.RunStagedCheck(tempNonGit)
	if err == nil {
		t.Errorf("expected error running on non-git directory")
	}
}
