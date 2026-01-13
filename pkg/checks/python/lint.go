package pythoncheck

import (
	"fmt"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// LintCheck runs Python linting.
type LintCheck struct {
	Config *config.PythonLanguageConfig
}

func (c *LintCheck) ID() string   { return "python:lint" }
func (c *LintCheck) Name() string { return "Python Lint" }

func (c *LintCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangPython)
	linter := c.detectLinter(path)

	var result *checkutil.CommandResult
	var cmdDesc string

	switch linter {
	case "ruff":
		result = checkutil.RunCommand(path, "ruff", "check", ".")
		cmdDesc = "ruff"
	case "flake8":
		result = checkutil.RunCommand(path, "flake8", ".")
		cmdDesc = "flake8"
	case "pylint":
		result = checkutil.RunCommand(path, "pylint", ".", "--output-format=text")
		cmdDesc = "pylint"
	default:
		// Try ruff first, fall back to flake8
		if checkutil.ToolAvailable("ruff") {
			result = checkutil.RunCommand(path, "ruff", "check", ".")
			cmdDesc = "ruff"
		} else if checkutil.ToolAvailable("flake8") {
			result = checkutil.RunCommand(path, "flake8", ".")
			cmdDesc = "flake8"
		} else {
			return rb.Pass("No linter installed (install ruff or flake8)"), nil
		}
	}

	if !result.Success() {
		if checkutil.ToolNotFoundError(result.Err) {
			return rb.ToolNotInstalled(linter, ""), nil
		}

		// Count issues
		output := result.Output()
		lines := strings.Split(output, "\n")
		issueCount := 0
		for _, line := range lines {
			if strings.TrimSpace(line) != "" && !strings.HasPrefix(line, "Found") {
				issueCount++
			}
		}

		return rb.Warn(fmt.Sprintf("%s found %d issues", cmdDesc, issueCount)), nil
	}

	return rb.Pass("No linting issues found"), nil
}

func (c *LintCheck) detectLinter(path string) string {
	if c.Config != nil && c.Config.Linter != "auto" && c.Config.Linter != "" {
		return c.Config.Linter
	}

	// Auto-detect based on config files
	if safepath.Exists(path, "ruff.toml") || safepath.Exists(path, ".ruff.toml") {
		return "ruff"
	}
	if safepath.Exists(path, ".flake8") || safepath.Exists(path, "setup.cfg") {
		return "flake8"
	}
	if safepath.Exists(path, ".pylintrc") || safepath.Exists(path, "pylintrc") {
		return "pylint"
	}

	// Check pyproject.toml for tool configs
	if safepath.Exists(path, "pyproject.toml") {
		data, err := safepath.ReadFile(path, "pyproject.toml")
		if err == nil {
			content := string(data)
			if strings.Contains(content, "[tool.ruff]") {
				return "ruff"
			}
			if strings.Contains(content, "[tool.flake8]") {
				return "flake8"
			}
			if strings.Contains(content, "[tool.pylint]") {
				return "pylint"
			}
		}
	}

	return "auto"
}
