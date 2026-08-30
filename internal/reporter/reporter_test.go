package reporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/joaooncode/envguard/internal/reporter"
	"github.com/joaooncode/envguard/internal/scanner"
)

func sampleResult() *scanner.Result {
	return &scanner.Result{
		Version:    "0.1.0",
		Timestamp:  time.Date(2026, 8, 29, 23, 0, 0, 0, time.UTC),
		ScannedDir: "/path/to/project",
		Findings: []scanner.Finding{
			{
				Path:        ".env",
				Severity:    scanner.SeverityCritical,
				Message:     "Environment file is tracked by Git (committed in repository history).",
				Suggestions: []string{"Remove file from git tracking: git rm --cached .env", "Add to .gitignore", "Rotate any leaked credentials"},
			},
			{
				Path:        ".env.local",
				Severity:    scanner.SeverityHigh,
				Message:     "Environment file is staged for commit in Git index.",
				Suggestions: []string{"Unstage file: git restore --staged .env.local", "Add to .gitignore"},
			},
			{
				Path:        ".env.backup",
				Severity:    scanner.SeverityWarning,
				Message:     "Environment file exists locally and is not ignored by .gitignore.",
				Suggestions: []string{"Add to .gitignore"},
			},
			{
				Path:        ".env.example",
				Severity:    scanner.SeverityInfo,
				Message:     "Safe environment template/example file allowed.",
				Suggestions: []string{},
				IsAllowed:   true,
			},
		},
		Summary: scanner.Summary{
			Total:    4,
			Critical: 1,
			High:     1,
			Warning:  1,
			Info:     1,
			Passed:   false,
		},
	}
}

func sampleCleanResult() *scanner.Result {
	return &scanner.Result{
		Version:    "0.1.0",
		Timestamp:  time.Date(2026, 8, 29, 23, 0, 0, 0, time.UTC),
		ScannedDir: "/path/to/clean-project",
		Findings:   []scanner.Finding{},
		Summary: scanner.Summary{
			Total:    0,
			Critical: 0,
			High:     0,
			Warning:  0,
			Info:     0,
			Passed:   true,
		},
	}
}

func TestNewReporter(t *testing.T) {
	t.Run("Terminal format", func(t *testing.T) {
		r, err := reporter.New(reporter.FormatTerminal)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if r == nil {
			t.Fatal("expected non-nil reporter")
		}
	})

	t.Run("JSON format", func(t *testing.T) {
		r, err := reporter.New(reporter.FormatJSON)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if r == nil {
			t.Fatal("expected non-nil reporter")
		}
	})

	t.Run("Unsupported format", func(t *testing.T) {
		r, err := reporter.New(reporter.Format("yaml"))
		if err == nil {
			t.Fatal("expected error for unsupported format, got nil")
		}
		if r != nil {
			t.Fatalf("expected nil reporter, got %v", r)
		}
	})
}

func TestTerminalReporter_Render(t *testing.T) {
	t.Run("Render nil result returns error", func(t *testing.T) {
		r, err := reporter.New(reporter.FormatTerminal)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var buf bytes.Buffer
		if err := r.Render(nil, &buf); err == nil {
			t.Fatal("expected error when rendering nil result, got nil")
		}
	})

	t.Run("Render clean result with NoColor", func(t *testing.T) {
		r, err := reporter.New(reporter.FormatTerminal, reporter.WithNoColor(true))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var buf bytes.Buffer
		res := sampleCleanResult()
		if err := r.Render(res, &buf); err != nil {
			t.Fatalf("render failed: %v", err)
		}

		out := buf.String()
		if strings.Contains(out, "\033[") {
			t.Errorf("expected no ANSI escape sequences with NoColor, found in: %q", out)
		}
		if !strings.Contains(out, "envguard") {
			t.Errorf("expected header with envguard, got: %s", out)
		}
		if !strings.Contains(out, "No environment file security issues detected") {
			t.Errorf("expected clean status message, got: %s", out)
		}
		if !strings.Contains(out, "PASSED") {
			t.Errorf("expected PASSED status in summary, got: %s", out)
		}
		if !strings.Contains(out, "Total Findings: 0") {
			t.Errorf("expected Total Findings: 0, got: %s", out)
		}
	})

	t.Run("Render findings with severity icons and details", func(t *testing.T) {
		r, err := reporter.New(reporter.FormatTerminal, reporter.WithNoColor(true))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var buf bytes.Buffer
		res := sampleResult()
		if err := r.Render(res, &buf); err != nil {
			t.Fatalf("render failed: %v", err)
		}

		out := buf.String()

		// Verify severity icons
		if !strings.Contains(out, "✗ [CRITICAL]") {
			t.Errorf("expected critical severity icon '✗ [CRITICAL]', got:\n%s", out)
		}
		if !strings.Contains(out, "✗ [HIGH]") {
			t.Errorf("expected high severity icon '✗ [HIGH]', got:\n%s", out)
		}
		if !strings.Contains(out, "⚠ [WARNING]") {
			t.Errorf("expected warning severity icon '⚠ [WARNING]', got:\n%s", out)
		}
		if !strings.Contains(out, "✓ [INFO]") {
			t.Errorf("expected info severity icon '✓ [INFO]', got:\n%s", out)
		}

		// Verify suggestions
		if !strings.Contains(out, "Remove file from git tracking: git rm --cached .env") {
			t.Errorf("expected suggestion in output, got:\n%s", out)
		}
		if !strings.Contains(out, "FAILED") {
			t.Errorf("expected FAILED status in summary, got:\n%s", out)
		}
	})

	t.Run("Render with ANSI colors enabled", func(t *testing.T) {
		r, err := reporter.New(reporter.FormatTerminal, reporter.WithNoColor(false))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var buf bytes.Buffer
		res := sampleResult()
		if err := r.Render(res, &buf); err != nil {
			t.Fatalf("render failed: %v", err)
		}

		out := buf.String()
		if !strings.Contains(out, "\033[") {
			t.Errorf("expected ANSI escape sequences with colors enabled, got:\n%s", out)
		}
	})
}

