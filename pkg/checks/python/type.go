package pythoncheck

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// TypeCheck runs mypy for Python type checking.
type TypeCheck struct {
	Config *config.PythonLanguageConfig
}

func (c *TypeCheck) ID() string   { return "python:type" }
func (c *TypeCheck) Name() string { return "Python Type Check" }

func (c *TypeCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangPython)

	// Check if this is a typed Python project
	if !c.isTypedProject(path) {
		return rb.Pass("Not a typed Python project (no py.typed marker or mypy config)"), nil
	}

	// Check if mypy is installed
	if !checkutil.ToolAvailable("mypy") {
		return rb.ToolNotInstalled("mypy", "pip install mypy"), nil
	}

	// Run mypy
	result := checkutil.RunCommand(path, "mypy", ".", "--ignore-missing-imports")
	output := result.CombinedOutput()

	if !result.Success() {
		// Parse error count
		errorCount := c.countTypeErrors(output)
		if errorCount > 0 {
			return rb.Warn(fmt.Sprintf("%s found. Run: mypy .", checkutil.PluralizeCount(errorCount, "type error", "type errors"))), nil
		}
		return rb.Warn("Type errors found. Run: mypy ."), nil
	}

	return rb.Pass("No type errors found"), nil
}

// isTypedProject checks if the project uses Python type hints.
func (c *TypeCheck) isTypedProject(path string) bool {
	// Check for mypy.ini
	if safepath.Exists(path, "mypy.ini") {
		return true
	}

	// Check for .mypy.ini
	if safepath.Exists(path, ".mypy.ini") {
		return true
	}

	// Check for py.typed marker (PEP 561)
	if safepath.Exists(path, "py.typed") {
		return true
	}

	// Check for setup.cfg with mypy section
	if safepath.Exists(path, "setup.cfg") {
		data, err := safepath.ReadFile(path, "setup.cfg")
		if err == nil && strings.Contains(string(data), "[mypy]") {
			return true
		}
	}

	// Check for pyproject.toml with mypy config
	if safepath.Exists(path, "pyproject.toml") {
		data, err := safepath.ReadFile(path, "pyproject.toml")
		if err == nil {
			content := string(data)
			if strings.Contains(content, "[tool.mypy]") {
				return true
			}
		}
	}

	return false
}

// countTypeErrors parses mypy output to count errors.
func (c *TypeCheck) countTypeErrors(output string) int {
	// mypy outputs "Found X errors in Y files" at the end
	re := regexp.MustCompile(`Found (\d+) errors? in`)
	matches := re.FindStringSubmatch(output)
	if len(matches) >= 2 {
		count, err := strconv.Atoi(matches[1])
		if err == nil {
			return count
		}
	}

	// Fallback: count lines that look like errors
	// Format: path/file.py:line: error: message
	count := 0
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, ": error:") {
			count++
		}
	}
	return count
}
