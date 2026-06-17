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
				Speed:       checker.SpeedSlow,
				Name:        "Python Build",
				Description: "Builds the Python package to verify it can be installed correctly.",
				Languages:   []checker.Language{checker.LangPython},
				Critical:    true,
				Order:       110,
				Suggestion:  "Fix build errors before continuing",
				Command:     "poetry check",
			},
		},
		{
			Checker: &TestsCheck{Config: pythonCfg},
			Meta: checker.CheckMeta{
				ID:          "python:tests",
				Speed:       checker.SpeedSlow,
				Name:        "Python Tests",
				Description: "Runs the test suite using pytest or unittest to verify all tests pass.",
				Languages:   []checker.Language{checker.LangPython},
				Critical:    true,
				Order:       120,
				Suggestion:  "Fix failing tests before continuing",
				Command:     "pytest -v --tb=short",
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
				Command:     "ruff format --check .",
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
				Command:     "ruff check .",
			},
		},
		{
			Checker: &TypeCheck{Config: pythonCfg},
			Meta: checker.CheckMeta{
				ID:          "python:type",
				Speed:       checker.SpeedSlow,
				Name:        "Python Type Check",
				Description: "Runs static type checking using mypy or pyright to catch type errors.",
				Languages:   []checker.Language{checker.LangPython},
				Critical:    false,
				Order:       215,
				Suggestion:  "Fix type errors reported by mypy/pyright",
				Command:     "mypy . --ignore-missing-imports",
			},
		},
		{
			Checker: &CoverageCheck{Config: pythonCfg},
			Meta: checker.CheckMeta{
				ID:          "python:coverage",
				Speed:       checker.SpeedSlow,
				Name:        "Python Coverage",
				Description: "Measures test coverage and verifies it meets the configured threshold.",
				Languages:   []checker.Language{checker.LangPython},
				Critical:    false,
				Order:       220,
				Suggestion:  "Add more tests to improve coverage",
				Command:     "pytest --cov=. --cov-report=term-missing -q",
			},
		},
		{
			Checker: &DepsCheck{},
			Meta: checker.CheckMeta{
				ID:          "python:deps",
				Speed:       checker.SpeedSlow,
				Name:        "Python Vulnerabilities",
				Description: "Scans dependencies for known vulnerabilities using pip-audit or safety.",
				Languages:   []checker.Language{checker.LangPython},
				Critical:    false,
				Order:       230,
				Suggestion:  "Update dependencies to fix vulnerabilities",
				Command:     "pip-audit",
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
				Command:     "radon cc -s .",
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
		{
			Checker: &DepsFreshnessCheck{},
			Meta: checker.CheckMeta{
				ID:          "python:deps_freshness",
				Speed:       checker.SpeedSlow,
				Name:        "Python Dependency Freshness",
				Description: "Checks for outdated Python packages using pip list --outdated.",
				Languages:   []checker.Language{checker.LangPython},
				Critical:    false,
				Order:       235,
				Suggestion:  "Run 'pip list --outdated' to review stale dependencies",
				Command:     "pip list --outdated --format=columns",
			},
		},
		{
			Checker: &DeadcodeCheck{},
			Meta: checker.CheckMeta{
				ID:          "python:deadcode",
				Speed:       checker.SpeedSlow,
				Name:        "Python Dead Code",
				Description: "Detects unused code (functions, variables, imports) using vulture.",
				Languages:   []checker.Language{checker.LangPython},
				Critical:    false,
				Order:       245,
				Suggestion:  "Remove unused code to improve maintainability",
				Command:     "vulture .",
			},
		},
	}
}
