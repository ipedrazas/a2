package rustcheck

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// LintCheck runs clippy to check for common mistakes and style issues.
type LintCheck struct{}

func (c *LintCheck) ID() string   { return "rust:lint" }
func (c *LintCheck) Name() string { return "Rust Clippy" }

// Run executes cargo clippy.
func (c *LintCheck) Run(path string) (checker.Result, error) {
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

	// Check for clippy config
	hasConfig := safepath.Exists(path, "clippy.toml") || safepath.Exists(path, ".clippy.toml")

	// Run cargo clippy
	cmd := exec.Command("cargo", "clippy", "--quiet", "--", "-D", "warnings")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		// Clippy found issues
		// Try to count warnings/errors
		warningRe := regexp.MustCompile(`warning:`)
		errorRe := regexp.MustCompile(`error\[`)

		warnings := len(warningRe.FindAllString(outputStr, -1))
		errors := len(errorRe.FindAllString(outputStr, -1))

		var msg strings.Builder
		msg.WriteString("Clippy issues found: ")
		if errors > 0 {
			msg.WriteString(string(rune('0'+errors)) + " error(s)")
			if warnings > 0 {
				msg.WriteString(", ")
			}
		}
		if warnings > 0 {
			msg.WriteString(string(rune('0'+warnings)) + " warning(s)")
		}
		if errors == 0 && warnings == 0 {
			msg.WriteString(err.Error())
		}

		result.Passed = false
		if errors > 0 {
			result.Status = checker.Fail
		} else {
			result.Status = checker.Warn
		}
		result.Message = msg.String()
		return result, nil
	}

	result.Passed = true
	result.Status = checker.Pass
	if hasConfig {
		result.Message = "Clippy passed (custom config)"
	} else {
		result.Message = "Clippy passed"
	}

	return result, nil
}
