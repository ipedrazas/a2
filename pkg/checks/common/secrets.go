package common

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// SecretsCheck verifies that secrets scanning is configured or no hardcoded secrets exist.
type SecretsCheck struct{}

func (c *SecretsCheck) ID() string   { return "common:secrets" }
func (c *SecretsCheck) Name() string { return "Secrets Detection" }

// secretPattern represents a pattern to detect potential secrets.
type secretPattern struct {
	name    string
	pattern *regexp.Regexp
}

// Run checks for secret scanning configuration or scans for hardcoded secrets.
func (c *SecretsCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangCommon,
	}

	// First, check if secret scanning tools are configured
	scannerConfigs := []struct {
		name string
		file string
	}{
		{"Gitleaks", ".gitleaks.toml"},
		{"Gitleaks", ".gitleaks.yaml"},
		{"Gitleaks", "gitleaks.toml"},
		{"TruffleHog", ".trufflehog.yml"},
		{"TruffleHog", "trufflehog.yml"},
		{"Secretlint", ".secretlintrc"},
		{"Secretlint", ".secretlintrc.json"},
		{"git-secrets", ".git-secrets"},
		{"detect-secrets", ".secrets.baseline"},
	}

	var foundScanners []string
	for _, scanner := range scannerConfigs {
		if safepath.Exists(path, scanner.file) {
			foundScanners = append(foundScanners, scanner.name)
		}
	}

	// Check for pre-commit hooks that might include secret scanning
	if c.hasSecretScanningInPreCommit(path) {
		foundScanners = append(foundScanners, "pre-commit hook")
	}

	// If secret scanning is configured, that's a pass
	if len(foundScanners) > 0 {
		result.Passed = true
		result.Status = checker.Pass
		// Deduplicate scanner names
		unique := uniqueStrings(foundScanners)
		result.Message = "Secret scanning configured: " + strings.Join(unique, ", ")
		return result, nil
	}

	// No scanner configured - scan for potential secrets
	findings := c.scanForSecrets(path)

	if len(findings) > 0 {
		result.Passed = false
		result.Status = checker.Warn
		if len(findings) == 1 {
			result.Message = fmt.Sprintf("Potential secret found: %s", findings[0])
		} else if len(findings) <= 3 {
			result.Message = fmt.Sprintf("Potential secrets found: %s", strings.Join(findings, ", "))
		} else {
			result.Message = fmt.Sprintf("%d potential secrets found (e.g., %s)", len(findings), strings.Join(findings[:3], ", "))
		}
		return result, nil
	}

	// No secrets found but no scanner configured
	result.Passed = false
	result.Status = checker.Warn
	result.Message = "No secret scanning configured (consider adding gitleaks or similar)"
	return result, nil
}

// hasSecretScanningInPreCommit checks if pre-commit config includes secret scanning.
func (c *SecretsCheck) hasSecretScanningInPreCommit(path string) bool {
	preCommitFile := filepath.Join(path, ".pre-commit-config.yaml")
	content, err := safepath.ReadFile(path, ".pre-commit-config.yaml")
	if err != nil {
		return false
	}

	// Look for common secret scanning hooks
	secretHooks := []string{
		"gitleaks",
		"trufflehog",
		"detect-secrets",
		"git-secrets",
		"secretlint",
	}

	contentStr := strings.ToLower(string(content))
	for _, hook := range secretHooks {
		if strings.Contains(contentStr, hook) {
			return true
		}
	}

	_ = preCommitFile // Suppress unused variable warning
	return false
}

