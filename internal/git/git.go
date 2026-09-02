package git

import (
	"fmt"
	"path/filepath"
	"strings"
)

// FileStatus represents the Git status of a file within a repository.
type FileStatus struct {
	// IsRepo indicates whether the target belongs to a valid Git repository.
	IsRepo bool `json:"is_repo"`
	// IsTracked indicates whether the file is tracked in Git (index or history).
	IsTracked bool `json:"is_tracked"`
	// IsStaged indicates whether the file has changes staged for commit.
	IsStaged bool `json:"is_staged"`
	// IsIgnored indicates whether the file is ignored according to .gitignore.
	IsIgnored bool `json:"is_ignored"`
}

// Client defines the operations for inspecting Git repositories and files.
type Client interface {
	// IsAvailable checks if the git executable is installed and available in PATH.
	IsAvailable() bool
	// IsGitRepo checks if the given directory is inside a Git repository work tree.
	IsGitRepo(dir string) bool
	// GetRepoRoot returns the top-level root directory of the Git repository.
	GetRepoRoot(dir string) (string, error)
	// IsTracked checks if a specific file path is tracked in the Git repository.
	IsTracked(dir string, filePath string) (bool, error)
	// IsStaged checks if a specific file path is staged in the Git index.
	IsStaged(dir string, filePath string) (bool, error)
	// IsIgnored checks if a specific file path matches .gitignore rules.
	IsIgnored(dir string, filePath string) (bool, error)
	// GetFileStatus returns the consolidated FileStatus for a given file.
	GetFileStatus(dir string, filePath string) (FileStatus, error)
	// GetStagedFiles returns the list of repository-relative paths currently staged in the index.
	GetStagedFiles(dir string) ([]string, error)
	// GetHooksDir returns the path to the Git hooks directory for the repository.
	GetHooksDir(dir string) (string, error)
}

// GitClient is the standard implementation of Client.
type GitClient struct {
	runner Runner
}

// NewClient creates a new GitClient with the default OS command runner.
func NewClient() *GitClient {
	return &GitClient{
		runner: NewOSCommandRunner(),
	}
}

// NewClientWithRunner creates a GitClient with a custom Runner (useful for testing).
func NewClientWithRunner(runner Runner) *GitClient {
	return &GitClient{
		runner: runner,
	}
}

// IsAvailable reports whether the git binary is found in PATH.
func (c *GitClient) IsAvailable() bool {
	_, err := c.runner.LookPath("git")
	return err == nil
}

// IsGitRepo reports whether dir is inside a Git repository work tree.
func (c *GitClient) IsGitRepo(dir string) bool {
	if !c.IsAvailable() {
		return false
	}
	stdout, _, exitCode, err := c.runner.Run(dir, "rev-parse", "--is-inside-work-tree")
	if err != nil || exitCode != 0 {
		return false
	}
	return strings.TrimSpace(string(stdout)) == "true"
}

// GetRepoRoot returns the root directory of the repository containing dir.
func (c *GitClient) GetRepoRoot(dir string) (string, error) {
	if !c.IsGitRepo(dir) {
		return "", fmt.Errorf("not a git repository: %s", dir)
	}
	stdout, stderr, exitCode, err := c.runner.Run(dir, "rev-parse", "--show-toplevel")
	if err != nil || exitCode != 0 {
		return "", fmt.Errorf("failed to get git repository root: %s (exit code %d)", strings.TrimSpace(string(stderr)), exitCode)
	}
	root := filepath.Clean(strings.TrimSpace(string(stdout)))
	if evalRoot, err := filepath.EvalSymlinks(root); err == nil {
		return evalRoot, nil
	}
	return root, nil
}

// IsTracked reports whether the file at filePath is tracked in Git.
func (c *GitClient) IsTracked(dir string, filePath string) (bool, error) {
	if !c.IsGitRepo(dir) {
		return false, nil
	}

	repoRoot, err := c.GetRepoRoot(dir)
	if err != nil {
		return false, err
	}

	relPath := normalizePath(repoRoot, dir, filePath)
	if relPath == "" {
		return false, nil
	}

	_, stderr, exitCode, _ := c.runner.Run(repoRoot, "ls-files", "--error-unmatch", "--", relPath)
	if exitCode == 0 {
		return true, nil
	}

	// Exit code 1 means file is not matched / untracked
	if exitCode == 1 || strings.Contains(string(stderr), "did not match any file(s) known to git") {
		return false, nil
	}

	return false, nil
}

// IsStaged reports whether the file at filePath has staged changes in the Git index.
func (c *GitClient) IsStaged(dir string, filePath string) (bool, error) {
	if !c.IsGitRepo(dir) {
		return false, nil
	}

	repoRoot, err := c.GetRepoRoot(dir)
	if err != nil {
		return false, err
	}

	relPath := normalizePath(repoRoot, dir, filePath)
	if relPath == "" {
		return false, nil
	}

	stdout, stderr, exitCode, err := c.runner.Run(repoRoot, "status", "--porcelain", "--", relPath)
	if err != nil || exitCode != 0 {
		return false, fmt.Errorf("failed to check staged status: %s (exit code %d)", strings.TrimSpace(string(stderr)), exitCode)
	}

	lines := strings.Split(strings.TrimSpace(string(stdout)), "\n")
	for _, line := range lines {
		if len(line) >= 2 {
			indexStatus := line[0]
			// First column represents the staging area status (A, M, R, C, D)
			if indexStatus != ' ' && indexStatus != '?' && indexStatus != '!' {
				return true, nil
			}
		}
	}

	return false, nil
}

