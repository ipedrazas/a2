package version

var (
	// Version is set at build time via ldflags
	Version = "dev"
	// GitSHA is set at build time via ldflags
	GitSHA = "unknown"
	// BuildDate is set at build time via ldflags
	BuildDate = "unknown"
)

