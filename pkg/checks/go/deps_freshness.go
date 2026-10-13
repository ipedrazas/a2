package gocheck

import (
	"context"
	"fmt"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
)

// DepsFreshnessCheck detects outdated direct Go module dependencies.
type DepsFreshnessCheck struct{}

// directUpdatesFormat makes `go list -m -u` print only direct requirements
// that have an update, as "path version [update]". Indirect modules are
// skipped: many are only in the module graph (dependencies of dependencies
// that are never compiled), and `go get -u` cannot update them.
const directUpdatesFormat = `{{if and .Update (not .Indirect) (not .Main)}}{{.Path}} {{.Version}} [{{.Update.Version}}]{{end}}`

func (c *DepsFreshnessCheck) ID() string   { return "go:deps_freshness" }
func (c *DepsFreshnessCheck) Name() string { return "Go Dependency Freshness" }

func (c *DepsFreshnessCheck) Run(ctx context.Context, path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangGo)

	result := rb.RunCommand(ctx, path, "go", "list", "-m", "-u", "-f", directUpdatesFormat, "all")
	output := result.CombinedOutput()

	if !result.Success() {
		return rb.WarnWithOutput("go list -m -u failed: "+checkutil.TruncateMessage(result.Output(), 200), output), nil
	}

	outdated := countOutdatedGoModules(result.Stdout)

	if outdated == 0 {
		return rb.PassWithOutput("All direct dependencies are up to date", output), nil
	}

	msg := fmt.Sprintf("%d outdated direct %s",
		outdated,
		checkutil.Pluralize(outdated, "dependency", "dependencies"),
	)

	if outdated > 20 {
		return rb.WarnWithOutput(msg+". Run: go get -u -t ./...", output), nil
	}

	return rb.PassWithOutput(msg+". Run: go get -u -t ./...", output), nil
}

// countOutdatedGoModules counts modules that have updates available.
// go list -m -u (with directUpdatesFormat) outputs lines like:
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
