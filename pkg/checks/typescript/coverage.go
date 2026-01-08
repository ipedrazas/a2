package typescriptcheck

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// CoverageCheck checks for test coverage configuration and reports.
type CoverageCheck struct {
	Config    *config.TypeScriptLanguageConfig
	Threshold float64
}

func (c *CoverageCheck) ID() string   { return "typescript:coverage" }
func (c *CoverageCheck) Name() string { return "TypeScript Coverage" }

// Run checks for coverage tooling and reports.
func (c *CoverageCheck) Run(path string) (checker.Result, error) {
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

	// Look for coverage configuration
	var coverageTools []string

	// Check for Jest coverage config
	if c.hasJestCoverage(path) {
		coverageTools = append(coverageTools, "Jest coverage")
	}

	// Check for Vitest coverage config
	if c.hasVitestCoverage(path) {
		coverageTools = append(coverageTools, "Vitest coverage")
	}

	// Check for c8 (Node.js coverage)
	if c.hasC8(path) {
		coverageTools = append(coverageTools, "c8")
	}

	// Check for nyc (Istanbul)
	if c.hasNYC(path) {
		coverageTools = append(coverageTools, "nyc/Istanbul")
	}

	// Check for coverage in CI configs
	ciCoverage := c.checkCICoverage(path)
	if ciCoverage != "" {
		coverageTools = append(coverageTools, ciCoverage)
	}

	// Try to find coverage reports
	coverage := c.findCoverageReports(path)

	// Build result
	if coverage >= 0 {
		result.Passed = coverage >= c.Threshold
		if result.Passed {
			result.Status = checker.Pass
		} else {
			result.Status = checker.Warn
		}
		result.Message = fmt.Sprintf("Coverage %.1f%% (threshold %.1f%%)", coverage, c.Threshold)
	} else if len(coverageTools) > 0 {
		result.Passed = true
		result.Status = checker.Pass
		result.Message = "Coverage configured: " + strings.Join(coverageTools, ", ")
	} else {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "No coverage tooling found (consider Jest or Vitest coverage)"
	}

	return result, nil
}

// hasJestCoverage checks if Jest coverage is configured.
func (c *CoverageCheck) hasJestCoverage(path string) bool {
	// Check jest.config files for coverage config
	jestConfigs := []string{"jest.config.js", "jest.config.ts", "jest.config.mjs", "jest.config.cjs"}
	for _, cfg := range jestConfigs {
		if content, err := safepath.ReadFile(path, cfg); err == nil {
			contentLower := strings.ToLower(string(content))
			if strings.Contains(contentLower, "coverage") || strings.Contains(contentLower, "collectcoverage") {
				return true
			}
		}
	}

	// Check package.json jest config
	if pkg, err := ParsePackageJSON(path); err == nil {
		if script, ok := pkg.Scripts["test:coverage"]; ok && strings.Contains(script, "jest") {
			return true
		}
		if script, ok := pkg.Scripts["coverage"]; ok && strings.Contains(script, "jest") {
			return true
		}
	}

	return false
}

// hasVitestCoverage checks if Vitest coverage is configured.
func (c *CoverageCheck) hasVitestCoverage(path string) bool {
	// Check vitest.config files for coverage config
	vitestConfigs := []string{"vitest.config.ts", "vitest.config.js", "vitest.config.mts"}
	for _, cfg := range vitestConfigs {
		if content, err := safepath.ReadFile(path, cfg); err == nil {
			if strings.Contains(string(content), "coverage") {
				return true
			}
		}
	}

	// Check for @vitest/coverage-* packages
	if pkg, err := ParsePackageJSON(path); err == nil {
		for dep := range pkg.DevDependencies {
			if strings.HasPrefix(dep, "@vitest/coverage") {
				return true
			}
		}
	}

	return false
}

