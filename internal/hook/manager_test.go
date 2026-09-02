package hook

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joaooncode/envguard/internal/git"
)

func setupTestGitRepo(t *testing.T) string {
	t.Helper()
	tempDir := t.TempDir()
	if evalDir, err := filepath.EvalSymlinks(tempDir); err == nil {
		tempDir = evalDir
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init failed: %v, out: %s", err, string(out))
	}
	return tempDir
}

func TestManagerInstallAndUninstall(t *testing.T) {
	repoDir := setupTestGitRepo(t)
	mgr := NewManager(git.NewClient())

	// 1. Install hook
	hookPath, err := mgr.Install(repoDir, false)
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	expectedPath := filepath.Join(repoDir, ".git", "hooks", HookFileName)
	if filepath.Clean(hookPath) != filepath.Clean(expectedPath) {
		t.Errorf("expected hook path %s, got %s", expectedPath, hookPath)
	}

	content, err := os.ReadFile(hookPath)
	if err != nil {
		t.Fatalf("failed to read hook file: %v", err)
	}
	if !strings.Contains(string(content), HookSignature) {
		t.Errorf("expected signature %q in hook script", HookSignature)
	}

	// 2. Re-installing over envguard hook succeeds without force
	_, err = mgr.Install(repoDir, false)
	if err != nil {
		t.Fatalf("re-installing over envguard hook failed: %v", err)
	}

	// 3. Uninstall hook
	uninstalledPath, err := mgr.Uninstall(repoDir, false)
	if err != nil {
		t.Fatalf("Uninstall failed: %v", err)
	}
	if filepath.Clean(uninstalledPath) != filepath.Clean(expectedPath) {
		t.Errorf("expected uninstalled path %s, got %s", expectedPath, uninstalledPath)
	}
	if _, err := os.Stat(hookPath); !os.IsNotExist(err) {
		t.Errorf("expected hook file to be removed")
	}

	// 4. Uninstalling when non-existent returns error
	_, err = mgr.Uninstall(repoDir, false)
	if err == nil {
		t.Errorf("expected error when uninstalling missing hook")
	}
}

func TestManagerForeignHookConflict(t *testing.T) {
	repoDir := setupTestGitRepo(t)
	mgr := NewManager(git.NewClient())

	hooksDir := filepath.Join(repoDir, ".git", "hooks")
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		t.Fatal(err)
	}
	hookPath := filepath.Join(hooksDir, HookFileName)

	foreignScript := "#!/usr/bin/env sh\necho 'custom linter'\n"
	if err := os.WriteFile(hookPath, []byte(foreignScript), 0755); err != nil {
		t.Fatal(err)
	}

	// Install without force should fail
	_, err := mgr.Install(repoDir, false)
	if err == nil {
		t.Errorf("expected conflict error when installing over foreign hook")
	}

	// Uninstall without force should fail
	_, err = mgr.Uninstall(repoDir, false)
	if err == nil {
		t.Errorf("expected error when uninstalling foreign hook without force")
	}

	// Install with force should succeed
	_, err = mgr.Install(repoDir, true)
	if err != nil {
		t.Fatalf("expected install with force to succeed, got: %v", err)
	}

	// Foreign hook should be replaced by envguard hook
	content, _ := os.ReadFile(hookPath)
	if !strings.Contains(string(content), HookSignature) {
		t.Errorf("expected hook to be overwritten with envguard script")
	}

	// Put foreign hook back to test uninstall with force
	if err := os.WriteFile(hookPath, []byte(foreignScript), 0755); err != nil {
		t.Fatal(err)
	}

	_, err = mgr.Uninstall(repoDir, true)
	if err != nil {
		t.Fatalf("expected uninstall with force to succeed on foreign hook: %v", err)
	}
}

func TestManagerNonGitRepo(t *testing.T) {
	tempNonGit := t.TempDir()
	mgr := NewManager(nil)

	_, err := mgr.Install(tempNonGit, false)
	if err == nil {
		t.Errorf("expected error installing in non-git directory")
	}

	_, err = mgr.Uninstall(tempNonGit, false)
	if err == nil {
		t.Errorf("expected error uninstalling in non-git directory")
	}
}
