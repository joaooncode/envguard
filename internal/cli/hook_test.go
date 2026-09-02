package cli_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joaooncode/envguard/internal/cli"
)

func setupTestGitRepo(t *testing.T) string {
	t.Helper()
	tempDir := t.TempDir()

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

func TestHookHelpCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	app := cli.New(&stdout, &stderr)

	code := app.Run([]string{"hook", "--help"})
	if code != cli.ExitCodeSuccess {
		t.Errorf("expected exit code %d, got %d", cli.ExitCodeSuccess, code)
	}
	if !strings.Contains(stdout.String(), "envguard hook - Manage Git pre-commit hooks") {
		t.Errorf("expected hook help in output, got: %s", stdout.String())
	}

	stdout.Reset()
	code = app.Run([]string{"hook"})
	if code != cli.ExitCodeSuccess {
		t.Errorf("expected exit code %d, got %d", cli.ExitCodeSuccess, code)
	}
	if !strings.Contains(stdout.String(), "envguard hook - Manage Git pre-commit hooks") {
		t.Errorf("expected hook help in output, got: %s", stdout.String())
	}
}

func TestHookInstallAndUninstallCLI(t *testing.T) {
	repoDir := setupTestGitRepo(t)

	var stdout, stderr bytes.Buffer
	app := cli.New(&stdout, &stderr)

	// 1. Install hook
	code := app.Run([]string{"hook", "install", "--path", repoDir})
	if code != cli.ExitCodeSuccess {
		t.Fatalf("hook install failed with exit code %d, stderr: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Successfully installed pre-commit hook") {
		t.Errorf("expected install confirmation in stdout, got: %s", stdout.String())
	}

	hookFile := filepath.Join(repoDir, ".git", "hooks", "pre-commit")
	if _, err := os.Stat(hookFile); os.IsNotExist(err) {
		t.Fatalf("hook file %s was not created", hookFile)
	}

	// 2. Uninstall hook
	stdout.Reset()
	stderr.Reset()
	code = app.Run([]string{"hook", "uninstall", "--path", repoDir})
	if code != cli.ExitCodeSuccess {
		t.Fatalf("hook uninstall failed with exit code %d, stderr: %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Successfully removed pre-commit hook") {
		t.Errorf("expected uninstall confirmation in stdout, got: %s", stdout.String())
	}
	if _, err := os.Stat(hookFile); !os.IsNotExist(err) {
		t.Errorf("hook file %s still exists after uninstall", hookFile)
	}
}

func TestHookRunCLI(t *testing.T) {
	repoDir := setupTestGitRepo(t)

	var stdout, stderr bytes.Buffer
	app := cli.New(&stdout, &stderr)

	// 1. Clean stage -> exit code 0
	code := app.Run([]string{"hook", "run", "--path", repoDir, "--no-color"})
	if code != cli.ExitCodeSuccess {
		t.Fatalf("expected exit code 0 on clean stage, got %d", code)
	}
	if !strings.Contains(stdout.String(), "No unprotected environment files staged") {
		t.Errorf("expected clean message in stdout, got: %s", stdout.String())
	}

	// 2. Stage safe file -> exit code 0
	exampleFile := filepath.Join(repoDir, ".env.example")
	_ = os.WriteFile(exampleFile, []byte("KEY=\n"), 0644)
	runGitCmd(t, repoDir, "add", ".env.example")

	stdout.Reset()
	stderr.Reset()
	code = app.Run([]string{"hook", "run", "--path", repoDir, "--no-color"})
	if code != cli.ExitCodeSuccess {
		t.Fatalf("expected exit code 0 with allowed staged file, got %d", code)
	}

	// 3. Stage sensitive file -> exit code 1
	secretFile := filepath.Join(repoDir, ".env.production")
	_ = os.WriteFile(secretFile, []byte("SECRET=true\n"), 0644)
	runGitCmd(t, repoDir, "add", ".env.production")

	stdout.Reset()
	stderr.Reset()
	code = app.Run([]string{"hook", "run", "--path", repoDir, "--no-color"})
	if code != cli.ExitCodeFindingsFound {
		t.Fatalf("expected exit code 1 with staged secret, got %d", code)
	}
	if !strings.Contains(stderr.String(), "Git Pre-Commit Check Failed") {
		t.Errorf("expected failure header in stderr, got: %s", stderr.String())
	}
	if !strings.Contains(stderr.String(), ".env.production") {
		t.Errorf("expected .env.production in stderr findings, got: %s", stderr.String())
	}
}

func TestHookCLIUnknownSubcommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	app := cli.New(&stdout, &stderr)

	code := app.Run([]string{"hook", "unknown"})
	if code != cli.ExitCodeUsageError {
		t.Errorf("expected ExitCodeUsageError (2), got %d", code)
	}
}

func TestHookCLINonGitRepo(t *testing.T) {
	tempNonGit := t.TempDir()

	var stdout, stderr bytes.Buffer
	app := cli.New(&stdout, &stderr)

	code := app.Run([]string{"hook", "install", "--path", tempNonGit})
	if code != cli.ExitCodeInternalError {
		t.Errorf("expected ExitCodeInternalError (3) for non-git repo, got %d", code)
	}
}
