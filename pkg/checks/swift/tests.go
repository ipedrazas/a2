package swiftcheck

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// TestsCheck runs swift test to verify tests pass.
type TestsCheck struct{}

func (c *TestsCheck) ID() string   { return "swift:tests" }
func (c *TestsCheck) Name() string { return "Swift Tests" }

// Run executes swift test.
func (c *TestsCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangSwift)

	// Check for Package.swift first
	if !safepath.Exists(path, "Package.swift") {
		return rb.Fail("No Package.swift found"), nil
	}

	// Run swift test
	cmd := exec.Command("swift", "test")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		// Try to extract test failure info
		if strings.Contains(outputStr, "failed") {
			// Count failures
			failRe := regexp.MustCompile(`(\d+) test[s]? failed`)
			if matches := failRe.FindStringSubmatch(outputStr); len(matches) > 1 {
				return rb.FailWithOutput("Tests failed: "+matches[1]+" test(s) failed", outputStr), nil
			}
			return rb.FailWithOutput("Tests failed", outputStr), nil
		}
		return rb.FailWithOutput("Tests failed: "+err.Error(), outputStr), nil
	}

	// Parse test results - swift test output format:
	// Test Suite 'All tests' passed at ...
	// Executed X tests, with 0 failures
	passedRe := regexp.MustCompile(`Executed (\d+) test`)
	matches := passedRe.FindStringSubmatch(outputStr)
	if len(matches) > 1 {
		return rb.PassWithOutput(matches[1]+" test(s) passed", outputStr), nil
	} else if strings.Contains(outputStr, "0 tests") || strings.Contains(outputStr, "no tests") {
		return rb.PassWithOutput("No tests found", outputStr), nil
	}
	return rb.PassWithOutput("Tests passed", outputStr), nil
}
