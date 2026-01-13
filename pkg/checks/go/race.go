package gocheck

import (
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
)

// RaceCheck runs tests with the race detector enabled to find data races.
type RaceCheck struct{}

func (c *RaceCheck) ID() string   { return "go:race" }
func (c *RaceCheck) Name() string { return "Go Race Detection" }

func (c *RaceCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangGo)

	result := checkutil.RunCommand(path, "go", "test", "-race", "-short", "./...")
	output := result.CombinedOutput()

	// Check for race conditions in output
	if strings.Contains(output, "WARNING: DATA RACE") {
		return rb.Warn("Race condition detected"), nil
	}

	// Handle no test files
	if strings.Contains(output, "no test files") {
		return rb.Pass("No test files to check for races"), nil
	}

	// Handle test failures (separate from race detection)
	if !result.Success() {
		return rb.Warn("Tests failed during race detection"), nil
	}

	return rb.Pass("No race conditions detected"), nil
}
