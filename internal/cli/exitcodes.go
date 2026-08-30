package cli

// Exit codes returned by the envguard CLI.
const (
	// ExitCodeSuccess indicates successful execution with no blocking findings (0).
	ExitCodeSuccess = 0

	// ExitCodeFindingsFound indicates one or more blocking findings were detected (1).
	ExitCodeFindingsFound = 1

	// ExitCodeUsageError indicates invalid arguments, unknown flags, or invalid subcommands (2).
	ExitCodeUsageError = 2

	// ExitCodeInternalError indicates an internal error occurred during execution (3).
	ExitCodeInternalError = 3
)
