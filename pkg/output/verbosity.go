package output

// VerbosityLevel controls how much detail to show in output.
type VerbosityLevel int

const (
	// VerbosityNormal shows status and message only (default).
	VerbosityNormal VerbosityLevel = 0
	// VerbosityFailures shows RawOutput for failing/warning checks (-v).
	VerbosityFailures VerbosityLevel = 1
	// VerbosityAll shows RawOutput for all checks (-vv).
	VerbosityAll VerbosityLevel = 2
)
