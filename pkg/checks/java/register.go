package javacheck

import (
	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
)

// Register returns all Java check registrations.
func Register(cfg *config.Config) []checker.CheckRegistration {
	// Use language-specific coverage threshold if set, otherwise use global
	coverageThreshold := cfg.Coverage.Threshold
	if cfg.Language.Java.CoverageThreshold > 0 {
		coverageThreshold = cfg.Language.Java.CoverageThreshold
	}

	return []checker.CheckRegistration{
		{
			Checker: &ProjectCheck{},
			Meta: checker.CheckMeta{
				ID:        "java:project",
				Name:      "Java Project",
				Languages: []checker.Language{checker.LangJava},
				Critical:  true,
				Order:     100,
			},
		},
		{
			Checker: &BuildCheck{Config: &cfg.Language.Java},
			Meta: checker.CheckMeta{
				ID:        "java:build",
				Name:      "Java Build",
				Languages: []checker.Language{checker.LangJava},
				Critical:  true,
				Order:     110,
			},
		},
		{
			Checker: &TestsCheck{Config: &cfg.Language.Java},
			Meta: checker.CheckMeta{
				ID:        "java:tests",
				Name:      "Java Tests",
				Languages: []checker.Language{checker.LangJava},
				Critical:  true,
				Order:     120,
			},
		},
		{
			Checker: &FormatCheck{},
			Meta: checker.CheckMeta{
				ID:        "java:format",
				Name:      "Java Format",
				Languages: []checker.Language{checker.LangJava},
				Critical:  false,
				Order:     200,
			},
		},
		{
			Checker: &LintCheck{},
			Meta: checker.CheckMeta{
				ID:        "java:lint",
				Name:      "Java Lint",
				Languages: []checker.Language{checker.LangJava},
				Critical:  false,
				Order:     210,
			},
		},
		{
			Checker: &CoverageCheck{Threshold: coverageThreshold},
			Meta: checker.CheckMeta{
				ID:        "java:coverage",
				Name:      "Java Coverage",
				Languages: []checker.Language{checker.LangJava},
				Critical:  false,
				Order:     220,
			},
		},
		{
			Checker: &DepsCheck{},
			Meta: checker.CheckMeta{
				ID:        "java:deps",
				Name:      "Java Dependencies",
				Languages: []checker.Language{checker.LangJava},
				Critical:  false,
				Order:     230,
			},
		},
		{
			Checker: &LoggingCheck{},
			Meta: checker.CheckMeta{
				ID:        "java:logging",
				Name:      "Java Logging",
				Languages: []checker.Language{checker.LangJava},
				Critical:  false,
				Order:     250,
			},
		},
	}
}
