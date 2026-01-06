package pythoncheck

import (
	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// ProjectCheck verifies that a Python project configuration exists.
type ProjectCheck struct{}

func (c *ProjectCheck) ID() string   { return "python:project" }
func (c *ProjectCheck) Name() string { return "Python Project" }

func (c *ProjectCheck) Run(path string) (checker.Result, error) {
	// Check for pyproject.toml (preferred)
	if safepath.Exists(path, "pyproject.toml") {
		return checker.Result{
			Name:     c.Name(),
			ID:       c.ID(),
			Passed:   true,
			Status:   checker.Pass,
			Message:  "Found pyproject.toml",
			Language: checker.LangPython,
		}, nil
	}

	// Check for setup.py (legacy)
	if safepath.Exists(path, "setup.py") {
		return checker.Result{
			Name:     c.Name(),
			ID:       c.ID(),
			Passed:   true,
			Status:   checker.Warn,
			Message:  "Found setup.py (consider migrating to pyproject.toml)",
			Language: checker.LangPython,
		}, nil
	}

	// Check for requirements.txt (minimal)
	if safepath.Exists(path, "requirements.txt") {
		return checker.Result{
			Name:     c.Name(),
			ID:       c.ID(),
			Passed:   true,
			Status:   checker.Warn,
			Message:  "Found requirements.txt only - consider adding pyproject.toml",
			Language: checker.LangPython,
		}, nil
	}

	return checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Passed:   false,
		Status:   checker.Fail,
		Message:  "No Python project configuration found (pyproject.toml, setup.py, or requirements.txt)",
		Language: checker.LangPython,
	}, nil
}
