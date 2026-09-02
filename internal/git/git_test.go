package git

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// setupTestGitRepo initializes a temporary Git repository for integration testing.
func setupTestGitRepo(t *testing.T) string {
	t.Helper()

	tempDir := t.TempDir()

	// Configure and initialize git repository
	runGitCmd(t, tempDir, "init")
	runGitCmd(t, tempDir, "config", "user.email", "test@envguard.dev")
	runGitCmd(t, tempDir, "config", "user.name", "Envguard Test")
	runGitCmd(t, tempDir, "config", "commit.gpgsign", "false")

	return tempDir
}

func runGitCmd(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_CONFIG_NOSYSTEM=1")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed in %s: %v\nOutput: %s", args, dir, err, string(out))
	}
	return string(out)
}

func TestNonGitDirectory(t *testing.T) {
	tempDir := t.TempDir()
	client := NewClient()

	if client.IsGitRepo(tempDir) {
		t.Errorf("expected IsGitRepo to be false for non-git directory: %s", tempDir)
	}

	_, err := client.GetRepoRoot(tempDir)
	if err == nil {
		t.Errorf("expected error from GetRepoRoot on non-git directory")
	}

	status, err := client.GetFileStatus(tempDir, ".env")
	if err != nil {
		t.Fatalf("GetFileStatus returned error on non-git directory: %v", err)
	}

	if status.IsRepo {
		t.Errorf("expected status.IsRepo to be false")
	}
	if status.IsTracked || status.IsStaged || status.IsIgnored {
		t.Errorf("expected all status flags to be false on non-git dir, got %+v", status)
	}
}

func TestGitRepoStatus(t *testing.T) {
	repoDir := setupTestGitRepo(t)
	client := NewClient()

	if !client.IsGitRepo(repoDir) {
		t.Fatalf("expected %s to be recognized as git repo", repoDir)
	}

	root, err := client.GetRepoRoot(repoDir)
	if err != nil {
		t.Fatalf("failed to get repo root: %v", err)
	}
	if root == "" {
		t.Errorf("expected non-empty repo root")
	}

	// 1. Untracked file (new file, not ignored, not staged, not committed)
	untrackedFile := filepath.Join(repoDir, ".env")
	if err := os.WriteFile(untrackedFile, []byte("SECRET=123\n"), 0644); err != nil {
		t.Fatal(err)
	}

	status, err := client.GetFileStatus(repoDir, ".env")
	if err != nil {
		t.Fatalf("GetFileStatus(.env) failed: %v", err)
	}
	if !status.IsRepo || status.IsTracked || status.IsStaged || status.IsIgnored {
		t.Errorf("expected untracked status (Repo: true, Tracked: false, Staged: false, Ignored: false), got: %+v", status)
	}

	// 2. Ignored file (.gitignore)
	gitignorePath := filepath.Join(repoDir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(".env.local\nconfig/*.env\n"), 0644); err != nil {
		t.Fatal(err)
	}

	ignoredFile := filepath.Join(repoDir, ".env.local")
	if err := os.WriteFile(ignoredFile, []byte("LOCAL=true\n"), 0644); err != nil {
		t.Fatal(err)
	}

	statusIgnored, err := client.GetFileStatus(repoDir, ".env.local")
	if err != nil {
		t.Fatalf("GetFileStatus(.env.local) failed: %v", err)
	}
	if !statusIgnored.IsRepo || statusIgnored.IsTracked || statusIgnored.IsStaged || !statusIgnored.IsIgnored {
		t.Errorf("expected ignored status (Repo: true, Tracked: false, Staged: false, Ignored: true), got: %+v", statusIgnored)
	}

	// Nested ignored file
	subDir := filepath.Join(repoDir, "config")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	nestedIgnored := filepath.Join(subDir, "app.env")
	if err := os.WriteFile(nestedIgnored, []byte("APP=1\n"), 0644); err != nil {
		t.Fatal(err)
	}
	statusNested, err := client.GetFileStatus(subDir, "app.env")
	if err != nil {
		t.Fatalf("GetFileStatus(config/app.env) failed: %v", err)
	}
	if !statusNested.IsIgnored {
		t.Errorf("expected nested file config/app.env to be ignored, got: %+v", statusNested)
	}

	// 3. Staged file (git add)
	stagedFile := filepath.Join(repoDir, ".env.staging")
	if err := os.WriteFile(stagedFile, []byte("STAGING=1\n"), 0644); err != nil {
		t.Fatal(err)
	}
	runGitCmd(t, repoDir, "add", ".env.staging")

	statusStaged, err := client.GetFileStatus(repoDir, ".env.staging")
	if err != nil {
		t.Fatalf("GetFileStatus(.env.staging) failed: %v", err)
	}
	if !statusStaged.IsStaged {
		t.Errorf("expected .env.staging to be Staged, got: %+v", statusStaged)
	}

	// 4. Tracked file (committed)
	runGitCmd(t, repoDir, "commit", "--no-gpg-sign", "-m", "chore: initial commit")

	statusCommitted, err := client.GetFileStatus(repoDir, ".env.staging")
	if err != nil {
		t.Fatalf("GetFileStatus(.env.staging) after commit failed: %v", err)
	}
	if !statusCommitted.IsTracked {
		t.Errorf("expected .env.staging to be Tracked after commit, got: %+v", statusCommitted)
	}
	if statusCommitted.IsStaged {
		t.Errorf("expected .env.staging to NOT be staged after commit, got: %+v", statusCommitted)
	}

	// 5. Tracked file that is later added to .gitignore
	if err := os.WriteFile(gitignorePath, []byte(".env.local\n.env.staging\n"), 0644); err != nil {
		t.Fatal(err)
	}
	statusTrackedIgnored, err := client.GetFileStatus(repoDir, ".env.staging")
	if err != nil {
		t.Fatalf("GetFileStatus(.env.staging) failed: %v", err)
	}
	if !statusTrackedIgnored.IsTracked {
		t.Errorf("expected .env.staging to still be Tracked, got: %+v", statusTrackedIgnored)
	}
	if !statusTrackedIgnored.IsIgnored {
		t.Errorf("expected .env.staging to be reported as Ignored by .gitignore, got: %+v", statusTrackedIgnored)
	}
}

