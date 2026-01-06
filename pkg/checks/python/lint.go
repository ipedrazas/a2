package pythoncheck

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
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
	linter := c.detectLinter(path)

	var cmd *exec.Cmd
	var cmdDesc string

	switch linter {
	case "ruff":
		cmd = exec.Command("ruff", "check", ".")
		cmdDesc = "ruff"
	case "flake8":
		cmd = exec.Command("flake8", ".")
		cmdDesc = "flake8"
	case "pylint":
		cmd = exec.Command("pylint", ".", "--output-format=text")
		cmdDesc = "pylint"
	default:
		// Try ruff first, fall back to flake8
		if _, err := exec.LookPath("ruff"); err == nil {
			cmd = exec.Command("ruff", "check", ".")
			cmdDesc = "ruff"
		} else if _, err := exec.LookPath("flake8"); err == nil {
			cmd = exec.Command("flake8", ".")
			cmdDesc = "flake8"
		} else {
			return checker.Result{
				Name:     c.Name(),
				ID:       c.ID(),
				Passed:   true,
				Status:   checker.Pass,
				Message:  "No linter installed (install ruff or flake8)",
				Language: checker.LangPython,
			}, nil
		}
	}

	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// Check if tool is not installed
		if strings.Contains(err.Error(), "executable file not found") {
			return checker.Result{
				Name:     c.Name(),
				ID:       c.ID(),
				Passed:   true,
				Status:   checker.Pass,
				Message:  linter + " not installed, skipping lint check",
				Language: checker.LangPython,
			}, nil
		}

		output := strings.TrimSpace(stdout.String())
		if output == "" {
			output = strings.TrimSpace(stderr.String())
		}

		// Count issues
		lines := strings.Split(output, "\n")
		issueCount := 0
		for _, line := range lines {
			if strings.TrimSpace(line) != "" && !strings.HasPrefix(line, "Found") {
				issueCount++
			}
		}

		return checker.Result{
			Name:     c.Name(),
			ID:       c.ID(),
			Passed:   false,
			Status:   checker.Warn,
			Message:  fmt.Sprintf("%s found %d issues", cmdDesc, issueCount),
			Language: checker.LangPython,
		}, nil
	}

	return checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Passed:   true,
		Status:   checker.Pass,
		Message:  "No linting issues found",
		Language: checker.LangPython,
	}, nil
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
