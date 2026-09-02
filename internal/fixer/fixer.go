package fixer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joaooncode/envguard/internal/scanner"
)

const SectionHeader = "# Added by envguard"

// Fixer orchestrates automated remediation of unprotected environment files.
type Fixer struct{}

// New creates a new Fixer instance.
func New() *Fixer {
	return &Fixer{}
}

// Options contains parameters for a remediation operation.
type Options struct {
	TargetDir string
	Findings  []scanner.Finding
	DryRun    bool
}

// Result summarizes the outcomes of a remediation execution.
type Result struct {
	TargetDir        string            `json:"target_dir"`
	GitignorePath    string            `json:"gitignore_path"`
	AddedRules       []string          `json:"added_rules"`
	SkippedRules     []string          `json:"skipped_rules"`
	CriticalFindings []scanner.Finding `json:"critical_findings"`
	GitignoreUpdated bool              `json:"gitignore_updated"`
	DryRun           bool              `json:"dry_run"`
	NewContent       string            `json:"new_content,omitempty"`
}

// Apply analyzes scan findings and updates the root .gitignore accordingly.
func (f *Fixer) Apply(opts Options) (*Result, error) {
	targetDir := opts.TargetDir
	if targetDir == "" {
		targetDir = "."
	}

	absTargetDir, err := filepath.Abs(targetDir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve target directory %s: %w", targetDir, err)
	}

	gitignorePath := filepath.Join(absTargetDir, ".gitignore")

	result := &Result{
		TargetDir:     absTargetDir,
		GitignorePath: gitignorePath,
		DryRun:        opts.DryRun,
	}

	var warningFindings []scanner.Finding
	for _, finding := range opts.Findings {
		switch finding.Severity {
		case scanner.SeverityCritical:
			result.CriticalFindings = append(result.CriticalFindings, finding)
		case scanner.SeverityWarning:
			warningFindings = append(warningFindings, finding)
		}
	}

	if len(warningFindings) == 0 {
		return result, nil
	}

	var existingContent string
	if data, err := os.ReadFile(gitignorePath); err == nil {
		existingContent = string(data)
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to read .gitignore at %s: %w", gitignorePath, err)
	}

	rulesToAdd, skippedRules := f.computeRules(absTargetDir, warningFindings, existingContent)
	result.SkippedRules = skippedRules
	result.AddedRules = rulesToAdd

	if len(rulesToAdd) == 0 {
		return result, nil
	}

	newContent := f.mergeGitignore(existingContent, rulesToAdd)
	result.NewContent = newContent
	result.GitignoreUpdated = true

	if !opts.DryRun {
		if err := os.WriteFile(gitignorePath, []byte(newContent), 0644); err != nil {
			return nil, fmt.Errorf("failed to write %s: %w", gitignorePath, err)
		}
	}

	return result, nil
}

func (f *Fixer) computeRules(baseDir string, findings []scanner.Finding, currentContent string) ([]string, []string) {
	existingRules := make(map[string]bool)
	scanner := bufio.NewScanner(strings.NewReader(currentContent))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			existingRules[line] = true
		}
	}

	var toAdd []string
	var skipped []string
	seenInBatch := make(map[string]bool)

	for _, finding := range findings {
		rule := f.CalculateRule(baseDir, finding.Path)
		if rule == "" {
			continue
		}

		if existingRules[rule] || seenInBatch[rule] {
			if !seenInBatch[rule] {
				skipped = append(skipped, rule)
				seenInBatch[rule] = true
			}
			continue
		}

		toAdd = append(toAdd, rule)
		seenInBatch[rule] = true
	}

	return toAdd, skipped
}

// CalculateRule returns the appropriate .gitignore pattern for a file path relative to baseDir.
func (f *Fixer) CalculateRule(baseDir, filePath string) string {
	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		absBase = baseDir
	}

	absPath := filePath
	if !filepath.IsAbs(filePath) {
		absPath = filepath.Join(absBase, filePath)
	}
	absPath = filepath.Clean(absPath)

	rel, err := filepath.Rel(absBase, absPath)
	if err != nil {
		rel = filepath.Base(filePath)
	}

	slashRel := filepath.ToSlash(rel)
	if slashRel == "." || slashRel == "" {
		return ""
	}

	if !strings.Contains(slashRel, "/") {
		return slashRel
	}

	return "/" + slashRel
}

func (f *Fixer) mergeGitignore(existingContent string, newRules []string) string {
	if len(newRules) == 0 {
		return existingContent
	}

	rulesText := strings.Join(newRules, "\n") + "\n"

	if existingContent == "" {
		return SectionHeader + "\n" + rulesText
	}

	lines := strings.Split(existingContent, "\n")
	headerIdx := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == SectionHeader {
			headerIdx = i
			break
		}
	}

	if headerIdx != -1 {
		// Insert directly under existing SectionHeader
		insertIdx := headerIdx + 1
		for insertIdx < len(lines) {
			trimmed := strings.TrimSpace(lines[insertIdx])
			if strings.HasPrefix(trimmed, "#") && trimmed != "" {
				// Reached next section
				break
			}
			if trimmed == "" && insertIdx+1 < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[insertIdx+1]), "#") {
				// Reached blank line before next section
				break
			}
			insertIdx++
		}

		var updatedLines []string
		updatedLines = append(updatedLines, lines[:insertIdx]...)
		for _, rule := range newRules {
			updatedLines = append(updatedLines, rule)
		}
		updatedLines = append(updatedLines, lines[insertIdx:]...)
		return strings.Join(updatedLines, "\n")
	}

	// Section header does not exist, append at the end
	var sb strings.Builder
	sb.WriteString(existingContent)
	if !strings.HasSuffix(existingContent, "\n") {
		sb.WriteString("\n")
	}
	if strings.TrimSpace(existingContent) != "" {
		sb.WriteString("\n")
	}
	sb.WriteString(SectionHeader + "\n")
	sb.WriteString(rulesText)

	return sb.String()
}
