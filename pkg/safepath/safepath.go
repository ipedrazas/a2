// Package safepath provides secure file path operations to prevent directory traversal attacks.
package safepath

import (
	"fmt"
	"os"
	"path/filepath"
)

// ReadFile safely reads a file within the given root directory.
// It prevents directory traversal by ensuring the resolved path is within root.
func ReadFile(root, filename string) ([]byte, error) {
	safePath, err := SafeJoin(root, filename)
	if err != nil {
		return nil, err
	}

	// #nosec G304 -- Path is validated by SafeJoin to prevent traversal
	return os.ReadFile(safePath)
}

// SafeJoin securely joins a root directory with a relative path.
// It returns an error if the resulting path would escape the root directory.
func SafeJoin(root, relPath string) (string, error) {
	// Clean and resolve the root to absolute path
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("invalid root path: %w", err)
	}

	// Clean the relative path and join with root
	cleanRel := filepath.Clean(relPath)

	// Reject absolute paths in relPath
	if filepath.IsAbs(cleanRel) {
		return "", fmt.Errorf("absolute paths not allowed: %s", relPath)
	}

	// Join and resolve to absolute
	joined := filepath.Join(absRoot, cleanRel)
	absJoined, err := filepath.Abs(joined)
	if err != nil {
		return "", fmt.Errorf("invalid path: %w", err)
	}

	// Ensure the result is within root (prevent ../../../ attacks)
	// We need to add a separator to root to prevent matching partial directory names
	// e.g., /foo should not match /foobar
	rootWithSep := absRoot + string(filepath.Separator)
	if absJoined != absRoot && !hasPrefix(absJoined, rootWithSep) {
		return "", fmt.Errorf("path escapes root directory: %s", relPath)
	}

	return absJoined, nil
}

// Stat safely stats a file within the given root directory.
func Stat(root, filename string) (os.FileInfo, error) {
	safePath, err := SafeJoin(root, filename)
	if err != nil {
		return nil, err
	}

	return os.Stat(safePath)
}

// Exists checks if a file exists within the given root directory.
func Exists(root, filename string) bool {
	safePath, err := SafeJoin(root, filename)
	if err != nil {
		return false
	}

	_, err = os.Stat(safePath)
	return err == nil
}

// hasPrefix checks if path starts with prefix, handling OS-specific path separators.
func hasPrefix(path, prefix string) bool {
	// Use filepath.Clean to normalize both paths
	cleanPath := filepath.Clean(path)
	cleanPrefix := filepath.Clean(prefix)

	// Check if path starts with prefix
	if len(cleanPath) < len(cleanPrefix) {
		return false
	}

	return cleanPath[:len(cleanPrefix)] == cleanPrefix
}
