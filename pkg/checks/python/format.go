package pythoncheck

import (
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
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
	rb := checkutil.NewResultBuilder(c, checker.LangPython)
	formatter := c.detectFormatter(path)

	var result *checkutil.CommandResult
	var cmdDesc string

	switch formatter {
	case "ruff":
		result = checkutil.RunCommand(path, "ruff", "format", "--check", ".")
		cmdDesc = "ruff format"
	case "black":
		result = checkutil.RunCommand(path, "black", "--check", ".")
		cmdDesc = "black"
	default:
		// Try ruff first, fall back to black
		if checkutil.ToolAvailable("ruff") {
			result = checkutil.RunCommand(path, "ruff", "format", "--check", ".")
			cmdDesc = "ruff format"
		} else if checkutil.ToolAvailable("black") {
			result = checkutil.RunCommand(path, "black", "--check", ".")
			cmdDesc = "black"
		} else {
			return rb.Pass("No formatter installed (install ruff or black)"), nil
		}
	}

	output := result.CombinedOutput()
	if !result.Success() {
		if checkutil.ToolNotFoundError(result.Err) {
			return rb.ToolNotInstalled(formatter, ""), nil
		}

		// Count unformatted files
		lines := strings.Split(result.Output(), "\n")
		fileCount := 0
		for _, line := range lines {
			if strings.Contains(line, "would reformat") || strings.Contains(line, "Would reformat") {
				fileCount++
			}
		}

		if fileCount > 0 {
			return rb.WarnWithOutput(cmdDesc+": "+checkutil.PluralizeCount(fileCount, "file", "files")+" need formatting", output), nil
		}

		return rb.WarnWithOutput(cmdDesc+" found issues: "+checkutil.TruncateMessage(result.Output(), 150), output), nil
	}

	return rb.Pass("All Python files are properly formatted"), nil
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
