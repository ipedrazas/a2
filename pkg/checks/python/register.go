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
				ID:          "python:project",
				Name:        "Python Project",
				Description: "Verifies that pyproject.toml or setup.py exists for proper project configuration.",
				Languages:   []checker.Language{checker.LangPython},
				Critical:    true,
				Order:       100,
				Suggestion:  "Ensure pyproject.toml or setup.py exists",
			},
		},
		{
			Checker: &BuildCheck{Config: pythonCfg},
			Meta: checker.CheckMeta{
				ID:          "python:build",
				Name:        "Python Build",
				Description: "Builds the Python package to verify it can be installed correctly.",
				Languages:   []checker.Language{checker.LangPython},
				Critical:    true,
				Order:       110,
				Suggestion:  "Fix build errors before continuing",
			},
		},
		{
			Checker: &TestsCheck{Config: pythonCfg},
			Meta: checker.CheckMeta{
				ID:          "python:tests",
				Name:        "Python Tests",
				Description: "Runs the test suite using pytest or unittest to verify all tests pass.",
				Languages:   []checker.Language{checker.LangPython},
				Critical:    true,
				Order:       120,
				Suggestion:  "Fix failing tests before continuing",
			},
		},
		{
			Checker: &FormatCheck{Config: pythonCfg},
			Meta: checker.CheckMeta{
				ID:          "python:format",
				Name:        "Python Format",
				Description: "Checks if code is formatted according to black or ruff standards.",
				Languages:   []checker.Language{checker.LangPython},
				Critical:    false,
				Order:       200,
				Suggestion:  "Run formatter (black/ruff) to format code",
			},
		},
		{
			Checker: &LintCheck{Config: pythonCfg},
			Meta: checker.CheckMeta{
				ID:          "python:lint",
				Name:        "Python Lint",
				Description: "Runs linting tools (ruff, flake8, pylint) to catch style and programming errors.",
				Languages:   []checker.Language{checker.LangPython},
				Critical:    false,
				Order:       210,
				Suggestion:  "Fix linting issues",
			},
		},
		{
			Checker: &TypeCheck{Config: pythonCfg},
			Meta: checker.CheckMeta{
				ID:          "python:type",
				Name:        "Python Type Check",
				Description: "Runs static type checking using mypy or pyright to catch type errors.",
				Languages:   []checker.Language{checker.LangPython},
				Critical:    false,
				Order:       215,
				Suggestion:  "Fix type errors reported by mypy/pyright",
			},
		},
		{
			Checker: &CoverageCheck{Config: pythonCfg},
			Meta: checker.CheckMeta{
				ID:          "python:coverage",
				Name:        "Python Coverage",
				Description: "Measures test coverage and verifies it meets the configured threshold.",
				Languages:   []checker.Language{checker.LangPython},
				Critical:    false,
				Order:       220,
				Suggestion:  "Add more tests to improve coverage",
			},
		},
		{
			Checker: &DepsCheck{},
			Meta: checker.CheckMeta{
				ID:          "python:deps",
				Name:        "Python Vulnerabilities",
				Description: "Scans dependencies for known vulnerabilities using pip-audit or safety.",
				Languages:   []checker.Language{checker.LangPython},
				Critical:    false,
				Order:       230,
				Suggestion:  "Update dependencies to fix vulnerabilities",
			},
		},
		{
			Checker: &ComplexityCheck{Config: pythonCfg},
			Meta: checker.CheckMeta{
				ID:          "python:complexity",
				Name:        "Python Complexity",
				Description: "Analyzes cyclomatic complexity of functions using radon to identify overly complex code.",
				Languages:   []checker.Language{checker.LangPython},
				Critical:    false,
				Order:       240,
				Suggestion:  "Refactor complex functions to reduce complexity",
			},
		},
		{
			Checker: &LoggingCheck{},
			Meta: checker.CheckMeta{
				ID:          "python:logging",
				Name:        "Python Logging",
				Description: "Checks for structured logging usage instead of print statements.",
				Languages:   []checker.Language{checker.LangPython},
				Critical:    false,
				Order:       250,
				Suggestion:  "Consider using structured logging",
			},
		},
	}
}
