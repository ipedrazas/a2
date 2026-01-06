package pythoncheck

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// TestsCheck runs Python tests.
type TestsCheck struct {
	Config *config.PythonLanguageConfig
}

func (c *TestsCheck) ID() string   { return "python:tests" }
func (c *TestsCheck) Name() string { return "Python Tests" }

func (c *TestsCheck) Run(path string) (checker.Result, error) {
	runner := c.detectTestRunner(path)

	var cmd *exec.Cmd
	switch runner {
	case "pytest":
		cmd = exec.Command("pytest", "--collect-only", "-q")
	case "unittest":
		cmd = exec.Command("python", "-m", "unittest", "discover", "--list")
	default:
		cmd = exec.Command("pytest", "--collect-only", "-q")
	}

	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := strings.TrimSpace(stdout.String())
	errOutput := strings.TrimSpace(stderr.String())

	// Check if test runner is not installed
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			return checker.Result{
				Name:     c.Name(),
				ID:       c.ID(),
				Passed:   true,
				Status:   checker.Pass,
				Message:  runner + " not installed, skipping tests",
				Language: checker.LangPython,
			}, nil
		}
	}

	// Check for no tests
	if strings.Contains(output, "no tests ran") ||
		strings.Contains(output, "0 tests") ||
		strings.Contains(output, "collected 0 items") ||
		strings.Contains(errOutput, "no tests") {
		return checker.Result{
			Name:     c.Name(),
			ID:       c.ID(),
			Passed:   true,
			Status:   checker.Pass,
			Message:  "No tests found",
			Language: checker.LangPython,
		}, nil
	}

	// Now run actual tests
	switch runner {
	case "pytest":
		cmd = exec.Command("pytest", "-v", "--tb=short")
	case "unittest":
		cmd = exec.Command("python", "-m", "unittest", "discover", "-v")
	default:
		cmd = exec.Command("pytest", "-v", "--tb=short")
	}

	cmd.Dir = path
	stdout.Reset()
	stderr.Reset()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		output = strings.TrimSpace(stderr.String())
		if output == "" {
			output = strings.TrimSpace(stdout.String())
		}

		return checker.Result{
			Name:     c.Name(),
			ID:       c.ID(),
			Passed:   false,
			Status:   checker.Fail,
			Message:  "Tests failed: " + truncateMessage(output, 200),
			Language: checker.LangPython,
		}, nil
	}

	return checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Passed:   true,
		Status:   checker.Pass,
		Message:  "All tests passed",
		Language: checker.LangPython,
	}, nil
}

func (c *TestsCheck) detectTestRunner(path string) string {
	if c.Config != nil && c.Config.TestRunner != "auto" && c.Config.TestRunner != "" {
		return c.Config.TestRunner
	}

	// Auto-detect: pytest if pytest.ini or conftest.py exists
	if safepath.Exists(path, "pytest.ini") ||
		safepath.Exists(path, "conftest.py") ||
		safepath.Exists(path, "pyproject.toml") {
		return "pytest"
	}

	return "pytest" // Default
}

func truncateMessage(msg string, maxLen int) string {
	if len(msg) <= maxLen {
		return msg
	}
	return msg[:maxLen] + "..."
}
