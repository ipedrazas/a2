package pythoncheck

import (
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
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
	rb := checkutil.NewResultBuilder(c, checker.LangPython)
	runner := c.detectTestRunner(path)

	// First check if tests exist
	var collectResult *checkutil.CommandResult
	switch runner {
	case "pytest":
		collectResult = runPythonCommand(path, "pytest", "--collect-only", "-q")
	case "unittest":
		collectResult = runPythonCommand(path, "python", "-m", "unittest", "discover", "--list")
	default:
		collectResult = runPythonCommand(path, "pytest", "--collect-only", "-q")
	}

	// Check if test runner is not installed
	if checkutil.ToolNotFoundError(collectResult.Err) {
		return rb.ToolNotInstalled(runner, ""), nil
	}

	// Check for no tests
	output := collectResult.Stdout
	errOutput := collectResult.Stderr
	if strings.Contains(output, "no tests ran") ||
		strings.Contains(output, "0 tests") ||
		strings.Contains(output, "collected 0 items") ||
		strings.Contains(errOutput, "no tests") {
		return rb.Pass("No tests found"), nil
	}

	// Now run actual tests
	var testResult *checkutil.CommandResult
	switch runner {
	case "pytest":
		testResult = runPythonCommand(path, "pytest", "-v", "--tb=short")
	case "unittest":
		testResult = runPythonCommand(path, "python", "-m", "unittest", "discover", "-v")
	default:
		testResult = runPythonCommand(path, "pytest", "-v", "--tb=short")
	}

	testOutput := testResult.CombinedOutput()
	if !testResult.Success() {
		return rb.FailWithOutput("Tests failed: "+checkutil.TruncateMessage(testResult.Output(), 200), testOutput), nil
	}

	return rb.PassWithOutput("All tests passed", testOutput), nil
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
