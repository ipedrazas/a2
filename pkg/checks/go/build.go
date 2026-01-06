package gocheck

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
)

// BuildCheck verifies that the project compiles successfully.
type BuildCheck struct{}

func (c *BuildCheck) ID() string   { return "go:build" }
func (c *BuildCheck) Name() string { return "Go Build" }

func (c *BuildCheck) Run(path string) (checker.Result, error) {
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		output := strings.TrimSpace(stderr.String())
		if output == "" {
			output = strings.TrimSpace(stdout.String())
		}

		return checker.Result{
			Name:     c.Name(),
			ID:       c.ID(),
			Passed:   false,
			Status:   checker.Fail, // Critical - stops execution
			Message:  "Build failed: " + output,
			Language: checker.LangGo,
		}, nil
	}

	return checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Passed:   true,
		Status:   checker.Pass,
		Message:  "Build successful",
		Language: checker.LangGo,
	}, nil
}
