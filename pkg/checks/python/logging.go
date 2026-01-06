package pythoncheck

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
)

// LoggingCheck verifies proper logging practices in Python code.
type LoggingCheck struct{}

func (c *LoggingCheck) ID() string   { return "python:logging" }
func (c *LoggingCheck) Name() string { return "Python Logging" }

func (c *LoggingCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangPython,
	}

	// Structured logging imports to detect
	loggingImports := []string{
		"import logging",
		"from logging import",
		"import structlog",
		"from structlog import",
		"from loguru import",
		"import loguru",
	}

	// Pattern to detect print() calls (excluding comments and strings)
	printPattern := regexp.MustCompile(`^\s*print\s*\(`)

	hasLoggingImport := false
	printCount := 0

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip directories
		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "venv" || name == "__pycache__" || name == ".venv" || name == "env" {
				return filepath.SkipDir
			}
			return nil
		}

		// Only check .py files
		if !strings.HasSuffix(filePath, ".py") {
			return nil
		}

		// Skip test files for print detection
		baseName := filepath.Base(filePath)
		isTestFile := strings.HasPrefix(baseName, "test_") ||
			strings.HasSuffix(baseName, "_test.py") ||
			baseName == "conftest.py"

		file, err := os.Open(filePath)
		if err != nil {
			return nil
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		inMultilineString := false

		for scanner.Scan() {
			line := scanner.Text()
			trimmed := strings.TrimSpace(line)

			// Skip empty lines and comments
			if trimmed == "" || strings.HasPrefix(trimmed, "#") {
				continue
			}

			// Track multiline strings (basic detection)
			if strings.Count(line, `"""`)%2 == 1 || strings.Count(line, `'''`)%2 == 1 {
				inMultilineString = !inMultilineString
			}
			if inMultilineString {
				continue
			}

			// Check for logging imports
			if !hasLoggingImport {
				for _, imp := range loggingImports {
					if strings.Contains(line, imp) {
						hasLoggingImport = true
						break
					}
				}
			}

			// Check for print() in non-test files
			if !isTestFile && printPattern.MatchString(trimmed) {
				// Skip if it looks like it's in a __name__ guard
				// (very basic check)
				printCount++
			}
		}

		return nil
	})

	if err != nil {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "Error scanning files: " + err.Error()
		return result, nil
	}

	// Determine result based on findings
	if hasLoggingImport && printCount == 0 {
		result.Passed = true
		result.Status = checker.Pass
		result.Message = "Uses logging module, no print() statements"
		return result, nil
	}

	if hasLoggingImport && printCount > 0 {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = fmt.Sprintf("Uses logging but found %d print() statement(s)", printCount)
		return result, nil
	}

	if !hasLoggingImport && printCount == 0 {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "No logging module detected (consider using logging or structlog)"
		return result, nil
	}

	// !hasLoggingImport && printCount > 0
	result.Passed = false
	result.Status = checker.Warn
	result.Message = fmt.Sprintf("No logging module and found %d print() statement(s)", printCount)
	return result, nil
}
