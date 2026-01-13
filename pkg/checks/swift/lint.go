package swiftcheck

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// LintCheck runs SwiftLint to check for common mistakes and style issues.
type LintCheck struct {
	Config *config.SwiftLanguageConfig
}

func (c *LintCheck) ID() string   { return "swift:lint" }
func (c *LintCheck) Name() string { return "Swift Lint" }

// Run executes swiftlint.
func (c *LintCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangSwift)

	// Check for Package.swift first
	if !safepath.Exists(path, "Package.swift") {
		return rb.Fail("No Package.swift found"), nil
	}

	// Check if swiftlint is available
	if _, err := exec.LookPath("swiftlint"); err != nil {
		return rb.Info("SwiftLint not installed (install with 'brew install swiftlint')"), nil
	}

	// Check for SwiftLint config
	hasConfig := safepath.Exists(path, ".swiftlint.yml") || safepath.Exists(path, ".swiftlint.yaml")

	// Run swiftlint
	cmd := exec.Command("swiftlint", "lint", "--quiet")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		// SwiftLint found issues
		// Try to count warnings/errors
		warningRe := regexp.MustCompile(`warning:`)
		errorRe := regexp.MustCompile(`error:`)

		warnings := len(warningRe.FindAllString(outputStr, -1))
		errors := len(errorRe.FindAllString(outputStr, -1))

		var msg strings.Builder
		msg.WriteString("SwiftLint issues found: ")
		if errors > 0 {
			msg.WriteString(formatIssueCount(errors) + " error(s)")
			if warnings > 0 {
				msg.WriteString(", ")
			}
		}
		if warnings > 0 {
			msg.WriteString(formatIssueCount(warnings) + " warning(s)")
		}
		if errors == 0 && warnings == 0 {
			msg.WriteString(err.Error())
		}

		if errors > 0 {
			return rb.Fail(msg.String()), nil
		}
		return rb.Warn(msg.String()), nil
	}

	if hasConfig {
		return rb.Pass("SwiftLint passed (custom config)"), nil
	}
	return rb.Pass("SwiftLint passed"), nil
}

// formatIssueCount converts an int to a string for display.
func formatIssueCount(n int) string {
	if n < 10 {
		return string(rune('0' + n))
	}
	return "multiple"
}
