package common

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// CICheck verifies that a CI/CD configuration exists.
type CICheck struct{}

func (c *CICheck) ID() string   { return "common:ci" }
func (c *CICheck) Name() string { return "CI Pipeline" }

func (c *CICheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangCommon,
	}

	// Check for various CI configurations
	ciConfigs := []struct {
		name  string
		check func(string) bool
	}{
		{"GitHub Actions", c.hasGitHubActions},
		{"GitLab CI", func(p string) bool { return safepath.Exists(p, ".gitlab-ci.yml") }},
		{"Jenkins", func(p string) bool { return safepath.Exists(p, "Jenkinsfile") }},
		{"CircleCI", func(p string) bool { return safepath.Exists(p, ".circleci/config.yml") }},
		{"Travis CI", func(p string) bool { return safepath.Exists(p, ".travis.yml") }},
		{"Azure Pipelines", func(p string) bool { return safepath.Exists(p, "azure-pipelines.yml") }},
		{"Bitbucket Pipelines", func(p string) bool { return safepath.Exists(p, "bitbucket-pipelines.yml") }},
		{"Drone CI", func(p string) bool { return safepath.Exists(p, ".drone.yml") }},
		{"Taskfile", func(p string) bool { return safepath.Exists(p, "Taskfile.yml") || safepath.Exists(p, "Taskfile.yaml") }},
	}

	var foundCIs []string
	for _, ci := range ciConfigs {
		if ci.check(path) {
			foundCIs = append(foundCIs, ci.name)
		}
	}

	if len(foundCIs) == 0 {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "No CI/CD configuration found"
		return result, nil
	}

	result.Passed = true
	result.Status = checker.Pass
	if len(foundCIs) == 1 {
		result.Message = foundCIs[0] + " configured"
	} else {
		result.Message = strings.Join(foundCIs, ", ") + " configured"
	}

	return result, nil
}

// hasGitHubActions checks for GitHub Actions workflow files.
func (c *CICheck) hasGitHubActions(path string) bool {
	workflowDir := filepath.Join(path, ".github", "workflows")
	entries, err := os.ReadDir(workflowDir)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml") {
			return true
		}
	}
	return false
}
