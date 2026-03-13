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

// Run runs coverage tools and checks that coverage meets the threshold.
func (c *CoverageCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangRust)

	threshold := c.Threshold
	if threshold == 0 {
		threshold = 80.0
	}

	// Check for Cargo.toml first
	if !safepath.Exists(path, "Cargo.toml") {
		return rb.Fail("No Cargo.toml found"), nil
	}

	// Try to run cargo-tarpaulin
	if checkutil.ToolAvailable("cargo-tarpaulin") {
		result := checkutil.RunCommand(path, "cargo", "tarpaulin", "--skip-clean", "--out", "stdout")
		output := result.CombinedOutput()

		if !result.Success() {
			return rb.WarnWithOutput("Could not measure coverage: cargo tarpaulin failed", output), nil
		}

		coverage := parseTarpaulinCoverage(result.Stdout)
		if coverage < threshold {
			return rb.WarnWithOutput(fmt.Sprintf("Coverage %.1f%% is below threshold %.1f%%", coverage, threshold), output), nil
		}
		return rb.PassWithOutput(fmt.Sprintf("Coverage: %.1f%%", coverage), output), nil
	}

	// Try to run cargo-llvm-cov
	if checkutil.ToolAvailable("cargo-llvm-cov") {
		result := checkutil.RunCommand(path, "cargo", "llvm-cov", "--summary-only")
		output := result.CombinedOutput()

		if !result.Success() {
			return rb.WarnWithOutput("Could not measure coverage: cargo llvm-cov failed", output), nil
		}

		coverage := parseLlvmCovCoverage(result.Stdout)
		if coverage < threshold {
			return rb.WarnWithOutput(fmt.Sprintf("Coverage %.1f%% is below threshold %.1f%%", coverage, threshold), output), nil
		}
		return rb.PassWithOutput(fmt.Sprintf("Coverage: %.1f%%", coverage), output), nil
	}

	return rb.Warn("No coverage tool available (install cargo-tarpaulin or cargo-llvm-cov)"), nil
}

// parseTarpaulinCoverage extracts coverage percentage from cargo tarpaulin output.
// Output typically ends with a line like: "85.00% coverage, 170/200 lines covered"
func parseTarpaulinCoverage(output string) float64 {
	re := regexp.MustCompile(`([\d.]+)%\s+coverage`)
	matches := re.FindAllStringSubmatch(output, -1)
	if len(matches) == 0 {
		return 0
	}
	// Take the last match (the summary line)
	last := matches[len(matches)-1]
	if cov, err := strconv.ParseFloat(last[1], 64); err == nil {
		return cov
	}
	return 0
}

// parseLlvmCovCoverage extracts coverage percentage from cargo llvm-cov output.
// Output contains lines like: "TOTAL  1000  850  85.00%"
func parseLlvmCovCoverage(output string) float64 {
	// Match TOTAL line with percentage
	re := regexp.MustCompile(`TOTAL\s+\d+\s+\d+\s+([\d.]+)%`)
	if matches := re.FindStringSubmatch(output); len(matches) > 1 {
		if cov, err := strconv.ParseFloat(matches[1], 64); err == nil {
			return cov
		}
	}
	// Fallback: match any percentage in the last lines
	rePercent := regexp.MustCompile(`([\d.]+)%`)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for i := len(lines) - 1; i >= 0 && i >= len(lines)-5; i-- {
		if matches := rePercent.FindStringSubmatch(lines[i]); len(matches) > 1 {
			if cov, err := strconv.ParseFloat(matches[1], 64); err == nil {
				return cov
			}
		}
	}
	return 0
}
