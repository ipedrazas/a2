package gocheck

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
)

// CoverageCheck verifies that test coverage meets a threshold.
type CoverageCheck struct {
	Threshold float64 // Minimum coverage percentage (default 80.0)
}

func (c *CoverageCheck) ID() string   { return "go:coverage" }
func (c *CoverageCheck) Name() string { return "Go Coverage" }

func (c *CoverageCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangGo)

	threshold := c.Threshold
	if threshold == 0 {
		threshold = 80.0 // Default threshold
	}

	result := checkutil.RunCommand(path, "go", "test", "-cover", "./...")

	// Check for no test files
	if strings.Contains(result.Stdout, "no test files") && !strings.Contains(result.Stdout, "coverage:") {
		return rb.Warn("No test files found - coverage is 0%"), nil
	}

	// If tests failed, report that
	if !result.Success() {
		return rb.Warn("Could not measure coverage: tests failed"), nil
	}

	// Parse coverage from output
	coverage := parseCoverage(result.Stdout)

	if coverage < threshold {
		return rb.Warn(fmt.Sprintf("Coverage %.1f%% is below threshold %.1f%%", coverage, threshold)), nil
	}

	return rb.Pass(fmt.Sprintf("Coverage: %.1f%%", coverage)), nil
}

// parseCoverage extracts the average coverage percentage from go test -cover output.
func parseCoverage(output string) float64 {
	// Match patterns like "coverage: 75.0% of statements"
	re := regexp.MustCompile(`coverage:\s*([\d.]+)%`)
	matches := re.FindAllStringSubmatch(output, -1)

	if len(matches) == 0 {
		return 0
	}

	// Calculate average coverage across all packages
	var total float64
	var count int
	for _, match := range matches {
		if len(match) >= 2 {
			if cov, err := strconv.ParseFloat(match[1], 64); err == nil {
				total += cov
				count++
			}
		}
	}

	if count == 0 {
		return 0
	}

	return total / float64(count)
}

// DefaultCoverageCheck returns a CoverageCheck with the default 80% threshold.
func DefaultCoverageCheck() *CoverageCheck {
	return &CoverageCheck{Threshold: 80.0}
}
