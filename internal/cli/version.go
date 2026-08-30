package cli

import "fmt"

// Version is the current semantic version of envguard.
const Version = "0.1.0"

// VersionString returns the formatted version string.
func VersionString() string {
	return fmt.Sprintf("envguard v%s", Version)
}
