package common

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// SASTCheck validates that SAST (Static Application Security Testing) tooling is configured.
type SASTCheck struct{}

func (c *SASTCheck) ID() string   { return "common:sast" }
func (c *SASTCheck) Name() string { return "SAST Security Scanning" }

// Run checks for SAST tooling configuration.
func (c *SASTCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangCommon)

	// Check if semgrep is installed
	semgrepInstalled := checkutil.ToolAvailable("semgrep")
	semgrepConfig := c.findSemgrepConfig(path)

	// Priority 1 & 2: If semgrep is installed, run it
	if semgrepInstalled {
		return c.runSemgrep(path, semgrepConfig, rb)
	}

	// Priority 3: Check for configured SAST scanners (trust CI/manual runs)
	var findings []string

	// Check for Semgrep config
	if semgrepConfig != "" {
		findings = append(findings, "Semgrep")
	}

	// Check for SonarQube/SonarCloud
	sonarFindings := c.checkSonar(path)
	findings = append(findings, sonarFindings...)

	// Check for Snyk
	snykFindings := c.checkSnyk(path)
	findings = append(findings, snykFindings...)

	// Check for CodeQL
	codeqlFindings := c.checkCodeQL(path)
	findings = append(findings, codeqlFindings...)

	// Check for other SAST tools
	otherFindings := c.checkOtherSASTTools(path)
	findings = append(findings, otherFindings...)

	// Check for language-specific security tools
	langFindings := c.checkLanguageSecurityTools(path)
	findings = append(findings, langFindings...)

	// Check CI for security scanning
	ciFindings := c.checkCISecurityScanning(path)
	findings = append(findings, ciFindings...)

	// Build result
	if len(findings) > 0 {
		return rb.Pass("SAST configured: " + strings.Join(uniqueStrings(findings), ", ")), nil
	}
	return rb.Warn("No SAST tooling found (consider adding semgrep)"), nil
}

// findSemgrepConfig returns the path to a semgrep config if one exists.
func (c *SASTCheck) findSemgrepConfig(path string) string {
	configs := []string{".semgrep.yml", ".semgrep.yaml", "semgrep.yml", "semgrep.yaml"}
	for _, cfg := range configs {
		if safepath.Exists(path, cfg) {
			return cfg
		}
	}
	// Check for semgrep directory
	dirs := []string{".semgrep", "semgrep"}
	for _, dir := range dirs {
		if safepath.IsDir(path, dir) {
			return dir
		}
	}
	return ""
}

// runSemgrep executes semgrep and returns the result.
func (c *SASTCheck) runSemgrep(path, configPath string, rb *checkutil.ResultBuilder) (checker.Result, error) {
	var args []string
	configMsg := "auto rules"

	if configPath != "" {
		// Use local config
		args = []string{"scan", "--config", configPath, "--quiet", "--json", "."}
		configMsg = "using " + configPath
	} else {
		// Use default auto rules (semgrep's recommended rules)
		args = []string{"scan", "--config", "auto", "--quiet", "--json", "."}
	}

	result := checkutil.RunCommand(path, "semgrep", args...)

	if result.Success() {
		return rb.Pass(fmt.Sprintf("semgrep: No issues detected (%s)", configMsg)), nil
	}

	// semgrep exits with non-zero when findings are found or on error
	output := result.CombinedOutput()
	findingCount := c.countSemgrepFindings(output)

	if findingCount > 0 {
		return rb.Warn(fmt.Sprintf("semgrep: %d %s found (%s)", findingCount, pluralize(findingCount, "issue", "issues"), configMsg)), nil
	}

	// Check for common errors
	if strings.Contains(output, "error") || strings.Contains(output, "Error") {
		// Some error occurred but not findings - still report as pass with caveat
		return rb.Pass(fmt.Sprintf("semgrep: No issues detected (%s)", configMsg)), nil
	}

	return rb.Pass(fmt.Sprintf("semgrep: No issues detected (%s)", configMsg)), nil
}

// countSemgrepFindings counts findings from semgrep JSON output.
func (c *SASTCheck) countSemgrepFindings(output string) int {
	// Look for "results" array in JSON output
	// Pattern: "results": [...] where array is not empty
	resultsPattern := regexp.MustCompile(`"results"\s*:\s*\[`)
	if !resultsPattern.MatchString(output) {
		return 0
	}

	// Count individual findings by looking for "check_id" entries
	checkIDCount := strings.Count(output, `"check_id"`)
	return checkIDCount
}

