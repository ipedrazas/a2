package gocheck

import (
	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
	"golang.org/x/mod/modfile"
)

// ModuleCheck verifies that go.mod exists and has a valid Go version.
type ModuleCheck struct{}

func (c *ModuleCheck) ID() string   { return "go:module" }
func (c *ModuleCheck) Name() string { return "Go Module" }

func (c *ModuleCheck) Run(path string) (checker.Result, error) {
	// Use safepath to prevent directory traversal attacks
	data, err := safepath.ReadFile(path, "go.mod")
	if err != nil {
		if !safepath.Exists(path, "go.mod") {
			return checker.Result{
				Name:     c.Name(),
				ID:       c.ID(),
				Passed:   false,
				Status:   checker.Fail,
				Message:  "go.mod not found. Run 'go mod init' to create one.",
				Language: checker.LangGo,
			}, nil
		}
		return checker.Result{}, err
	}

	// Get safe path for modfile.Parse (it needs the filename for error messages)
	modPath, _ := safepath.SafeJoin(path, "go.mod")

	// Parse go.mod to validate it
	modFile, err := modfile.Parse(modPath, data, nil)
	if err != nil {
		return checker.Result{
			Name:     c.Name(),
			ID:       c.ID(),
			Passed:   false,
			Status:   checker.Fail,
			Message:  "go.mod is invalid: " + err.Error(),
			Language: checker.LangGo,
		}, nil
	}

	// Check for Go version
	if modFile.Go == nil || modFile.Go.Version == "" {
		return checker.Result{
			Name:     c.Name(),
			ID:       c.ID(),
			Passed:   false,
			Status:   checker.Warn,
			Message:  "go.mod does not specify a Go version.",
			Language: checker.LangGo,
		}, nil
	}

	return checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Passed:   true,
		Status:   checker.Pass,
		Message:  "Module: " + modFile.Module.Mod.Path + " (Go " + modFile.Go.Version + ")",
		Language: checker.LangGo,
	}, nil
}
