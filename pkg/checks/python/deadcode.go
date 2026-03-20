package pythoncheck

import (
	"fmt"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// DeadcodeCheck detects unused code using vulture.
type DeadcodeCheck struct{}

func (c *DeadcodeCheck) ID() string   { return "python:deadcode" }
func (c *DeadcodeCheck) Name() string { return "Python Dead Code" }

func (c *DeadcodeCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangPython)

	// Verify Python project
	hasPython := safepath.Exists(path, "pyproject.toml") ||
		safepath.Exists(path, "setup.py") ||
		safepath.Exists(path, "requirements.txt")
	if !hasPython {
		return rb.Fail("Python project not found"), nil
	}

	if !pythonToolAvailable(path, "vulture") {
		return rb.ToolNotInstalled("vulture", "pip install vulture"), nil
	}

	result := runPythonCommand(path, "vulture", ".")
	output := result.CombinedOutput()

	if result.Success() {
		return rb.PassWithOutput("No unused code detected", output), nil
	}

	// vulture exits with code 1 when it finds dead code, 3 for syntax errors
	if result.ExitCode == 3 {
		return rb.WarnWithOutput("vulture: syntax error in source files", output), nil
	}

	// Count findings — each line is one unused item
	lines := countFindings(result.Stdout)
	if lines == 0 {
		lines = countFindings(output)
	}

	if lines == 0 {
		return rb.PassWithOutput("No unused code detected", output), nil
	}

	msg := fmt.Sprintf("%d unused code %s found",
		lines,
		checkutil.Pluralize(lines, "item", "items"),
	)

	return rb.WarnWithOutput(msg, output), nil
}

// countFindings counts non-empty lines that look like vulture findings.
func countFindings(output string) int {
	count := 0
	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && strings.Contains(trimmed, "unused") {
			count++
		}
	}
	return count
}