// mockRunner allows simulating error conditions and absence of Git.
type mockRunner struct {
	lookPathErr error
	runFn       func(dir string, args ...string) ([]byte, []byte, int, error)
}

func (m *mockRunner) LookPath(file string) (string, error) {
	if m.lookPathErr != nil {
		return "", m.lookPathErr
	}
	return "/usr/bin/git", nil
}

func (m *mockRunner) Run(dir string, args ...string) ([]byte, []byte, int, error) {
	if m.runFn != nil {
		return m.runFn(dir, args...)
	}
	return nil, nil, 0, nil
}

func TestGitUnavailable(t *testing.T) {
	mock := &mockRunner{
		lookPathErr: errors.New("git not found"),
	}
	client := NewClientWithRunner(mock)

	if client.IsAvailable() {
		t.Errorf("expected IsAvailable to be false when git is missing")
	}

	if client.IsGitRepo("/any/path") {
		t.Errorf("expected IsGitRepo to be false when git is missing")
	}

	status, err := client.GetFileStatus("/any/path", ".env")
	if err != nil {
		t.Fatalf("expected no error from fallback GetFileStatus, got: %v", err)
	}
	if status.IsRepo {
		t.Errorf("expected status.IsRepo to be false")
	}
}

func TestGetStagedFiles(t *testing.T) {
	repoDir := setupTestGitRepo(t)
	client := NewClient()

	// Initial state: no staged files
	staged, err := client.GetStagedFiles(repoDir)
	if err != nil {
		t.Fatalf("unexpected error on empty stage: %v", err)
	}
	if len(staged) != 0 {
		t.Errorf("expected 0 staged files, got %d", len(staged))
	}

	// Create and stage files
	file1 := filepath.Join(repoDir, ".env")
	file2 := filepath.Join(repoDir, "sub", ".env.local")
	if err := os.MkdirAll(filepath.Dir(file2), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file1, []byte("FOO=1\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte("BAR=2\n"), 0644); err != nil {
		t.Fatal(err)
	}

	runGitCmd(t, repoDir, "add", ".env", "sub/.env.local")

	staged, err = client.GetStagedFiles(repoDir)
	if err != nil {
		t.Fatalf("failed to get staged files: %v", err)
	}
	if len(staged) != 2 {
		t.Fatalf("expected 2 staged files, got %d: %v", len(staged), staged)
	}

	// Verify non-git directory returns error
	tempNonGit := t.TempDir()
	_, err = client.GetStagedFiles(tempNonGit)
	if err == nil {
		t.Errorf("expected error on non-git dir for GetStagedFiles")
	}
}

func TestGetHooksDir(t *testing.T) {
	repoDir := setupTestGitRepo(t)
	client := NewClient()

	hooksDir, err := client.GetHooksDir(repoDir)
	if err != nil {
		t.Fatalf("failed to get hooks dir: %v", err)
	}
	expectedDefault := filepath.Clean(filepath.Join(repoDir, ".git", "hooks"))
	if hooksDir != expectedDefault {
		t.Errorf("expected hooks dir %s, got %s", expectedDefault, hooksDir)
	}

	// Test custom core.hooksPath
	customHooksRel := ".githooks"
	runGitCmd(t, repoDir, "config", "core.hooksPath", customHooksRel)

	hooksDirCustom, err := client.GetHooksDir(repoDir)
	if err != nil {
		t.Fatalf("failed to get custom hooks dir: %v", err)
	}
	expectedCustom := filepath.Clean(filepath.Join(repoDir, customHooksRel))
	if hooksDirCustom != expectedCustom {
		t.Errorf("expected custom hooks dir %s, got %s", expectedCustom, hooksDirCustom)
	}

	// Test non-git directory
	tempNonGit := t.TempDir()
	_, err = client.GetHooksDir(tempNonGit)
	if err == nil {
		t.Errorf("expected error on non-git dir for GetHooksDir")
	}
}

func TestPackageLevelDefaults(t *testing.T) {
	// Test package-level helper methods
	_ = IsAvailable()
	_ = IsGitRepo(".")
	_, _ = GetRepoRoot(".")
	_, _ = IsTracked(".", ".env")
	_, _ = IsStaged(".", ".env")
	_, _ = IsIgnored(".", ".env")
	_, _ = GetFileStatus(".", ".env")
	_, _ = GetStagedFiles(".")
	_, _ = GetHooksDir(".")
}
