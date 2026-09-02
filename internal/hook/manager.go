package hook

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joaooncode/envguard/internal/git"
)

// HookSignature is the unique header comment used to identify envguard-managed hook scripts.
const HookSignature = "# Installed by envguard"

// HookFileName is the Git hook file name for pre-commit.
const HookFileName = "pre-commit"

// PreCommitHookScript is the POSIX shell script content installed to .git/hooks/pre-commit.
const PreCommitHookScript = `#!/usr/bin/env sh
# Installed by envguard
# Pre-commit hook to prevent uncommitted or leaked environment files

if command -v envguard >/dev/null 2>&1; then
    envguard hook run
else
    echo "Warning: envguard is not installed or not found in PATH. Skipping hook check." >&2
fi
`

// Manager handles the lifecycle of Git hook scripts in local repositories.
type Manager struct {
	gitClient git.Client
}

// NewManager creates a Manager instance with the given Git client.
func NewManager(client git.Client) *Manager {
	if client == nil {
		client = git.NewClient()
	}
	return &Manager{
		gitClient: client,
	}
}

// Install writes the envguard pre-commit hook script to the repository's Git hooks directory.
func (m *Manager) Install(repoPath string, force bool) (string, error) {
	if repoPath == "" {
		repoPath = "."
	}

	if !m.gitClient.IsGitRepo(repoPath) {
		return "", fmt.Errorf("not a git repository: %s", repoPath)
	}

	hooksDir, err := m.gitClient.GetHooksDir(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to locate hooks directory: %w", err)
	}

	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create hooks directory: %w", err)
	}

	hookPath := filepath.Join(hooksDir, HookFileName)

	// Check if hook already exists
	if data, err := os.ReadFile(hookPath); err == nil {
		content := string(data)
		isEnvguardHook := strings.Contains(content, HookSignature)
		if !isEnvguardHook && !force {
			return "", fmt.Errorf("existing pre-commit hook found at %s. Use --force to overwrite", hookPath)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("failed to read existing hook: %w", err)
	}

	// Write executable hook script (0755)
	if err := os.WriteFile(hookPath, []byte(PreCommitHookScript), 0755); err != nil {
		return "", fmt.Errorf("failed to write pre-commit hook: %w", err)
	}

	return hookPath, nil
}

// Uninstall removes the envguard pre-commit hook script from the repository's Git hooks directory.
func (m *Manager) Uninstall(repoPath string, force bool) (string, error) {
	if repoPath == "" {
		repoPath = "."
	}

	if !m.gitClient.IsGitRepo(repoPath) {
		return "", fmt.Errorf("not a git repository: %s", repoPath)
	}

	hooksDir, err := m.gitClient.GetHooksDir(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to locate hooks directory: %w", err)
	}

	hookPath := filepath.Join(hooksDir, HookFileName)

	data, err := os.ReadFile(hookPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("no pre-commit hook found at %s", hookPath)
		}
		return "", fmt.Errorf("failed to inspect pre-commit hook: %w", err)
	}

	isEnvguardHook := strings.Contains(string(data), HookSignature)
	if !isEnvguardHook && !force {
		return "", fmt.Errorf("existing pre-commit hook at %s was not installed by envguard. Use --force to remove", hookPath)
	}

	if err := os.Remove(hookPath); err != nil {
		return "", fmt.Errorf("failed to remove pre-commit hook: %w", err)
	}

	return hookPath, nil
}