// hasC8 checks if c8 is configured.
func (c *CoverageCheck) hasC8(path string) bool {
	if pkg, err := ParsePackageJSON(path); err == nil {
		if _, ok := pkg.DevDependencies["c8"]; ok {
			return true
		}
	}

	if safepath.Exists(path, ".c8rc.json") || safepath.Exists(path, ".c8rc") {
		return true
	}

	return false
}

// hasNYC checks if nyc/Istanbul is configured.
func (c *CoverageCheck) hasNYC(path string) bool {
	if pkg, err := ParsePackageJSON(path); err == nil {
		if _, ok := pkg.DevDependencies["nyc"]; ok {
			return true
		}
	}

	if safepath.Exists(path, ".nycrc") || safepath.Exists(path, ".nycrc.json") || safepath.Exists(path, "nyc.config.js") {
		return true
	}

	return false
}

// checkCICoverage checks for coverage reporting in CI configs.
func (c *CoverageCheck) checkCICoverage(path string) string {
	ciFiles := []string{
		".github/workflows/ci.yml", ".github/workflows/ci.yaml",
		".github/workflows/test.yml", ".github/workflows/test.yaml",
		".gitlab-ci.yml",
	}

	for _, ciFile := range ciFiles {
		if content, err := safepath.ReadFile(path, ciFile); err == nil {
			contentLower := strings.ToLower(string(content))
			if strings.Contains(contentLower, "codecov") {
				return "Codecov (CI)"
			}
			if strings.Contains(contentLower, "coveralls") {
				return "Coveralls (CI)"
			}
		}
	}

	return ""
}

// findCoverageReports looks for coverage report files and extracts percentage.
func (c *CoverageCheck) findCoverageReports(path string) float64 {
	// Check common coverage report locations
	reportPaths := []string{
		"coverage/lcov.info",
		"coverage/coverage-final.json",
		"coverage/cobertura-coverage.xml",
		"coverage/clover.xml",
		"lcov.info",
	}

	for _, reportPath := range reportPaths {
		if content, err := safepath.ReadFile(path, reportPath); err == nil {
			if cov := extractCoverage(string(content)); cov >= 0 {
				return cov
			}
		}
	}

	return -1
}

// extractCoverage tries to extract coverage percentage from various formats.
func extractCoverage(content string) float64 {
	// Cobertura XML format: line-rate="0.85"
	coberturaRe := regexp.MustCompile(`line-rate="([0-9.]+)"`)
	if matches := coberturaRe.FindStringSubmatch(content); len(matches) > 1 {
		if rate, err := strconv.ParseFloat(matches[1], 64); err == nil {
			return rate * 100
		}
	}

	// LCOV summary format: LF:100 LH:85 (lines found, lines hit)
	lfRe := regexp.MustCompile(`LF:(\d+)`)
	lhRe := regexp.MustCompile(`LH:(\d+)`)
	lfMatches := lfRe.FindAllStringSubmatch(content, -1)
	lhMatches := lhRe.FindAllStringSubmatch(content, -1)

	if len(lfMatches) > 0 && len(lhMatches) > 0 {
		var totalLF, totalLH float64
		for _, m := range lfMatches {
			if len(m) > 1 {
				lf, _ := strconv.ParseFloat(m[1], 64)
				totalLF += lf
			}
		}
		for _, m := range lhMatches {
			if len(m) > 1 {
				lh, _ := strconv.ParseFloat(m[1], 64)
				totalLH += lh
			}
		}
		if totalLF > 0 {
			return (totalLH / totalLF) * 100
		}
	}

	// JSON coverage format (Jest): check for coverage summary
	if strings.Contains(content, "\"pct\"") {
		re := regexp.MustCompile(`"lines":\s*\{[^}]*"pct":\s*([0-9.]+)`)
		if matches := re.FindStringSubmatch(content); len(matches) > 1 {
			if pct, err := strconv.ParseFloat(matches[1], 64); err == nil {
				return pct
			}
		}
	}

	return -1
}
