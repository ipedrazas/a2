package gocheck

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
)

// CoverageCheck verifies that test coverage meets a threshold.
type CoverageCheck struct {
	Threshold float64 // Minimum coverage percentage (default 80.0)
}

func (c *CoverageCheck) ID() string   { return "go:coverage" }
func (c *CoverageCheck) Name() string { return "Go Coverage" }

func (c *CoverageCheck) Run(path string) (checker.Result, error) {
	threshold := c.Threshold
	if threshold == 0 {
		threshold = 80.0 // Default threshold
	}

	cmd := exec.Command("go", "test", "-cover", "./...")
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String()

	// Check for no test files
	if strings.Contains(output, "no test files") && !strings.Contains(output, "coverage:") {
		return checker.Result{
			Name:     c.Name(),
			ID:       c.ID(),
			Passed:   false,
			Status:   checker.Warn,
			Message:  "No test files found - coverage is 0%",
			Language: checker.LangGo,
		}, nil
	}

	// If tests failed, report that
	if err != nil {
		return checker.Result{
			Name:     c.Name(),
			ID:       c.ID(),
			Passed:   false,
			Status:   checker.Warn,
			Message:  "Could not measure coverage: tests failed",
			Language: checker.LangGo,
		}, nil
	}

	// Parse coverage from output
	coverage := parseCoverage(output)

	if coverage < threshold {
		return checker.Result{
			Name:     c.Name(),
			ID:       c.ID(),
			Passed:   false,
			Status:   checker.Warn,
			Message:  fmt.Sprintf("Coverage %.1f%% is below threshold %.1f%%", coverage, threshold),
			Language: checker.LangGo,
		}, nil
	}

	return checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Passed:   true,
		Status:   checker.Pass,
		Message:  fmt.Sprintf("Coverage: %.1f%%", coverage),
		Language: checker.LangGo,
	}, nil
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
