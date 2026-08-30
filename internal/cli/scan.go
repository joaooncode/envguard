package cli

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/joaooncode/envguard/internal/reporter"
	"github.com/joaooncode/envguard/internal/scanner"
)

// scanConfig holds the parsed options for scan/check commands.
type scanConfig struct {
	path      string
	format    string
	severity  string
	noColor   bool
}

func parseSeverity(s string) (scanner.Severity, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "info":
		return scanner.SeverityInfo, true
	case "warning", "warn":
		return scanner.SeverityWarning, true
	case "high":
		return scanner.SeverityHigh, true
	case "critical":
		return scanner.SeverityCritical, true
	case "all", "":
		return "", true
	default:
		return "", false
	}
}

func severityRank(sev scanner.Severity) int {
	switch sev {
	case scanner.SeverityInfo:
		return 0
	case scanner.SeverityWarning:
		return 1
	case scanner.SeverityHigh:
		return 2
	case scanner.SeverityCritical:
		return 3
	default:
		return 0
	}
}

func runScanCommand(args []string, stdout, stderr io.Writer, scannerInstance *scanner.Scanner) int {
	fs := flag.NewFlagSet("scan", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var cfg scanConfig
	fs.StringVar(&cfg.path, "path", ".", "Target directory path to scan")
	fs.StringVar(&cfg.path, "p", ".", "Target directory path to scan (shorthand)")
	fs.StringVar(&cfg.format, "format", "text", "Output format: text|terminal|json")
	fs.StringVar(&cfg.format, "f", "text", "Output format (shorthand)")
	fs.StringVar(&cfg.severity, "severity", "all", "Minimum severity level to report (info, warning, high, critical)")
	fs.StringVar(&cfg.severity, "s", "all", "Minimum severity level (shorthand)")
	fs.BoolVar(&cfg.noColor, "no-color", false, "Disable ANSI color escape codes in terminal output")

	if err := fs.Parse(args); err != nil {
		return ExitCodeUsageError
	}

	// Validate format
	var repFormat reporter.Format
	switch strings.ToLower(strings.TrimSpace(cfg.format)) {
	case "text", "terminal", "":
		repFormat = reporter.FormatTerminal
	case "json":
		repFormat = reporter.FormatJSON
	default:
		fmt.Fprintf(stderr, "Error: invalid format %q. Supported formats: text, terminal, json\n", cfg.format)
		return ExitCodeUsageError
	}

	// Validate severity
	minSev, ok := parseSeverity(cfg.severity)
	if !ok {
		fmt.Fprintf(stderr, "Error: invalid severity %q. Supported levels: info, warning, high, critical, all\n", cfg.severity)
		return ExitCodeUsageError
	}

	if scannerInstance == nil {
		scannerInstance = scanner.DefaultScanner
	}

	result, err := scannerInstance.Scan(cfg.path)
	if err != nil {
		fmt.Fprintf(stderr, "Error: failed to scan directory %s: %v\n", cfg.path, err)
		return ExitCodeInternalError
	}

	// Set version metadata
	result.Version = Version

	// Filter by minimum severity if specified
	if minSev != "" {
		minRank := severityRank(minSev)
		filtered := make([]scanner.Finding, 0, len(result.Findings))
		for _, f := range result.Findings {
			if severityRank(f.Severity) >= minRank {
				filtered = append(filtered, f)
			}
		}
		result.Findings = filtered
		result.CalculateSummary()
	}

	// Create and render reporter
	rep, err := reporter.New(
		repFormat,
		reporter.WithNoColor(cfg.noColor),
		reporter.WithVersion(Version),
		reporter.WithPretty(true),
	)
	if err != nil {
		fmt.Fprintf(stderr, "Error: failed to initialize reporter: %v\n", err)
		return ExitCodeInternalError
	}

	if err := rep.Render(result, stdout); err != nil {
		fmt.Fprintf(stderr, "Error: failed to render scan results: %v\n", err)
		return ExitCodeInternalError
	}

	if !result.Summary.Passed {
		return ExitCodeFindingsFound
	}

	return ExitCodeSuccess
}
