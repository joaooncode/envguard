package fixer_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joaooncode/envguard/internal/fixer"
	"github.com/joaooncode/envguard/internal/git"
	"github.com/joaooncode/envguard/internal/scanner"
)

func TestFixer_Apply_NewGitignore(t *testing.T) {
	tempDir := t.TempDir()

	findings := []scanner.Finding{
		{
			Path:      filepath.Join(tempDir, ".env"),
			Severity:  scanner.SeverityWarning,
			Message:   "unprotected local file",
			GitStatus: git.FileStatus{IsTracked: false, IsIgnored: false},
		},
		{
			Path:      filepath.Join(tempDir, "backend", ".env.local"),
			Severity:  scanner.SeverityWarning,
			Message:   "unprotected local file",
			GitStatus: git.FileStatus{IsTracked: false, IsIgnored: false},
		},
	}

	f := fixer.New()
	res, err := f.Apply(fixer.Options{
		TargetDir: tempDir,
		Findings:  findings,
		DryRun:    false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !res.GitignoreUpdated {
		t.Fatalf("expected GitignoreUpdated to be true")
	}

	if len(res.AddedRules) != 2 {
		t.Fatalf("expected 2 added rules, got %d", len(res.AddedRules))
	}

	gitignorePath := filepath.Join(tempDir, ".gitignore")
	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("failed to read created .gitignore: %v", err)
	}

	expectedSubstrings := []string{
		"# Added by envguard",
		".env",
		"/backend/.env.local",
	}

	for _, sub := range expectedSubstrings {
		if !strings.Contains(string(content), sub) {
			t.Errorf("expected .gitignore to contain %q, got:\n%s", sub, string(content))
		}
	}
}

func TestFixer_Apply_ExistingGitignore_AppendSection(t *testing.T) {
	tempDir := t.TempDir()
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	initialContent := "node_modules/\ndist/\n"
	if err := os.WriteFile(gitignorePath, []byte(initialContent), 0644); err != nil {
		t.Fatalf("failed to write initial .gitignore: %v", err)
	}

	findings := []scanner.Finding{
		{
			Path:      filepath.Join(tempDir, ".env"),
			Severity:  scanner.SeverityWarning,
			GitStatus: git.FileStatus{IsTracked: false, IsIgnored: false},
		},
	}

	f := fixer.New()
	res, err := f.Apply(fixer.Options{
		TargetDir: tempDir,
		Findings:  findings,
		DryRun:    false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !res.GitignoreUpdated {
		t.Fatalf("expected GitignoreUpdated to be true")
	}

	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("failed to read .gitignore: %v", err)
	}

	strContent := string(content)
	if !strings.HasPrefix(strContent, "node_modules/\ndist/\n") {
		t.Errorf("expected original content to be preserved at start, got:\n%s", strContent)
	}
	if !strings.Contains(strContent, "# Added by envguard\n.env") {
		t.Errorf("expected envguard section appended, got:\n%s", strContent)
	}
}

func TestFixer_Apply_ExistingEnvguardSection(t *testing.T) {
	tempDir := t.TempDir()
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	initialContent := "# Existing rules\nnode_modules/\n\n# Added by envguard\n.env\n"
	if err := os.WriteFile(gitignorePath, []byte(initialContent), 0644); err != nil {
		t.Fatalf("failed to write initial .gitignore: %v", err)
	}

	findings := []scanner.Finding{
		{
			Path:      filepath.Join(tempDir, ".env"), // duplicate
			Severity:  scanner.SeverityWarning,
			GitStatus: git.FileStatus{IsTracked: false, IsIgnored: false},
		},
		{
			Path:      filepath.Join(tempDir, ".env.local"), // new
			Severity:  scanner.SeverityWarning,
			GitStatus: git.FileStatus{IsTracked: false, IsIgnored: false},
		},
	}

	f := fixer.New()
	res, err := f.Apply(fixer.Options{
		TargetDir: tempDir,
		Findings:  findings,
		DryRun:    false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.AddedRules) != 1 || res.AddedRules[0] != ".env.local" {
		t.Fatalf("expected 1 added rule (.env.local), got: %v", res.AddedRules)
	}
	if len(res.SkippedRules) != 1 || res.SkippedRules[0] != ".env" {
		t.Fatalf("expected 1 skipped rule (.env), got: %v", res.SkippedRules)
	}

	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("failed to read .gitignore: %v", err)
	}

	strContent := string(content)
	// Should not duplicate the section header
	if strings.Count(strContent, "# Added by envguard") != 1 {
		t.Errorf("expected exactly 1 '# Added by envguard' header, got: %d", strings.Count(strContent, "# Added by envguard"))
	}
	if !strings.Contains(strContent, ".env.local") {
		t.Errorf("expected .env.local in .gitignore, got:\n%s", strContent)
	}
}

func TestFixer_Apply_DryRun(t *testing.T) {
	tempDir := t.TempDir()
	gitignorePath := filepath.Join(tempDir, ".gitignore")

	findings := []scanner.Finding{
		{
			Path:      filepath.Join(tempDir, ".env"),
			Severity:  scanner.SeverityWarning,
			GitStatus: git.FileStatus{IsTracked: false, IsIgnored: false},
		},
	}

	f := fixer.New()
	res, err := f.Apply(fixer.Options{
		TargetDir: tempDir,
		Findings:  findings,
		DryRun:    true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !res.DryRun {
		t.Fatalf("expected DryRun to be true")
	}
	if len(res.AddedRules) != 1 {
		t.Fatalf("expected 1 added rule in dry-run result, got %d", len(res.AddedRules))
	}

	// File should NOT exist in dry-run mode
	if _, err := os.Stat(gitignorePath); !os.IsNotExist(err) {
		t.Errorf(".gitignore should not have been created during dry-run")
	}
}

func TestFixer_Apply_CriticalFindings(t *testing.T) {
	tempDir := t.TempDir()

	findings := []scanner.Finding{
		{
			Path:      filepath.Join(tempDir, ".env"),
			Severity:  scanner.SeverityCritical,
			GitStatus: git.FileStatus{IsTracked: true},
		},
		{
			Path:      filepath.Join(tempDir, "config", ".env.local"),
			Severity:  scanner.SeverityWarning,
			GitStatus: git.FileStatus{IsTracked: false, IsIgnored: false},
		},
	}

	f := fixer.New()
	res, err := f.Apply(fixer.Options{
		TargetDir: tempDir,
		Findings:  findings,
		DryRun:    false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.CriticalFindings) != 1 {
		t.Fatalf("expected 1 critical finding, got %d", len(res.CriticalFindings))
	}
	if len(res.AddedRules) != 1 {
		t.Fatalf("expected 1 added rule for warning, got %d", len(res.AddedRules))
	}
}

func TestFixer_Apply_NoTrailingNewlineHandling(t *testing.T) {
	tempDir := t.TempDir()
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	initialContent := "node_modules" // no trailing newline
	if err := os.WriteFile(gitignorePath, []byte(initialContent), 0644); err != nil {
		t.Fatalf("failed to write initial .gitignore: %v", err)
	}

	findings := []scanner.Finding{
		{
			Path:      filepath.Join(tempDir, ".env"),
			Severity:  scanner.SeverityWarning,
			GitStatus: git.FileStatus{IsTracked: false, IsIgnored: false},
		},
	}

	f := fixer.New()
	res, err := f.Apply(fixer.Options{
		TargetDir: tempDir,
		Findings:  findings,
		DryRun:    false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !res.GitignoreUpdated {
		t.Fatalf("expected GitignoreUpdated to be true")
	}

	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("failed to read .gitignore: %v", err)
	}

	strContent := string(content)
	if !strings.HasPrefix(strContent, "node_modules\n\n# Added by envguard\n.env\n") {
		t.Errorf("unexpected content structure:\n%s", strContent)
	}
}
