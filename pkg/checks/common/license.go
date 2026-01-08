package common

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// LicenseCheck validates dependency license compliance tooling.
type LicenseCheck struct{}

func (c *LicenseCheck) ID() string   { return "common:license" }
func (c *LicenseCheck) Name() string { return "License Compliance" }

// Run checks for license compliance tooling configuration.
func (c *LicenseCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangCommon,
	}

	var findings []string

	// Check for license audit config files
	configFindings := c.checkConfigFiles(path)
	findings = append(findings, configFindings...)

	// Check for FOSSA
	fossaFindings := c.checkFOSSA(path)
	findings = append(findings, fossaFindings...)

	// Check for SPDX
	spdxFindings := c.checkSPDX(path)
	findings = append(findings, spdxFindings...)

	// Check for license tools in Go dependencies
	goFindings := c.checkGoLicenseTools(path)
	findings = append(findings, goFindings...)

	// Check for license tools in Python dependencies
	pythonFindings := c.checkPythonLicenseTools(path)
	findings = append(findings, pythonFindings...)

	// Check for license tools in Node.js dependencies
	nodeFindings := c.checkNodeLicenseTools(path)
	findings = append(findings, nodeFindings...)

	// Check for license tools in Java dependencies
	javaFindings := c.checkJavaLicenseTools(path)
	findings = append(findings, javaFindings...)

	// Check CI for license scanning
	ciFindings := c.checkCILicenseScanning(path)
	findings = append(findings, ciFindings...)

	// Build result
	if len(findings) > 0 {
		result.Passed = true
		result.Status = checker.Pass
		result.Message = "License compliance: " + strings.Join(findings, ", ")
	} else {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "No license compliance tooling found (consider adding license scanning)"
	}

	return result, nil
}

// checkConfigFiles checks for license audit configuration files.
func (c *LicenseCheck) checkConfigFiles(path string) []string {
	var findings []string

	// License checker configs
	configFiles := map[string]string{
		".licensrc":             "licensrc config",
		".licensrc.json":        "licensrc config",
		".licensrc.yaml":        "licensrc config",
		".licensrc.yml":         "licensrc config",
		"license-checker.json":  "license-checker config",
		".license-checker.json": "license-checker config",
		"license.json":          "license config",
		".licenserc":            "license config",
		".licenserc.json":       "license config",
		".licenserc.yaml":       "license config",
	}

	for file, desc := range configFiles {
		if safepath.Exists(path, file) {
			findings = append(findings, desc)
			break // Only report once per type
		}
	}

	return findings
}

// checkFOSSA checks for FOSSA configuration.
func (c *LicenseCheck) checkFOSSA(path string) []string {
	var findings []string

	fossaFiles := []string{".fossa.yml", ".fossa.yaml", "fossa.yml", "fossa.yaml"}
	for _, file := range fossaFiles {
		if safepath.Exists(path, file) {
			findings = append(findings, "FOSSA")
			return findings
		}
	}

	return findings
}

// checkSPDX checks for SPDX license files.
func (c *LicenseCheck) checkSPDX(path string) []string {
	var findings []string

	// Check for SPDX files
	spdxPatterns := []string{
		"*.spdx", "*.spdx.json", "*.spdx.yaml", "*.spdx.yml",
		"spdx.json", "spdx.yaml", "spdx.yml",
	}

	for _, pattern := range spdxPatterns {
		if pattern[0] == '*' {
			// Glob pattern - check common names
			if safepath.Exists(path, "sbom.spdx") ||
				safepath.Exists(path, "sbom.spdx.json") ||
				safepath.Exists(path, "bom.spdx.json") {
				findings = append(findings, "SPDX SBOM")
				return findings
			}
		} else {
			if safepath.Exists(path, pattern) {
				findings = append(findings, "SPDX")
				return findings
			}
		}
	}

	return findings
}

// checkGoLicenseTools checks for Go license checking tools.
func (c *LicenseCheck) checkGoLicenseTools(path string) []string {
	var findings []string

	if !safepath.Exists(path, "go.mod") {
		return findings
	}

	content, err := safepath.ReadFile(path, "go.mod")
	if err != nil {
		return findings
	}

	goLicenseTools := map[string]string{
		"github.com/google/go-licenses":        "go-licenses",
		"github.com/mitchellh/golicense":       "golicense",
		"github.com/uw-labs/lichen":            "lichen",
		"github.com/ribice/glice":              "glice",
		"github.com/fossa/fossa-cli":           "FOSSA CLI",
		"github.com/anchore/syft":              "Syft",
		"github.com/CycloneDX/cyclonedx-gomod": "CycloneDX",
	}

	for tool, name := range goLicenseTools {
		if strings.Contains(string(content), tool) {
			findings = append(findings, name)
		}
	}

	return findings
}

