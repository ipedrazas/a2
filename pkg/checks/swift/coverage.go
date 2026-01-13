package swiftcheck

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// CoverageCheck verifies test coverage meets the threshold.
type CoverageCheck struct {
	Threshold float64
}

func (c *CoverageCheck) ID() string   { return "swift:coverage" }
func (c *CoverageCheck) Name() string { return "Swift Coverage" }

// Run checks test coverage using swift test with coverage enabled.
func (c *CoverageCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangSwift)

	// Check for Package.swift first
	if !safepath.Exists(path, "Package.swift") {
		return rb.Fail("No Package.swift found"), nil
	}

	// Run swift test with coverage
	cmd := exec.Command("swift", "test", "--enable-code-coverage")
	cmd.Dir = path
	_, err := cmd.CombinedOutput()

	if err != nil {
		return rb.Warn("Cannot generate coverage: " + err.Error()), nil
	}

	// Get coverage data using llvm-cov
	// First, find the test binary and profdata
	coverage, err := c.extractCoverage(path)
	if err != nil {
		return rb.Info("Coverage data not available (llvm-cov not configured)"), nil
	}

	threshold := c.Threshold
	if threshold == 0 {
		threshold = 80.0
	}

	if coverage >= threshold {
		return rb.Pass(fmt.Sprintf("Coverage: %.1f%% (threshold: %.1f%%)", coverage, threshold)), nil
	}
	return rb.Warn(fmt.Sprintf("Coverage: %.1f%% (below threshold: %.1f%%)", coverage, threshold)), nil
}

// extractCoverage attempts to get coverage percentage from swift test output.
func (c *CoverageCheck) extractCoverage(path string) (float64, error) {
	// Use swift package show-codecov-path to find the coverage file
	cmd := exec.Command("swift", "package", "show-codecov-path")
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	// The codecov path points to a JSON file with coverage data
	codecovPath := string(output)
	if codecovPath == "" {
		return 0, fmt.Errorf("no codecov path")
	}

	// Read the codecov JSON
	data, err := safepath.ReadFile(path, codecovPath)
	if err != nil {
		return 0, err
	}

	// Parse the JSON to extract coverage percentage
	var codecov struct {
		Data []struct {
			Totals struct {
				Lines struct {
					Percent float64 `json:"percent"`
				} `json:"lines"`
			} `json:"totals"`
		} `json:"data"`
	}

	if err := json.Unmarshal(data, &codecov); err != nil {
		return 0, err
	}

	if len(codecov.Data) > 0 {
		return codecov.Data[0].Totals.Lines.Percent, nil
	}

	return 0, fmt.Errorf("no coverage data")
}
