package common

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// HealthCheck verifies that health endpoint patterns exist in the codebase.
type HealthCheck struct{}

func (c *HealthCheck) ID() string   { return "common:health" }
func (c *HealthCheck) Name() string { return "Health Endpoint" }

func (c *HealthCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangCommon,
	}

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

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
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
		return result, err
	}

	if !found {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "No health endpoint pattern detected"
		return result, nil
	}

	result.Passed = true
	result.Status = checker.Pass
	result.Message = "Health endpoint pattern found: " + foundPattern

	return result, nil
}

// searchFileForPatterns searches a file for any of the given patterns.
func (c *HealthCheck) searchFileForPatterns(root, filePath string, patterns []string) (string, bool) {
	file, err := safepath.OpenPath(root, filePath)
	if err != nil {
		return "", false
	}
	defer file.Close()

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
