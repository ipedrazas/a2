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

	// Ensure the result is within root (prevent ../../../ and symlink escapes)
	if !contained(absRoot, absJoined) {
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

// Open safely opens a file within the given root directory.
// It prevents directory traversal by ensuring the resolved path is within root.
func Open(root, filename string) (*os.File, error) {
	safePath, err := SafeJoin(root, filename)
	if err != nil {
		return nil, err
	}

	// #nosec G304 -- Path is validated by SafeJoin to prevent traversal
	return os.Open(safePath)
}

// OpenPath safely opens a file path that should be within the given root directory.
// Unlike Open, this validates an already-joined path against the root.
// This is useful when the path comes from filepath.Walk.
func OpenPath(root, filePath string) (*os.File, error) {
	// Resolve both paths to absolute
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("invalid root path: %w", err)
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("invalid file path: %w", err)
	}

	// Ensure the path is within root, following symlinks
	if !contained(absRoot, absPath) {
		return nil, fmt.Errorf("path escapes root directory: %s", filePath)
	}

	// #nosec G304 -- Path is validated above to prevent traversal
	return os.Open(absPath)
}

// within reports whether absPath is absRoot or lexically inside it. The
// separator suffix stops /foo from matching /foobar.
func within(absRoot, absPath string) bool {
	return absPath == absRoot || hasPrefix(absPath, absRoot+string(filepath.Separator))
}

// contained reports whether absPath stays inside absRoot both lexically and
// after resolving symlinks, so a link inside a scanned repository cannot point
// a2 at files outside it. A path that does not exist has no links to follow,
// so it is only checked lexically.
func contained(absRoot, absPath string) bool {
	if !within(absRoot, absPath) {
		return false
	}
	realPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		return true
	}
	realRoot, err := filepath.EvalSymlinks(absRoot)
	if err != nil {
		realRoot = absRoot
	}
	return within(realRoot, realPath)
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

// IsDir checks if a path is a directory within the given root directory.
func IsDir(root, dirname string) bool {
	safePath, err := SafeJoin(root, dirname)
	if err != nil {
		return false
	}

	info, err := os.Stat(safePath)
	if err != nil {
		return false
	}

	return info.IsDir()
}

// Glob performs a glob pattern match within the given root directory.
// Returns absolute paths of matching files.
func Glob(root, pattern string) ([]string, error) {
	safePath, err := SafeJoin(root, pattern)
	if err != nil {
		// If the pattern itself contains traversal, fail
		return nil, err
	}

	// Use filepath.Glob on the safe path
	matches, err := filepath.Glob(safePath)
	if err != nil {
		return nil, err
	}

	// Validate each match is within root
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	var validMatches []string
	for _, match := range matches {
		absMatch, err := filepath.Abs(match)
		if err != nil {
			continue
		}
		if contained(absRoot, absMatch) {
			validMatches = append(validMatches, absMatch)
		}
	}

	return validMatches, nil
}

// ReadFileAbs reads a file given an absolute path, validating it doesn't escape the cwd.
// This is useful for reading files returned by Glob.
func ReadFileAbs(absPath string) ([]byte, error) {
	// Basic validation that it's an absolute path
	if !filepath.IsAbs(absPath) {
		return nil, fmt.Errorf("path must be absolute: %s", absPath)
	}

	// #nosec G304 -- Caller is responsible for ensuring path is safe (e.g., from Glob)
	return os.ReadFile(absPath)
}
