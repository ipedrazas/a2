package gocheck

import (
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
)

// TestsCheck runs go test to verify all tests pass.
type TestsCheck struct{}

func (c *TestsCheck) ID() string   { return "go:tests" }
func (c *TestsCheck) Name() string { return "Go Tests" }

func (c *TestsCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangGo)

	result := checkutil.RunCommand(path, "go", "test", "./...")

	// Check if it's "no test files" which is not a failure
	if strings.Contains(result.Stdout, "no test files") || strings.Contains(result.Stderr, "no test files") {
		return rb.Pass("No test files found"), nil
	}

	if !result.Success() {
		return rb.Fail("Tests failed: " + result.Output()), nil
	}

	return rb.Pass("All tests passed"), nil
}
