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
				ID:        "typescript:project",
				Name:      "TypeScript Project",
				Languages: []checker.Language{checker.LangTypeScript},
				Critical:  true,
				Order:     100,
			},
		},
		{
			Checker: &BuildCheck{Config: tsCfg},
			Meta: checker.CheckMeta{
				ID:        "typescript:build",
				Name:      "TypeScript Build",
				Languages: []checker.Language{checker.LangTypeScript},
				Critical:  true,
				Order:     110,
			},
		},
		{
			Checker: &TestsCheck{Config: tsCfg},
			Meta: checker.CheckMeta{
				ID:        "typescript:tests",
				Name:      "TypeScript Tests",
				Languages: []checker.Language{checker.LangTypeScript},
				Critical:  true,
				Order:     120,
			},
		},
		{
			Checker: &FormatCheck{Config: tsCfg},
			Meta: checker.CheckMeta{
				ID:        "typescript:format",
				Name:      "TypeScript Format",
				Languages: []checker.Language{checker.LangTypeScript},
				Critical:  false,
				Order:     200,
			},
		},
		{
			Checker: &LintCheck{Config: tsCfg},
			Meta: checker.CheckMeta{
				ID:        "typescript:lint",
				Name:      "TypeScript Lint",
				Languages: []checker.Language{checker.LangTypeScript},
				Critical:  false,
				Order:     210,
			},
		},
		{
			Checker: &TypeCheck{Config: tsCfg},
			Meta: checker.CheckMeta{
				ID:        "typescript:type",
				Name:      "TypeScript Type Check",
				Languages: []checker.Language{checker.LangTypeScript},
				Critical:  true, // Type checking is critical for TypeScript
				Order:     215,
			},
		},
		{
			Checker: &CoverageCheck{Config: tsCfg, Threshold: coverageThreshold},
			Meta: checker.CheckMeta{
				ID:        "typescript:coverage",
				Name:      "TypeScript Coverage",
				Languages: []checker.Language{checker.LangTypeScript},
				Critical:  false,
				Order:     220,
			},
		},
		{
			Checker: &DepsCheck{Config: tsCfg},
			Meta: checker.CheckMeta{
				ID:        "typescript:deps",
				Name:      "TypeScript Vulnerabilities",
				Languages: []checker.Language{checker.LangTypeScript},
				Critical:  false,
				Order:     230,
			},
		},
		{
			Checker: &LoggingCheck{},
			Meta: checker.CheckMeta{
				ID:        "typescript:logging",
				Name:      "TypeScript Logging",
				Languages: []checker.Language{checker.LangTypeScript},
				Critical:  false,
				Order:     250,
			},
		},
	}
}
