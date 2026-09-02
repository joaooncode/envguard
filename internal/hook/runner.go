package hook

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/joaooncode/envguard/internal/config"
	"github.com/joaooncode/envguard/internal/detector"
	"github.com/joaooncode/envguard/internal/git"
	"github.com/joaooncode/envguard/internal/scanner"
)

// Runner executes pre-commit inspection specifically targeting staged Git files.
type Runner struct {
	gitClient git.Client
	detector  *detector.Detector
	cfg       *config.Config
}

// NewRunner creates a new pre-commit Runner instance.
func NewRunner(client git.Client, det *detector.Detector, cfg *config.Config) *Runner {
	if client == nil {
		client = git.NewClient()
	}
	if cfg == nil {
		cfg = config.NewDefault()
	}
	if det == nil {
		det = detector.NewWithPatterns(cfg.Detector.CustomPatterns, cfg.Detector.Allowlist)
	}
	return &Runner{
		gitClient: client,
		detector:  det,
		cfg:       cfg,
	}
}

// RunStagedCheck inspects all files staged in the Git index and returns findings for unprotected environment files.
func (r *Runner) RunStagedCheck(repoPath string) ([]scanner.Finding, error) {
	if repoPath == "" {
		repoPath = "."
	}

	if !r.gitClient.IsGitRepo(repoPath) {
		return nil, fmt.Errorf("not a git repository: %s", repoPath)
	}

	stagedFiles, err := r.gitClient.GetStagedFiles(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve staged files: %w", err)
	}

	var findings []scanner.Finding
	for _, stagedPath := range stagedFiles {
		cleanPath := filepath.ToSlash(stagedPath)
		isEnv, isAllowed := r.detector.Detect(cleanPath)
		if !isEnv {
			continue
		}

		// Safe templates matching allowlist are ignored in pre-commit blocking check
		if isAllowed {
			continue
		}

		// File is staged and not allowed
		severity := scanner.SeverityHigh
		message := "Environment file is staged for commit in Git index."
		suggestions := []string{
			fmt.Sprintf("Unstage file: git restore --staged %s", cleanPath),
			"Add to .gitignore",
		}

		// Apply severity overrides if configured
		baseName := filepath.Base(cleanPath)
		for _, override := range r.cfg.Detector.SeverityOverrides {
			patternLower := strings.ToLower(override.Pattern)
			baseLower := strings.ToLower(baseName)
			relLower := strings.ToLower(cleanPath)

			matched := false
			if patternLower == baseLower || patternLower == relLower {
				matched = true
			} else if m, err := filepath.Match(patternLower, baseLower); err == nil && m {
				matched = true
			} else if m, err := filepath.Match(patternLower, relLower); err == nil && m {
				matched = true
			}

			if matched {
				switch strings.ToLower(override.Severity) {
				case "info":
					severity = scanner.SeverityInfo
				case "warning", "warn":
					severity = scanner.SeverityWarning
				case "high":
					severity = scanner.SeverityHigh
				case "critical":
					severity = scanner.SeverityCritical
				}
				break
			}
		}

		findings = append(findings, scanner.Finding{
			Path:        cleanPath,
			Severity:    severity,
			Message:     message,
			Suggestions: suggestions,
			GitStatus: git.FileStatus{
				IsRepo:    true,
				IsTracked: false,
				IsStaged:  true,
				IsIgnored: false,
			},
			IsAllowed: false,
		})
	}

	return findings, nil
}