// scanForSecrets scans code files for potential hardcoded secrets.
func (c *SecretsCheck) scanForSecrets(path string) []string {
	var findings []string

	// Patterns to detect potential secrets
	patterns := []secretPattern{
		{"AWS Access Key", regexp.MustCompile(`(?i)AKIA[0-9A-Z]{16}`)},
		{"AWS Secret Key", regexp.MustCompile(`(?i)aws.{0,20}secret.{0,20}['"][0-9a-zA-Z/+]{40}['"]`)},
		{"Generic API Key", regexp.MustCompile(`(?i)api[_-]?key\s*[=:]\s*['"][a-zA-Z0-9]{20,}['"]`)},
		{"Generic Secret", regexp.MustCompile(`(?i)secret[_-]?key\s*[=:]\s*['"][a-zA-Z0-9]{20,}['"]`)},
		{"Private Key", regexp.MustCompile(`-----BEGIN\s+(RSA|DSA|EC|OPENSSH|PGP)\s+PRIVATE\s+KEY-----`)},
		{"GitHub Token", regexp.MustCompile(`(?i)gh[pousr]_[A-Za-z0-9_]{36,}`)},
		{"Generic Password", regexp.MustCompile(`(?i)password\s*[=:]\s*['"][^'"]{8,}['"]`)},
		{"Database URL", regexp.MustCompile(`(?i)(mysql|postgres|mongodb|redis):\/\/[^:]+:[^@]+@`)},
		{"JWT Token", regexp.MustCompile(`eyJ[A-Za-z0-9-_]+\.eyJ[A-Za-z0-9-_]+\.[A-Za-z0-9-_.+/]*`)},
		{"Slack Token", regexp.MustCompile(`xox[baprs]-[0-9]{10,13}-[0-9]{10,13}[a-zA-Z0-9-]*`)},
		{"Stripe Key", regexp.MustCompile(`(?i)sk_live_[0-9a-zA-Z]{24,}`)},
		{"SendGrid Key", regexp.MustCompile(`SG\.[a-zA-Z0-9_-]{22}\.[a-zA-Z0-9_-]{43}`)},
	}

	// File extensions to scan
	codeExtensions := map[string]bool{
		".go":     true,
		".py":     true,
		".js":     true,
		".ts":     true,
		".jsx":    true,
		".tsx":    true,
		".java":   true,
		".rb":     true,
		".php":    true,
		".cs":     true,
		".yaml":   true,
		".yml":    true,
		".json":   true,
		".xml":    true,
		".config": true,
		".env":    true,
		".sh":     true,
		".bash":   true,
	}

	// Directories to skip
	skipDirs := map[string]bool{
		"node_modules": true,
		"vendor":       true,
		".git":         true,
		"__pycache__":  true,
		".venv":        true,
		"venv":         true,
		"dist":         true,
		"build":        true,
		".idea":        true,
		".vscode":      true,
	}

	// Files to skip (templates, examples)
	skipFiles := map[string]bool{
		".env.example":  true,
		".env.sample":   true,
		".env.template": true,
		"example.env":   true,
	}

	_ = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip directories
		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") && name != "." {
				// Skip hidden directories except current
				if name != ".github" && name != ".circleci" {
					return filepath.SkipDir
				}
			}
			if skipDirs[name] {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if file should be skipped
		fileName := info.Name()
		if skipFiles[strings.ToLower(fileName)] {
			return nil
		}

		// Check file extension
		ext := strings.ToLower(filepath.Ext(filePath))
		isEnvFile := strings.HasPrefix(fileName, ".env") && !strings.Contains(fileName, "example") && !strings.Contains(fileName, "sample")
		if !codeExtensions[ext] && !isEnvFile {
			return nil
		}

		// Skip test files for generic patterns (they often have fake secrets)
		isTestFile := strings.Contains(fileName, "_test.") ||
			strings.Contains(fileName, ".test.") ||
			strings.Contains(fileName, ".spec.") ||
			strings.HasPrefix(fileName, "test_")

		// Scan file for secrets
		fileFindings := c.scanFile(path, filePath, patterns, isTestFile)
		findings = append(findings, fileFindings...)

		// Limit findings to avoid excessive output
		if len(findings) >= 10 {
			return filepath.SkipAll
		}

		return nil
	})

	return findings
}

// scanFile scans a single file for secret patterns.
func (c *SecretsCheck) scanFile(root, filePath string, patterns []secretPattern, isTestFile bool) []string {
	var findings []string

	file, err := safepath.OpenPath(root, filePath)
	if err != nil {
		return nil
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

	// Get relative path for reporting
	relPath, _ := filepath.Rel(root, filePath)
	if relPath == "" {
		relPath = filepath.Base(filePath)
	}

	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip comment lines (common false positives)
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "*") {
			// But still check for actual secrets in comments
			if !strings.Contains(strings.ToLower(trimmed), "example") &&
				!strings.Contains(strings.ToLower(trimmed), "todo") &&
				!strings.Contains(strings.ToLower(trimmed), "fixme") {
				// Check high-confidence patterns only in comments
				for _, p := range patterns {
					if p.name == "AWS Access Key" || p.name == "Private Key" || p.name == "JWT Token" {
						if p.pattern.MatchString(line) {
							findings = append(findings, fmt.Sprintf("%s in %s:%d", p.name, relPath, lineNum))
							break
						}
					}
				}
			}
			continue
		}

		// For test files, only check high-confidence patterns
		if isTestFile {
			for _, p := range patterns {
				if p.name == "AWS Access Key" || p.name == "Private Key" {
					if p.pattern.MatchString(line) {
						findings = append(findings, fmt.Sprintf("%s in %s:%d", p.name, relPath, lineNum))
						break
					}
				}
			}
			continue
		}

		// Check all patterns
		for _, p := range patterns {
			if p.pattern.MatchString(line) {
				findings = append(findings, fmt.Sprintf("%s in %s:%d", p.name, relPath, lineNum))
				break // One finding per line is enough
			}
		}

		if len(findings) >= 10 {
			break
		}
	}

	return findings
}

// uniqueStrings returns unique strings from a slice while preserving order.
func uniqueStrings(strs []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, s := range strs {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}
