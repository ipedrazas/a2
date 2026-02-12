package devops

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// TerraformCheck verifies that Terraform configurations exist and validates them.
type TerraformCheck struct{}

func (c *TerraformCheck) ID() string   { return "devops:terraform" }
func (c *TerraformCheck) Name() string { return "Terraform Configuration" }

// Run checks for Terraform files and validates them if terraform is installed.
func (c *TerraformCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangCommon)

	// Check if terraform files exist
	foundFiles := c.findTerraformFiles(path)
	if len(foundFiles) == 0 {
		return rb.Info("No Terraform files found"), nil
	}

	// Check if terraform tool is installed
	if !checkutil.ToolAvailable("terraform") {
		return rb.ToolNotInstalled("terraform", "install from https://terraform.io"), nil
	}

	// Initialize terraform if needed (for modules)
	// Check if .terraform directory exists
	hasDotTerraform := safepath.Exists(path, ".terraform")
	if !hasDotTerraform {
		// Try to run terraform init, but don't fail if it fails
		// Some repos might need backend config or credentials
		_ = checkutil.RunCommand(path, "terraform", "init", "-backend=false")
		// Silently ignore init failures - terraform validate might still work
	}

	// Run terraform validate
	result := checkutil.RunCommand(path, "terraform", "validate")

	if result.Success() {
		return rb.Pass("terraform validate passed"), nil
	}

	// Check if the error is about terraform not being initialized
	if strings.Contains(result.Output(), "run terraform init") {
		return rb.WarnWithOutput("terraform configuration found, run 'terraform init' first", result.CombinedOutput()), nil
	}

	return rb.WarnWithOutput("terraform validation failed", result.CombinedOutput()), nil
}

// findTerraformFiles searches for Terraform configuration files.
func (c *TerraformCheck) findTerraformFiles(path string) []string {
	var foundFiles []string

	// Check for .tf files recursively
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue walking on error
		}
		if info.IsDir() {
			// Skip hidden directories and common non-terraform dirs
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			if info.Name() == "node_modules" || info.Name() == "vendor" || info.Name() == "dist" || info.Name() == "build" {
				return filepath.SkipDir
			}
			return nil
		}

		// Check for .tf files
		if strings.HasSuffix(info.Name(), ".tf") {
			foundFiles = append(foundFiles, filePath)
		}

		return nil
	})

	if err != nil {
		return foundFiles
	}

	// Also check for .tfvars files in root
	tfvarsFiles := []string{"terraform.tfvars", "*.auto.tfvars"}
	for _, pattern := range tfvarsFiles {
		matches, _ := filepath.Glob(filepath.Join(path, pattern))
		foundFiles = append(foundFiles, matches...)
	}

	// Check for versions.tf file (common indicator)
	if safepath.Exists(path, "versions.tf") {
		foundFiles = append(foundFiles, filepath.Join(path, "versions.tf"))
	}

	return foundFiles
}
