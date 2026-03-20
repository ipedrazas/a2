package nodecheck

import (
	"fmt"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// DeadcodeCheck detects unused exports and dependencies using knip.
type DeadcodeCheck struct {
	Config *config.NodeLanguageConfig
}

func (c *DeadcodeCheck) ID() string   { return "node:deadcode" }
func (c *DeadcodeCheck) Name() string { return "Node Dead Code" }

func (c *DeadcodeCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangNode)

	if !safepath.Exists(path, "package.json") {
		return rb.Fail("package.json not found"), nil
	}

	// Check for knip (preferred) or ts-prune
	if c.hasKnip(path) {
		return c.runKnip(path, rb)
	}

	if checkutil.ToolAvailable("ts-prune") {
		return c.runTsPrune(path, rb)
	}

	// Check if knip is a devDependency but npx is available
	if checkutil.ToolAvailable("npx") && c.hasKnipDep(path) {
		return c.runKnip(path, rb)
	}

	return rb.ToolNotInstalled("knip", "npm install -D knip"), nil
}

// hasKnip checks if knip is available as a command.
func (c *DeadcodeCheck) hasKnip(path string) bool {
	return checkutil.ToolAvailable("knip")
}

// hasKnipDep checks if knip is listed as a devDependency.
func (c *DeadcodeCheck) hasKnipDep(path string) bool {
	pkg, err := ParsePackageJSON(path)
	if err != nil || pkg == nil {
		return false
	}
	_, ok := pkg.DevDependencies["knip"]
	return ok
}

// runKnip executes knip and parses its output.
func (c *DeadcodeCheck) runKnip(path string, rb *checkutil.ResultBuilder) (checker.Result, error) {
	var result *checkutil.CommandResult

	if checkutil.ToolAvailable("knip") {
		result = checkutil.RunCommand(path, "knip", "--no-progress")
	} else {
		result = checkutil.RunCommand(path, "npx", "knip", "--no-progress")
	}

	output := result.CombinedOutput()

	if result.Success() {
		return rb.PassWithOutput("No unused code detected (knip)", output), nil
	}

	// knip exits non-zero when issues are found
	issues := countKnipIssues(result.Stdout)
	if issues == 0 {
		issues = countKnipIssues(output)
	}

	if issues == 0 {
		// Might be an error rather than findings
		if strings.Contains(output, "error") || strings.Contains(output, "Error") {
			return rb.WarnWithOutput("knip error: "+checkutil.TruncateMessage(result.Output(), 200), output), nil
		}
		return rb.PassWithOutput("No unused code detected (knip)", output), nil
	}

	msg := fmt.Sprintf("knip: %d unused %s found",
		issues,
		checkutil.Pluralize(issues, "item", "items"),
	)

	return rb.WarnWithOutput(msg, output), nil
}

// runTsPrune executes ts-prune as a fallback.
func (c *DeadcodeCheck) runTsPrune(path string, rb *checkutil.ResultBuilder) (checker.Result, error) {
	result := checkutil.RunCommand(path, "ts-prune")
	output := result.CombinedOutput()

	if !result.Success() {
		return rb.WarnWithOutput("ts-prune error: "+checkutil.TruncateMessage(result.Output(), 200), output), nil
	}

	// Count unused exports — each line with "used in module" excluded
	unused := 0
	for _, line := range strings.Split(result.Stdout, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.Contains(trimmed, "used in module") {
			unused++
		}
	}

	if unused == 0 {
		return rb.PassWithOutput("No unused exports detected (ts-prune)", output), nil
	}

	msg := fmt.Sprintf("ts-prune: %d unused %s found",
		unused,
		checkutil.Pluralize(unused, "export", "exports"),
	)

	return rb.WarnWithOutput(msg, output), nil
}

// countKnipIssues counts issues in knip output.
// knip outputs sections like "Unused files (3)" or "Unused exports (5)".
func countKnipIssues(output string) int {
	count := 0
	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		// Count lines that look like findings (file paths or export references)
		// Skip section headers and blank lines
		if !strings.HasPrefix(trimmed, "Unused") &&
			!strings.HasPrefix(trimmed, "Unlisted") &&
			!strings.HasPrefix(trimmed, "---") &&
			len(trimmed) > 0 {
			count++
		}
	}
	return count
}
