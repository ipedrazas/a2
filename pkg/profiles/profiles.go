// Package profiles provides application type profiles for different project types.
package profiles

import (
	"sort"
	"sync"
)

// Source indicates where a profile is defined.
type Source string

const (
	SourceBuiltIn Source = "built-in"
	SourceUser    Source = "user"
)

// Profile defines a set of checks to disable for a specific application type.
type Profile struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Disabled    []string `yaml:"disabled"`
	Source      Source   `yaml:"-"` // Not serialized to YAML
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

var (
	userProfiles map[string]Profile
	initOnce     sync.Once
	initErr      error
)

// Init loads user profiles from the config directory.
// This should be called once at application startup.
// It's safe to call multiple times; subsequent calls are no-ops.
func Init() error {
	initOnce.Do(func() {
		userProfiles, initErr = LoadUserProfiles()
	})
	return initErr
}

// Get returns a profile by name and whether it was found.
// User profiles take precedence over built-in profiles.
func Get(name string) (Profile, bool) {
	// User profiles take precedence
	if p, ok := userProfiles[name]; ok {
		return p, true
	}
	// Fall back to built-in
	if p, ok := BuiltInProfiles[name]; ok {
		p.Source = SourceBuiltIn
		return p, true
	}
	return Profile{}, false
}

// List returns all available profiles in a consistent order.
// User profiles are listed after built-in profiles, both sorted by name.
func List() []Profile {
	var profiles []Profile

	// Add built-in profiles
	for _, p := range BuiltInProfiles {
		// Skip if overridden by user
		if _, ok := userProfiles[p.Name]; ok {
			continue
		}
		p.Source = SourceBuiltIn
		profiles = append(profiles, p)
	}

	// Add user profiles
	for _, p := range userProfiles {
		profiles = append(profiles, p)
	}

	// Sort by source (built-in first), then by name
	sort.Slice(profiles, func(i, j int) bool {
		if profiles[i].Source != profiles[j].Source {
			return profiles[i].Source == SourceBuiltIn
		}
		return profiles[i].Name < profiles[j].Name
	})

	return profiles
}

// Names returns the names of all available profiles.
func Names() []string {
	profiles := List()
	names := make([]string, len(profiles))
	for i, p := range profiles {
		names[i] = p.Name
	}
	return names
}
