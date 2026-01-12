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
				ID:        "swift:project",
				Name:      "Swift Project",
				Languages: []checker.Language{checker.LangSwift},
				Critical:  true,
				Order:     100,
			},
		},
		{
			Checker: &BuildCheck{},
			Meta: checker.CheckMeta{
				ID:        "swift:build",
				Name:      "Swift Build",
				Languages: []checker.Language{checker.LangSwift},
				Critical:  true,
				Order:     110,
			},
		},
		{
			Checker: &TestsCheck{},
			Meta: checker.CheckMeta{
				ID:        "swift:tests",
				Name:      "Swift Tests",
				Languages: []checker.Language{checker.LangSwift},
				Critical:  true,
				Order:     120,
			},
		},
		{
			Checker: &FormatCheck{Config: &cfg.Language.Swift},
			Meta: checker.CheckMeta{
				ID:        "swift:format",
				Name:      "Swift Format",
				Languages: []checker.Language{checker.LangSwift},
				Critical:  false,
				Order:     200,
			},
		},
		{
			Checker: &LintCheck{Config: &cfg.Language.Swift},
			Meta: checker.CheckMeta{
				ID:        "swift:lint",
				Name:      "Swift Lint",
				Languages: []checker.Language{checker.LangSwift},
				Critical:  false,
				Order:     210,
			},
		},
		{
			Checker: &CoverageCheck{Threshold: coverageThreshold},
			Meta: checker.CheckMeta{
				ID:        "swift:coverage",
				Name:      "Swift Coverage",
				Languages: []checker.Language{checker.LangSwift},
				Critical:  false,
				Order:     220,
			},
		},
		{
			Checker: &DepsCheck{},
			Meta: checker.CheckMeta{
				ID:        "swift:deps",
				Name:      "Swift Vulnerabilities",
				Languages: []checker.Language{checker.LangSwift},
				Critical:  false,
				Order:     230,
			},
		},
	}
}
