package rustcheck

import (
	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
)

// Register returns all Rust check registrations.
func Register(cfg *config.Config) []checker.CheckRegistration {
	// Use language-specific coverage threshold if set, otherwise use global
	coverageThreshold := cfg.Coverage.Threshold
	if cfg.Language.Rust.CoverageThreshold > 0 {
		coverageThreshold = cfg.Language.Rust.CoverageThreshold
	}

	return []checker.CheckRegistration{
		{
			Checker: &ProjectCheck{},
			Meta: checker.CheckMeta{
				ID:         "rust:project",
				Name:       "Rust Project",
				Languages:  []checker.Language{checker.LangRust},
				Critical:   true,
				Order:      100,
				Suggestion: "Ensure Cargo.toml exists and is valid",
			},
		},
		{
			Checker: &BuildCheck{},
			Meta: checker.CheckMeta{
				ID:         "rust:build",
				Name:       "Rust Build",
				Languages:  []checker.Language{checker.LangRust},
				Critical:   true,
				Order:      110,
				Suggestion: "Fix build errors before continuing",
			},
		},
		{
			Checker: &TestsCheck{},
			Meta: checker.CheckMeta{
				ID:         "rust:tests",
				Name:       "Rust Tests",
				Languages:  []checker.Language{checker.LangRust},
				Critical:   true,
				Order:      120,
				Suggestion: "Fix failing tests before continuing",
			},
		},
		{
			Checker: &FormatCheck{},
			Meta: checker.CheckMeta{
				ID:         "rust:format",
				Name:       "Rust Format",
				Languages:  []checker.Language{checker.LangRust},
				Critical:   false,
				Order:      200,
				Suggestion: "Run 'cargo fmt' to format code",
			},
		},
		{
			Checker: &LintCheck{},
			Meta: checker.CheckMeta{
				ID:         "rust:lint",
				Name:       "Rust Clippy",
				Languages:  []checker.Language{checker.LangRust},
				Critical:   false,
				Order:      210,
				Suggestion: "Fix clippy warnings",
			},
		},
		{
			Checker: &CoverageCheck{Threshold: coverageThreshold},
			Meta: checker.CheckMeta{
				ID:         "rust:coverage",
				Name:       "Rust Coverage",
				Languages:  []checker.Language{checker.LangRust},
				Critical:   false,
				Order:      220,
				Suggestion: "Add more tests to improve coverage",
			},
		},
		{
			Checker: &DepsCheck{},
			Meta: checker.CheckMeta{
				ID:         "rust:deps",
				Name:       "Rust Vulnerabilities",
				Languages:  []checker.Language{checker.LangRust},
				Critical:   false,
				Order:      230,
				Suggestion: "Update dependencies to fix vulnerabilities",
			},
		},
		{
			Checker: &LoggingCheck{},
			Meta: checker.CheckMeta{
				ID:         "rust:logging",
				Name:       "Rust Logging",
				Languages:  []checker.Language{checker.LangRust},
				Critical:   false,
				Order:      250,
				Suggestion: "Consider using structured logging (e.g., tracing, log)",
			},
		},
	}
}
