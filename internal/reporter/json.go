package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/joaooncode/envguard/internal/scanner"
)

// JSONReport represents the structured report payload serialized to JSON.
type JSONReport struct {
	Version    string            `json:"version"`
	Timestamp  string            `json:"timestamp"`
	ScannedDir string            `json:"scanned_dir"`
	Findings   []scanner.Finding `json:"findings"`
	Summary    scanner.Summary   `json:"summary"`
}

// JSONReporter renders scan results into structured, typed JSON.
type JSONReporter struct {
	opts Options
}

// NewJSONReporter creates a new JSONReporter with specified options.
func NewJSONReporter(opts Options) *JSONReporter {
	return &JSONReporter{
		opts: opts,
	}
}

// Render serializes the scan results into JSON and writes to w.
func (j *JSONReporter) Render(result *scanner.Result, w io.Writer) error {
	if result == nil {
		return fmt.Errorf("cannot render nil result")
	}

	version := result.Version
	if version == "" {
		version = j.opts.Version
	}
	if version == "" {
		version = "dev"
	}

	ts := result.Timestamp
	if ts.IsZero() {
		ts = time.Now().UTC()
	}

	findings := result.Findings
	if findings == nil {
		findings = make([]scanner.Finding, 0)
	}

	// Validate and sanitize findings to ensure strictly structured metadata without leaking contents
	sanitizedFindings := make([]scanner.Finding, len(findings))
	for i, f := range findings {
		sanitizedFindings[i] = scanner.Finding{
			Path:        f.Path,
			Severity:    f.Severity,
			Message:     f.Message,
			Suggestions: f.Suggestions,
			GitStatus:   f.GitStatus,
			IsAllowed:   f.IsAllowed,
		}
	}

	report := JSONReport{
		Version:    version,
		Timestamp:  ts.Format(time.RFC3339),
		ScannedDir: result.ScannedDir,
		Findings:   sanitizedFindings,
		Summary:    result.Summary,
	}

	var data []byte
	var err error
	if j.opts.Pretty {
		data, err = json.MarshalIndent(report, "", "  ")
	} else {
		data, err = json.Marshal(report)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal json report: %w", err)
	}

	if _, err := w.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write json report: %w", err)
	}

	return nil
}