// IsIgnored reports whether the file at filePath is ignored according to .gitignore.
func (c *GitClient) IsIgnored(dir string, filePath string) (bool, error) {
	if !c.IsGitRepo(dir) {
		return false, nil
	}

	repoRoot, err := c.GetRepoRoot(dir)
	if err != nil {
		return false, err
	}

	relPath := normalizePath(repoRoot, dir, filePath)
	if relPath == "" {
		return false, nil
	}

	// Use --no-index to check ignore rules even if the file is tracked in Git index
	_, stderr, exitCode, err := c.runner.Run(repoRoot, "check-ignore", "-q", "--no-index", "--", relPath)
	if exitCode == 0 {
		return true, nil
	}
	if exitCode == 1 {
		return false, nil
	}

	// Any other exit code is a git error
	return false, fmt.Errorf("failed to check git ignore status: %s (exit code %d)", strings.TrimSpace(string(stderr)), exitCode)
}

// GetFileStatus returns the consolidated FileStatus for the specified filePath.
func (c *GitClient) GetFileStatus(dir string, filePath string) (FileStatus, error) {
	if !c.IsGitRepo(dir) {
		return FileStatus{
			IsRepo:    false,
			IsTracked: false,
			IsStaged:  false,
			IsIgnored: false,
		}, nil
	}

	tracked, err := c.IsTracked(dir, filePath)
	if err != nil {
		return FileStatus{}, err
	}

	staged, err := c.IsStaged(dir, filePath)
	if err != nil {
		return FileStatus{}, err
	}

	ignored, err := c.IsIgnored(dir, filePath)
	if err != nil {
		return FileStatus{}, err
	}

	return FileStatus{
		IsRepo:    true,
		IsTracked: tracked,
		IsStaged:  staged,
		IsIgnored: ignored,
	}, nil
}

// GetStagedFiles returns the list of repository-relative paths currently staged in the index.
func (c *GitClient) GetStagedFiles(dir string) ([]string, error) {
	if !c.IsGitRepo(dir) {
		return nil, fmt.Errorf("not a git repository: %s", dir)
	}

	repoRoot, err := c.GetRepoRoot(dir)
	if err != nil {
		return nil, err
	}

	stdout, stderr, exitCode, err := c.runner.Run(repoRoot, "diff", "--name-only", "--cached", "--diff-filter=ACM")
	if err != nil || exitCode != 0 {
		return nil, fmt.Errorf("failed to get staged files: %s (exit code %d)", strings.TrimSpace(string(stderr)), exitCode)
	}

	trimmed := strings.TrimSpace(string(stdout))
	if trimmed == "" {
		return []string{}, nil
	}

	lines := strings.Split(trimmed, "\n")
	var files []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			files = append(files, filepath.ToSlash(line))
		}
	}
	return files, nil
}

// GetHooksDir returns the path to the Git hooks directory for the repository.
func (c *GitClient) GetHooksDir(dir string) (string, error) {
	if !c.IsGitRepo(dir) {
		return "", fmt.Errorf("not a git repository: %s", dir)
	}

	repoRoot, err := c.GetRepoRoot(dir)
	if err != nil {
		return "", err
	}

	stdout, _, exitCode, _ := c.runner.Run(repoRoot, "config", "--get", "core.hooksPath")
	if exitCode == 0 {
		customPath := strings.TrimSpace(string(stdout))
		if customPath != "" {
			if filepath.IsAbs(customPath) {
				return filepath.Clean(customPath), nil
			}
			return filepath.Clean(filepath.Join(repoRoot, customPath)), nil
		}
	}

	return filepath.Clean(filepath.Join(repoRoot, ".git", "hooks")), nil
}

// DefaultClient is the package-level default Git client.
var DefaultClient = NewClient()

// IsAvailable reports whether the git binary is available in PATH using the default client.
func IsAvailable() bool {
	return DefaultClient.IsAvailable()
}

// IsGitRepo reports whether the given directory is inside a Git repository work tree using the default client.
func IsGitRepo(dir string) bool {
	return DefaultClient.IsGitRepo(dir)
}

// GetRepoRoot returns the top-level repository root using the default client.
func GetRepoRoot(dir string) (string, error) {
	return DefaultClient.GetRepoRoot(dir)
}

// IsTracked reports whether the file is tracked in Git using the default client.
func IsTracked(dir string, filePath string) (bool, error) {
	return DefaultClient.IsTracked(dir, filePath)
}

// IsStaged reports whether the file has staged changes in the Git index using the default client.
func IsStaged(dir string, filePath string) (bool, error) {
	return DefaultClient.IsStaged(dir, filePath)
}

// IsIgnored reports whether the file is ignored by .gitignore using the default client.
func IsIgnored(dir string, filePath string) (bool, error) {
	return DefaultClient.IsIgnored(dir, filePath)
}

// GetFileStatus returns the consolidated FileStatus using the default client.
func GetFileStatus(dir string, filePath string) (FileStatus, error) {
	return DefaultClient.GetFileStatus(dir, filePath)
}

// GetStagedFiles returns the list of staged files using the default client.
func GetStagedFiles(dir string) ([]string, error) {
	return DefaultClient.GetStagedFiles(dir)
}

// GetHooksDir returns the Git hooks directory using the default client.
func GetHooksDir(dir string) (string, error) {
	return DefaultClient.GetHooksDir(dir)
}
