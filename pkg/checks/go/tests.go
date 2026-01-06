package gocheck

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
)

// TestsCheck runs go test to verify all tests pass.
type TestsCheck struct{}

func (c *TestsCheck) ID() string   { return "go:tests" }
func (c *TestsCheck) Name() string { return "Go Tests" }

func (c *TestsCheck) Run(path string) (checker.Result, error) {
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
				Name:     c.Name(),
				ID:       c.ID(),
				Passed:   true,
				Status:   checker.Pass,
				Message:  "No test files found",
				Language: checker.LangGo,
			}, nil
		}

		return checker.Result{
			Name:     c.Name(),
			ID:       c.ID(),
			Passed:   false,
			Status:   checker.Fail, // Critical - stops execution
			Message:  "Tests failed: " + output,
			Language: checker.LangGo,
		}, nil
	}

	// Parse output to count tests
	output := stdout.String()
	if strings.Contains(output, "no test files") {
		return checker.Result{
			Name:     c.Name(),
			ID:       c.ID(),
			Passed:   true,
			Status:   checker.Pass,
			Message:  "No test files found",
			Language: checker.LangGo,
		}, nil
	}

	return checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Passed:   true,
		Status:   checker.Pass,
		Message:  "All tests passed",
		Language: checker.LangGo,
	}, nil
}
