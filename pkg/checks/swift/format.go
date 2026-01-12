package swiftcheck

import (
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
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

	// Determine which formatter to use
	formatter := c.detectFormatter(path)
	if formatter == "" {
		result.Passed = true
		result.Status = checker.Info
		result.Message = "No Swift formatter found (install swift-format or swiftformat)"
		return result, nil
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

	if err != nil {
		outputStr := strings.TrimSpace(string(output))
		// Count issues
		lines := strings.Split(outputStr, "\n")
		issueCount := 0
		for _, line := range lines {
			if strings.Contains(line, "warning:") || strings.Contains(line, "error:") {
				issueCount++
			}
		}
		result.Passed = false
		result.Status = checker.Warn
		if issueCount > 0 {
			result.Message = "Code not formatted: " + formatCount(issueCount) + " issue(s) found"
		} else {
			result.Message = "Code not formatted (run '" + formatter + "')"
		}
		return result, nil
	}

	result.Passed = true
	result.Status = checker.Pass
	if hasConfig {
		result.Message = "Code formatted (" + formatter + ", custom config)"
	} else {
		result.Message = "Code formatted (" + formatter + ")"
	}

	return result, nil
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