func TestJSONReporter_Render(t *testing.T) {
	t.Run("Render nil result returns error", func(t *testing.T) {
		r, err := reporter.New(reporter.FormatJSON)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var buf bytes.Buffer
		if err := r.Render(nil, &buf); err == nil {
			t.Fatal("expected error when rendering nil result, got nil")
		}
	})

	t.Run("Render clean result valid schema", func(t *testing.T) {
		r, err := reporter.New(reporter.FormatJSON, reporter.WithPretty(true))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var buf bytes.Buffer
		res := sampleCleanResult()
		if err := r.Render(res, &buf); err != nil {
			t.Fatalf("render failed: %v", err)
		}

		var report reporter.JSONReport
		if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
			t.Fatalf("failed to unmarshal JSON output: %v\nJSON:\n%s", err, buf.String())
		}

		if report.Version != "0.1.0" {
			t.Errorf("expected version 0.1.0, got %s", report.Version)
		}
		if report.ScannedDir != "/path/to/clean-project" {
			t.Errorf("expected scanned_dir /path/to/clean-project, got %s", report.ScannedDir)
		}
		if report.Findings == nil || len(report.Findings) != 0 {
			t.Errorf("expected empty non-nil findings slice, got %v", report.Findings)
		}
		if !report.Summary.Passed {
			t.Errorf("expected summary.passed to be true, got false")
		}
	})

	t.Run("Render findings with full metadata and unmarshal validation", func(t *testing.T) {
		r, err := reporter.New(reporter.FormatJSON, reporter.WithPretty(false))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var buf bytes.Buffer
		res := sampleResult()
		if err := r.Render(res, &buf); err != nil {
			t.Fatalf("render failed: %v", err)
		}

		var report reporter.JSONReport
		if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
			t.Fatalf("failed to unmarshal JSON output: %v\nJSON:\n%s", err, buf.String())
		}

		if len(report.Findings) != 4 {
			t.Fatalf("expected 4 findings, got %d", len(report.Findings))
		}

		if report.Findings[0].Path != ".env" || report.Findings[0].Severity != scanner.SeverityCritical {
			t.Errorf("unexpected first finding: %+v", report.Findings[0])
		}

		if report.Summary.Critical != 1 || report.Summary.High != 1 || report.Summary.Warning != 1 || report.Summary.Info != 1 {
			t.Errorf("unexpected summary counts: %+v", report.Summary)
		}

		// Verify raw map to ensure no extraneous or secret-leaking fields are emitted
		var rawMap map[string]interface{}
		if err := json.Unmarshal(buf.Bytes(), &rawMap); err != nil {
			t.Fatalf("failed to unmarshal raw map: %v", err)
		}

		allowedKeys := map[string]bool{
			"version":     true,
			"timestamp":   true,
			"scanned_dir": true,
			"findings":    true,
			"summary":     true,
		}
		for k := range rawMap {
			if !allowedKeys[k] {
				t.Errorf("unexpected top-level JSON key: %s", k)
			}
		}
	})
}

