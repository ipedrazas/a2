package nodecheck

import (
	"encoding/json"
	"fmt"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// DepsFreshnessCheck detects outdated Node.js dependencies.
type DepsFreshnessCheck struct {
	Config *config.NodeLanguageConfig
}

func (c *DepsFreshnessCheck) ID() string   { return "node:deps_freshness" }
func (c *DepsFreshnessCheck) Name() string { return "Node Dependency Freshness" }

func (c *DepsFreshnessCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangNode)

	if !safepath.Exists(path, "package.json") {
		return rb.Fail("package.json not found"), nil
	}

	pm := c.detectPackageManager(path)

	// npm outdated --json returns a JSON object of outdated packages
	result := checkutil.RunCommand(path, pm, "outdated", "--json")
	output := result.CombinedOutput()

	// npm outdated exits with code 1 when outdated packages exist, so don't check Success()
	outdated, err := countOutdatedNodePackages(result.Stdout)
	if err != nil {
		// Failed to parse; try combined output
		outdated, err = countOutdatedNodePackages(output)
		if err != nil {
			if result.Stdout == "" && result.Stderr == "" {
				return rb.PassWithOutput("All dependencies are up to date", output), nil
			}
			return rb.WarnWithOutput(pm+" outdated failed: "+checkutil.TruncateMessage(result.Output(), 200), output), nil
		}
	}

	if outdated == 0 {
		return rb.PassWithOutput("All dependencies are up to date", output), nil
	}

	msg := fmt.Sprintf("%d outdated %s",
		outdated,
		checkutil.Pluralize(outdated, "package", "packages"),
	)

	if outdated > 20 {
		return rb.WarnWithOutput(msg+". Run: "+pm+" outdated", output), nil
	}

	return rb.PassWithOutput(msg+". Run: "+pm+" update", output), nil
}

// detectPackageManager determines which package manager to use.
func (c *DepsFreshnessCheck) detectPackageManager(path string) string {
	if c.Config != nil && c.Config.PackageManager != "auto" && c.Config.PackageManager != "" {
		return c.Config.PackageManager
	}

	if safepath.Exists(path, "pnpm-lock.yaml") {
		return "pnpm"
	}
	if safepath.Exists(path, "yarn.lock") {
		return "yarn"
	}
	if safepath.Exists(path, "bun.lockb") {
		return "bun"
	}

	return "npm"
}

// countOutdatedNodePackages parses npm/yarn/pnpm outdated --json output.
// npm outdated --json returns: {"package-name": {"current": "1.0.0", "wanted": "1.1.0", "latest": "2.0.0"}, ...}
func countOutdatedNodePackages(output string) (int, error) {
	if output == "" {
		return 0, nil
	}

	var packages map[string]json.RawMessage
	if err := json.Unmarshal([]byte(output), &packages); err != nil {
		return 0, err
	}

	return len(packages), nil
}
