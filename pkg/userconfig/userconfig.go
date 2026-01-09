// Package userconfig provides utilities for accessing user configuration directories.
package userconfig

import (
	"os"
	"path/filepath"
)

const appName = "a2"

// GetConfigDir returns the user configuration directory for a2.
// On Unix: ~/.config/a2
// On Windows: %APPDATA%\a2
// Returns an error if the home directory cannot be determined.
func GetConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, appName), nil
}

// GetSubDir returns a subdirectory within the user config directory.
// For example, GetSubDir("profiles") returns ~/.config/a2/profiles on Unix.
func GetSubDir(subdir string) (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, subdir), nil
}

// EnsureDir creates a subdirectory within the config directory if it doesn't exist.
// Returns the full path to the directory.
func EnsureDir(subdir string) (string, error) {
	dir, err := GetSubDir(subdir)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return dir, nil
}

// DirExists checks if a subdirectory exists within the config directory.
func DirExists(subdir string) bool {
	dir, err := GetSubDir(subdir)
	if err != nil {
		return false
	}
	info, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return info.IsDir()
}
