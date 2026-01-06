package nodecheck

import (
	"encoding/json"
	"fmt"

	"github.com/ipedrazas/a2/pkg/checker"
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
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangNode,
	}

	// Check if package.json exists
	if !safepath.Exists(path, "package.json") {
		result.Status = checker.Fail
		result.Passed = false
		result.Message = "package.json not found"
		return result, nil
	}

	// Read and parse package.json
	data, err := safepath.ReadFile(path, "package.json")
	if err != nil {
		result.Status = checker.Fail
		result.Passed = false
		result.Message = fmt.Sprintf("Failed to read package.json: %v", err)
		return result, nil
	}

	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		result.Status = checker.Fail
		result.Passed = false
		result.Message = fmt.Sprintf("package.json is invalid JSON: %v", err)
		return result, nil
	}

	// Validate required fields
	if pkg.Name == "" {
		result.Status = checker.Fail
		result.Passed = false
		result.Message = "package.json is missing required 'name' field"
		return result, nil
	}

	if pkg.Version == "" {
		result.Status = checker.Warn
		result.Passed = false
		result.Message = fmt.Sprintf("Package %s is missing 'version' field", pkg.Name)
		return result, nil
	}

	result.Status = checker.Pass
	result.Passed = true
	result.Message = fmt.Sprintf("Package: %s v%s", pkg.Name, pkg.Version)
	return result, nil
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
