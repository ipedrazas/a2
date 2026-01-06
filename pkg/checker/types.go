package checker

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
	Name    string // Human-readable name of the check
	ID      string // Unique identifier for the check
	Passed  bool   // Whether the check passed
	Status  Status // Severity level (Pass, Warn, Fail)
	Message string // Descriptive message about the result
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
