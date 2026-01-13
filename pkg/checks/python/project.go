package pythoncheck

import (
	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// ProjectCheck verifies that a Python project configuration exists.
type ProjectCheck struct{}

func (c *ProjectCheck) ID() string   { return "python:project" }
func (c *ProjectCheck) Name() string { return "Python Project" }

func (c *ProjectCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangPython)

	// Check for pyproject.toml (preferred)
	if safepath.Exists(path, "pyproject.toml") {
		return rb.Pass("Found pyproject.toml"), nil
	}

	// Check for setup.py (legacy)
	if safepath.Exists(path, "setup.py") {
		return rb.Warn("Found setup.py (consider migrating to pyproject.toml)"), nil
	}

	// Check for requirements.txt (minimal)
	if safepath.Exists(path, "requirements.txt") {
		return rb.Warn("Found requirements.txt only - consider adding pyproject.toml"), nil
	}

	return rb.Fail("No Python project configuration found (pyproject.toml, setup.py, or requirements.txt)"), nil
}
