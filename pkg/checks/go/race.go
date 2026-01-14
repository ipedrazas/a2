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
		return rb.WarnWithOutput("Race condition detected", output), nil
	}

	// Handle no test files
	if strings.Contains(output, "no test files") {
		return rb.PassWithOutput("No test files to check for races", output), nil
	}

	// Handle test failures (separate from race detection)
	if !result.Success() {
		return rb.WarnWithOutput("Tests failed during race detection", output), nil
	}

	return rb.PassWithOutput("No race conditions detected", output), nil
}
