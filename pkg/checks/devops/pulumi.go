package devops

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// PulumiCheck verifies that Pulumi configurations exist and validates them.
type PulumiCheck struct{}

func (c *PulumiCheck) ID() string   { return "devops:pulumi" }
func (c *PulumiCheck) Name() string { return "Pulumi Configuration" }

// Run checks for Pulumi files and validates them if pulumi is installed.
func (c *PulumiCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangCommon)

	// Check if pulumi files exist
	foundFiles := c.findPulumiFiles(path)
	if len(foundFiles) == 0 {
		return rb.Info("No Pulumi files found"), nil
	}

	// Check if pulumi tool is installed
	if !checkutil.ToolAvailable("pulumi") {
		return rb.ToolNotInstalled("pulumi", "install from https://pulumi.io"), nil
	}

	// Run pulumi validate
	result := checkutil.RunCommand(path, "pulumi", "validate", "--offline")

	if result.Success() {
		return rb.Pass("pulumi validate passed"), nil
	}

	// Check if error is about not being a pulumi project
	if strings.Contains(result.Output(), "not a Pulumi project") || strings.Contains(result.Output(), "no Pulumi.yaml") {
		return rb.Info("Pulumi files found but not a valid Pulumi project"), nil
	}

	return rb.WarnWithOutput("pulumi validation failed", result.CombinedOutput()), nil
}

// findPulumiFiles searches for Pulumi configuration files.
func (c *PulumiCheck) findPulumiFiles(path string) []string {
	var foundFiles []string

	// Check for Pulumi.yaml in root
	if safepath.Exists(path, "Pulumi.yaml") {
		foundFiles = append(foundFiles, filepath.Join(path, "Pulumi.yaml"))
	}

	// Check for Pulumi.yml in root
	if safepath.Exists(path, "Pulumi.yml") {
		foundFiles = append(foundFiles, filepath.Join(path, "Pulumi.yml"))
	}

	// Check for language-specific Pulumi files
	pulumiLangFiles := []string{
		"Pulumi.go.yaml",
		"Pulumi.python.yaml",
		"Pulumi.ts.yaml",
		"Pulumi.javascript.yaml",
		"Pulumi.cs.yaml",
		"Pulumi.java.yaml",
		"Pulumi.yaml.yaml",
	}

	for _, file := range pulumiLangFiles {
		if safepath.Exists(path, file) {
			foundFiles = append(foundFiles, filepath.Join(path, file))
		}
	}

	// Check for Pulumi directory
	pulumiDir := filepath.Join(path, "pulumi")
	if info, err := os.Stat(pulumiDir); err == nil && info.IsDir() {
		foundFiles = append(foundFiles, pulumiDir)
	}

	// Check for .pulumi directory
	dotPulumiDir := filepath.Join(path, ".pulumi")
	if info, err := os.Stat(dotPulumiDir); err == nil && info.IsDir() {
		foundFiles = append(foundFiles, dotPulumiDir)
	}

	// Recursively check for Pulumi.yaml in subdirectories
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue walking on error
		}
		if info.IsDir() {
			// Skip hidden directories and common non-pulumi dirs
			if strings.HasPrefix(info.Name(), ".") && info.Name() != ".pulumi" {
				return filepath.SkipDir
			}
			if info.Name() == "node_modules" || info.Name() == "vendor" || info.Name() == "dist" || info.Name() == "build" {
				return filepath.SkipDir
			}
			return nil
		}

		// Check for Pulumi.yaml files in subdirectories
		if info.Name() == "Pulumi.yaml" || info.Name() == "Pulumi.yml" {
			// Only add if not already added from root
			if filePath != filepath.Join(path, info.Name()) {
				foundFiles = append(foundFiles, filePath)
			}
		}

		return nil
	})

	if err != nil {
		return foundFiles
	}

	return foundFiles
}
