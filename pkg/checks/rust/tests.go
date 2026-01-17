package rustcheck

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// TestsCheck runs cargo test to verify tests pass.
type TestsCheck struct{}

func (c *TestsCheck) ID() string   { return "rust:tests" }
func (c *TestsCheck) Name() string { return "Rust Tests" }

// Run executes cargo test.
func (c *TestsCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangRust)

	// Check for Cargo.toml first
	if !safepath.Exists(path, "Cargo.toml") {
		return rb.Fail("No Cargo.toml found"), nil
	}

	// Run cargo test
	cmd := exec.Command("cargo", "test", "--quiet")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		// Try to extract test failure info
		if strings.Contains(outputStr, "FAILED") {
			// Count failures
			failRe := regexp.MustCompile(`(\d+) failed`)
			if matches := failRe.FindStringSubmatch(outputStr); len(matches) > 1 {
				return rb.FailWithOutput("Tests failed: "+matches[1]+" test(s) failed", outputStr), nil
			}
			return rb.FailWithOutput("Tests failed", outputStr), nil
		}
		return rb.FailWithOutput("Tests failed: "+err.Error(), outputStr), nil
	}

	// Parse test results - cargo test output format:
	// running X tests
	// test result: ok. X passed; 0 failed; 0 ignored
	passedRe := regexp.MustCompile(`(\d+) passed`)
	matches := passedRe.FindStringSubmatch(outputStr)
	if len(matches) > 1 {
		return rb.PassWithOutput(matches[1]+" test(s) passed", outputStr), nil
	}
	if strings.Contains(outputStr, "running 0 tests") || strings.Contains(outputStr, "0 passed") {
		return rb.Pass("No tests found"), nil
	}
	return rb.PassWithOutput("Tests passed", outputStr), nil
}
