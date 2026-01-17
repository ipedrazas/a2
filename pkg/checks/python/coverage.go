package pythoncheck

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/config"
)

// CoverageCheck verifies that test coverage meets a threshold.
type CoverageCheck struct {
	Config *config.PythonLanguageConfig
}

func (c *CoverageCheck) ID() string   { return "python:coverage" }
func (c *CoverageCheck) Name() string { return "Python Coverage" }

func (c *CoverageCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangPython)

	threshold := 80.0
	if c.Config != nil && c.Config.CoverageThreshold > 0 {
		threshold = c.Config.CoverageThreshold
	}

	result := checkutil.RunCommand(path, "pytest", "--cov=.", "--cov-report=term-missing", "-q")
	output := result.CombinedOutput()

	// Check if pytest or pytest-cov is not installed
	if strings.Contains(result.Stderr, "unrecognized arguments: --cov") ||
		strings.Contains(result.Stderr, "No module named pytest") {
		return rb.ToolNotInstalled("pytest-cov", "pip install pytest-cov"), nil
	}

	if checkutil.ToolNotFoundError(result.Err) {
		return rb.ToolNotInstalled("pytest", "pip install pytest"), nil
	}

	// Check for no tests
	if strings.Contains(result.Stdout, "no tests ran") || strings.Contains(result.Stdout, "collected 0 items") {
		return rb.WarnWithOutput("No tests found - coverage is 0%", output), nil
	}

	// Parse coverage from output
	coverage := parsePythonCoverage(result.Stdout)

	if coverage < threshold {
		return rb.WarnWithOutput(fmt.Sprintf("Coverage %.1f%% is below threshold %.1f%%", coverage, threshold), output), nil
	}

	return rb.PassWithOutput(fmt.Sprintf("Coverage: %.1f%%", coverage), output), nil
}

// parsePythonCoverage extracts the total coverage percentage from pytest-cov output.
func parsePythonCoverage(output string) float64 {
	// Look for "TOTAL ... XX%" pattern
	re := regexp.MustCompile(`TOTAL\s+\d+\s+\d+\s+(\d+)%`)
	matches := re.FindStringSubmatch(output)
	if len(matches) >= 2 {
		if cov, err := strconv.ParseFloat(matches[1], 64); err == nil {
			return cov
		}
	}

	// Alternative pattern: "Coverage: XX%"
	re2 := regexp.MustCompile(`(?:coverage|Coverage):\s*(\d+(?:\.\d+)?)%`)
	matches2 := re2.FindStringSubmatch(output)
	if len(matches2) >= 2 {
		if cov, err := strconv.ParseFloat(matches2[1], 64); err == nil {
			return cov
		}
	}

	return 0
}
