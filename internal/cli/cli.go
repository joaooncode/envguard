package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/joaooncode/envguard/internal/git"
	"github.com/joaooncode/envguard/internal/scanner"
)

// App encapsulates the CLI execution environment and dependencies.
type App struct {
	stdout    io.Writer
	stderr    io.Writer
	scanner   *scanner.Scanner
	gitClient git.Client
}

// Option configures an App instance.
type Option func(*App)

// WithScanner sets a custom Scanner instance (useful for testing).
func WithScanner(s *scanner.Scanner) Option {
	return func(a *App) {
		a.scanner = s
	}
}

// WithGitClient sets a custom Git Client instance (useful for testing).
func WithGitClient(g git.Client) Option {
	return func(a *App) {
		a.gitClient = g
	}
}

// New creates a new App configured with stdout and stderr writers.
func New(stdout, stderr io.Writer, opts ...Option) *App {
	app := &App{
		stdout: stdout,
		stderr: stderr,
	}
	for _, opt := range opts {
		opt(app)
	}
	return app
}

// Run executes the command line application with the provided arguments and returns an exit code.
func Run(args []string, stdout, stderr io.Writer, opts ...Option) int {
	app := New(stdout, stderr, opts...)
	return app.Run(args)
}

// Run parses arguments, dispatches commands, and returns deterministic exit codes.
func (a *App) Run(args []string) int {
	if len(args) == 0 {
		a.printHelp()
		return ExitCodeSuccess
	}

	cmd := strings.ToLower(strings.TrimSpace(args[0]))

	switch cmd {
	case "version", "-v", "--version", "-version":
		fmt.Fprintf(a.stdout, "%s\n", VersionString())
		return ExitCodeSuccess

	case "help", "-h", "--help", "-help":
		a.printHelp()
		return ExitCodeSuccess

	case "scan":
		return runScanCommand(args[1:], a.stdout, a.stderr, a.scanner)

	case "check":
		return runCheckCommand(args[1:], a.stdout, a.stderr, a.scanner)

	case "init":
		return runInitCommand(args[1:], a.stdout, a.stderr)

	case "fix":
		return runFixCommand(args[1:], a.stdout, a.stderr, a.scanner)

	case "hook":
		return runHookCommand(args[1:], a.stdout, a.stderr, a.gitClient)

	default:
		fmt.Fprintf(a.stderr, "Error: unknown command or flag %q\n\n", args[0])
		a.printHelpTo(a.stderr)
		return ExitCodeUsageError
	}
}

func (a *App) printHelp() {
	a.printHelpTo(a.stdout)
}

func (a *App) printHelpTo(w io.Writer) {
	help := `🛡️  envguard - Prevent .env files and environment secrets from reaching Git

Usage:
  envguard <command> [flags]

Available Commands:
  scan       Scan a directory for unprotected environment files
  check      Run verification optimized for CI/CD pipelines
  fix        Automatically add unprotected environment files to .gitignore
  hook       Manage Git pre-commit hooks and perform staged inspections
  init       Initialize configuration file and safe template files
  version    Show current envguard version
  help       Show help for envguard commands

Flags:
  -h, --help       Show help for envguard
  -v, --version    Show envguard version

Scan & Check Flags:
  -p, --path       Target directory path to scan (default: ".")
  -f, --format     Output format: text|terminal|json (default: "text")
  -s, --severity   Minimum severity level: info|warning|high|critical|all (default: "all")
      --no-color   Disable ANSI color escape codes in terminal output

Fix Flags:
  -p, --path       Target directory path to scan and fix (default: ".")
  -d, --dry-run    Preview proposed .gitignore additions without modifying files
  -c, --config     Path to custom configuration file
      --no-color   Disable ANSI color escape codes in terminal output

Hook Flags (see 'envguard hook --help' for details):
  install          Install executable pre-commit hook into .git/hooks/pre-commit
  run              Inspect currently staged files and block commits with sensitive env files
  uninstall        Remove envguard pre-commit hook from .git/hooks/pre-commit

Init Flags:
  -p, --path           Target directory path to initialize (default: ".")
  -f, --force          Overwrite existing configuration or template files
  -t, --template       Generate a safe .env.example template file
      --template-from  Source .env file to sanitize and create template from

Examples:
  envguard scan
  envguard scan --path ./my-project --format json
  envguard scan --severity warning
  envguard check --path . --severity high
  envguard fix
  envguard fix --dry-run
  envguard hook install
  envguard hook run
  envguard hook uninstall
  envguard init
  envguard init --template
  envguard init --path ./my-project --force
  envguard version
`
	fmt.Fprint(w, help)
}
