package gocheck

import (
	"context"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
)

// RaceCheck runs tests with the race detector enabled to find data races.
type RaceCheck struct{}

func (c *RaceCheck) ID() string   { return "go:race" }
func (c *RaceCheck) Name() string { return "Go Race Detection" }

func (c *RaceCheck) Run(ctx context.Context, path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangGo)

	result := rb.RunCommand(ctx, path, "go", "test", "-race", "-short", "./...")
	output := result.CombinedOutput()

	// Check for race conditions in output
	if strings.Contains(output, "WARNING: DATA RACE") {
		return rb.WarnWithOutput("Race condition detected", output), nil
	}

	// Handle test failures (separate from race detection). This must come
	// before the no-test-files case so an untested package cannot mask it.
	if !result.Success() {
		return rb.WarnWithOutput("Tests failed during race detection", output), nil
	}

	if !hasTestedPackages(result.Stdout) {
		return rb.PassWithOutput("No test files to check for races", output), nil
	}

	return rb.PassWithOutput("No race conditions detected", output), nil
}
