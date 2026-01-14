package gocheck

import (
	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
)

// Register returns all Go check registrations.
func Register(cfg *config.Config) []checker.CheckRegistration {
	// Use language-specific coverage threshold if set, otherwise use global
	coverageThreshold := cfg.Coverage.Threshold
	if cfg.Language.Go.CoverageThreshold > 0 {
		coverageThreshold = cfg.Language.Go.CoverageThreshold
	}

	// Use language-specific cyclomatic threshold or default
	cyclomaticThreshold := cfg.Language.Go.CyclomaticThreshold
	if cyclomaticThreshold <= 0 {
		cyclomaticThreshold = 15 // Default threshold
	}

	return []checker.CheckRegistration{
		{
			Checker: &ModuleCheck{},
			Meta: checker.CheckMeta{
				ID:          "go:module",
				Name:        "Go Module",
				Description: "Verifies that go.mod exists and the project uses Go modules for dependency management.",
				Languages:   []checker.Language{checker.LangGo},
				Critical:    true,
				Order:       100,
				Suggestion:  "Ensure go.mod file exists and is valid",
			},
		},
		{
			Checker: &BuildCheck{},
			Meta: checker.CheckMeta{
				ID:          "go:build",
				Name:        "Go Build",
				Description: "Compiles the project using 'go build ./...' to verify all packages build without errors.",
				Languages:   []checker.Language{checker.LangGo},
				Critical:    true,
				Order:       110,
				Suggestion:  "Fix build errors before continuing",
			},
		},
		{
			Checker: &TestsCheck{},
			Meta: checker.CheckMeta{
				ID:          "go:tests",
				Name:        "Go Tests",
				Description: "Runs the test suite using 'go test ./...' to verify all tests pass.",
				Languages:   []checker.Language{checker.LangGo},
				Critical:    true,
				Order:       120,
				Suggestion:  "Fix failing tests before continuing",
			},
		},
		{
			Checker: &RaceCheck{},
			Meta: checker.CheckMeta{
				ID:          "go:race",
				Name:        "Go Race Detection",
				Description: "Runs tests with the -race flag to detect data races between goroutines.",
				Languages:   []checker.Language{checker.LangGo},
				Critical:    false,
				Order:       125,
				Suggestion:  "Fix race conditions detected by -race flag",
			},
		},
		{
			Checker: &FormatCheck{},
			Meta: checker.CheckMeta{
				ID:          "go:format",
				Name:        "Go Format",
				Description: "Checks if code is formatted according to gofmt standards.",
				Languages:   []checker.Language{checker.LangGo},
				Critical:    false,
				Order:       200,
				Suggestion:  "Run 'gofmt -w .' to format code",
			},
		},
		{
			Checker: &VetCheck{},
			Meta: checker.CheckMeta{
				ID:          "go:vet",
				Name:        "Go Vet",
				Description: "Runs 'go vet' static analysis to find common programming mistakes.",
				Languages:   []checker.Language{checker.LangGo},
				Critical:    false,
				Order:       210,
				Suggestion:  "Fix issues reported by 'go vet ./...'",
			},
		},
		{
			Checker: &CoverageCheck{Threshold: coverageThreshold},
			Meta: checker.CheckMeta{
				ID:          "go:coverage",
				Name:        "Go Coverage",
				Description: "Measures test coverage and verifies it meets the configured threshold.",
				Languages:   []checker.Language{checker.LangGo},
				Critical:    false,
				Order:       220,
				Suggestion:  "Add more tests to improve coverage",
			},
		},
		{
			Checker: &DepsCheck{},
			Meta: checker.CheckMeta{
				ID:          "go:deps",
				Name:        "Go Vulnerabilities",
				Description: "Scans dependencies for known vulnerabilities using 'govulncheck'.",
				Languages:   []checker.Language{checker.LangGo},
				Critical:    false,
				Order:       230,
				Suggestion:  "Update dependencies to fix vulnerabilities",
			},
		},
		{
			Checker: &CyclomaticCheck{Threshold: cyclomaticThreshold},
			Meta: checker.CheckMeta{
				ID:          "go:cyclomatic",
				Name:        "Go Complexity",
				Description: "Analyzes cyclomatic complexity of functions to identify overly complex code.",
				Languages:   []checker.Language{checker.LangGo},
				Critical:    false,
				Order:       240,
				Suggestion:  "Refactor complex functions to reduce cyclomatic complexity",
			},
		},
		{
			Checker: &LoggingCheck{},
			Meta: checker.CheckMeta{
				ID:          "go:logging",
				Name:        "Go Logging",
				Description: "Checks for structured logging usage instead of fmt.Print or log.Print.",
				Languages:   []checker.Language{checker.LangGo},
				Critical:    false,
				Order:       250,
				Suggestion:  "Consider using structured logging (e.g., slog, zap)",
			},
		},
	}
}
