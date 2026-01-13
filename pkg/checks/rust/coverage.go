package rustcheck

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// CoverageCheck checks for test coverage configuration and reports.
type CoverageCheck struct {
	Threshold float64
}

func (c *CoverageCheck) ID() string   { return "rust:coverage" }
func (c *CoverageCheck) Name() string { return "Rust Coverage" }

// Run checks for coverage tooling and reports.
func (c *CoverageCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangRust)

	// Check for Cargo.toml first
	if !safepath.Exists(path, "Cargo.toml") {
		return rb.Fail("No Cargo.toml found"), nil
	}

	// Look for coverage tools/config
	var coverageTools []string

	// Check for tarpaulin (cargo-tarpaulin)
	if safepath.Exists(path, "tarpaulin.toml") || safepath.Exists(path, ".tarpaulin.toml") {
		coverageTools = append(coverageTools, "tarpaulin")
	}

	// Check for llvm-cov config in Cargo.toml
	if content, err := safepath.ReadFile(path, "Cargo.toml"); err == nil {
		if strings.Contains(string(content), "cargo-llvm-cov") ||
			strings.Contains(string(content), "llvm-cov") {
			coverageTools = append(coverageTools, "llvm-cov")
		}
	}

	// Check for coverage in CI configs
	ciFiles := []string{".github/workflows/ci.yml", ".github/workflows/ci.yaml",
		".github/workflows/rust.yml", ".github/workflows/rust.yaml",
		".gitlab-ci.yml"}
	for _, ciFile := range ciFiles {
		if content, err := safepath.ReadFile(path, ciFile); err == nil {
			contentLower := strings.ToLower(string(content))
			if strings.Contains(contentLower, "tarpaulin") {
				if !contains(coverageTools, "tarpaulin") {
					coverageTools = append(coverageTools, "tarpaulin (CI)")
				}
			}
			if strings.Contains(contentLower, "llvm-cov") || strings.Contains(contentLower, "cargo-llvm-cov") {
				if !contains(coverageTools, "llvm-cov") {
					coverageTools = append(coverageTools, "llvm-cov (CI)")
				}
			}
			if strings.Contains(contentLower, "codecov") || strings.Contains(contentLower, "coveralls") {
				coverageTools = append(coverageTools, "coverage reporting")
			}
		}
	}

	// Try to find coverage reports
	coverage := c.findCoverageReports(path)

	// Build result
	if coverage >= 0 {
		msg := fmt.Sprintf("Coverage %.1f%% (threshold %.1f%%)", coverage, c.Threshold)
		if coverage >= c.Threshold {
			return rb.Pass(msg), nil
		}
		return rb.Warn(msg), nil
	}
	if len(coverageTools) > 0 {
		return rb.Pass("Coverage configured: " + strings.Join(coverageTools, ", ")), nil
	}
	return rb.Warn("No coverage tooling found (consider cargo-tarpaulin or cargo-llvm-cov)"), nil
}

// findCoverageReports looks for coverage report files and extracts percentage.
func (c *CoverageCheck) findCoverageReports(path string) float64 {
	// Check common coverage report locations
	reportPaths := []string{
		"target/tarpaulin/cobertura.xml",
		"target/llvm-cov/html/index.html",
		"coverage/cobertura.xml",
		"coverage.xml",
		"lcov.info",
	}

	for _, reportPath := range reportPaths {
		if content, err := safepath.ReadFile(path, reportPath); err == nil {
			// Try to extract coverage percentage
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
	lfMatches := lfRe.FindStringSubmatch(content)
	lhMatches := lhRe.FindStringSubmatch(content)
	if len(lfMatches) > 1 && len(lhMatches) > 1 {
		lf, _ := strconv.ParseFloat(lfMatches[1], 64)
		lh, _ := strconv.ParseFloat(lhMatches[1], 64)
		if lf > 0 {
			return (lh / lf) * 100
		}
	}

	return -1
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
