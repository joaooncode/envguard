package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/joaooncode/envguard/internal/config"
	"github.com/joaooncode/envguard/internal/fixer"
	"github.com/joaooncode/envguard/internal/scanner"
)

type fixConfig struct {
	path       string
	configPath string
	dryRun     bool
	noColor    bool
}

func runFixCommand(args []string, stdout, stderr io.Writer, scannerInstance *scanner.Scanner) int {
	fs := flag.NewFlagSet("fix", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var cfg fixConfig
	fs.StringVar(&cfg.path, "path", ".", "Target directory path to scan and fix")
	fs.StringVar(&cfg.path, "p", ".", "Target directory path (shorthand)")
	fs.StringVar(&cfg.configPath, "config", "", "Path to custom configuration file")
	fs.StringVar(&cfg.configPath, "c", "", "Path to custom configuration file (shorthand)")
	fs.BoolVar(&cfg.dryRun, "dry-run", false, "Preview proposed .gitignore changes without writing to disk")
	fs.BoolVar(&cfg.dryRun, "d", false, "Preview proposed changes (shorthand)")
	fs.BoolVar(&cfg.noColor, "no-color", false, "Disable ANSI color escape codes in terminal output")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return ExitCodeSuccess
		}
		return ExitCodeUsageError
	}

	// Load configuration
	appConfig, _, err := config.DiscoverAndLoad(cfg.path, cfg.configPath)
	if err != nil {
		fmt.Fprintf(stderr, "Error: %v\n", err)
		return ExitCodeUsageError
	}

	if scannerInstance == nil {
		scannerInstance = scanner.NewWithConfig(nil, nil, appConfig)
	}

	scanResult, err := scannerInstance.Scan(cfg.path)
	if err != nil {
		fmt.Fprintf(stderr, "Error: failed to scan directory %s: %v\n", cfg.path, err)
		return ExitCodeInternalError
	}

	f := fixer.New()
	fixRes, err := f.Apply(fixer.Options{
		TargetDir: cfg.path,
		Findings:  scanResult.Findings,
		DryRun:    cfg.dryRun,
	})
	if err != nil {
		fmt.Fprintf(stderr, "Error: failed to apply remediation: %v\n", err)
		return ExitCodeInternalError
	}

	renderFixOutput(stdout, fixRes, cfg.noColor)

	if len(fixRes.CriticalFindings) > 0 {
		return ExitCodeFindingsFound
	}

	return ExitCodeSuccess
}

func renderFixOutput(w io.Writer, res *fixer.Result, noColor bool) {
	bold := func(s string) string {
		if noColor || s == "" {
			return s
		}
		return "\033[1m" + s + "\033[0m"
	}
	green := func(s string) string {
		if noColor || s == "" {
			return s
		}
		return "\033[32m" + s + "\033[0m"
	}
	yellow := func(s string) string {
		if noColor || s == "" {
			return s
		}
		return "\033[33m" + s + "\033[0m"
	}
	red := func(s string) string {
		if noColor || s == "" {
			return s
		}
		return "\033[31m" + s + "\033[0m"
	}
	dim := func(s string) string {
		if noColor || s == "" {
			return s
		}
		return "\033[90m" + s + "\033[0m"
	}

	var sb strings.Builder

	if res.DryRun {
		sb.WriteString(bold("🔍 Dry run mode: changes will not be written to disk\n\n"))
		if len(res.AddedRules) > 0 {
			sb.WriteString(bold("Proposed .gitignore additions:\n"))
			for _, r := range res.AddedRules {
				sb.WriteString(fmt.Sprintf("  %s %s\n", green("+"), r))
			}
			sb.WriteString("\n")
		} else {
			sb.WriteString(green("No unprotected environment files found to remediate.\n\n"))
		}
	} else {
		if len(res.AddedRules) > 0 {
			sb.WriteString(fmt.Sprintf("%s Successfully updated %s with %d rule(s):\n",
				green("✓"),
				bold(".gitignore"),
				len(res.AddedRules),
			))
			for _, r := range res.AddedRules {
				sb.WriteString(fmt.Sprintf("  %s %s\n", green("+"), r))
			}
			if len(res.SkippedRules) > 0 {
				sb.WriteString(dim(fmt.Sprintf("  (Skipped %d rule(s) already present in .gitignore)\n", len(res.SkippedRules))))
			}
			sb.WriteString("\n")
		} else if len(res.SkippedRules) > 0 {
			sb.WriteString(fmt.Sprintf("%s All detected environment files are already present in %s.\n\n",
				green("✓"),
				bold(".gitignore"),
			))
		} else {
			sb.WriteString(fmt.Sprintf("%s No unprotected environment files found to remediate.\n\n",
				green("✓"),
			))
		}
	}

	if len(res.CriticalFindings) > 0 {
		sb.WriteString(yellow(fmt.Sprintf("⚠️  Warning: %d tracked environment file(s) detected (CRITICAL severity)\n", len(res.CriticalFindings))))
		sb.WriteString(dim("Tracked files in Git cannot be ignored via .gitignore alone:\n"))
		for _, f := range res.CriticalFindings {
			sb.WriteString(fmt.Sprintf("  %s %s\n", red("✗"), f.Path))
		}
		sb.WriteString("\n" + bold("Remediation steps:") + "\n")
		sb.WriteString("  1. Remove file from Git index without deleting local copy:\n")
		for _, f := range res.CriticalFindings {
			sb.WriteString(fmt.Sprintf("     git rm --cached %s\n", f.Path))
		}
		sb.WriteString("  2. Commit the removal and rotate any exposed credentials immediately.\n\n")
	}

	fmt.Fprint(w, sb.String())
}
