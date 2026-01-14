package swiftcheck

import (
	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
)

// Register returns all Swift check registrations.
func Register(cfg *config.Config) []checker.CheckRegistration {
	// Use language-specific coverage threshold if set, otherwise use global
	coverageThreshold := cfg.Coverage.Threshold
	if cfg.Language.Swift.CoverageThreshold > 0 {
		coverageThreshold = cfg.Language.Swift.CoverageThreshold
	}

	return []checker.CheckRegistration{
		{
			Checker: &ProjectCheck{},
			Meta: checker.CheckMeta{
				ID:          "swift:project",
				Name:        "Swift Project",
				Description: "Verifies that Package.swift exists and is valid for proper project configuration.",
				Languages:   []checker.Language{checker.LangSwift},
				Critical:    true,
				Order:       100,
				Suggestion:  "Ensure Package.swift exists and is valid",
			},
		},
		{
			Checker: &BuildCheck{},
			Meta: checker.CheckMeta{
				ID:          "swift:build",
				Name:        "Swift Build",
				Description: "Compiles the project using 'swift build' to verify it builds without errors.",
				Languages:   []checker.Language{checker.LangSwift},
				Critical:    true,
				Order:       110,
				Suggestion:  "Fix build errors before continuing",
			},
		},
		{
			Checker: &TestsCheck{},
			Meta: checker.CheckMeta{
				ID:          "swift:tests",
				Name:        "Swift Tests",
				Description: "Runs the test suite using 'swift test' to verify all tests pass.",
				Languages:   []checker.Language{checker.LangSwift},
				Critical:    true,
				Order:       120,
				Suggestion:  "Fix failing tests before continuing",
			},
		},
		{
			Checker: &FormatCheck{Config: &cfg.Language.Swift},
			Meta: checker.CheckMeta{
				ID:          "swift:format",
				Name:        "Swift Format",
				Description: "Checks if code is formatted according to SwiftFormat or swift-format standards.",
				Languages:   []checker.Language{checker.LangSwift},
				Critical:    false,
				Order:       200,
				Suggestion:  "Run SwiftFormat or swift-format to format code",
			},
		},
		{
			Checker: &LintCheck{Config: &cfg.Language.Swift},
			Meta: checker.CheckMeta{
				ID:          "swift:lint",
				Name:        "Swift Lint",
				Description: "Runs SwiftLint to catch style and programming issues.",
				Languages:   []checker.Language{checker.LangSwift},
				Critical:    false,
				Order:       210,
				Suggestion:  "Fix SwiftLint warnings",
			},
		},
		{
			Checker: &CoverageCheck{Threshold: coverageThreshold},
			Meta: checker.CheckMeta{
				ID:          "swift:coverage",
				Name:        "Swift Coverage",
				Description: "Measures test coverage and verifies it meets the configured threshold.",
				Languages:   []checker.Language{checker.LangSwift},
				Critical:    false,
				Order:       220,
				Suggestion:  "Add more tests to improve coverage",
			},
		},
		{
			Checker: &DepsCheck{},
			Meta: checker.CheckMeta{
				ID:          "swift:deps",
				Name:        "Swift Vulnerabilities",
				Description: "Reviews Swift Package Manager dependencies for potential issues.",
				Languages:   []checker.Language{checker.LangSwift},
				Critical:    false,
				Order:       230,
				Suggestion:  "Review and update dependencies",
			},
		},
	}
}
