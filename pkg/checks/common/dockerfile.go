package common

import (
	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// DockerfileCheck verifies that a Dockerfile or Containerfile exists.
type DockerfileCheck struct{}

func (c *DockerfileCheck) ID() string   { return "common:dockerfile" }
func (c *DockerfileCheck) Name() string { return "Container Ready" }

func (c *DockerfileCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangCommon)

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
		return rb.Warn("No Dockerfile or Containerfile found"), nil
	}

	// Check for .dockerignore (bonus)
	hasIgnore := safepath.Exists(path, ".dockerignore")

	if hasIgnore {
		return rb.Pass(foundFile + " found with .dockerignore"), nil
	}
	return rb.Pass(foundFile + " found (consider adding .dockerignore)"), nil
}
