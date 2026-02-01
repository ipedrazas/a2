package common

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// HealthCheck verifies that health endpoint patterns exist in the codebase.
type HealthCheck struct{}

func (c *HealthCheck) ID() string   { return "common:health" }
func (c *HealthCheck) Name() string { return "Health Endpoint" }

func (c *HealthCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangCommon)

	// Resolve to absolute path so Walk and OpenPath use the same root (avoids
	// failures when path is relative or when running a2 from another directory).
	absPath, err := filepath.Abs(path)
	if err != nil {
		return rb.Fail(err.Error()), err
	}
	path = absPath

	// Health endpoint patterns to search for
	patterns := []string{
		"/health",
		"/healthz",
		"/ready",
		"/readiness",
		"/live",
		"/liveness",
		"/ping",
		"HealthCheck",
		"healthCheck",
		"health_check",
		"healthcheck",
		"readinessProbe",
		"livenessProbe",
	}

	// Code file extensions to search
	codeExtensions := map[string]bool{
		".go":   true,
		".py":   true,
		".js":   true,
		".ts":   true,
		".jsx":  true,
		".tsx":  true,
		".java": true,
		".rb":   true,
		".rs":   true,
		".yaml": true,
		".yml":  true,
		".json": true,
	}

	found := false
	var foundPattern string

	err = filepath.Walk(path, func(filePath string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return nil // Skip files we can't access
		}

		// Skip directories and hidden paths
		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" || name == "__pycache__" {
				return filepath.SkipDir
			}
			return nil
		}

		// Check file extension
		ext := strings.ToLower(filepath.Ext(filePath))
		if !codeExtensions[ext] {
			return nil
		}

		// Search file for patterns
		pattern, hasPattern := c.searchFileForPatterns(path, filePath, patterns)
		if hasPattern {
			found = true
			foundPattern = pattern
			return filepath.SkipAll // Stop walking once found
		}

		return nil
	})

	if err != nil {
		return rb.Fail(err.Error()), err
	}

	if !found {
		return rb.Warn("No health endpoint pattern detected"), nil
	}

	return rb.Pass("Health endpoint pattern found: " + foundPattern), nil
}

// searchFileForPatterns searches a file for any of the given patterns.
func (c *HealthCheck) searchFileForPatterns(root, filePath string, patterns []string) (string, bool) {
	file, err := safepath.OpenPath(root, filePath)
	if err != nil {
		return "", false
	}
	defer func() {
		if err := file.Close(); err != nil {
			// File close errors are typically not critical in read-only scenarios
			fmt.Println("Error closing file:", err)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		for _, pattern := range patterns {
			if strings.Contains(line, pattern) {
				return pattern, true
			}
		}
	}

	return "", false
}
