package pythoncheck

import (
	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
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
	rb := checkutil.NewResultBuilder(c, checker.LangPython)
	pm := c.detectPackageManager(path)

	var result *checkutil.CommandResult
	var cmdDesc string

	switch pm {
	case "poetry":
		result = checkutil.RunCommand(path, "poetry", "check")
		cmdDesc = "poetry check"
	case "pipenv":
		result = checkutil.RunCommand(path, "pipenv", "check")
		cmdDesc = "pipenv check"
	default:
		result = checkutil.RunCommand(path, "pip", "--version")
		cmdDesc = "pip --version"
	}

	output := result.CombinedOutput()
	if !result.Success() {
		if checkutil.ToolNotFoundError(result.Err) {
			return rb.ToolNotInstalled(pm, ""), nil
		}
		return rb.FailWithOutput(cmdDesc+" failed: "+checkutil.TruncateMessage(result.Output(), 200), output), nil
	}

	return rb.Pass("Build check passed (" + pm + ")"), nil
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
