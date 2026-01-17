package swiftcheck

import (
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// FormatCheck verifies Swift code is formatted.
type FormatCheck struct {
	Config *config.SwiftLanguageConfig
}

func (c *FormatCheck) ID() string   { return "swift:format" }
func (c *FormatCheck) Name() string { return "Swift Format" }

// Run checks if code is properly formatted using swift-format or swiftformat.
func (c *FormatCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangSwift)

	// Check for Package.swift first
	if !safepath.Exists(path, "Package.swift") {
		return rb.Fail("No Package.swift found"), nil
	}

	// Determine which formatter to use
	formatter := c.detectFormatter(path)
	if formatter == "" {
		return rb.Info("No Swift formatter found (install swift-format or swiftformat)"), nil
	}

	// Check for formatter config
	hasConfig := c.hasFormatterConfig(path, formatter)

	// Run formatter check
	var cmd *exec.Cmd
	switch formatter {
	case "swift-format":
		// swift-format lint --recursive .
		cmd = exec.Command("swift-format", "lint", "--recursive", ".")
	case "swiftformat":
		// swiftformat --lint .
		cmd = exec.Command("swiftformat", "--lint", ".")
	}

	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		trimmedOutput := strings.TrimSpace(outputStr)
		// Count issues
		lines := strings.Split(trimmedOutput, "\n")
		issueCount := 0
		for _, line := range lines {
			if strings.Contains(line, "warning:") || strings.Contains(line, "error:") {
				issueCount++
			}
		}
		if issueCount > 0 {
			return rb.WarnWithOutput("Code not formatted: "+formatCount(issueCount)+" issue(s) found", outputStr), nil
		}
		return rb.WarnWithOutput("Code not formatted (run '"+formatter+"')", outputStr), nil
	}

	if hasConfig {
		return rb.PassWithOutput("Code formatted ("+formatter+", custom config)", outputStr), nil
	}
	return rb.PassWithOutput("Code formatted ("+formatter+")", outputStr), nil
}

// detectFormatter determines which formatter is available.
func (c *FormatCheck) detectFormatter(path string) string {
	// Check config preference first
	if c.Config != nil && c.Config.Formatter != "" && c.Config.Formatter != "auto" {
		// Verify the configured formatter is available
		if c.Config.Formatter == "swift-format" {
			if _, err := exec.LookPath("swift-format"); err == nil {
				return "swift-format"
			}
		} else if c.Config.Formatter == "swiftformat" {
			if _, err := exec.LookPath("swiftformat"); err == nil {
				return "swiftformat"
			}
		}
	}

	// Auto-detect: prefer swift-format (Apple's official)
	if _, err := exec.LookPath("swift-format"); err == nil {
		return "swift-format"
	}
	if _, err := exec.LookPath("swiftformat"); err == nil {
		return "swiftformat"
	}

	return ""
}

// hasFormatterConfig checks if a formatter config file exists.
func (c *FormatCheck) hasFormatterConfig(path string, formatter string) bool {
	switch formatter {
	case "swift-format":
		return safepath.Exists(path, ".swift-format")
	case "swiftformat":
		return safepath.Exists(path, ".swiftformat")
	}
	return false
}

// formatCount converts an int to a string for display.
func formatCount(n int) string {
	if n < 10 {
		return string(rune('0' + n))
	}
	return "multiple"
}
