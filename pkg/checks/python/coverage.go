package pythoncheck

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
)

// CoverageCheck verifies that test coverage meets a threshold.
type CoverageCheck struct {
	Config *config.PythonLanguageConfig
}

func (c *CoverageCheck) ID() string   { return "python:coverage" }
func (c *CoverageCheck) Name() string { return "Python Coverage" }

func (c *CoverageCheck) Run(path string) (checker.Result, error) {
	threshold := 80.0
	if c.Config != nil && c.Config.CoverageThreshold > 0 {
		threshold = c.Config.CoverageThreshold
	}

	// Check if pytest-cov is available by trying to run pytest with coverage
	cmd := exec.Command("pytest", "--cov=.", "--cov-report=term-missing", "-q")
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String()
	errOutput := stderr.String()

	// Check if pytest or pytest-cov is not installed
	if strings.Contains(errOutput, "unrecognized arguments: --cov") ||
		strings.Contains(errOutput, "No module named pytest") {
		return checker.Result{
			Name:     c.Name(),
			ID:       c.ID(),
			Passed:   true,
			Status:   checker.Pass,
			Message:  "pytest-cov not installed (run: pip install pytest-cov)",
			Language: checker.LangPython,
		}, nil
	}

	if err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			return checker.Result{
				Name:     c.Name(),
				ID:       c.ID(),
				Passed:   true,
				Status:   checker.Pass,
				Message:  "pytest not installed, skipping coverage",
				Language: checker.LangPython,
			}, nil
		}
	}

	// Check for no tests
	if strings.Contains(output, "no tests ran") || strings.Contains(output, "collected 0 items") {
		return checker.Result{
			Name:     c.Name(),
			ID:       c.ID(),
			Passed:   false,
			Status:   checker.Warn,
			Message:  "No tests found - coverage is 0%",
			Language: checker.LangPython,
		}, nil
	}

	// Parse coverage from output
	coverage := parsePythonCoverage(output)

	if coverage < threshold {
		return checker.Result{
			Name:     c.Name(),
			ID:       c.ID(),
			Passed:   false,
			Status:   checker.Warn,
			Message:  fmt.Sprintf("Coverage %.1f%% is below threshold %.1f%%", coverage, threshold),
			Language: checker.LangPython,
		}, nil
	}

	return checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Passed:   true,
		Status:   checker.Pass,
		Message:  fmt.Sprintf("Coverage: %.1f%%", coverage),
		Language: checker.LangPython,
	}, nil
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
