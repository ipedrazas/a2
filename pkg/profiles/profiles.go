// Package profiles provides application type profiles for different project types.
package profiles

// Profile defines a set of checks to disable for a specific application type.
type Profile struct {
	Name        string
	Description string
	Disabled    []string // Check IDs to disable
}

// BuiltInProfiles contains the predefined application profiles.
var BuiltInProfiles = map[string]Profile{
	"cli": {
		Name:        "cli",
		Description: "Command-line tool - skip server-related checks",
		Disabled: []string{
			"common:health",      // No health endpoints
			"common:k8s",         // Not containerized typically
			"common:metrics",     // No Prometheus metrics
			"common:api_docs",    // No API documentation
			"common:integration", // CLI doesn't need integration tests
			"common:shutdown",    // No graceful shutdown
			"common:errors",      // No error tracking service
			"common:e2e",         // No browser E2E tests
			"common:tracing",     // No distributed tracing
		},
	},
	"api": {
		Name:        "api",
		Description: "Web service/API - all operational checks enabled",
		Disabled: []string{
			"common:e2e", // API tests via integration tests instead
		},
	},
	"library": {
		Name:        "library",
		Description: "Reusable package/library - focus on code quality",
		Disabled: []string{
			"common:dockerfile",  // Libraries aren't containerized
			"common:health",      // No health endpoints
			"common:k8s",         // Not deployed
			"common:shutdown",    // No server to shutdown
			"common:metrics",     // No runtime metrics
			"common:errors",      // No error tracking
			"common:integration", // Unit tests suffice
			"common:tracing",     // No distributed tracing
			"common:e2e",         // No E2E tests
			"common:api_docs",    // API docs via code docs
		},
	},
	"desktop": {
		Name:        "desktop",
		Description: "Desktop application - focus on user-facing quality",
		Disabled: []string{
			"common:health",   // No health endpoints
			"common:k8s",      // Not containerized
			"common:api_docs", // No REST API
			"common:tracing",  // No distributed tracing
			"common:metrics",  // Different metrics approach
			"common:shutdown", // OS handles shutdown
		},
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
		BuiltInProfiles["cli"],
		BuiltInProfiles["api"],
		BuiltInProfiles["library"],
		BuiltInProfiles["desktop"],
	}
}

// Names returns the names of all available profiles.
func Names() []string {
	return []string{"cli", "api", "library", "desktop"}
}
