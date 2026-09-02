package scanner

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joaooncode/envguard/internal/config"
	"github.com/joaooncode/envguard/internal/detector"
	"github.com/joaooncode/envguard/internal/git"
)

// setupGitRepo creates a temporary initialized git repository.
func setupGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "test@envguard.dev")
	runGit(t, dir, "config", "user.name", "Envguard Test")
	runGit(t, dir, "config", "commit.gpgsign", "false")

	return dir
}

func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_CONFIG_NOSYSTEM=1")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\nOutput: %s", args, err, string(out))
	}
	return string(out)
}

func TestCalculateSummary(t *testing.T) {
	findings := []Finding{
		{Severity: SeverityCritical},
		{Severity: SeverityCritical},
		{Severity: SeverityHigh},
		{Severity: SeverityWarning},
		{Severity: SeverityWarning},
		{Severity: SeverityWarning},
		{Severity: SeverityInfo},
	}

	summary := CalculateSummary(findings)

	if summary.Critical != 2 {
		t.Errorf("expected Critical = 2, got %d", summary.Critical)
	}
	if summary.High != 1 {
		t.Errorf("expected High = 1, got %d", summary.High)
	}
	if summary.Warning != 3 {
		t.Errorf("expected Warning = 3, got %d", summary.Warning)
	}
	if summary.Info != 1 {
		t.Errorf("expected Info = 1, got %d", summary.Info)
	}
	if summary.Total != 7 {
		t.Errorf("expected Total = 7, got %d", summary.Total)
	}

	// Test Result method CalculateSummary
	res := &Result{Findings: findings}
	resSummary := res.CalculateSummary()
	if resSummary.Total != 7 || res.Summary.Total != 7 {
		t.Errorf("expected Result.CalculateSummary to update Summary, got %d", res.Summary.Total)
	}
}

func TestClassifyFinding(t *testing.T) {
	s := NewDefault()

	tests := []struct {
		name                string
		path                string
		status              git.FileStatus
		isAllowed           bool
		expectedSeverity    Severity
		expectedMessage     string
		expectedSuggestions []string
	}{
		{
			name: "Tracked env file (CRITICAL)",
			path: ".env",
			status: git.FileStatus{
				IsRepo:    true,
				IsTracked: true,
			},
			isAllowed:        false,
			expectedSeverity: SeverityCritical,
			expectedMessage:  "Environment file is tracked by Git (committed in repository history).",
			expectedSuggestions: []string{
				"Remove file from git tracking: git rm --cached .env",
				"Add to .gitignore",
				"Rotate any leaked credentials",
			},
		},
		{
			name: "Staged env file (HIGH)",
			path: ".env.staging",
			status: git.FileStatus{
				IsRepo:    true,
				IsTracked: false,
				IsStaged:  true,
			},
			isAllowed:        false,
			expectedSeverity: SeverityHigh,
			expectedMessage:  "Environment file is staged for commit in Git index.",
			expectedSuggestions: []string{
				"Unstage file: git restore --staged .env.staging",
				"Add to .gitignore",
			},
		},
		{
			name: "Untracked, unignored env file (WARNING)",
			path: ".env.local",
			status: git.FileStatus{
				IsRepo:    true,
				IsTracked: false,
				IsStaged:  false,
				IsIgnored: false,
			},
			isAllowed:        false,
			expectedSeverity: SeverityWarning,
			expectedMessage:  "Environment file exists locally and is not ignored by .gitignore.",
			expectedSuggestions: []string{
				"Add to .gitignore",
			},
		},
		{
			name: "Properly ignored env file (INFO)",
			path: ".env.local",
			status: git.FileStatus{
				IsRepo:    true,
				IsTracked: false,
				IsStaged:  false,
				IsIgnored: true,
			},
			isAllowed:           false,
			expectedSeverity:    SeverityInfo,
			expectedMessage:     "Environment file is properly ignored by .gitignore.",
			expectedSuggestions: []string{},
		},
		{
			name: "Allowed template file (INFO)",
			path: ".env.example",
			status: git.FileStatus{
				IsRepo:    true,
				IsTracked: true,
			},
			isAllowed:           true,
			expectedSeverity:    SeverityInfo,
			expectedMessage:     "Safe environment template/example file allowed.",
			expectedSuggestions: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finding := s.classifyFinding(tt.path, tt.status, tt.isAllowed)

			if finding.Severity != tt.expectedSeverity {
				t.Errorf("Severity = %s, want %s", finding.Severity, tt.expectedSeverity)
			}
			if finding.Message != tt.expectedMessage {
				t.Errorf("Message = %q, want %q", finding.Message, tt.expectedMessage)
			}
			if len(finding.Suggestions) != len(tt.expectedSuggestions) {
				t.Fatalf("Suggestions length = %d, want %d", len(finding.Suggestions), len(tt.expectedSuggestions))
			}
			for i, sug := range tt.expectedSuggestions {
				if finding.Suggestions[i] != sug {
					t.Errorf("Suggestion[%d] = %q, want %q", i, finding.Suggestions[i], sug)
				}
			}
		})
	}
}

