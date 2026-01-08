package common

import (
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// ContributingCheck verifies contribution guidelines exist.
type ContributingCheck struct{}

func (c *ContributingCheck) ID() string   { return "common:contributing" }
func (c *ContributingCheck) Name() string { return "Contributing Guidelines" }

// Run checks for contribution guidelines and templates.
func (c *ContributingCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangCommon,
	}

	var found []string

	// Check for CONTRIBUTING file
	contributingFiles := []string{
		"CONTRIBUTING.md",
		"CONTRIBUTING.txt",
		"CONTRIBUTING",
		".github/CONTRIBUTING.md",
		"docs/CONTRIBUTING.md",
	}
	for _, file := range contributingFiles {
		if safepath.Exists(path, file) {
			found = append(found, file)
			break
		}
	}

	// Check for PR template
	prTemplates := []string{
		".github/PULL_REQUEST_TEMPLATE.md",
		".github/pull_request_template.md",
		".github/PULL_REQUEST_TEMPLATE/",
		"docs/pull_request_template.md",
	}
	for _, template := range prTemplates {
		if safepath.Exists(path, template) {
			found = append(found, "PR template")
			break
		}
	}

	// Check for issue templates
	issueTemplates := []string{
		".github/ISSUE_TEMPLATE.md",
		".github/issue_template.md",
		".github/ISSUE_TEMPLATE/",
	}
	for _, template := range issueTemplates {
		if safepath.Exists(path, template) {
			found = append(found, "issue templates")
			break
		}
	}

	// Check for code owners
	codeOwners := []string{
		"CODEOWNERS",
		".github/CODEOWNERS",
		"docs/CODEOWNERS",
	}
	for _, file := range codeOwners {
		if safepath.Exists(path, file) {
			found = append(found, "CODEOWNERS")
			break
		}
	}

	// Build result
	if len(found) > 0 {
		result.Passed = true
		result.Status = checker.Pass
		result.Message = "Found: " + strings.Join(found, ", ")
	} else {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "No contribution guidelines found (consider adding CONTRIBUTING.md)"
	}

	return result, nil
}
