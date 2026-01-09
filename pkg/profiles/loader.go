package profiles

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ipedrazas/a2/pkg/userconfig"
	"gopkg.in/yaml.v3"
)

const profilesSubdir = "profiles"

// LoadUserProfiles discovers and loads user profiles from the config directory.
// Returns an empty map if the directory doesn't exist or is empty.
func LoadUserProfiles() (map[string]Profile, error) {
	dir, err := userconfig.GetSubDir(profilesSubdir)
	if err != nil {
		return make(map[string]Profile), nil
	}

	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return make(map[string]Profile), nil
	}

	return DiscoverProfiles(dir)
}

// DiscoverProfiles finds all YAML files in a directory and loads them as profiles.
func DiscoverProfiles(dir string) (map[string]Profile, error) {
	profiles := make(map[string]Profile)

	pattern := filepath.Join(dir, "*.yaml")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return profiles, err
	}

	// Also check for .yml extension
	ymlPattern := filepath.Join(dir, "*.yml")
	ymlFiles, err := filepath.Glob(ymlPattern)
	if err != nil {
		return profiles, err
	}
	files = append(files, ymlFiles...)

	for _, file := range files {
		profile, err := LoadFromFile(file)
		if err != nil {
			// Log warning but continue loading other profiles
			fmt.Fprintf(os.Stderr, "Warning: failed to load profile from %s: %v\n", file, err)
			continue
		}
		profiles[profile.Name] = *profile
	}

	return profiles, nil
}

// LoadFromFile loads a single profile from a YAML file.
func LoadFromFile(path string) (*Profile, error) {
	// #nosec G304 -- Path comes from trusted source (user config directory via filepath.Glob)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var profile Profile
	if err := yaml.Unmarshal(data, &profile); err != nil {
		return nil, err
	}

	// If name is not set, use the filename
	if profile.Name == "" {
		base := filepath.Base(path)
		profile.Name = strings.TrimSuffix(strings.TrimSuffix(base, ".yaml"), ".yml")
	}

	profile.Source = SourceUser
	return &profile, nil
}

// WriteBuiltInProfiles writes all built-in profiles to the user config directory.
// This is used by the "profiles init" command.
func WriteBuiltInProfiles() error {
	dir, err := userconfig.EnsureDir(profilesSubdir)
	if err != nil {
		return fmt.Errorf("failed to create profiles directory: %w", err)
	}

	for name, profile := range BuiltInProfiles {
		path := filepath.Join(dir, name+".yaml")

		// Don't overwrite existing files
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("Skipping %s (already exists)\n", path)
			continue
		}

		if err := writeProfile(path, profile); err != nil {
			return fmt.Errorf("failed to write profile %s: %w", name, err)
		}
		fmt.Printf("Created %s\n", path)
	}

	return nil
}

// writeProfile writes a profile to a YAML file.
func writeProfile(path string, profile Profile) error {
	// Create a copy without the Source field for serialization
	p := struct {
		Name        string   `yaml:"name"`
		Description string   `yaml:"description"`
		Disabled    []string `yaml:"disabled"`
	}{
		Name:        profile.Name,
		Description: profile.Description,
		Disabled:    profile.Disabled,
	}

	data, err := yaml.Marshal(&p)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
