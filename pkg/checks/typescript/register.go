// Package typescriptcheck provides code quality checks for TypeScript projects.
package typescriptcheck

import (
	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
)

// Register returns all TypeScript checks with their metadata.
func Register(cfg *config.Config) []checker.CheckRegistration {
	tsCfg := &cfg.Language.TypeScript

	// Get coverage threshold with fallback to global
	coverageThreshold := cfg.Coverage.Threshold
	if cfg.Language.TypeScript.CoverageThreshold > 0 {
		coverageThreshold = cfg.Language.TypeScript.CoverageThreshold
	}

	return []checker.CheckRegistration{
		{
			Checker: &ProjectCheck{},
			Meta: checker.CheckMeta{
				ID:          "typescript:project",
				Name:        "TypeScript Project",
				Description: "Verifies that tsconfig.json exists for proper TypeScript configuration.",
				Languages:   []checker.Language{checker.LangTypeScript},
				Critical:    true,
				Order:       100,
				Suggestion:  "Ensure tsconfig.json exists",
			},
		},
		{
			Checker: &BuildCheck{Config: tsCfg},
			Meta: checker.CheckMeta{
				ID:          "typescript:build",
				Name:        "TypeScript Build",
				Description: "Compiles the project to verify it builds without errors.",
				Languages:   []checker.Language{checker.LangTypeScript},
				Critical:    true,
				Order:       110,
				Suggestion:  "Fix build errors before continuing",
			},
		},
		{
			Checker: &TestsCheck{Config: tsCfg},
			Meta: checker.CheckMeta{
				ID:          "typescript:tests",
				Name:        "TypeScript Tests",
				Description: "Runs the test suite using the configured test runner to verify all tests pass.",
				Languages:   []checker.Language{checker.LangTypeScript},
				Critical:    true,
				Order:       120,
				Suggestion:  "Fix failing tests before continuing",
			},
		},
		{
			Checker: &FormatCheck{Config: tsCfg},
			Meta: checker.CheckMeta{
				ID:          "typescript:format",
				Name:        "TypeScript Format",
				Description: "Checks if code is formatted according to prettier or biome standards.",
				Languages:   []checker.Language{checker.LangTypeScript},
				Critical:    false,
				Order:       200,
				Suggestion:  "Run formatter (prettier/biome) to format code",
			},
		},
		{
			Checker: &LintCheck{Config: tsCfg},
			Meta: checker.CheckMeta{
				ID:          "typescript:lint",
				Name:        "TypeScript Lint",
				Description: "Runs ESLint or biome to catch style and programming errors.",
				Languages:   []checker.Language{checker.LangTypeScript},
				Critical:    false,
				Order:       210,
				Suggestion:  "Fix linting issues",
			},
		},
		{
			Checker: &TypeCheck{Config: tsCfg},
			Meta: checker.CheckMeta{
				ID:          "typescript:type",
				Name:        "TypeScript Type Check",
				Description: "Runs the TypeScript compiler to verify all types are correct.",
				Languages:   []checker.Language{checker.LangTypeScript},
				Critical:    true, // Type checking is critical for TypeScript
				Order:       215,
				Suggestion:  "Fix TypeScript type errors",
			},
		},
		{
			Checker: &CoverageCheck{Config: tsCfg, Threshold: coverageThreshold},
			Meta: checker.CheckMeta{
				ID:          "typescript:coverage",
				Name:        "TypeScript Coverage",
				Description: "Measures test coverage and verifies it meets the configured threshold.",
				Languages:   []checker.Language{checker.LangTypeScript},
				Critical:    false,
				Order:       220,
				Suggestion:  "Add more tests to improve coverage",
			},
		},
		{
			Checker: &DepsCheck{Config: tsCfg},
			Meta: checker.CheckMeta{
				ID:          "typescript:deps",
				Name:        "TypeScript Vulnerabilities",
				Description: "Scans dependencies for known vulnerabilities using npm audit.",
				Languages:   []checker.Language{checker.LangTypeScript},
				Critical:    false,
				Order:       230,
				Suggestion:  "Update dependencies to fix vulnerabilities",
			},
		},
		{
			Checker: &LoggingCheck{},
			Meta: checker.CheckMeta{
				ID:          "typescript:logging",
				Name:        "TypeScript Logging",
				Description: "Checks for structured logging usage instead of console.log.",
				Languages:   []checker.Language{checker.LangTypeScript},
				Critical:    false,
				Order:       250,
				Suggestion:  "Consider using structured logging",
			},
		},
	}
}
