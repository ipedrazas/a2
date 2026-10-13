package gocheck

import (
	"context"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
)

// TestsCheck runs go test to verify all tests pass.
type TestsCheck struct{}

func (c *TestsCheck) ID() string   { return "go:tests" }
func (c *TestsCheck) Name() string { return "Go Tests" }

func (c *TestsCheck) Run(ctx context.Context, path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangGo)

	result := rb.RunCommand(ctx, path, "go", "test", "./...")
	output := result.CombinedOutput()

	// Failures come first: a package without tests ("[no test files]") must
	// not mask a failing package elsewhere in the module.
	if !result.Success() {
		return rb.FailWithOutput("Tests failed: "+checkutil.TruncateMessage(result.Output(), 200), output), nil
	}

	if !hasTestedPackages(result.Stdout) {
		return rb.PassWithOutput("No test files found", output), nil
	}

	return rb.PassWithOutput("All tests passed", output), nil
}

// hasTestedPackages reports whether go test output contains at least one
// package that actually ran tests ("ok  <pkg>"), as opposed to only
// "?   <pkg> [no test files]" lines.
func hasTestedPackages(output string) bool {
	for line := range strings.Lines(output) {
		if strings.HasPrefix(line, "ok ") || strings.HasPrefix(line, "ok\t") {
			return true
		}
	}
	return false
}
