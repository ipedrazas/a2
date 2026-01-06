package nodecheck

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// CoverageCheck measures test coverage for Node.js projects.
type CoverageCheck struct {
	Config    *config.NodeLanguageConfig
	Threshold float64
}

// ID returns the unique identifier for this check.
func (c *CoverageCheck) ID() string {
	return "node:coverage"
}

// Name returns the human-readable name for this check.
func (c *CoverageCheck) Name() string {
	return "Node Coverage"
}

// Run executes the coverage check.
func (c *CoverageCheck) Run(path string) (checker.Result, error) {
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

	// Parse package.json to check for test script
	pkg, err := ParsePackageJSON(path)
	if err != nil {
		result.Status = checker.Fail
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to parse package.json: %v", err)
		return result, nil
	}

	// Check if test script exists
	testScript, hasTest := pkg.Scripts["test"]
	if !hasTest || testScript == "" || strings.Contains(testScript, "no test specified") {
		result.Status = checker.Warn
		result.Passed = false
		result.Message = "No test script configured, coverage is 0%"
		return result, nil
	}

	threshold := c.Threshold
	if threshold == 0 {
		threshold = 80.0
	}

	// Detect test runner and run with coverage
	runner := c.detectTestRunner(path, pkg)
	coverage, err := c.runCoverage(path, runner)

	if err != nil {
		result.Status = checker.Warn
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to measure coverage: %v", err)
		return result, nil
	}

	if coverage < 0 {
		result.Status = checker.Warn
		result.Passed = false
		result.Message = "Could not determine coverage. Ensure coverage tools are installed."
		return result, nil
	}

	if coverage < threshold {
		result.Status = checker.Warn
		result.Passed = false
		result.Message = fmt.Sprintf("Coverage %.1f%% is below threshold %.1f%%", coverage, threshold)
		return result, nil
	}

	result.Status = checker.Pass
	result.Passed = true
	result.Message = fmt.Sprintf("Coverage: %.1f%% (threshold: %.1f%%)", coverage, threshold)
	return result, nil
}

// detectTestRunner determines which test runner to use.
func (c *CoverageCheck) detectTestRunner(path string, pkg *PackageJSON) string {
	// Check config override first
	if c.Config != nil && c.Config.TestRunner != "auto" && c.Config.TestRunner != "" {
		return c.Config.TestRunner
	}

	// Check for Jest config files
	jestConfigs := []string{"jest.config.js", "jest.config.ts", "jest.config.mjs", "jest.config.cjs", "jest.config.json"}
	for _, cfg := range jestConfigs {
		if safepath.Exists(path, cfg) {
			return "jest"
		}
	}

	// Check for Vitest config files
	vitestConfigs := []string{"vitest.config.js", "vitest.config.ts", "vitest.config.mjs", "vitest.config.mts"}
	for _, cfg := range vitestConfigs {
		if safepath.Exists(path, cfg) {
			return "vitest"
		}
	}

	// Check devDependencies
	if pkg != nil {
		if _, ok := pkg.DevDependencies["jest"]; ok {
			return "jest"
		}
		if _, ok := pkg.DevDependencies["vitest"]; ok {
			return "vitest"
		}
		if _, ok := pkg.DevDependencies["c8"]; ok {
			return "c8"
		}
		if _, ok := pkg.DevDependencies["nyc"]; ok {
			return "nyc"
		}
	}

	return "jest" // Default to jest
}

// runCoverage runs the appropriate coverage command and returns the coverage percentage.
func (c *CoverageCheck) runCoverage(path, runner string) (float64, error) {
	var cmd *exec.Cmd

	switch runner {
	case "jest":
		cmd = exec.Command("npx", "jest", "--coverage", "--coverageReporters=text-summary", "--passWithNoTests")
	case "vitest":
		cmd = exec.Command("npx", "vitest", "run", "--coverage", "--reporter=text")
	case "c8":
		cmd = exec.Command("npx", "c8", "npm", "test")
	case "nyc":
		cmd = exec.Command("npx", "nyc", "npm", "test")
	default:
		// Try jest with coverage
		cmd = exec.Command("npx", "jest", "--coverage", "--coverageReporters=text-summary", "--passWithNoTests")
	}

	cmd.Dir = path
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()

	// Parse coverage from output even if command fails (tests may fail but coverage reported)
	coverage := parseNodeCoverage(output)
	if coverage >= 0 {
		return coverage, nil
	}

	if err != nil {
		return -1, fmt.Errorf("coverage command failed: %w", err)
	}

	return -1, nil
}

// parseNodeCoverage extracts coverage percentage from various coverage tool outputs.
func parseNodeCoverage(output string) float64 {
	// Jest text-summary format: "All files | XX.XX | XX.XX | XX.XX | XX.XX"
	// or "Statements   : XX.XX% ( X/Y )"
	jestPattern := regexp.MustCompile(`(?:Statements|All files)\s*[|:]\s*(\d+(?:\.\d+)?)\s*%?`)
	if matches := jestPattern.FindStringSubmatch(output); len(matches) > 1 {
		if coverage, err := strconv.ParseFloat(matches[1], 64); err == nil {
			return coverage
		}
	}

	// c8/nyc format: "Statements   : XX.XX% ( X/Y )"
	c8Pattern := regexp.MustCompile(`Statements\s*:\s*(\d+(?:\.\d+)?)\s*%`)
	if matches := c8Pattern.FindStringSubmatch(output); len(matches) > 1 {
		if coverage, err := strconv.ParseFloat(matches[1], 64); err == nil {
			return coverage
		}
	}

	// Vitest format: "All files | XX.XX | XX.XX | XX.XX | XX.XX"
	vitestPattern := regexp.MustCompile(`All files\s*\|\s*(\d+(?:\.\d+)?)\s*\|`)
	if matches := vitestPattern.FindStringSubmatch(output); len(matches) > 1 {
		if coverage, err := strconv.ParseFloat(matches[1], 64); err == nil {
			return coverage
		}
	}

	// Generic pattern: "Coverage: XX.XX%" or "XX.XX% coverage"
	genericPattern := regexp.MustCompile(`(?:Coverage|coverage)[:\s]*(\d+(?:\.\d+)?)\s*%|(\d+(?:\.\d+)?)\s*%\s*(?:Coverage|coverage)`)
	if matches := genericPattern.FindStringSubmatch(output); len(matches) > 1 {
		for _, match := range matches[1:] {
			if match != "" {
				if coverage, err := strconv.ParseFloat(match, 64); err == nil {
					return coverage
				}
			}
		}
	}

	return -1
}
