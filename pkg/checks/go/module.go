package gocheck

import (
	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
	"golang.org/x/mod/modfile"
)

// ModuleCheck verifies that go.mod exists and has a valid Go version.
type ModuleCheck struct{}

func (c *ModuleCheck) ID() string   { return "go:module" }
func (c *ModuleCheck) Name() string { return "Go Module" }

func (c *ModuleCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangGo)

	// Use safepath to prevent directory traversal attacks
	data, err := safepath.ReadFile(path, "go.mod")
	if err != nil {
		if !safepath.Exists(path, "go.mod") {
			return rb.Fail("go.mod not found. Run 'go mod init' to create one."), nil
		}
		return checker.Result{}, err
	}

	// Get safe path for modfile.Parse (it needs the filename for error messages)
	modPath, _ := safepath.SafeJoin(path, "go.mod")

	// Parse go.mod to validate it
	modFile, err := modfile.Parse(modPath, data, nil)
	if err != nil {
		return rb.Fail("go.mod is invalid: " + err.Error()), nil
	}

	// Check for Go version
	if modFile.Go == nil || modFile.Go.Version == "" {
		return rb.WarnWithOutput("go.mod does not specify a Go version.", string(data)), nil
	}

	return rb.Pass("Module: " + modFile.Module.Mod.Path + " (Go " + modFile.Go.Version + ")"), nil
}