// checkPythonLicenseTools checks for Python license checking tools.
func (c *LicenseCheck) checkPythonLicenseTools(path string) []string {
	var findings []string

	pythonConfigs := []string{"pyproject.toml", "requirements.txt", "requirements-dev.txt", "setup.py", "Pipfile"}
	pythonLicenseTools := map[string]string{
		"pip-licenses":     "pip-licenses",
		"liccheck":         "liccheck",
		"license-check":    "license-check",
		"piplicenses":      "pip-licenses",
		"scancode-toolkit": "ScanCode",
		"cyclonedx-bom":    "CycloneDX",
		"cyclonedx-py":     "CycloneDX",
		"fossa-cli":        "FOSSA CLI",
	}

	for _, cfg := range pythonConfigs {
		if safepath.Exists(path, cfg) {
			content, err := safepath.ReadFile(path, cfg)
			if err != nil {
				continue
			}
			contentLower := strings.ToLower(string(content))
			for tool, name := range pythonLicenseTools {
				if strings.Contains(contentLower, tool) {
					findings = append(findings, name)
				}
			}
		}
	}

	// Check for liccheck config
	if safepath.Exists(path, ".liccheck.ini") || safepath.Exists(path, "liccheck.ini") {
		findings = append(findings, "liccheck config")
	}

	// Deduplicate findings
	return uniqueStrings(findings)
}

// checkNodeLicenseTools checks for Node.js license checking tools.
func (c *LicenseCheck) checkNodeLicenseTools(path string) []string {
	var findings []string

	if !safepath.Exists(path, "package.json") {
		return findings
	}

	content, err := safepath.ReadFile(path, "package.json")
	if err != nil {
		return findings
	}

	nodeLicenseTools := map[string]string{
		"license-checker":          "license-checker",
		"license-compliance":       "license-compliance",
		"license-webpack-plugin":   "license-webpack-plugin",
		"legally":                  "legally",
		"nlf":                      "nlf",
		"@cyclonedx/bom":           "CycloneDX",
		"@cyclonedx/cyclonedx-npm": "CycloneDX",
		"snyk":                     "Snyk",
		"fossa-cli":                "FOSSA CLI",
	}

	for tool, name := range nodeLicenseTools {
		if strings.Contains(string(content), "\""+tool+"\"") {
			findings = append(findings, name)
		}
	}

	return findings
}

// checkJavaLicenseTools checks for Java license checking tools.
func (c *LicenseCheck) checkJavaLicenseTools(path string) []string {
	var findings []string

	javaConfigs := []string{"pom.xml", "build.gradle", "build.gradle.kts"}
	javaLicenseTools := map[string]string{
		"license-maven-plugin":            "license-maven-plugin",
		"com.mycila:license-maven-plugin": "Mycila License Plugin",
		"license-gradle-plugin":           "license-gradle-plugin",
		"com.github.hierynomus.license":   "Gradle License Plugin",
		"org.cyclonedx":                   "CycloneDX",
		"cyclonedx-maven-plugin":          "CycloneDX",
		"cyclonedx-gradle-plugin":         "CycloneDX",
		"dependency-license-report":       "Dependency License Report",
	}

	for _, cfg := range javaConfigs {
		if safepath.Exists(path, cfg) {
			content, err := safepath.ReadFile(path, cfg)
			if err != nil {
				continue
			}
			contentLower := strings.ToLower(string(content))
			for tool, name := range javaLicenseTools {
				if strings.Contains(contentLower, strings.ToLower(tool)) {
					findings = append(findings, name)
				}
			}
		}
	}

	return uniqueStrings(findings)
}

// checkCILicenseScanning checks for license scanning in CI configuration.
func (c *LicenseCheck) checkCILicenseScanning(path string) []string {
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
				file := entry.Name()
				if strings.HasSuffix(file, ".yml") || strings.HasSuffix(file, ".yaml") {
					content, err := safepath.ReadFile(path, filepath.Join(workflowsDir, file))
					if err != nil {
						continue
					}
					contentLower := strings.ToLower(string(content))

					// Check for license scanning actions/tools
					licensePatterns := []string{
						"fossa-action",
						"license-checker",
						"go-licenses",
						"pip-licenses",
						"license-scan",
						"snyk/actions",
						"cyclonedx",
						"sbom",
					}

					for _, pattern := range licensePatterns {
						if strings.Contains(contentLower, pattern) {
							findings = append(findings, "CI license scanning")
							return findings
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
			if strings.Contains(contentLower, "license") &&
				(strings.Contains(contentLower, "scan") || strings.Contains(contentLower, "check")) {
				findings = append(findings, "CI license scanning")
			}
		}
	}

	return findings
}
