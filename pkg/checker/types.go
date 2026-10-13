package checker

import (
	"context"
	"time"
)

// Language represents a supported programming language.
type Language string

const (
	LangGo         Language = "go"
	LangPython     Language = "python"
	LangNode       Language = "node"
	LangJava       Language = "java"
	LangRust       Language = "rust"
	LangTypeScript Language = "typescript"
	LangSwift      Language = "swift"
	LangN8n        Language = "n8n"    // n8n workflow exports (domain, not a programming language)
	LangCommon     Language = "common" // Language-agnostic checks
)

// AllLanguages returns all supported programming-language identifiers
// (excluding common and the n8n domain, which has no toolchain).
func AllLanguages() []Language {
	return []Language{LangGo, LangPython, LangNode, LangJava, LangRust, LangTypeScript, LangSwift}
}

// Status represents the severity level of a check result.
type Status int

const (
	Pass    Status = iota // Pass: Check passed, no issues
	Warn                  // Warn: Something is wrong, but not critical
	Fail                  // Fail: Critical failure, stops execution (the Veto)
	Info                  // Info: Informational only, does not affect maturity score
	Errored               // Errored: a2 could not evaluate the check (tool crash, panic, timeout). Excluded from score.
	Skipped               // Skipped: the check was not run (e.g. cancelled by fail-fast). Excluded from score.
)

func (s Status) String() string {
	switch s {
	case Pass:
		return "PASS"
	case Warn:
		return "WARN"
	case Fail:
		return "FAIL"
	case Info:
		return "INFO"
	case Errored:
		return "ERROR"
	case Skipped:
		return "SKIP"
	default:
		return "UNKNOWN"
	}
}

// Result represents the outcome of running a check.
type Result struct {
	Name        string        // Human-readable name of the check
	ID          string        // Unique identifier for the check
	Passed      bool          // Whether the check passed
	Status      Status        // Severity level (Pass, Warn, Fail)
	Message     string        // What happened (short, user-facing summary)
	Reason      string        // Why it happened (evidence or rationale)
	Language    Language      // Which language this check applies to
	Duration    time.Duration // How long the check took to execute
	RawOutput   string        // Full command output for verbose display
	Critical    bool          // Set by the runner from CheckMeta; weights the maturity score
	SourceDir   string        // Directory the check ran in (empty means repo root); set for source_dir-scoped checks
	Command     string        // The actual command a2 executed (e.g. "cd nodeagent && govulncheck ./..."); empty for checks that don't shell out
	ToolMissing bool          // The check could not run because a required tool is not installed; strict mode fails it
}

// Speed indicates the relative cost of a check.
// Fast checks are static/IO-light (format, lint, file existence, regex scans);
// Slow checks spawn builds, tests, coverage/race runs, or network/heavy scans.
type Speed int

const (
	SpeedFast Speed = iota // Static / IO-light. Included in --quick runs. (zero value)
	SpeedSlow              // Spawns builds/tests/network. Excluded from --quick runs.
)

// CheckMeta provides metadata about a check for registration.
type CheckMeta struct {
	ID          string     // Unique identifier (e.g., "go:build", "python:tests")
	Name        string     // Human-readable name
	Description string     // Detailed explanation of what this check does
	Languages   []Language // Which languages this check applies to
	Critical    bool       // If true, failure = veto/abort
	Optional    bool       // If true, Warn results are converted to Info (excluded from score)
	Order       int        // Execution priority (lower = first)
	Suggestion  string     // Recommendation shown when check fails (e.g., "Run 'go fmt' to fix")
	Speed       Speed      // Relative cost; SpeedSlow checks are skipped by --quick (default SpeedFast)
	Command     string     // Representative command template the check runs (e.g. "govulncheck ./..."); shown by `a2 explain`. Empty for checks that don't shell out.
	Detail      string     // Long-form: what it detects, what patterns/conditions trigger it, and why it matters. Shown by `a2 explain`.
	FixPrompt   string     // Copy-paste remediation prompt for an AI, for checks with no fix command. Shown by `a2 explain`.
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
	// ctx is cancelled on timeout or fail-fast; checks must pass it to any
	// command they run (see checkutil.Command) so the process is killed.
	Run(ctx context.Context, path string) (Result, error)
}

// CoverageThresholdSetter is implemented by checks that support
// a per-source-directory coverage threshold override.
type CoverageThresholdSetter interface {
	SetCoverageThreshold(float64)
}
