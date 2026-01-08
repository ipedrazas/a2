package rustcheck

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// TestsCheck runs cargo test to verify tests pass.
type TestsCheck struct{}

func (c *TestsCheck) ID() string   { return "rust:tests" }
func (c *TestsCheck) Name() string { return "Rust Tests" }

// Run executes cargo test.
func (c *TestsCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangRust,
	}

	// Check for Cargo.toml first
	if !safepath.Exists(path, "Cargo.toml") {
		result.Passed = false
		result.Status = checker.Fail
		result.Message = "No Cargo.toml found"
		return result, nil
	}

	// Run cargo test
	cmd := exec.Command("cargo", "test", "--quiet")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		result.Passed = false
		result.Status = checker.Fail
		// Try to extract test failure info
		if strings.Contains(outputStr, "FAILED") {
			// Count failures
			failRe := regexp.MustCompile(`(\d+) failed`)
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

	// Parse test results - cargo test output format:
	// running X tests
	// test result: ok. X passed; 0 failed; 0 ignored
	passedRe := regexp.MustCompile(`(\d+) passed`)
	matches := passedRe.FindStringSubmatch(outputStr)
	if len(matches) > 1 {
		result.Passed = true
		result.Status = checker.Pass
		result.Message = matches[1] + " test(s) passed"
	} else if strings.Contains(outputStr, "running 0 tests") || strings.Contains(outputStr, "0 passed") {
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
