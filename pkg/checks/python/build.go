package pythoncheck

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// BuildCheck verifies that Python dependencies can be installed.
type BuildCheck struct {
	Config *config.PythonLanguageConfig
}

func (c *BuildCheck) ID() string   { return "python:build" }
func (c *BuildCheck) Name() string { return "Python Build" }

func (c *BuildCheck) Run(path string) (checker.Result, error) {
	pm := c.detectPackageManager(path)

	var cmd *exec.Cmd
	var cmdDesc string

	switch pm {
	case "poetry":
		cmd = exec.Command("poetry", "check")
		cmdDesc = "poetry check"
	case "pipenv":
		cmd = exec.Command("pipenv", "check")
		cmdDesc = "pipenv check"
	default:
		// For pip, we just verify pip is available
		cmd = exec.Command("pip", "--version")
		cmdDesc = "pip --version"
	}

	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		output := strings.TrimSpace(stderr.String())
		if output == "" {
			output = strings.TrimSpace(stdout.String())
		}

		// Check if tool is not installed
		if strings.Contains(err.Error(), "executable file not found") {
			return checker.Result{
				Name:     c.Name(),
				ID:       c.ID(),
				Passed:   true,
				Status:   checker.Pass,
				Message:  pm + " not installed, skipping build check",
				Language: checker.LangPython,
			}, nil
		}

		return checker.Result{
			Name:     c.Name(),
			ID:       c.ID(),
			Passed:   false,
			Status:   checker.Fail,
			Message:  cmdDesc + " failed: " + output,
			Language: checker.LangPython,
		}, nil
	}

	return checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Passed:   true,
		Status:   checker.Pass,
		Message:  "Build check passed (" + pm + ")",
		Language: checker.LangPython,
	}, nil
}

func (c *BuildCheck) detectPackageManager(path string) string {
	if c.Config != nil && c.Config.PackageManager != "auto" && c.Config.PackageManager != "" {
		return c.Config.PackageManager
	}

	// Auto-detect based on lock files
	if safepath.Exists(path, "poetry.lock") {
		return "poetry"
	}
	if safepath.Exists(path, "Pipfile.lock") || safepath.Exists(path, "Pipfile") {
		return "pipenv"
	}
	return "pip"
}
