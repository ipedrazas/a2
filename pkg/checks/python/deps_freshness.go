package pythoncheck

import (
	"fmt"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// DepsFreshnessCheck detects outdated Python dependencies.
type DepsFreshnessCheck struct{}

func (c *DepsFreshnessCheck) ID() string   { return "python:deps_freshness" }
func (c *DepsFreshnessCheck) Name() string { return "Python Dependency Freshness" }

func (c *DepsFreshnessCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangPython)

	hasPython := safepath.Exists(path, "pyproject.toml") ||
		safepath.Exists(path, "setup.py") ||
		safepath.Exists(path, "requirements.txt")
	if !hasPython {
		return rb.Fail("Python project not found"), nil
	}

	// pip list --outdated shows packages with newer versions
	if !pythonToolAvailable(path, "pip") {
		return rb.ToolNotInstalled("pip", ""), nil
	}

	result := runPythonCommand(path, "pip", "list", "--outdated", "--format=columns")
	output := result.CombinedOutput()

	if !result.Success() {
		// Some pip versions exit non-zero for various reasons; try parsing anyway
		if result.Stdout == "" {
			return rb.WarnWithOutput("pip list --outdated failed: "+checkutil.TruncateMessage(result.Output(), 200), output), nil
		}
	}

	outdated := countOutdatedPipPackages(result.Stdout)

	if outdated == 0 {
		return rb.PassWithOutput("All dependencies are up to date", output), nil
	}

	msg := fmt.Sprintf("%d outdated %s",
		outdated,
		checkutil.Pluralize(outdated, "package", "packages"),
	)

	if outdated > 20 {
		return rb.WarnWithOutput(msg+". Run: pip list --outdated", output), nil
	}

	return rb.PassWithOutput(msg+". Run: pip list --outdated", output), nil
}

// countOutdatedPipPackages counts outdated packages from pip list --outdated output.
// Output format (columns):
// Package    Version  Latest   Type
// --------   -------  ------   ----
// requests   2.28.0   2.31.0   wheel
func countOutdatedPipPackages(output string) int {
	count := 0
	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		// Skip header lines
		if strings.HasPrefix(trimmed, "Package") || strings.HasPrefix(trimmed, "---") {
			continue
		}
		// Each non-header, non-empty line is an outdated package
		if len(strings.Fields(trimmed)) >= 3 {
			count++
		}
	}
	return count
}
