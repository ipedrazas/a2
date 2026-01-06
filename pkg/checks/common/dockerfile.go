package common

import (
	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// DockerfileCheck verifies that a Dockerfile or Containerfile exists.
type DockerfileCheck struct{}

func (c *DockerfileCheck) ID() string   { return "common:dockerfile" }
func (c *DockerfileCheck) Name() string { return "Container Ready" }

func (c *DockerfileCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangCommon,
	}

	// Check for various Dockerfile variants
	dockerfiles := []string{
		"Dockerfile",
		"dockerfile",
		"Containerfile",
		"containerfile",
	}

	var foundFile string
	for _, df := range dockerfiles {
		if safepath.Exists(path, df) {
			foundFile = df
			break
		}
	}

	if foundFile == "" {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "No Dockerfile or Containerfile found"
		return result, nil
	}

	// Check for .dockerignore (bonus)
	hasIgnore := safepath.Exists(path, ".dockerignore")

	if hasIgnore {
		result.Passed = true
		result.Status = checker.Pass
		result.Message = foundFile + " found with .dockerignore"
	} else {
		result.Passed = true
		result.Status = checker.Pass
		result.Message = foundFile + " found (consider adding .dockerignore)"
	}

	return result, nil
}
