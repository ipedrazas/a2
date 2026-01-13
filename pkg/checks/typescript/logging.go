package typescriptcheck

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// LoggingCheck detects logging libraries and warns about console.log usage.
type LoggingCheck struct{}

func (c *LoggingCheck) ID() string   { return "typescript:logging" }
func (c *LoggingCheck) Name() string { return "TypeScript Logging" }

// Run checks for logging configuration.
func (c *LoggingCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangTypeScript,
	}

	// Check for tsconfig.json
	if !safepath.Exists(path, "tsconfig.json") && !safepath.Exists(path, "tsconfig.base.json") {
		result.Passed = false
		result.Status = checker.Fail
		result.Message = "No tsconfig.json found"
		return result, nil
	}

	// Detect logging libraries
	var loggingLibs []string

	pkg, err := ParsePackageJSON(path)
	if err == nil {
		// Check for common logging libraries
		loggingDeps := map[string]string{
			"winston":       "Winston",
			"pino":          "Pino",
			"bunyan":        "Bunyan",
			"log4js":        "Log4js",
			"loglevel":      "LogLevel",
			"signale":       "Signale",
			"tslog":         "tslog",
			"@sentry/node":  "Sentry",
			"@sentry/react": "Sentry",
			"dd-trace":      "Datadog",
			"newrelic":      "New Relic",
		}

		for dep, name := range loggingDeps {
			if _, ok := pkg.Dependencies[dep]; ok {
				loggingLibs = append(loggingLibs, name)
			}
			if _, ok := pkg.DevDependencies[dep]; ok {
				loggingLibs = append(loggingLibs, name)
			}
		}
	}

	// Check for console.log usage in source files
	consoleLogCount := c.countConsoleLog(path)

	// Build result
	if len(loggingLibs) > 0 {
		if consoleLogCount > 0 {
			result.Passed = true
			result.Status = checker.Warn
			result.Message = fmt.Sprintf("Logging: %s; found %d console.log %s (consider removing)",
				strings.Join(loggingLibs, ", "),
				consoleLogCount,
				checkutil.Pluralize(consoleLogCount, "call", "calls"))
		} else {
			result.Passed = true
			result.Status = checker.Pass
			result.Message = "Logging configured: " + strings.Join(loggingLibs, ", ")
		}
	} else if consoleLogCount > 0 {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = fmt.Sprintf("Found %d console.log %s; consider a logging library (winston, pino)",
			consoleLogCount,
			checkutil.Pluralize(consoleLogCount, "call", "calls"))
	} else {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "No logging library detected (consider winston, pino, or tslog)"
	}

	return result, nil
}

// countConsoleLog counts console.log calls in TypeScript source files.
func (c *LoggingCheck) countConsoleLog(path string) int {
	count := 0
	consoleLogRe := regexp.MustCompile(`console\.(log|info|warn|error|debug)\s*\(`)

	// Walk through source directories
	srcDirs := []string{"src", "lib", "app", "."}
	for _, srcDir := range srcDirs {
		srcPath := filepath.Join(path, srcDir)
		if _, err := os.Stat(srcPath); err != nil {
			continue
		}

		err := filepath.Walk(srcPath, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			// Skip node_modules, dist, build, and test directories
			if info.IsDir() {
				name := info.Name()
				if name == "node_modules" || name == "dist" || name == "build" ||
					name == ".git" || name == "coverage" || name == "__tests__" {
					return filepath.SkipDir
				}
				return nil
			}

			// Only check TypeScript/JavaScript files
			ext := filepath.Ext(filePath)
			if ext != ".ts" && ext != ".tsx" && ext != ".js" && ext != ".jsx" {
				return nil
			}

			// Skip test files
			base := filepath.Base(filePath)
			if strings.Contains(base, ".test.") || strings.Contains(base, ".spec.") ||
				strings.HasSuffix(base, "_test.ts") || strings.HasSuffix(base, "_test.tsx") {
				return nil
			}

			// Read and scan file
			content, err := os.ReadFile(filePath) // #nosec G304 - filePath is from controlled walk
			if err != nil {
				return nil
			}

			matches := consoleLogRe.FindAllString(string(content), -1)
			count += len(matches)

			return nil
		})

		if err != nil {
			continue
		}

		// Only scan the first valid source directory
		if srcDir != "." {
			break
		}
	}

	return count
}
