// Package targets provides maturity target definitions for different project stages.
package targets

// Target defines a set of checks to disable for a specific maturity level.
type Target struct {
	Name        string
	Description string
	Disabled    []string // Check IDs to disable
}

// BuiltInTargets contains the predefined maturity targets.
var BuiltInTargets = map[string]Target{
	"poc": {
		Name:        "poc",
		Description: "Proof of Concept - minimal checks for early development",
		Disabled: []string{
			// Common checks to disable for PoC
			"common:license", "common:sast", "common:changelog",
			"common:precommit", "common:env", "common:contributing",
			"common:secrets", "common:config_validation", "common:retry",
			"common:editorconfig",
			// Go optional checks
			"go:coverage", "go:deps", "go:cyclomatic", "go:logging", "go:race",
			// Python optional checks
			"python:coverage", "python:deps", "python:complexity", "python:logging",
			// Node.js optional checks
			"node:coverage", "node:deps", "node:logging",
			// Java optional checks
			"java:coverage", "java:deps", "java:logging",
			// Rust optional checks
			"rust:coverage", "rust:deps", "rust:logging",
			// TypeScript optional checks
			"typescript:coverage", "typescript:deps", "typescript:logging",
		},
	},
	"production": {
		Name:        "production",
		Description: "Production application - all checks enabled",
		Disabled:    []string{}, // Enable everything
	},
}

// Get returns a target by name and whether it was found.
func Get(name string) (Target, bool) {
	t, ok := BuiltInTargets[name]
	return t, ok
}

// List returns all available targets in a consistent order.
func List() []Target {
	return []Target{
		BuiltInTargets["poc"],
		BuiltInTargets["production"],
	}
}

// Names returns the names of all available targets.
func Names() []string {
	return []string{"poc", "production"}
}
