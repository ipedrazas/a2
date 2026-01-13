package gocheck

import (
	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
)

// BuildCheck verifies that the project compiles successfully.
type BuildCheck struct{}

func (c *BuildCheck) ID() string   { return "go:build" }
func (c *BuildCheck) Name() string { return "Go Build" }

func (c *BuildCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangGo)

	result := checkutil.RunCommand(path, "go", "build", "./...")
	if !result.Success() {
		return rb.Fail("Build failed: " + checkutil.TruncateMessage(result.Output(), 200)), nil
	}

	return rb.Pass("Build successful"), nil
}
