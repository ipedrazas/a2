package gocheck

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
)

// RaceCheck runs tests with the race detector enabled to find data races.
type RaceCheck struct{}

func (c *RaceCheck) ID() string   { return "go:race" }
func (c *RaceCheck) Name() string { return "Go Race Detection" }

func (c *RaceCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangGo,
	}

	cmd := exec.Command("go", "test", "-race", "-short", "./...")
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()

	// Check for race conditions in output
	if strings.Contains(output, "WARNING: DATA RACE") {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "Race condition detected"
		return result, nil
	}

	// Handle no test files
	if strings.Contains(output, "no test files") {
		result.Passed = true
		result.Status = checker.Pass
		result.Message = "No test files to check for races"
		return result, nil
	}

	// Handle test failures (separate from race detection)
	if err != nil {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "Tests failed during race detection"
		return result, nil
	}

	result.Passed = true
	result.Status = checker.Pass
	result.Message = "No race conditions detected"
	return result, nil
}
