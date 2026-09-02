package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joaooncode/envguard/internal/config"
	"github.com/joaooncode/envguard/internal/detector"
	"github.com/joaooncode/envguard/internal/git"
)

// IgnoredDirectories contains the directory names skipped during recursive scanning.
var IgnoredDirectories = map[string]bool{
	".git":         true,
	"node_modules": true,
	"vendor":       true,
	"dist":         true,
	"build":        true,
	".idea":        true,
	".vscode":      true,
}

// Scanner coordinates filesystem traversal, environment detection, and Git status inspection.
type Scanner struct {
	gitClient  git.Client
	detector   *detector.Detector
	cfg        *config.Config
	ignoreDirs map[string]bool
}

// New creates a new Scanner instance with the provided git client and detector.
func New(gitClient git.Client, det *detector.Detector) *Scanner {
	return NewWithConfig(gitClient, det, nil)
}

// NewWithConfig creates a Scanner initialized with configuration settings.
func NewWithConfig(gitClient git.Client, det *detector.Detector, cfg *config.Config) *Scanner {
	if gitClient == nil {
		gitClient = git.NewClient()
	}
	if cfg == nil {
		cfg = config.NewDefault()
	}
	if det == nil {
		det = detector.NewWithPatterns(cfg.Detector.CustomPatterns, cfg.Detector.Allowlist)
	}

	ignoreMap := make(map[string]bool)
	for k, v := range IgnoredDirectories {
		ignoreMap[k] = v
	}
	for _, dir := range cfg.Scanner.IgnoreDirs {
		ignoreMap[dir] = true
	}

	return &Scanner{
		gitClient:  gitClient,
		detector:   det,
		cfg:        cfg,
		ignoreDirs: ignoreMap,
	}
}

// NewDefault creates a Scanner configured with default Git and Detector implementations.
func NewDefault() *Scanner {
	return NewWithConfig(nil, nil, nil)
}


// Scan recursively walks the directory and classifies any detected environment files.
// It never opens or reads file contents.
func (s *Scanner) Scan(dir string) (*Result, error) {
	if dir == "" {
		dir = "."
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve directory path: %w", err)
	}

	info, err := os.Stat(absDir)
	if err != nil {
		return nil, fmt.Errorf("failed to stat directory %s: %w", dir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %s", dir)
	}

	result := &Result{
		ScannedDir: dir,
		Findings:   make([]Finding, 0),
	}

	err = filepath.WalkDir(absDir, func(currentPath string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if d.IsDir() {
			if currentPath != absDir && s.ignoreDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}


		// Evaluate path against detector rules without opening or reading the file
		isEnv, isAllowed := s.detector.Detect(d.Name())
		if !isEnv {
			return nil
		}

		relPath, err := filepath.Rel(absDir, currentPath)
		if err != nil {
			relPath = currentPath
		}
		relPath = filepath.ToSlash(relPath)

		status, err := s.gitClient.GetFileStatus(absDir, relPath)
		if err != nil {
			status = git.FileStatus{
				IsRepo:    false,
				IsTracked: false,
				IsStaged:  false,
				IsIgnored: false,
			}
		}

		finding := s.classifyFinding(relPath, status, isAllowed)
		result.Findings = append(result.Findings, finding)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory %s: %w", dir, err)
	}

	result.CalculateSummary()
	return result, nil
}

// classifyFinding evaluates file metadata and git status to assign a severity level, message, and recommendations.
func (s *Scanner) classifyFinding(relPath string, status git.FileStatus, isAllowed bool) Finding {
	var severity Severity
	var message string
	var suggestions []string

	if isAllowed {
		severity = SeverityInfo
		message = "Safe environment template/example file allowed."
		suggestions = []string{}
	} else if status.IsStaged {
		severity = SeverityHigh
		message = "Environment file is staged for commit in Git index."
		suggestions = []string{
			fmt.Sprintf("Unstage file: git restore --staged %s", relPath),
			"Add to .gitignore",
		}
	} else if status.IsTracked {
		severity = SeverityCritical
		message = "Environment file is tracked by Git (committed in repository history)."
		suggestions = []string{
			fmt.Sprintf("Remove file from git tracking: git rm --cached %s", relPath),
			"Add to .gitignore",
			"Rotate any leaked credentials",
		}
	} else if status.IsIgnored {
		severity = SeverityInfo
		message = "Environment file is properly ignored by .gitignore."
		suggestions = []string{}
	} else {
		severity = SeverityWarning
		message = "Environment file exists locally and is not ignored by .gitignore."
		suggestions = []string{
			"Add to .gitignore",
		}
	}

	// Check if any severity override matches this file path or base name
	baseName := filepath.Base(relPath)
	for _, override := range s.cfg.Detector.SeverityOverrides {
		patternLower := strings.ToLower(override.Pattern)
		baseLower := strings.ToLower(baseName)
		relLower := strings.ToLower(relPath)

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
				severity = SeverityInfo
			case "warning", "warn":
				severity = SeverityWarning
			case "high":
				severity = SeverityHigh
			case "critical":
				severity = SeverityCritical
			}
			break
		}
	}

	return Finding{
		Path:        relPath,
		Severity:    severity,
		Message:     message,
		Suggestions: suggestions,
		GitStatus:   status,
		IsAllowed:   isAllowed,
	}
}


// DefaultScanner is the package-level default scanner instance.
var DefaultScanner = NewDefault()

// Scan walks the specified directory using the default scanner.
func Scan(dir string) (*Result, error) {
	return DefaultScanner.Scan(dir)
}
