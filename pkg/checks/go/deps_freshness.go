package gocheck

import (
	"fmt"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
)

// DepsFreshnessCheck detects outdated Go module dependencies.
type DepsFreshnessCheck struct{}

func (c *DepsFreshnessCheck) ID() string   { return "go:deps_freshness" }
func (c *DepsFreshnessCheck) Name() string { return "Go Dependency Freshness" }

func (c *DepsFreshnessCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangGo)

	// go list -m -u all shows modules with available updates
	result := checkutil.RunCommand(path, "go", "list", "-m", "-u", "all")
	output := result.CombinedOutput()

	if !result.Success() {
		return rb.WarnWithOutput("go list -m -u all failed: "+checkutil.TruncateMessage(result.Output(), 200), output), nil
	}

	outdated := countOutdatedGoModules(result.Stdout)

	if outdated == 0 {
		return rb.PassWithOutput("All dependencies are up to date", output), nil
	}

	msg := fmt.Sprintf("%d outdated %s",
		outdated,
		checkutil.Pluralize(outdated, "dependency", "dependencies"),
	)

	if outdated > 20 {
		return rb.WarnWithOutput(msg+". Run: go list -m -u all", output), nil
	}

	return rb.PassWithOutput(msg+". Run: go get -u ./...", output), nil
}

// countOutdatedGoModules counts modules that have updates available.
// go list -m -u all outputs lines like:
// github.com/foo/bar v1.0.0 [v1.1.0]
// The [v1.1.0] suffix indicates an update is available.
func countOutdatedGoModules(output string) int {
	count := 0
	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		// Lines with [vX.Y.Z] indicate available updates
		if strings.Contains(trimmed, "[") && strings.Contains(trimmed, "]") {
			count++
		}
	}
	return count
}
