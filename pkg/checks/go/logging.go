package gocheck

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// LoggingCheck verifies proper logging practices in Go code.
type LoggingCheck struct{}

func (c *LoggingCheck) ID() string   { return "go:logging" }
func (c *LoggingCheck) Name() string { return "Go Logging" }

func (c *LoggingCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangGo,
	}

	// Structured logging imports to detect
	structuredLoggers := []string{
		"log/slog",
		"go.uber.org/zap",
		"github.com/rs/zerolog",
		"github.com/sirupsen/logrus",
	}

	// Bad patterns (fmt.Print* in non-test files)
	printPattern := regexp.MustCompile(`\bfmt\.(Print|Println|Printf)\s*\(`)

	hasStructuredLogger := false
	printCount := 0
	var filesWithPrints []string

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip directories
		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "vendor" || name == "testdata" {
				return filepath.SkipDir
			}
			return nil
		}

		// Only check .go files
		if !strings.HasSuffix(filePath, ".go") {
			return nil
		}

		// Skip test files for print detection
		isTestFile := strings.HasSuffix(filePath, "_test.go")

		file, err := safepath.OpenPath(path, filePath)
		if err != nil {
			return nil
		}
		defer func() {
			if err := file.Close(); err != nil {
				// File close errors are typically not critical in read-only scenarios
				fmt.Println("Error closing file:", err)
			}
		}()

		scanner := bufio.NewScanner(file)
		lineNum := 0
		filePrintCount := 0

		for scanner.Scan() {
			line := scanner.Text()
			lineNum++

			// Check for structured logger imports
			if !hasStructuredLogger {
				for _, logger := range structuredLoggers {
					if strings.Contains(line, `"`+logger+`"`) {
						hasStructuredLogger = true
						break
					}
				}
			}

			// Check for fmt.Print* in non-test files
			if !isTestFile && printPattern.MatchString(line) {
				// Skip if it's in a comment
				trimmed := strings.TrimSpace(line)
				if !strings.HasPrefix(trimmed, "//") {
					printCount++
					filePrintCount++
				}
			}
		}

		if filePrintCount > 0 {
			relPath, _ := filepath.Rel(path, filePath)
			filesWithPrints = append(filesWithPrints, relPath)
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
	if hasStructuredLogger && printCount == 0 {
		result.Passed = true
		result.Status = checker.Pass
		result.Message = "Uses structured logging, no fmt.Print statements"
		return result, nil
	}

	if hasStructuredLogger && printCount > 0 {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "Uses structured logging but found " + checkutil.PluralizeCount(printCount, "fmt.Print statement", "fmt.Print statements")
		return result, nil
	}

	if !hasStructuredLogger && printCount == 0 {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "No structured logging detected (consider slog, zap, or zerolog)"
		return result, nil
	}

	// !hasStructuredLogger && printCount > 0
	result.Passed = false
	result.Status = checker.Warn
	result.Message = "No structured logging and found " + checkutil.PluralizeCount(printCount, "fmt.Print statement", "fmt.Print statements")
	return result, nil
}
