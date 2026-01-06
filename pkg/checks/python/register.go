package pythoncheck

import (
	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
)

// Register returns all Python check registrations.
func Register(cfg *config.Config) []checker.CheckRegistration {
	pythonCfg := &cfg.Language.Python

	return []checker.CheckRegistration{
		{
			Checker: &ProjectCheck{},
			Meta: checker.CheckMeta{
				ID:        "python:project",
				Name:      "Python Project",
				Languages: []checker.Language{checker.LangPython},
				Critical:  true,
				Order:     100,
			},
		},
		{
			Checker: &BuildCheck{Config: pythonCfg},
			Meta: checker.CheckMeta{
				ID:        "python:build",
				Name:      "Python Build",
				Languages: []checker.Language{checker.LangPython},
				Critical:  true,
				Order:     110,
			},
		},
		{
			Checker: &TestsCheck{Config: pythonCfg},
			Meta: checker.CheckMeta{
				ID:        "python:tests",
				Name:      "Python Tests",
				Languages: []checker.Language{checker.LangPython},
				Critical:  true,
				Order:     120,
			},
		},
		{
			Checker: &FormatCheck{Config: pythonCfg},
			Meta: checker.CheckMeta{
				ID:        "python:format",
				Name:      "Python Format",
				Languages: []checker.Language{checker.LangPython},
				Critical:  false,
				Order:     200,
			},
		},
		{
			Checker: &LintCheck{Config: pythonCfg},
			Meta: checker.CheckMeta{
				ID:        "python:lint",
				Name:      "Python Lint",
				Languages: []checker.Language{checker.LangPython},
				Critical:  false,
				Order:     210,
			},
		},
		{
			Checker: &TypeCheck{Config: pythonCfg},
			Meta: checker.CheckMeta{
				ID:        "python:type",
				Name:      "Python Type Check",
				Languages: []checker.Language{checker.LangPython},
				Critical:  false,
				Order:     215,
			},
		},
		{
			Checker: &CoverageCheck{Config: pythonCfg},
			Meta: checker.CheckMeta{
				ID:        "python:coverage",
				Name:      "Python Coverage",
				Languages: []checker.Language{checker.LangPython},
				Critical:  false,
				Order:     220,
			},
		},
		{
			Checker: &DepsCheck{},
			Meta: checker.CheckMeta{
				ID:        "python:deps",
				Name:      "Python Vulnerabilities",
				Languages: []checker.Language{checker.LangPython},
				Critical:  false,
				Order:     230,
			},
		},
	}
}
