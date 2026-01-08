package rustcheck

import (
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// ProjectCheck verifies that a Rust project exists (Cargo.toml).
type ProjectCheck struct{}

func (c *ProjectCheck) ID() string   { return "rust:project" }
func (c *ProjectCheck) Name() string { return "Rust Project" }

// Run checks for Cargo.toml and extracts project information.
func (c *ProjectCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangRust,
	}

	// Check for Cargo.toml
	if !safepath.Exists(path, "Cargo.toml") {
		result.Passed = false
		result.Status = checker.Fail
		result.Message = "No Cargo.toml found"
		return result, nil
	}

	// Read Cargo.toml to extract package info
	content, err := safepath.ReadFile(path, "Cargo.toml")
	if err != nil {
		result.Passed = false
		result.Status = checker.Fail
		result.Message = "Cannot read Cargo.toml: " + err.Error()
		return result, nil
	}

	// Parse basic info from Cargo.toml
	name := extractTomlValue(string(content), "name")
	version := extractTomlValue(string(content), "version")

	var info []string
	if name != "" {
		info = append(info, name)
	}
	if version != "" {
		info = append(info, "v"+version)
	}

	// Check for workspace
	isWorkspace := strings.Contains(string(content), "[workspace]")
	if isWorkspace {
		info = append(info, "(workspace)")
	}

	result.Passed = true
	result.Status = checker.Pass
	if len(info) > 0 {
		result.Message = "Package: " + strings.Join(info, " ")
	} else {
		result.Message = "Cargo.toml found"
	}

	return result, nil
}

// extractTomlValue extracts a simple string value from TOML content.
func extractTomlValue(content, key string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, key+" =") || strings.HasPrefix(line, key+"=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				value := strings.TrimSpace(parts[1])
				value = strings.Trim(value, "\"'")
				return value
			}
		}
	}
	return ""
}
