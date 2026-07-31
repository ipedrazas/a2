package version

import "runtime/debug"

var (
	// Version is set at build time via ldflags
	Version = "dev"
	// GitSHA is set at build time via ldflags
	GitSHA = "unknown"
	// BuildDate is set at build time via ldflags
	BuildDate = "unknown"
)

// init fills in any value that was not stamped via ldflags from the build
// information Go embeds in the binary. This keeps `go install
// github.com/ipedrazas/a2@latest` reporting the module version rather than
// "dev", and keeps plain `go build` from a checkout reporting the commit.
func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	if Version == "dev" {
		if v := info.Main.Version; v != "" && v != "(devel)" {
			Version = v
		}
	}

	var modified, shaFromBuildInfo bool
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			if GitSHA == "unknown" && s.Value != "" {
				GitSHA = shortSHA(s.Value)
				shaFromBuildInfo = true
			}
		case "vcs.time":
			if BuildDate == "unknown" && s.Value != "" {
				BuildDate = s.Value
			}
		case "vcs.modified":
			modified = s.Value == "true"
		}
	}

	// Only annotate a SHA we derived ourselves; one stamped via ldflags is
	// whatever the release build decided it should be.
	if shaFromBuildInfo && modified {
		GitSHA += "-dirty"
	}
}

func shortSHA(sha string) string {
	if len(sha) > 7 {
		return sha[:7]
	}
	return sha
}