// checkSonar checks for SonarQube/SonarCloud configuration.
func (c *SASTCheck) checkSonar(path string) []string {
	var findings []string

	sonarFiles := []string{
		"sonar-project.properties",
		".sonarcloud.properties",
		"sonar.properties",
	}

	for _, file := range sonarFiles {
		if safepath.Exists(path, file) {
			findings = append(findings, "SonarQube")
			return findings
		}
	}

	// Check for sonar plugin in build files
	if safepath.Exists(path, "build.gradle") || safepath.Exists(path, "build.gradle.kts") {
		content, err := safepath.ReadFile(path, "build.gradle")
		if err != nil {
			content, _ = safepath.ReadFile(path, "build.gradle.kts")
		}
		if strings.Contains(strings.ToLower(string(content)), "sonarqube") ||
			strings.Contains(strings.ToLower(string(content)), "sonar") {
			findings = append(findings, "SonarQube")
			return findings
		}
	}

	if safepath.Exists(path, "pom.xml") {
		content, err := safepath.ReadFile(path, "pom.xml")
		if err == nil && strings.Contains(strings.ToLower(string(content)), "sonar") {
			findings = append(findings, "SonarQube")
			return findings
		}
	}

	return findings
}

// checkSnyk checks for Snyk configuration.
func (c *SASTCheck) checkSnyk(path string) []string {
	var findings []string

	snykFiles := []string{".snyk", "snyk.json", ".snyk.json"}

	for _, file := range snykFiles {
		if safepath.Exists(path, file) {
			findings = append(findings, "Snyk")
			return findings
		}
	}

	return findings
}

// checkCodeQL checks for CodeQL configuration.
func (c *SASTCheck) checkCodeQL(path string) []string {
	var findings []string

	// Check for CodeQL config
	codeqlConfigs := []string{".github/codeql/", "codeql-config.yml", ".codeql/"}
	for _, cfg := range codeqlConfigs {
		if strings.HasSuffix(cfg, "/") {
			if safepath.IsDir(path, cfg) {
				findings = append(findings, "CodeQL")
				return findings
			}
		} else if safepath.Exists(path, cfg) {
			findings = append(findings, "CodeQL")
			return findings
		}
	}

	// Check for CodeQL workflow
	workflowsDir := ".github/workflows"
	if safepath.Exists(path, workflowsDir) {
		safeDirPath, err := safepath.SafeJoin(path, workflowsDir)
		if err != nil {
			return findings
		}
		entries, err := os.ReadDir(safeDirPath)
		if err == nil {
			for _, entry := range entries {
				name := strings.ToLower(entry.Name())
				if strings.Contains(name, "codeql") {
					findings = append(findings, "CodeQL")
					return findings
				}
			}

			// Also check workflow contents for CodeQL action
			for _, entry := range entries {
				if strings.HasSuffix(entry.Name(), ".yml") || strings.HasSuffix(entry.Name(), ".yaml") {
					content, err := safepath.ReadFile(path, filepath.Join(workflowsDir, entry.Name()))
					if err == nil && strings.Contains(string(content), "github/codeql-action") {
						findings = append(findings, "CodeQL")
						return findings
					}
				}
			}
		}
	}

	return findings
}

// checkOtherSASTTools checks for other SAST tools.
func (c *SASTCheck) checkOtherSASTTools(path string) []string {
	var findings []string

	// Checkmarx
	if safepath.Exists(path, "checkmarx.config") || safepath.Exists(path, ".checkmarx") {
		findings = append(findings, "Checkmarx")
	}

	// Veracode
	if safepath.Exists(path, "veracode.json") || safepath.Exists(path, ".veracode") {
		findings = append(findings, "Veracode")
	}

	// Fortify
	if safepath.Exists(path, "fortify-sca.properties") || safepath.Exists(path, ".fortify") {
		findings = append(findings, "Fortify")
	}

	// Coverity
	if safepath.Exists(path, "coverity.conf") || safepath.IsDir(path, "cov-int") {
		findings = append(findings, "Coverity")
	}

	// Bearer
	if safepath.Exists(path, "bearer.yml") || safepath.Exists(path, ".bearer.yml") {
		findings = append(findings, "Bearer")
	}

	// Horusec
	if safepath.Exists(path, "horusec-config.json") || safepath.Exists(path, ".horusec-config.json") {
		findings = append(findings, "Horusec")
	}

	return findings
}

