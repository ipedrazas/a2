// Package nodecheck provides code quality checks for Node.js projects.
package nodecheck

import (
	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
)

// Register returns all Node.js checks with their metadata.
func Register(cfg *config.Config) []checker.CheckRegistration {
	nodeCfg := &cfg.Language.Node

	// Get coverage threshold with fallback to global
	coverageThreshold := cfg.Coverage.Threshold
	if cfg.Language.Node.CoverageThreshold > 0 {
		coverageThreshold = cfg.Language.Node.CoverageThreshold
	}

	return []checker.CheckRegistration{
		{
			Checker: &ProjectCheck{},
			Meta: checker.CheckMeta{
				ID:          "node:project",
				Name:        "Node Project",
				Description: "Verifies that package.json exists and is valid for proper project configuration.",
				Languages:   []checker.Language{checker.LangNode},
				Critical:    true,
				Order:       100,
				Suggestion:  "Ensure package.json exists and is valid",
			},
		},
		{
			Checker: &BuildCheck{Config: nodeCfg},
			Meta: checker.CheckMeta{
				ID:          "node:build",
				Name:        "Node Build",
				Description: "Runs the build script to verify the project compiles without errors.",
				Languages:   []checker.Language{checker.LangNode},
				Critical:    true,
				Order:       110,
				Suggestion:  "Fix build errors before continuing",
			},
		},
		{
			Checker: &TestsCheck{Config: nodeCfg},
			Meta: checker.CheckMeta{
				ID:          "node:tests",
				Name:        "Node Tests",
				Description: "Runs the test suite using the configured test runner to verify all tests pass.",
				Languages:   []checker.Language{checker.LangNode},
				Critical:    true,
				Order:       120,
				Suggestion:  "Fix failing tests before continuing",
			},
		},
		{
			Checker: &FormatCheck{Config: nodeCfg},
			Meta: checker.CheckMeta{
				ID:          "node:format",
				Name:        "Node Format",
				Description: "Checks if code is formatted according to prettier or biome standards.",
				Languages:   []checker.Language{checker.LangNode},
				Critical:    false,
				Order:       200,
				Suggestion:  "Run formatter (prettier/biome) to format code",
			},
		},
		{
			Checker: &LintCheck{Config: nodeCfg},
			Meta: checker.CheckMeta{
				ID:          "node:lint",
				Name:        "Node Lint",
				Description: "Runs ESLint or biome to catch style and programming errors.",
				Languages:   []checker.Language{checker.LangNode},
				Critical:    false,
				Order:       210,
				Suggestion:  "Fix linting issues",
			},
		},
		{
			Checker: &TypeCheck{Config: nodeCfg},
			Meta: checker.CheckMeta{
				ID:          "node:type",
				Name:        "TypeScript Type Check",
				Description: "Runs TypeScript compiler to check for type errors in the codebase.",
				Languages:   []checker.Language{checker.LangNode},
				Critical:    false,
				Order:       215,
				Suggestion:  "Fix TypeScript type errors",
			},
		},
		{
			Checker: &CoverageCheck{Config: nodeCfg, Threshold: coverageThreshold},
			Meta: checker.CheckMeta{
				ID:          "node:coverage",
				Name:        "Node Coverage",
				Description: "Measures test coverage and verifies it meets the configured threshold.",
				Languages:   []checker.Language{checker.LangNode},
				Critical:    false,
				Order:       220,
				Suggestion:  "Add more tests to improve coverage",
			},
		},
		{
			Checker: &DepsCheck{Config: nodeCfg},
			Meta: checker.CheckMeta{
				ID:          "node:deps",
				Name:        "Node Vulnerabilities",
				Description: "Scans dependencies for known vulnerabilities using npm audit.",
				Languages:   []checker.Language{checker.LangNode},
				Critical:    false,
				Order:       230,
				Suggestion:  "Update dependencies to fix vulnerabilities",
			},
		},
		{
			Checker: &LoggingCheck{},
			Meta: checker.CheckMeta{
				ID:          "node:logging",
				Name:        "Node Logging",
				Description: "Checks for structured logging usage instead of console.log.",
				Languages:   []checker.Language{checker.LangNode},
				Critical:    false,
				Order:       250,
				Suggestion:  "Consider using structured logging (e.g., pino, winston)",
			},
		},
	}
}
