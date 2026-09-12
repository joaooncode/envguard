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

// SecretMatch represents a single detected secret instance within a Finding's file.
// Its Severity is independent of the parent Finding's severity.
type SecretMatch struct {
	Line     int      `json:"line"`
	Key      string   `json:"key,omitempty"`
	Method   string   `json:"method"`
	Provider string   `json:"provider,omitempty"`
	Severity Severity `json:"severity"`
}

// Finding represents a detected environment file and its Git security posture.
type Finding struct {
	Path          string         `json:"path"`
	Severity      Severity       `json:"severity"`
	Message       string         `json:"message"`
	Suggestions   []string       `json:"suggestions,omitempty"`
	GitStatus     git.FileStatus `json:"git_status"`
	IsAllowed     bool           `json:"is_allowed"`
	SecretMatches []SecretMatch  `json:"secret_matches,omitempty"`
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

// EffectiveSeverity returns the maximum severity between the Finding itself and any
// of its SecretMatches. Per ADR 0005 item 4, a Secret Match's severity is independent
// of its parent Finding's severity, and exit-code aggregation must take the max of
// the two — otherwise a real secret in a properly ignored file (Finding severity
// Info) would be silently dropped from the summary.
func (f Finding) EffectiveSeverity() Severity {
	max := f.Severity
	for _, m := range f.SecretMatches {
		if severityRank(m.Severity) > severityRank(max) {
			max = m.Severity
		}
	}
	return max
}

// CalculateSummary computes metrics for a slice of findings.
func CalculateSummary(findings []Finding) Summary {
	var s Summary
	for _, f := range findings {
		switch f.EffectiveSeverity() {
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