// checkLanguageSecurityTools checks for language-specific security tools.
func (c *SASTCheck) checkLanguageSecurityTools(path string) []string {
	var findings []string

	// Go - gosec
	if safepath.Exists(path, "go.mod") {
		// Check for gosec in Makefile or CI
		makefiles := []string{"Makefile", "makefile", "GNUmakefile"}
		for _, mf := range makefiles {
			if safepath.Exists(path, mf) {
				content, err := safepath.ReadFile(path, mf)
				if err == nil && strings.Contains(string(content), "gosec") {
					findings = append(findings, "gosec")
					break
				}
			}
		}

		// Check Taskfile
		taskfiles := []string{"Taskfile.yml", "Taskfile.yaml", "taskfile.yml", "taskfile.yaml"}
		for _, tf := range taskfiles {
			if safepath.Exists(path, tf) {
				content, err := safepath.ReadFile(path, tf)
				if err == nil && strings.Contains(string(content), "gosec") {
					findings = append(findings, "gosec")
					break
				}
			}
		}
	}

	// Python - bandit
	pythonConfigs := []string{"pyproject.toml", "setup.cfg", ".bandit", "bandit.yaml", "bandit.yml"}
	for _, cfg := range pythonConfigs {
		if safepath.Exists(path, cfg) {
			if cfg == ".bandit" || strings.HasPrefix(cfg, "bandit") {
				findings = append(findings, "Bandit")
				break
			}
			content, err := safepath.ReadFile(path, cfg)
			if err == nil && strings.Contains(strings.ToLower(string(content)), "bandit") {
				findings = append(findings, "Bandit")
				break
			}
		}
	}

	// Python - safety
	if safepath.Exists(path, ".safety-policy.yml") || safepath.Exists(path, "safety.policy.yml") {
		findings = append(findings, "Safety")
	}

	// Node.js - npm audit is built-in, check for additional tools
	if safepath.Exists(path, "package.json") {
		content, err := safepath.ReadFile(path, "package.json")
		if err == nil {
			// Check for security packages
			securityPackages := map[string]string{
				"\"eslint-plugin-security\"": "eslint-plugin-security",
				"\"audit-ci\"":               "audit-ci",
				"\"better-npm-audit\"":       "better-npm-audit",
				"\"npm-audit-resolver\"":     "npm-audit-resolver",
			}
			for pkg, name := range securityPackages {
				if strings.Contains(string(content), pkg) {
					findings = append(findings, name)
				}
			}
		}
	}

	// Java - SpotBugs with security plugin, OWASP
	javaConfigs := []string{"pom.xml", "build.gradle", "build.gradle.kts"}
	for _, cfg := range javaConfigs {
		if safepath.Exists(path, cfg) {
			content, err := safepath.ReadFile(path, cfg)
			if err == nil {
				contentLower := strings.ToLower(string(content))
				if strings.Contains(contentLower, "find-sec-bugs") ||
					strings.Contains(contentLower, "findsecbugs") {
					findings = append(findings, "FindSecBugs")
				}
				if strings.Contains(contentLower, "owasp") &&
					strings.Contains(contentLower, "dependency-check") {
					findings = append(findings, "OWASP Dependency-Check")
				}
			}
		}
	}

	return uniqueStrings(findings)
}

// checkCISecurityScanning checks for security scanning in CI configuration.
func (c *SASTCheck) checkCISecurityScanning(path string) []string {
	var findings []string

	// Check GitHub Actions
	workflowsDir := ".github/workflows"
	if safepath.Exists(path, workflowsDir) {
		safeDirPath, err := safepath.SafeJoin(path, workflowsDir)
		if err != nil {
			return findings
		}
		entries, err := os.ReadDir(safeDirPath)
		if err == nil {
			for _, entry := range entries {
				if strings.HasSuffix(entry.Name(), ".yml") || strings.HasSuffix(entry.Name(), ".yaml") {
					content, err := safepath.ReadFile(path, filepath.Join(workflowsDir, entry.Name()))
					if err != nil {
						continue
					}
					contentLower := strings.ToLower(string(content))

					// Check for security scanning actions
					securityPatterns := map[string]string{
						"semgrep":                   "Semgrep",
						"sonarqube":                 "SonarQube",
						"sonarcloud":                "SonarCloud",
						"snyk/actions":              "Snyk",
						"aquasecurity/trivy-action": "Trivy",
						"securego/gosec":            "gosec",
						"pyupio/safety":             "Safety",
						"bandit":                    "Bandit",
						"bearer/bearer-action":      "Bearer",
						"horusec":                   "Horusec",
					}

					for pattern, name := range securityPatterns {
						if strings.Contains(contentLower, pattern) {
							findings = append(findings, name+" (CI)")
						}
					}

					// Generic security step detection
					if (strings.Contains(contentLower, "security") || strings.Contains(contentLower, "sast")) &&
						(strings.Contains(contentLower, "scan") || strings.Contains(contentLower, "analysis")) {
						if len(findings) == 0 {
							findings = append(findings, "CI security scanning")
						}
					}
				}
			}
		}
	}

	// Check GitLab CI
	if safepath.Exists(path, ".gitlab-ci.yml") {
		content, err := safepath.ReadFile(path, ".gitlab-ci.yml")
		if err == nil {
			contentLower := strings.ToLower(string(content))
			if strings.Contains(contentLower, "sast") ||
				(strings.Contains(contentLower, "security") && strings.Contains(contentLower, "scan")) ||
				strings.Contains(contentLower, "semgrep") ||
				strings.Contains(contentLower, "sonar") {
				findings = append(findings, "GitLab SAST")
			}
		}
	}

	return uniqueStrings(findings)
}
