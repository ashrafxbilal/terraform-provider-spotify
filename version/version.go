package version

import (
	"fmt"
)

// These variables are set during the build process via ldflags
var (
	// Version is the main version number that is being run at the moment.
	Version = "0.1.0"

	// VersionPrerelease is a pre-release marker for the version. If this is ""
	// (empty string) then it means that it is a final release. Otherwise, this is
	// a pre-release such as "dev" (in development), "beta", "rc1", etc.
	VersionPrerelease = "dev"

	// GitCommit is the git commit that was compiled. This will be filled in by
	// the compiler.
	GitCommit = ""

	// BuildDate is the date when the binary was built
	BuildDate = ""
)

// GetVersion returns the full version string
func GetVersion() string {
	if VersionPrerelease != "" {
		return fmt.Sprintf("%s-%s", Version, VersionPrerelease)
	}
	return Version
}

// GetVersionInfo returns a formatted string with the full version information
func GetVersionInfo() string {
	version := GetVersion()
	result := fmt.Sprintf("Terraform Provider Spotify v%s", version)

	if GitCommit != "" {
		result += fmt.Sprintf(" (%s)", GitCommit)
	}

	if BuildDate != "" {
		result += fmt.Sprintf(", built on %s", BuildDate)
	}

	return result
}