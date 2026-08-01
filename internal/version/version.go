package version

var (
	// Version is the semantic version of the application
	Version = "1.0.0"
	// GitCommit is the git commit hash (set at build time)
	GitCommit = "unknown"
	// BuildDate is the date of the build (set at build time)
	BuildDate = "unknown"
)

// GetVersion returns the full version string
func GetVersion() string {
	return Version
}

// GetFullVersion returns version with git commit and build date
func GetFullVersion() string {
	return Version + " (commit: " + GitCommit + ", built: " + BuildDate + ")"
}
