package checks

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
)

// TestRunnerCheck runs go test to verify all tests pass.
type TestRunnerCheck struct{}

func (c *TestRunnerCheck) ID() string   { return "tests" }
func (c *TestRunnerCheck) Name() string { return "Unit Tests" }

func (c *TestRunnerCheck) Run(path string) (checker.Result, error) {
	cmd := exec.Command("go", "test", "./...")
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// Combine output for error message
		output := strings.TrimSpace(stderr.String())
		if output == "" {
			output = strings.TrimSpace(stdout.String())
		}

		// Check if it's "no test files" which is not a failure
		if strings.Contains(output, "no test files") || strings.Contains(stdout.String(), "no test files") {
			return checker.Result{
				Name:    c.Name(),
				ID:      c.ID(),
				Passed:  true,
				Status:  checker.Pass,
				Message: "No test files found",
			}, nil
		}

		return checker.Result{
			Name:    c.Name(),
			ID:      c.ID(),
			Passed:  false,
			Status:  checker.Fail, // Critical - stops execution
			Message: "Tests failed: " + output,
		}, nil
	}

	// Parse output to count tests
	output := stdout.String()
	if strings.Contains(output, "no test files") {
		return checker.Result{
			Name:    c.Name(),
			ID:      c.ID(),
			Passed:  true,
			Status:  checker.Pass,
			Message: "No test files found",
		}, nil
	}

	return checker.Result{
		Name:    c.Name(),
		ID:      c.ID(),
		Passed:  true,
		Status:  checker.Pass,
		Message: "All tests passed",
	}, nil
}
