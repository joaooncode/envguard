package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/joaooncode/envguard/internal/config"
	"github.com/joaooncode/envguard/internal/detector"
	"github.com/joaooncode/envguard/internal/git"
	"github.com/joaooncode/envguard/internal/hook"
	"github.com/joaooncode/envguard/internal/scanner"
)

type hookInstallConfig struct {
	path  string
	force bool
}

type hookUninstallConfig struct {
	path  string
	force bool
}

type hookRunConfig struct {
	path       string
	configPath string
	noColor    bool
}

func runHookCommand(args []string, stdout, stderr io.Writer, gitClient git.Client) int {
	if len(args) == 0 {
		printHookHelp(stdout)
		return ExitCodeSuccess
	}

	subcmd := strings.ToLower(strings.TrimSpace(args[0]))
	subargs := args[1:]

	switch subcmd {
	case "install":
		return runHookInstallCommand(subargs, stdout, stderr, gitClient)
	case "uninstall":
		return runHookUninstallCommand(subargs, stdout, stderr, gitClient)
	case "run":
		return runHookRunCommand(subargs, stdout, stderr, gitClient)
	case "help", "-h", "--help", "-help":
		printHookHelp(stdout)
		return ExitCodeSuccess
	default:
		fmt.Fprintf(stderr, "Error: unknown hook sub-command %q\n\n", args[0])
		printHookHelp(stderr)
		return ExitCodeUsageError
	}
}

func runHookInstallCommand(args []string, stdout, stderr io.Writer, gitClient git.Client) int {
	fs := flag.NewFlagSet("hook install", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var cfg hookInstallConfig
	fs.StringVar(&cfg.path, "path", ".", "Target repository directory path")
	fs.StringVar(&cfg.path, "p", ".", "Target repository directory path (shorthand)")
	fs.BoolVar(&cfg.force, "force", false, "Overwrite existing pre-commit hooks")
	fs.BoolVar(&cfg.force, "f", false, "Overwrite existing pre-commit hooks (shorthand)")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return ExitCodeSuccess
		}
		return ExitCodeUsageError
	}

	mgr := hook.NewManager(gitClient)
	hookPath, err := mgr.Install(cfg.path, cfg.force)
	if err != nil {
		fmt.Fprintf(stderr, "Error: %v\n", err)
		return ExitCodeInternalError
	}

	fmt.Fprintf(stdout, "✓ Successfully installed pre-commit hook at %s\n", hookPath)
	return ExitCodeSuccess
}

func runHookUninstallCommand(args []string, stdout, stderr io.Writer, gitClient git.Client) int {
	fs := flag.NewFlagSet("hook uninstall", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var cfg hookUninstallConfig
	fs.StringVar(&cfg.path, "path", ".", "Target repository directory path")
	fs.StringVar(&cfg.path, "p", ".", "Target repository directory path (shorthand)")
	fs.BoolVar(&cfg.force, "force", false, "Force removal of non-envguard pre-commit hooks")
	fs.BoolVar(&cfg.force, "f", false, "Force removal of hooks (shorthand)")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return ExitCodeSuccess
		}
		return ExitCodeUsageError
	}

	mgr := hook.NewManager(gitClient)
	hookPath, err := mgr.Uninstall(cfg.path, cfg.force)
	if err != nil {
		fmt.Fprintf(stderr, "Error: %v\n", err)
		return ExitCodeInternalError
	}

	fmt.Fprintf(stdout, "✓ Successfully removed pre-commit hook from %s\n", hookPath)
	return ExitCodeSuccess
}

