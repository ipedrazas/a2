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
				return rb.Fail("Tests failed: " + matches[1] + " test(s) failed"), nil
			}
			return rb.Fail("Tests failed"), nil
		}
		return rb.Fail("Tests failed: " + err.Error()), nil
	}

	// Parse test results - swift test output format:
	// Test Suite 'All tests' passed at ...
	// Executed X tests, with 0 failures
	passedRe := regexp.MustCompile(`Executed (\d+) test`)
	matches := passedRe.FindStringSubmatch(outputStr)
	if len(matches) > 1 {
		return rb.Pass(matches[1] + " test(s) passed"), nil
	} else if strings.Contains(outputStr, "0 tests") || strings.Contains(outputStr, "no tests") {
		return rb.Pass("No tests found"), nil
	}
	return rb.Pass("Tests passed"), nil
}
