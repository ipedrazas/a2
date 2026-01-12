package swiftcheck

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// TestsCheck runs swift test to verify tests pass.
type TestsCheck struct{}

func (c *TestsCheck) ID() string   { return "swift:tests" }
func (c *TestsCheck) Name() string { return "Swift Tests" }

// Run executes swift test.
func (c *TestsCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangSwift,
	}

	// Check for Package.swift first
	if !safepath.Exists(path, "Package.swift") {
		result.Passed = false
		result.Status = checker.Fail
		result.Message = "No Package.swift found"
		return result, nil
	}

	// Run swift test
	cmd := exec.Command("swift", "test")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		result.Passed = false
		result.Status = checker.Fail
		// Try to extract test failure info
		if strings.Contains(outputStr, "failed") {
			// Count failures
			failRe := regexp.MustCompile(`(\d+) test[s]? failed`)
			if matches := failRe.FindStringSubmatch(outputStr); len(matches) > 1 {
				result.Message = "Tests failed: " + matches[1] + " test(s) failed"
			} else {
				result.Message = "Tests failed"
			}
		} else {
			result.Message = "Tests failed: " + err.Error()
		}
		return result, nil
	}

	// Parse test results - swift test output format:
	// Test Suite 'All tests' passed at ...
	// Executed X tests, with 0 failures
	passedRe := regexp.MustCompile(`Executed (\d+) test`)
	matches := passedRe.FindStringSubmatch(outputStr)
	if len(matches) > 1 {
		result.Passed = true
		result.Status = checker.Pass
		result.Message = matches[1] + " test(s) passed"
	} else if strings.Contains(outputStr, "0 tests") || strings.Contains(outputStr, "no tests") {
		result.Passed = true
		result.Status = checker.Pass
		result.Message = "No tests found"
	} else {
		result.Passed = true
		result.Status = checker.Pass
		result.Message = "Tests passed"
	}

	return result, nil
}
