package nodecheck

import (
	"encoding/json"
	"fmt"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// PackageJSON represents the structure of a package.json file.
type PackageJSON struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Description     string            `json:"description"`
	Main            string            `json:"main"`
	Scripts         map[string]string `json:"scripts"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

// ProjectCheck verifies that package.json exists and is valid.
type ProjectCheck struct{}

// ID returns the unique identifier for this check.
func (c *ProjectCheck) ID() string {
	return "node:project"
}

// Name returns the human-readable name for this check.
func (c *ProjectCheck) Name() string {
	return "Node Project"
}

// Run executes the project check.
func (c *ProjectCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangNode)

	// Check if package.json exists
	if !safepath.Exists(path, "package.json") {
		return rb.Fail("package.json not found"), nil
	}

	// Read and parse package.json
	data, err := safepath.ReadFile(path, "package.json")
	if err != nil {
		return rb.Fail(fmt.Sprintf("Failed to read package.json: %v", err)), nil
	}

	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return rb.Fail(fmt.Sprintf("package.json is invalid JSON: %v", err)), nil
	}

	// Validate required fields
	if pkg.Name == "" {
		return rb.Fail("package.json is missing required 'name' field"), nil
	}

	if pkg.Version == "" {
		return rb.Warn(fmt.Sprintf("Package %s is missing 'version' field", pkg.Name)), nil
	}

	return rb.Pass(fmt.Sprintf("Package: %s v%s", pkg.Name, pkg.Version)), nil
}

// ParsePackageJSON reads and parses package.json from the given path.
func ParsePackageJSON(path string) (*PackageJSON, error) {
	data, err := safepath.ReadFile(path, "package.json")
	if err != nil {
		return nil, err
	}

	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}

	return &pkg, nil
}
