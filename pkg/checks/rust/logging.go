package rustcheck

import (
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// LoggingCheck verifies structured logging is configured.
type LoggingCheck struct{}

func (c *LoggingCheck) ID() string   { return "rust:logging" }
func (c *LoggingCheck) Name() string { return "Rust Logging" }

// Run checks for structured logging libraries.
func (c *LoggingCheck) Run(path string) (checker.Result, error) {
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

	// Read Cargo.toml to check dependencies
	content, err := safepath.ReadFile(path, "Cargo.toml")
	if err != nil {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "Cannot read Cargo.toml"
		return result, nil
	}

	contentStr := strings.ToLower(string(content))
	var loggingLibs []string

	// Check for popular logging crates
	loggingCrates := map[string]string{
		"tracing":           "tracing",
		"log":               "log",
		"env_logger":        "env_logger",
		"fern":              "fern",
		"slog":              "slog",
		"flexi_logger":      "flexi_logger",
		"log4rs":            "log4rs",
		"pretty_env_logger": "pretty_env_logger",
	}

	// Structured logging crates (preferred)
	structuredCrates := map[string]string{
		"tracing":            "tracing",
		"tracing-subscriber": "tracing-subscriber",
		"slog":               "slog",
		"slog-json":          "slog-json",
	}

	for crate, name := range structuredCrates {
		if strings.Contains(contentStr, crate) {
			loggingLibs = append(loggingLibs, name)
		}
	}

	// If no structured logging, check for basic logging
	if len(loggingLibs) == 0 {
		for crate, name := range loggingCrates {
			if strings.Contains(contentStr, crate) {
				loggingLibs = append(loggingLibs, name)
			}
		}
	}

	// Check for println! usage in source files (anti-pattern)
	hasPrintln := c.checkForPrintln(path)

	// Build result
	if len(loggingLibs) > 0 {
		result.Passed = true
		result.Status = checker.Pass
		result.Message = "Logging configured: " + strings.Join(uniqueStrings(loggingLibs), ", ")
		if hasPrintln {
			result.Status = checker.Warn
			result.Message += " (but println! found in source)"
		}
	} else if hasPrintln {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "Using println! instead of structured logging (consider tracing or log crate)"
	} else {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "No logging library detected (consider tracing or log crate)"
	}

	return result, nil
}

// checkForPrintln looks for println! macro usage in Rust source files.
func (c *LoggingCheck) checkForPrintln(path string) bool {
	// Check main.rs and lib.rs
	sourceFiles := []string{"src/main.rs", "src/lib.rs"}

	for _, file := range sourceFiles {
		if content, err := safepath.ReadFile(path, file); err == nil {
			// Look for println! or print! macros (excluding comments)
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				// Skip comments
				if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") {
					continue
				}
				// Check for println! or print! macros
				if strings.Contains(trimmed, "println!") || strings.Contains(trimmed, "print!") {
					// Make sure it's not in a test module
					if !strings.Contains(trimmed, "#[test]") {
						return true
					}
				}
			}
		}
	}

	return false
}

// uniqueStrings removes duplicates from a string slice.
func uniqueStrings(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, s := range slice {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}
