package checker

// Language represents a supported programming language.
type Language string

const (
	LangGo     Language = "go"
	LangPython Language = "python"
	LangCommon Language = "common" // Language-agnostic checks
)

// AllLanguages returns all supported language identifiers (excluding common).
func AllLanguages() []Language {
	return []Language{LangGo, LangPython}
}

// Status represents the severity level of a check result.
type Status int

const (
	Pass Status = iota // Pass: Check passed, no issues
	Warn               // Warn: Something is wrong, but not critical
	Fail               // Fail: Critical failure, stops execution (the Veto)
)

func (s Status) String() string {
	switch s {
	case Pass:
		return "PASS"
	case Warn:
		return "WARN"
	case Fail:
		return "FAIL"
	default:
		return "UNKNOWN"
	}
}

// Result represents the outcome of running a check.
type Result struct {
	Name     string   // Human-readable name of the check
	ID       string   // Unique identifier for the check
	Passed   bool     // Whether the check passed
	Status   Status   // Severity level (Pass, Warn, Fail)
	Message  string   // Descriptive message about the result
	Language Language // Which language this check applies to
}

// CheckMeta provides metadata about a check for registration.
type CheckMeta struct {
	ID        string     // Unique identifier (e.g., "go:build", "python:tests")
	Name      string     // Human-readable name
	Languages []Language // Which languages this check applies to
	Critical  bool       // If true, failure = veto/abort
	Order     int        // Execution priority (lower = first)
}

// CheckRegistration combines a Checker with its metadata.
type CheckRegistration struct {
	Checker Checker
	Meta    CheckMeta
}

// Checker is the interface that all checks must implement.
type Checker interface {
	// ID returns a unique identifier for this check.
	ID() string

	// Name returns a human-readable name for this check.
	Name() string

	// Run executes the check against the given path.
	// Returns a Result and any error encountered during execution.
	Run(path string) (Result, error)
}
