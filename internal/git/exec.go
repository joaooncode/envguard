package git

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
)

// Runner abstracts running Git commands and checking for Git executable availability.
type Runner interface {
	Run(dir string, args ...string) (stdout []byte, stderr []byte, exitCode int, err error)
	LookPath(file string) (string, error)
}

// OSCommandRunner is the standard Runner that executes commands on the OS.
type OSCommandRunner struct {
	GitBinary string
}

// NewOSCommandRunner creates a new OSCommandRunner defaulting to "git".
func NewOSCommandRunner() *OSCommandRunner {
	return &OSCommandRunner{
		GitBinary: "git",
	}
}

// LookPath checks whether the git binary exists in system PATH.
func (r *OSCommandRunner) LookPath(file string) (string, error) {
	if r.GitBinary != "" && file == "git" {
		return exec.LookPath(r.GitBinary)
	}
	return exec.LookPath(file)
}

// Run executes a git command in the specified directory and returns stdout, stderr, exit code, and error.
func (r *OSCommandRunner) Run(dir string, args ...string) ([]byte, []byte, int, error) {
	bin := r.GitBinary
	if bin == "" {
		bin = "git"
	}

	cmd := exec.Command(bin, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Env = append(cmd.Environ(), "GIT_TERMINAL_PROMPT=0")

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}

	return stdoutBuf.Bytes(), stderrBuf.Bytes(), exitCode, err
}

// normalizePath converts a path (relative to dir or absolute) to a slash-delimited relative path to repoRoot.
func normalizePath(repoRoot, dir, targetPath string) string {
	if evalRoot, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = evalRoot
	}

	var fullPath string
	if filepath.IsAbs(targetPath) {
		fullPath = targetPath
	} else if dir != "" {
		fullPath = filepath.Join(dir, targetPath)
	} else {
		fullPath = targetPath
	}

	if evalFull, err := filepath.EvalSymlinks(fullPath); err == nil {
		fullPath = evalFull
	} else {
		parent := filepath.Dir(fullPath)
		if evalParent, err := filepath.EvalSymlinks(parent); err == nil {
			fullPath = filepath.Join(evalParent, filepath.Base(fullPath))
		}
	}

	if repoRoot != "" {
		repoRootClean := filepath.Clean(repoRoot)
		fullPathClean := filepath.Clean(fullPath)

		volRepo := filepath.VolumeName(repoRootClean)
		volFull := filepath.VolumeName(fullPathClean)
		if volRepo != "" && strings.EqualFold(volRepo, volFull) && volRepo != volFull {
			fullPathClean = volRepo + fullPathClean[len(volFull):]
		}

		if rel, err := filepath.Rel(repoRootClean, fullPathClean); err == nil && !strings.HasPrefix(rel, "..") {
			cleanPath := filepath.ToSlash(rel)
			cleanPath = strings.TrimPrefix(cleanPath, "./")
			return cleanPath
		}
	}

	if !filepath.IsAbs(targetPath) {
		clean := filepath.ToSlash(filepath.Clean(targetPath))
		clean = strings.TrimPrefix(clean, "./")
		return clean
	}

	cleanPath := filepath.ToSlash(fullPath)
	cleanPath = strings.TrimPrefix(cleanPath, "./")
	return cleanPath
}
