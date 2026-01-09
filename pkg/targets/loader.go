package targets

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ipedrazas/a2/pkg/userconfig"
	"gopkg.in/yaml.v3"
)

const targetsSubdir = "targets"

// LoadUserTargets discovers and loads user targets from the config directory.
// Returns an empty map if the directory doesn't exist or is empty.
func LoadUserTargets() (map[string]Target, error) {
	dir, err := userconfig.GetSubDir(targetsSubdir)
	if err != nil {
		return make(map[string]Target), nil
	}

	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return make(map[string]Target), nil
	}

	return DiscoverTargets(dir)
}

// DiscoverTargets finds all YAML files in a directory and loads them as targets.
func DiscoverTargets(dir string) (map[string]Target, error) {
	targets := make(map[string]Target)

	pattern := filepath.Join(dir, "*.yaml")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return targets, err
	}

	// Also check for .yml extension
	ymlPattern := filepath.Join(dir, "*.yml")
	ymlFiles, err := filepath.Glob(ymlPattern)
	if err != nil {
		return targets, err
	}
	files = append(files, ymlFiles...)

	for _, file := range files {
		target, err := LoadFromFile(file)
		if err != nil {
			// Log warning but continue loading other targets
			fmt.Fprintf(os.Stderr, "Warning: failed to load target from %s: %v\n", file, err)
			continue
		}
		targets[target.Name] = *target
	}

	return targets, nil
}

// LoadFromFile loads a single target from a YAML file.
func LoadFromFile(path string) (*Target, error) {
	// #nosec G304 -- Path comes from trusted source (user config directory via filepath.Glob)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var target Target
	if err := yaml.Unmarshal(data, &target); err != nil {
		return nil, err
	}

	// If name is not set, use the filename
	if target.Name == "" {
		base := filepath.Base(path)
		target.Name = strings.TrimSuffix(strings.TrimSuffix(base, ".yaml"), ".yml")
	}

	target.Source = SourceUser
	return &target, nil
}

// WriteBuiltInTargets writes all built-in targets to the user config directory.
// This is used by the "targets init" command.
func WriteBuiltInTargets() error {
	dir, err := userconfig.EnsureDir(targetsSubdir)
	if err != nil {
		return fmt.Errorf("failed to create targets directory: %w", err)
	}

	for name, target := range BuiltInTargets {
		path := filepath.Join(dir, name+".yaml")

		// Don't overwrite existing files
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("Skipping %s (already exists)\n", path)
			continue
		}

		if err := writeTarget(path, target); err != nil {
			return fmt.Errorf("failed to write target %s: %w", name, err)
		}
		fmt.Printf("Created %s\n", path)
	}

	return nil
}

// writeTarget writes a target to a YAML file.
func writeTarget(path string, target Target) error {
	// Create a copy without the Source field for serialization
	t := struct {
		Name        string   `yaml:"name"`
		Description string   `yaml:"description"`
		Disabled    []string `yaml:"disabled"`
	}{
		Name:        target.Name,
		Description: target.Description,
		Disabled:    target.Disabled,
	}

	data, err := yaml.Marshal(&t)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
