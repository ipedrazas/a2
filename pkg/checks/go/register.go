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
				ID:        "go:module",
				Name:      "Go Module",
				Languages: []checker.Language{checker.LangGo},
				Critical:  true,
				Order:     100,
			},
		},
		{
			Checker: &BuildCheck{},
			Meta: checker.CheckMeta{
				ID:        "go:build",
				Name:      "Go Build",
				Languages: []checker.Language{checker.LangGo},
				Critical:  true,
				Order:     110,
			},
		},
		{
			Checker: &TestsCheck{},
			Meta: checker.CheckMeta{
				ID:        "go:tests",
				Name:      "Go Tests",
				Languages: []checker.Language{checker.LangGo},
				Critical:  true,
				Order:     120,
			},
		},
		{
			Checker: &RaceCheck{},
			Meta: checker.CheckMeta{
				ID:        "go:race",
				Name:      "Go Race Detection",
				Languages: []checker.Language{checker.LangGo},
				Critical:  false,
				Order:     125,
			},
		},
		{
			Checker: &FormatCheck{},
			Meta: checker.CheckMeta{
				ID:        "go:format",
				Name:      "Go Format",
				Languages: []checker.Language{checker.LangGo},
				Critical:  false,
				Order:     200,
			},
		},
		{
			Checker: &VetCheck{},
			Meta: checker.CheckMeta{
				ID:        "go:vet",
				Name:      "Go Vet",
				Languages: []checker.Language{checker.LangGo},
				Critical:  false,
				Order:     210,
			},
		},
		{
			Checker: &CoverageCheck{Threshold: coverageThreshold},
			Meta: checker.CheckMeta{
				ID:        "go:coverage",
				Name:      "Go Coverage",
				Languages: []checker.Language{checker.LangGo},
				Critical:  false,
				Order:     220,
			},
		},
		{
			Checker: &DepsCheck{},
			Meta: checker.CheckMeta{
				ID:        "go:deps",
				Name:      "Go Vulnerabilities",
				Languages: []checker.Language{checker.LangGo},
				Critical:  false,
				Order:     230,
			},
		},
		{
			Checker: &CyclomaticCheck{Threshold: cyclomaticThreshold},
			Meta: checker.CheckMeta{
				ID:        "go:cyclomatic",
				Name:      "Go Complexity",
				Languages: []checker.Language{checker.LangGo},
				Critical:  false,
				Order:     240,
			},
		},
		{
			Checker: &LoggingCheck{},
			Meta: checker.CheckMeta{
				ID:        "go:logging",
				Name:      "Go Logging",
				Languages: []checker.Language{checker.LangGo},
				Critical:  false,
				Order:     250,
			},
		},
	}
}