func TestScanIgnoredDirectories(t *testing.T) {
	tempDir := t.TempDir()

	// Create root .env
	if err := os.WriteFile(filepath.Join(tempDir, ".env"), []byte("ROOT=1"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create ignored directories with .env files
	ignoredDirs := []string{
		".git",
		"node_modules",
		"vendor",
		"dist",
		"build",
		".idea",
		".vscode",
		filepath.Join("nested", "node_modules"),
	}

	for _, d := range ignoredDirs {
		dirPath := filepath.Join(tempDir, d)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			t.Fatal(err)
		}
		envPath := filepath.Join(dirPath, ".env")
		if err := os.WriteFile(envPath, []byte("SECRET=1"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	s := NewDefault()
	result, err := s.Scan(tempDir)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// Only root .env should be found
	if len(result.Findings) != 1 {
		t.Fatalf("expected exactly 1 finding (root .env), got %d findings: %+v", len(result.Findings), result.Findings)
	}

	if result.Findings[0].Path != ".env" {
		t.Errorf("expected finding path to be '.env', got %s", result.Findings[0].Path)
	}
}

func TestScanGitRepositoryIntegration(t *testing.T) {
	repoDir := setupGitRepo(t)

	// 1. CRITICAL: committed .env
	trackedFile := filepath.Join(repoDir, ".env.production")
	if err := os.WriteFile(trackedFile, []byte("PROD=secret"), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, repoDir, "add", ".env.production")
	runGit(t, repoDir, "commit", "-m", "chore: add prod env")

	// 2. HIGH: staged .env
	stagedFile := filepath.Join(repoDir, ".env.staging")
	if err := os.WriteFile(stagedFile, []byte("STAGE=secret"), 0644); err != nil {
		t.Fatal(err)
	}
	runGit(t, repoDir, "add", ".env.staging")

	// 3. INFO: .gitignore + ignored .env
	gitignoreFile := filepath.Join(repoDir, ".gitignore")
	if err := os.WriteFile(gitignoreFile, []byte(".env.local\n"), 0644); err != nil {
		t.Fatal(err)
	}
	ignoredFile := filepath.Join(repoDir, ".env.local")
	if err := os.WriteFile(ignoredFile, []byte("LOCAL=1"), 0644); err != nil {
		t.Fatal(err)
	}

	// 4. WARNING: untracked .env not in .gitignore
	warningFile := filepath.Join(repoDir, ".env.test")
	if err := os.WriteFile(warningFile, []byte("TEST=1"), 0644); err != nil {
		t.Fatal(err)
	}

	// 5. INFO: allowed template
	allowedFile := filepath.Join(repoDir, ".env.example")
	if err := os.WriteFile(allowedFile, []byte("EXAMPLE=1"), 0644); err != nil {
		t.Fatal(err)
	}

	s := NewDefault()
	res, err := s.Scan(repoDir)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if res.Summary.Total != 5 {
		t.Fatalf("expected 5 total findings, got %d (findings: %+v)", res.Summary.Total, res.Findings)
	}
	if res.Summary.Critical != 1 {
		t.Errorf("expected 1 Critical finding, got %d", res.Summary.Critical)
	}
	if res.Summary.High != 1 {
		t.Errorf("expected 1 High finding, got %d", res.Summary.High)
	}
	if res.Summary.Warning != 1 {
		t.Errorf("expected 1 Warning finding, got %d", res.Summary.Warning)
	}
	if res.Summary.Info != 2 {
		t.Errorf("expected 2 Info findings, got %d", res.Summary.Info)
	}

	// Validate finding details
	findingsMap := make(map[string]Finding)
	for _, f := range res.Findings {
		findingsMap[f.Path] = f
	}

	if f, ok := findingsMap[".env.production"]; !ok || f.Severity != SeverityCritical {
		t.Errorf("expected .env.production to be CRITICAL, got %+v", f)
	}
	if f, ok := findingsMap[".env.staging"]; !ok || f.Severity != SeverityHigh {
		t.Errorf("expected .env.staging to be HIGH, got %+v", f)
	}
	if f, ok := findingsMap[".env.test"]; !ok || f.Severity != SeverityWarning {
		t.Errorf("expected .env.test to be WARNING, got %+v", f)
	}
	if f, ok := findingsMap[".env.local"]; !ok || f.Severity != SeverityInfo {
		t.Errorf("expected .env.local to be INFO, got %+v", f)
	}
	if f, ok := findingsMap[".env.example"]; !ok || f.Severity != SeverityInfo {
		t.Errorf("expected .env.example to be INFO, got %+v", f)
	}
}

func TestScanNonGitDirectory(t *testing.T) {
	tempDir := t.TempDir()

	envFile := filepath.Join(tempDir, ".env")
	if err := os.WriteFile(envFile, []byte("FOO=bar"), 0644); err != nil {
		t.Fatal(err)
	}
	exampleFile := filepath.Join(tempDir, ".env.example")
	if err := os.WriteFile(exampleFile, []byte("FOO=example"), 0644); err != nil {
		t.Fatal(err)
	}

	res, err := Scan(tempDir)
	if err != nil {
		t.Fatalf("Scan on non-git dir failed: %v", err)
	}

	if res.Summary.Total != 2 {
		t.Errorf("expected 2 findings, got %d", res.Summary.Total)
	}
	if res.Summary.Warning != 1 {
		t.Errorf("expected 1 Warning finding for .env, got %d", res.Summary.Warning)
	}
	if res.Summary.Info != 1 {
		t.Errorf("expected 1 Info finding for .env.example, got %d", res.Summary.Info)
	}
}

func TestScanInvalidDirectory(t *testing.T) {
	s := NewDefault()

	// Non-existent directory
	_, err := s.Scan("/path/that/does/not/exist/9999")
	if err == nil {
		t.Errorf("expected error scanning non-existent directory")
	}

	// File instead of directory
	tempFile := filepath.Join(t.TempDir(), "file.txt")
	if err := os.WriteFile(tempFile, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err = s.Scan(tempFile)
	if err == nil {
		t.Errorf("expected error scanning regular file instead of directory")
	}
}

// mockGitClient simulates custom git client behaviors.
type mockGitClient struct {
	statusFn func(dir, filePath string) (git.FileStatus, error)
}

func (m *mockGitClient) IsAvailable() bool                            { return true }
func (m *mockGitClient) IsGitRepo(dir string) bool                    { return true }
func (m *mockGitClient) GetRepoRoot(dir string) (string, error)       { return dir, nil }
func (m *mockGitClient) IsTracked(dir, filePath string) (bool, error) { return false, nil }
func (m *mockGitClient) IsStaged(dir, filePath string) (bool, error)  { return false, nil }
func (m *mockGitClient) IsIgnored(dir, filePath string) (bool, error) { return false, nil }
func (m *mockGitClient) GetStagedFiles(dir string) ([]string, error)  { return nil, nil }
func (m *mockGitClient) GetHooksDir(dir string) (string, error) {
	return filepath.Join(dir, ".git", "hooks"), nil
}
func (m *mockGitClient) GetFileStatus(dir, filePath string) (git.FileStatus, error) {
	if m.statusFn != nil {
		return m.statusFn(dir, filePath)
	}
	return git.FileStatus{}, nil
}

func TestScannerCustomComponents(t *testing.T) {
	mock := &mockGitClient{
		statusFn: func(dir, filePath string) (git.FileStatus, error) {
			if strings.Contains(filePath, "tracked") {
				return git.FileStatus{IsRepo: true, IsTracked: true}, nil
			}
			return git.FileStatus{IsRepo: true}, nil
		},
	}

	customDetector := detector.NewWithAllowlist([]string{".env.custom.allowed"})
	s := New(mock, customDetector)

	tempDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(tempDir, ".env.tracked"), []byte(""), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, ".env.custom.allowed"), []byte(""), 0644)

	res, err := s.Scan(tempDir)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	if res.Summary.Critical != 1 {
		t.Errorf("expected 1 Critical finding, got %d", res.Summary.Critical)
	}
	if res.Summary.Info != 1 {
		t.Errorf("expected 1 Info finding for custom allowed file, got %d", res.Summary.Info)
	}
}

func TestScannerWithConfig(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Setup files:
	// a. Custom ignored dir
	customIgnoredDir := filepath.Join(tempDir, "custom_vendor")
	if err := os.MkdirAll(customIgnoredDir, 0755); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(customIgnoredDir, ".env"), []byte("IGNORED=1"), 0644)

	// b. Custom env pattern
	_ = os.WriteFile(filepath.Join(tempDir, "app.env.vault"), []byte("SECRET=1"), 0644)

	// c. Custom allowlist
	_ = os.WriteFile(filepath.Join(tempDir, ".env.dist"), []byte("DIST=1"), 0644)

	// d. Severity override (.env.test normally Warning in non-git, override to Critical)
	_ = os.WriteFile(filepath.Join(tempDir, ".env.test"), []byte("TEST=1"), 0644)

	cfg := &config.Config{
		Scanner: config.ScannerConfig{
			IgnoreDirs: []string{"custom_vendor"},
		},
		Detector: config.DetectorConfig{
			CustomPatterns: []string{"*.env.vault"},
			Allowlist:      []string{".env.dist"},
			SeverityOverrides: []config.SeverityOverride{
				{Pattern: ".env.test", Severity: "critical"},
			},
		},
	}

	s := NewWithConfig(nil, nil, cfg)
	res, err := s.Scan(tempDir)
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	findingsMap := make(map[string]Finding)
	for _, f := range res.Findings {
		findingsMap[f.Path] = f
	}

	// custom_vendor/.env should NOT be found
	if _, ok := findingsMap["custom_vendor/.env"]; ok {
		t.Errorf("expected custom_vendor/.env to be ignored, but was found")
	}

	// app.env.vault should be detected
	if _, ok := findingsMap["app.env.vault"]; !ok {
		t.Errorf("expected app.env.vault to be detected as env file")
	}

	// .env.dist should be allowed (Info)
	if f, ok := findingsMap[".env.dist"]; !ok || f.Severity != SeverityInfo {
		t.Errorf("expected .env.dist to be allowed (Info), got %+v", f)
	}

	// .env.test should have overridden severity Critical
	if f, ok := findingsMap[".env.test"]; !ok || f.Severity != SeverityCritical {
		t.Errorf("expected .env.test to be overridden to Critical, got %+v", f)
	}
}
