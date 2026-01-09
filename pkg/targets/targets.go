// Package targets provides maturity target definitions for different project stages.
package targets

import (
	"sort"
	"sync"
)

// Source indicates where a target is defined.
type Source string

const (
	SourceBuiltIn Source = "built-in"
	SourceUser    Source = "user"
)

// Target defines a set of checks to disable for a specific maturity level.
type Target struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Disabled    []string `yaml:"disabled"`
	Source      Source   `yaml:"-"` // Not serialized to YAML
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

var (
	userTargets map[string]Target
	initOnce    sync.Once
	initErr     error
)

// Init loads user targets from the config directory.
// This should be called once at application startup.
// It's safe to call multiple times; subsequent calls are no-ops.
func Init() error {
	initOnce.Do(func() {
		userTargets, initErr = LoadUserTargets()
	})
	return initErr
}

// Get returns a target by name and whether it was found.
// User targets take precedence over built-in targets.
func Get(name string) (Target, bool) {
	// User targets take precedence
	if t, ok := userTargets[name]; ok {
		return t, true
	}
	// Fall back to built-in
	if t, ok := BuiltInTargets[name]; ok {
		t.Source = SourceBuiltIn
		return t, true
	}
	return Target{}, false
}

// List returns all available targets in a consistent order.
// User targets are listed after built-in targets, both sorted by name.
func List() []Target {
	var targets []Target

	// Add built-in targets
	for _, t := range BuiltInTargets {
		// Skip if overridden by user
		if _, ok := userTargets[t.Name]; ok {
			continue
		}
		t.Source = SourceBuiltIn
		targets = append(targets, t)
	}

	// Add user targets
	for _, t := range userTargets {
		targets = append(targets, t)
	}

	// Sort by source (built-in first), then by name
	sort.Slice(targets, func(i, j int) bool {
		if targets[i].Source != targets[j].Source {
			return targets[i].Source == SourceBuiltIn
		}
		return targets[i].Name < targets[j].Name
	})

	return targets
}

// Names returns the names of all available targets.
func Names() []string {
	targets := List()
	names := make([]string, len(targets))
	for i, t := range targets {
		names[i] = t.Name
	}
	return names
}
