package cli

import "fmt"

// Version is the current semantic version of envguard.
// Overridden at build time via -ldflags in official release builds (see .goreleaser.yaml).
var Version = "0.3.6"

// Commit is the git commit SHA envguard was built from.
// Overridden at build time via -ldflags in official release builds (see .goreleaser.yaml).
var Commit = "none"

// Date is the UTC build timestamp of the release.
// Overridden at build time via -ldflags in official release builds (see .goreleaser.yaml).
var Date = "unknown"

// VersionString returns the formatted version string.
func VersionString() string {
	return fmt.Sprintf("envguard v%s", Version)
}
