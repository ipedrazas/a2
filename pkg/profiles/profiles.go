// Package profiles provides built-in check profiles for different project types.
package profiles

// Profile defines a set of checks to disable for a specific use case.
type Profile struct {
	Name        string
	Description string
	Disabled    []string // Check IDs to disable
}

// BuiltInProfiles contains the predefined profiles.
var BuiltInProfiles = map[string]Profile{
	"poc": {
		Name:        "poc",
		Description: "Proof of Concept - minimal checks for early development",
		Disabled: []string{
			// Common checks to disable for PoC
			"common:license", "common:sast", "common:k8s", "common:shutdown",
			"common:health", "common:api_docs", "common:changelog",
			"common:integration", "common:metrics", "common:errors",
			"common:precommit", "common:env",
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
	"library": {
		Name:        "library",
		Description: "Library/package - focus on code quality, skip deployment checks",
		Disabled: []string{
			"common:dockerfile", "common:health", "common:k8s", "common:shutdown",
			"common:metrics", "common:errors", "common:integration",
		},
	},
	"production": {
		Name:        "production",
		Description: "Production application - all checks enabled",
		Disabled:    []string{}, // Enable everything
	},
}

// Get returns a profile by name and whether it was found.
func Get(name string) (Profile, bool) {
	p, ok := BuiltInProfiles[name]
	return p, ok
}

// List returns all available profiles in a consistent order.
func List() []Profile {
	return []Profile{
		BuiltInProfiles["poc"],
		BuiltInProfiles["library"],
		BuiltInProfiles["production"],
	}
}

// Names returns the names of all available profiles.
func Names() []string {
	return []string{"poc", "library", "production"}
}
