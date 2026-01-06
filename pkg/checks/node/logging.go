package nodecheck

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// LoggingCheck verifies proper logging practices in Node.js code.
type LoggingCheck struct{}

func (c *LoggingCheck) ID() string   { return "node:logging" }
func (c *LoggingCheck) Name() string { return "Node Logging" }

func (c *LoggingCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangNode,
	}

	// Check if package.json exists
	if !safepath.Exists(path, "package.json") {
		result.Status = checker.Fail
		result.Passed = false
		result.Message = "package.json not found"
		return result, nil
	}

	// Structured logging libraries to detect
	loggingLibs := []string{
		"winston",
		"pino",
		"bunyan",
		"log4js",
		"loglevel",
		"signale",
	}

	// Patterns for console.* calls
	consolePattern := regexp.MustCompile(`\bconsole\.(log|error|warn|info|debug)\s*\(`)

	hasLoggingLib := false
	consoleCount := 0

	// First check package.json for logging libraries
	pkg, err := ParsePackageJSON(path)
	if err == nil && pkg != nil {
		for _, lib := range loggingLibs {
			if _, ok := pkg.Dependencies[lib]; ok {
				hasLoggingLib = true
				break
			}
			if _, ok := pkg.DevDependencies[lib]; ok {
				hasLoggingLib = true
				break
			}
		}
	}

	// Walk through JS/TS files
	err = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip directories
		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "dist" || name == "build" {
				return filepath.SkipDir
			}
			return nil
		}

		// Only check JS/TS files
		ext := strings.ToLower(filepath.Ext(filePath))
		if ext != ".js" && ext != ".ts" && ext != ".jsx" && ext != ".tsx" && ext != ".mjs" && ext != ".cjs" {
			return nil
		}

		// Skip test files and config files
		baseName := filepath.Base(filePath)
		isTestFile := strings.Contains(baseName, ".test.") ||
			strings.Contains(baseName, ".spec.") ||
			strings.HasSuffix(baseName, ".config.js") ||
			strings.HasSuffix(baseName, ".config.ts")

		file, err := os.Open(filePath)
		if err != nil {
			return nil
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			trimmed := strings.TrimSpace(line)

			// Skip comments
			if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*") {
				continue
			}

			// Check for logging library imports (if not already found)
			if !hasLoggingLib {
				for _, lib := range loggingLibs {
					if strings.Contains(line, `"`+lib+`"`) || strings.Contains(line, `'`+lib+`'`) {
						hasLoggingLib = true
						break
					}
				}
			}

			// Check for console.* in non-test files
			if !isTestFile && consolePattern.MatchString(line) {
				consoleCount++
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
	if hasLoggingLib && consoleCount == 0 {
		result.Passed = true
		result.Status = checker.Pass
		result.Message = "Uses structured logging, no console.* statements"
		return result, nil
	}

	if hasLoggingLib && consoleCount > 0 {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = fmt.Sprintf("Uses structured logging but found %d console.* statement(s)", consoleCount)
		return result, nil
	}

	if !hasLoggingLib && consoleCount == 0 {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "No structured logging detected (consider using winston or pino)"
		return result, nil
	}

	// !hasLoggingLib && consoleCount > 0
	result.Passed = false
	result.Status = checker.Warn
	result.Message = fmt.Sprintf("No structured logging and found %d console.* statement(s)", consoleCount)
	return result, nil
}
