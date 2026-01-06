package pythoncheck

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// FormatCheck verifies that Python code is properly formatted.
type FormatCheck struct {
	Config *config.PythonLanguageConfig
}

func (c *FormatCheck) ID() string   { return "python:format" }
func (c *FormatCheck) Name() string { return "Python Format" }

func (c *FormatCheck) Run(path string) (checker.Result, error) {
	formatter := c.detectFormatter(path)

	var cmd *exec.Cmd
	var cmdDesc string

	switch formatter {
	case "ruff":
		cmd = exec.Command("ruff", "format", "--check", ".")
		cmdDesc = "ruff format"
	case "black":
		cmd = exec.Command("black", "--check", ".")
		cmdDesc = "black"
	default:
		// Try ruff first, fall back to black
		if _, err := exec.LookPath("ruff"); err == nil {
			cmd = exec.Command("ruff", "format", "--check", ".")
			cmdDesc = "ruff format"
		} else if _, err := exec.LookPath("black"); err == nil {
			cmd = exec.Command("black", "--check", ".")
			cmdDesc = "black"
		} else {
			return checker.Result{
				Name:     c.Name(),
				ID:       c.ID(),
				Passed:   true,
				Status:   checker.Pass,
				Message:  "No formatter installed (install ruff or black)",
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
				Message:  formatter + " not installed, skipping format check",
				Language: checker.LangPython,
			}, nil
		}

		output := strings.TrimSpace(stderr.String())
		if output == "" {
			output = strings.TrimSpace(stdout.String())
		}

		// Count unformatted files
		lines := strings.Split(output, "\n")
		fileCount := 0
		for _, line := range lines {
			if strings.Contains(line, "would reformat") || strings.Contains(line, "Would reformat") {
				fileCount++
			}
		}

		if fileCount > 0 {
			return checker.Result{
				Name:     c.Name(),
				ID:       c.ID(),
				Passed:   false,
				Status:   checker.Warn,
				Message:  cmdDesc + ": " + pluralize(fileCount, "file", "files") + " need formatting",
				Language: checker.LangPython,
			}, nil
		}

		return checker.Result{
			Name:     c.Name(),
			ID:       c.ID(),
			Passed:   false,
			Status:   checker.Warn,
			Message:  cmdDesc + " found issues: " + truncateMessage(output, 150),
			Language: checker.LangPython,
		}, nil
	}

	return checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Passed:   true,
		Status:   checker.Pass,
		Message:  "All Python files are properly formatted",
		Language: checker.LangPython,
	}, nil
}

func (c *FormatCheck) detectFormatter(path string) string {
	if c.Config != nil && c.Config.Formatter != "auto" && c.Config.Formatter != "" {
		return c.Config.Formatter
	}

	// Auto-detect based on config files
	if safepath.Exists(path, "ruff.toml") || safepath.Exists(path, ".ruff.toml") {
		return "ruff"
	}
	if safepath.Exists(path, ".black.toml") || safepath.Exists(path, "pyproject.toml") {
		// Check if pyproject.toml has black config
		data, err := safepath.ReadFile(path, "pyproject.toml")
		if err == nil && strings.Contains(string(data), "[tool.black]") {
			return "black"
		}
		if err == nil && strings.Contains(string(data), "[tool.ruff]") {
			return "ruff"
		}
	}

	return "auto"
}

func pluralize(count int, singular, plural string) string {
	if count == 1 {
		return "1 " + singular
	}
	return strings.Replace(plural, "files", string(rune(count+'0'))+" files", 1)
}
