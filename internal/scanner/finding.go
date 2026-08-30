package scanner

import (
	"time"

	"github.com/joaooncode/envguard/internal/git"
)

// Severity represents the severity level of a detected finding.
type Severity string

const (
	// SeverityCritical indicates an environment file tracked in Git history.
	SeverityCritical Severity = "critical"
	// SeverityHigh indicates an environment file staged in Git index.
	SeverityHigh Severity = "high"
	// SeverityWarning indicates an unprotected local environment file not covered by .gitignore.
	SeverityWarning Severity = "warning"
	// SeverityInfo indicates an allowed template/example file or properly ignored file.
	SeverityInfo Severity = "info"
)

// Finding represents a detected environment file and its Git security posture.
type Finding struct {
	Path        string         `json:"path"`
	Severity    Severity       `json:"severity"`
	Message     string         `json:"message"`
	Suggestions []string       `json:"suggestions,omitempty"`
	GitStatus   git.FileStatus `json:"git_status"`
	IsAllowed   bool           `json:"is_allowed"`
}

// Summary aggregates finding counts categorized by severity.
type Summary struct {
	Total    int  `json:"total"`
	Critical int  `json:"critical"`
	High     int  `json:"high"`
	Warning  int  `json:"warning"`
	Info     int  `json:"info"`
	Passed   bool `json:"passed"`
}

// Result contains the complete scan output including findings and summary metrics.
type Result struct {
	Version    string    `json:"version,omitempty"`
	Timestamp  time.Time `json:"timestamp,omitempty"`
	ScannedDir string    `json:"scanned_dir"`
	Findings   []Finding `json:"findings"`
	Summary    Summary   `json:"summary"`
}

// CalculateSummary computes metrics for a slice of findings.
func CalculateSummary(findings []Finding) Summary {
	var s Summary
	for _, f := range findings {
		switch f.Severity {
		case SeverityCritical:
			s.Critical++
		case SeverityHigh:
			s.High++
		case SeverityWarning:
			s.Warning++
		case SeverityInfo:
			s.Info++
		}
	}
	s.Total = len(findings)
	s.Passed = (s.Critical == 0 && s.High == 0 && s.Warning == 0)
	return s
}

// CalculateSummary updates and returns the Result's summary.
func (r *Result) CalculateSummary() Summary {
	r.Summary = CalculateSummary(r.Findings)
	return r.Summary
}
