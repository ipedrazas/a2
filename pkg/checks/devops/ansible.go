package devops

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// AnsibleCheck verifies that Ansible configurations exist and validates them.
type AnsibleCheck struct{}

func (c *AnsibleCheck) ID() string   { return "devops:ansible" }
func (c *AnsibleCheck) Name() string { return "Ansible Configuration" }

// Run checks for Ansible files and validates them if ansible-lint is installed.
func (c *AnsibleCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangCommon)

	// Check if ansible files exist
	foundFiles := c.findAnsibleFiles(path)
	if len(foundFiles) == 0 {
		return rb.Info("No Ansible files found"), nil
	}

	// Check if ansible-lint is installed
	if !checkutil.ToolAvailable("ansible-lint") {
		return rb.ToolNotInstalled("ansible-lint", "install from https://ansible.com"), nil
	}

	// Run ansible-lint
	result := checkutil.RunCommand(path, "ansible-lint", ".")

	if result.Success() {
		return rb.Pass("ansible-lint passed"), nil
	}

	return rb.WarnWithOutput("ansible-lint found issues", result.CombinedOutput()), nil
}

// findAnsibleFiles searches for Ansible configuration files.
func (c *AnsibleCheck) findAnsibleFiles(path string) []string {
	var foundFiles []string

	// Check for ansible.cfg
	if safepath.Exists(path, "ansible.cfg") {
		foundFiles = append(foundFiles, filepath.Join(path, "ansible.cfg"))
	}

	// Check for common playbook filenames
	playbookNames := []string{
		"playbook.yml",
		"playbook.yaml",
		"site.yml",
		"site.yaml",
		"main.yml",
		"main.yaml",
	}

	for _, name := range playbookNames {
		if safepath.Exists(path, name) {
			foundFiles = append(foundFiles, filepath.Join(path, name))
		}
	}

	// Check for roles directory
	rolesDir := filepath.Join(path, "roles")
	if info, err := os.Stat(rolesDir); err == nil && info.IsDir() {
		foundFiles = append(foundFiles, rolesDir)
	}

	// Check for .ansible directory
	ansibleDir := filepath.Join(path, ".ansible")
	if info, err := os.Stat(ansibleDir); err == nil && info.IsDir() {
		foundFiles = append(foundFiles, ansibleDir)
	}

	// Check for ansible directory
	ansibleDir2 := filepath.Join(path, "ansible")
	if info, err := os.Stat(ansibleDir2); err == nil && info.IsDir() {
		foundFiles = append(foundFiles, ansibleDir2)
	}

	// Recursively check for .yml/.yaml files that might be playbooks
	// Only check in common ansible directories
	commonDirs := []string{".", "ansible", "playbooks", "roles"}
	for _, dir := range commonDirs {
		checkPath := path
		if dir != "." {
			checkPath = filepath.Join(path, dir)
		}

		if info, err := os.Stat(checkPath); err != nil || !info.IsDir() {
			continue
		}

		err := filepath.Walk(checkPath, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				// Skip hidden directories and common non-ansible dirs
				if strings.HasPrefix(info.Name(), ".") && info.Name() != ".ansible" {
					return filepath.SkipDir
				}
				if info.Name() == "node_modules" || info.Name() == "vendor" || info.Name() == "dist" || info.Name() == "build" {
					return filepath.SkipDir
				}
				// Limit depth to avoid scanning too deep
				depth := strings.Count(strings.TrimPrefix(filePath, checkPath), string(filepath.Separator))
				if depth > 3 {
					return filepath.SkipDir
				}
				return nil
			}

			// Check for playbook-like YAML files
			if strings.HasSuffix(info.Name(), ".yml") || strings.HasSuffix(info.Name(), ".yaml") {
				// Skip certain files that are unlikely to be playbooks
				name := strings.ToLower(info.Name())
				if strings.Contains(name, "requirements") || strings.Contains(name, "galaxy") {
					return nil
				}
				// Check if file might be a playbook (contains 'hosts:' or 'roles:')
				if c.isAnsiblePlaybook(checkPath, filePath) {
					foundFiles = append(foundFiles, filePath)
				}
			}

			return nil
		})

		if err == nil {
			break // Found at least one directory to check
		}
	}

	return foundFiles
}

// isAnsiblePlaybook checks if a YAML file appears to be an Ansible playbook.
func (c *AnsibleCheck) isAnsiblePlaybook(root, filePath string) bool {
	file, err := safepath.OpenPath(root, filePath)
	if err != nil {
		return false
	}

	defer func() {
		if err := file.Close(); err != nil {
			// File close errors are typically not critical in read-only scenarios
			fmt.Println("Error closing file:", err)
		}
	}()

	// Read first 50 lines looking for ansible indicators
	buf := make([]byte, 4096)
	n, _ := file.Read(buf)
	content := string(buf[:n])

	// Common ansible playbook indicators
	indicators := []string{
		"hosts:",
		"roles:",
		"tasks:",
		"- name:",
		"ansible.builtin.",
		"become:",
		"gather_facts:",
	}

	contentLower := strings.ToLower(content)
	for _, indicator := range indicators {
		if strings.Contains(contentLower, indicator) {
			return true
		}
	}

	return false
}