func runHookRunCommand(args []string, stdout, stderr io.Writer, gitClient git.Client) int {
	fs := flag.NewFlagSet("hook run", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var cfg hookRunConfig
	fs.StringVar(&cfg.path, "path", ".", "Target repository directory path")
	fs.StringVar(&cfg.path, "p", ".", "Target repository directory path (shorthand)")
	fs.StringVar(&cfg.configPath, "config", "", "Path to custom configuration file")
	fs.StringVar(&cfg.configPath, "c", "", "Path to custom configuration file (shorthand)")
	fs.BoolVar(&cfg.noColor, "no-color", false, "Disable ANSI color escape codes in terminal output")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return ExitCodeSuccess
		}
		return ExitCodeUsageError
	}

	appConfig, _, err := config.DiscoverAndLoad(cfg.path, cfg.configPath)
	if err != nil {
		fmt.Fprintf(stderr, "Error: %v\n", err)
		return ExitCodeUsageError
	}

	det := detector.NewWithPatterns(appConfig.Detector.CustomPatterns, appConfig.Detector.Allowlist)
	runner := hook.NewRunner(gitClient, det, appConfig)

	findings, err := runner.RunStagedCheck(cfg.path)
	if err != nil {
		fmt.Fprintf(stderr, "Error: %v\n", err)
		return ExitCodeInternalError
	}

	if len(findings) == 0 {
		renderHookRunSuccess(stdout, cfg.noColor)
		return ExitCodeSuccess
	}

	renderHookRunViolations(stderr, findings, cfg.noColor)
	return ExitCodeFindingsFound
}

func renderHookRunSuccess(w io.Writer, noColor bool) {
	green := func(s string) string {
		if noColor || s == "" {
			return s
		}
		return "\033[32m" + s + "\033[0m"
	}
	fmt.Fprintf(w, "%s No unprotected environment files staged for commit.\n", green("✓"))
}

func renderHookRunViolations(w io.Writer, findings []scanner.Finding, noColor bool) {
	bold := func(s string) string {
		if noColor || s == "" {
			return s
		}
		return "\033[1m" + s + "\033[0m"
	}
	red := func(s string) string {
		if noColor || s == "" {
			return s
		}
		return "\033[31m" + s + "\033[0m"
	}
	boldRed := func(s string) string {
		if noColor || s == "" {
			return s
		}
		return "\033[1;31m" + s + "\033[0m"
	}
	cyan := func(s string) string {
		if noColor || s == "" {
			return s
		}
		return "\033[36m" + s + "\033[0m"
	}
	dim := func(s string) string {
		if noColor || s == "" {
			return s
		}
		return "\033[90m" + s + "\033[0m"
	}

	var sb strings.Builder
	sb.WriteString(boldRed("🚨 Git Pre-Commit Check Failed!\n"))
	sb.WriteString(red(fmt.Sprintf("Found %d unprotected environment file(s) staged for commit:\n\n", len(findings))))

	for _, f := range findings {
		sb.WriteString(fmt.Sprintf("  %s %s\n", boldRed("✗ [STAGED]"), bold(f.Path)))
		if f.Message != "" {
			sb.WriteString(fmt.Sprintf("    %s %s\n", dim("Message:    "), f.Message))
		}
		if len(f.Suggestions) > 0 {
			sb.WriteString(fmt.Sprintf("    %s\n", cyan("Suggestions:")))
			for _, sug := range f.Suggestions {
				sb.WriteString(fmt.Sprintf("      • %s\n", sug))
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString(bold("Commit blocked to prevent sensitive credentials from reaching Git.\n"))
	sb.WriteString(dim("To remediate, unstage the files or add them to .gitignore (or run `envguard fix`).\n\n"))

	fmt.Fprint(w, sb.String())
}

func printHookHelp(w io.Writer) {
	help := `🛡️  envguard hook - Manage Git pre-commit hooks and perform staged inspections

Usage:
  envguard hook <subcommand> [flags]

Available Subcommands:
  install      Install executable pre-commit hook into .git/hooks/pre-commit
  run          Inspect currently staged files and block commits with sensitive env files
  uninstall    Remove envguard pre-commit hook from .git/hooks/pre-commit
  help         Show help for hook commands

Flags:
  -p, --path       Target repository directory path (default: ".")
  -f, --force      Force installation or removal over non-envguard hooks
  -c, --config     Path to custom configuration file (for run subcommand)
      --no-color   Disable ANSI color escape codes in output
`
	fmt.Fprint(w, help)
}
