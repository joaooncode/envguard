package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/joaooncode/envguard/internal/scanner"
)

type colorizer struct {
	noColor bool
}

func (c colorizer) color(code, text string) string {
	if c.noColor || text == "" {
		return text
	}
	return code + text + "\033[0m"
}

func (c colorizer) bold(text string) string {
	return c.color("\033[1m", text)
}

func (c colorizer) red(text string) string {
	return c.color("\033[31m", text)
}

func (c colorizer) boldRed(text string) string {
	return c.color("\033[1;31m", text)
}

func (c colorizer) yellow(text string) string {
	return c.color("\033[33m", text)
}

func (c colorizer) boldYellow(text string) string {
	return c.color("\033[1;33m", text)
}

func (c colorizer) green(text string) string {
	return c.color("\033[32m", text)
}

func (c colorizer) boldGreen(text string) string {
	return c.color("\033[1;32m", text)
}

func (c colorizer) cyan(text string) string {
	return c.color("\033[36m", text)
}

func (c colorizer) dim(text string) string {
	return c.color("\033[90m", text)
}

// TerminalReporter renders scan results formatted with icons and colors for CLI.
type TerminalReporter struct {
	opts Options
	c    colorizer
}

// NewTerminalReporter creates a new TerminalReporter with specified options.
func NewTerminalReporter(opts Options) *TerminalReporter {
	return &TerminalReporter{
		opts: opts,
		c:    colorizer{noColor: opts.NoColor},
	}
}

// Render writes human-readable formatted scan results to the writer.
func (t *TerminalReporter) Render(result *scanner.Result, w io.Writer) error {
	if result == nil {
		return fmt.Errorf("cannot render nil result")
	}

	var sb strings.Builder

	// Header
	version := result.Version
	if version == "" {
		version = t.opts.Version
	}
	if version != "" {
		sb.WriteString(fmt.Sprintf("%s %s\n", t.c.bold("🛡️  envguard"), t.c.dim("v"+version)))
	} else {
		sb.WriteString(fmt.Sprintf("%s\n", t.c.bold("🛡️  envguard - Environment Security Scanner")))
	}

	dir := result.ScannedDir
	if dir != "" {
		sb.WriteString(fmt.Sprintf("Target: %s\n", t.c.dim(dir)))
	}
	sb.WriteString(t.c.dim("──────────────────────────────────────────────────") + "\n")

	// Findings
	if len(result.Findings) == 0 {
		sb.WriteString(fmt.Sprintf("\n%s %s\n\n", t.c.boldGreen("✓"), t.c.green("No environment file security issues detected.")))
	} else {
		sb.WriteString(fmt.Sprintf("\n%s\n\n", t.c.bold("Findings:")))

		for _, f := range result.Findings {
			icon, badge := t.severityBadge(f.Severity)

			sb.WriteString(fmt.Sprintf("  %s %s %s\n", icon, badge, t.c.bold(f.Path)))
			if f.Message != "" {
				sb.WriteString(fmt.Sprintf("    %s %s\n", t.c.dim("Message:    "), f.Message))
			}
			if len(f.Suggestions) > 0 {
				sb.WriteString(fmt.Sprintf("    %s\n", t.c.cyan("Suggestions:")))
				for _, sug := range f.Suggestions {
					sb.WriteString(fmt.Sprintf("      • %s\n", sug))
				}
			}
			sb.WriteString("\n")
		}
	}

	// Summary
	sb.WriteString(t.c.dim("──────────────────────────────────────────────────") + "\n")
	sb.WriteString(fmt.Sprintf("%s\n", t.c.bold("Summary:")))
	sb.WriteString(fmt.Sprintf("  Total Findings: %d (Critical: %d, High: %d, Warning: %d, Info: %d)\n",
		result.Summary.Total,
		result.Summary.Critical,
		result.Summary.High,
		result.Summary.Warning,
		result.Summary.Info,
	))

	if result.Summary.Passed {
		sb.WriteString(fmt.Sprintf("  Status:         %s %s\n", t.c.boldGreen("✓"), t.c.boldGreen("PASSED")))
	} else {
		sb.WriteString(fmt.Sprintf("  Status:         %s %s\n", t.c.boldRed("✗"), t.c.boldRed("FAILED")))
	}

	_, err := io.WriteString(w, sb.String())
	return err
}

func (t *TerminalReporter) severityBadge(sev scanner.Severity) (icon string, badge string) {
	switch sev {
	case scanner.SeverityCritical:
		return t.c.boldRed("✗"), t.c.boldRed("[CRITICAL]")
	case scanner.SeverityHigh:
		return t.c.red("✗"), t.c.red("[HIGH]")
	case scanner.SeverityWarning:
		return t.c.boldYellow("⚠"), t.c.boldYellow("[WARNING]")
	case scanner.SeverityInfo:
		return t.c.boldGreen("✓"), t.c.cyan("[INFO]")
	default:
		sevUpper := strings.ToUpper(string(sev))
		if sevUpper == "" {
			sevUpper = "UNKNOWN"
		}
		return t.c.dim("•"), t.c.dim(fmt.Sprintf("[%s]", sevUpper))
	}
}
