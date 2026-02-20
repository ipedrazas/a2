package common

import (
	"path/filepath"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// RetryCheck verifies retry logic exists for external calls.
type RetryCheck struct{}

func (c *RetryCheck) ID() string   { return "common:retry" }
func (c *RetryCheck) Name() string { return "Retry Logic" }

// Run checks for retry and backoff libraries.
func (c *RetryCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangCommon)

	var found []string

	// Check Go dependencies for retry libraries
	if safepath.Exists(path, "go.mod") {
		if content, err := safepath.ReadFile(path, "go.mod"); err == nil {
			goRetryLibs := map[string]string{
				"github.com/cenkalti/backoff":           "backoff",
				"github.com/avast/retry-go":             "retry-go",
				"github.com/hashicorp/go-retryablehttp": "go-retryablehttp",
				"github.com/sethvargo/go-retry":         "go-retry",
				"github.com/eapache/go-resiliency":      "go-resiliency",
				"github.com/sony/gobreaker":             "gobreaker",
				"github.com/afex/hystrix-go":            "hystrix-go",
				"github.com/failsafe-go/failsafe-go":    "failsafe-go",
			}
			for dep, name := range goRetryLibs {
				if strings.Contains(string(content), dep) {
					if !containsString(found, name) {
						found = append(found, name)
					}
				}
			}
		}
	}

	// Check Node.js dependencies
	if safepath.Exists(path, "package.json") {
		if content, err := safepath.ReadFile(path, "package.json"); err == nil {
			nodeRetryLibs := map[string]string{
				"async-retry":         "async-retry",
				"axios-retry":         "axios-retry",
				"p-retry":             "p-retry",
				"retry":               "retry",
				"got":                 "got (built-in retry)",
				"ky":                  "ky (built-in retry)",
				"exponential-backoff": "exponential-backoff",
				"cockatiel":           "Cockatiel",
				"opossum":             "Opossum (circuit breaker)",
				"@nestjs/terminus":    "NestJS Terminus",
			}
			for dep, name := range nodeRetryLibs {
				if strings.Contains(string(content), `"`+dep+`"`) {
					if !containsString(found, name) {
						found = append(found, name)
					}
				}
			}
		}
	}

	// Check Python dependencies
	pythonFiles := []string{"pyproject.toml", "requirements.txt", "setup.py"}
	for _, file := range pythonFiles {
		if content, err := safepath.ReadFile(path, file); err == nil {
			pythonRetryLibs := map[string]string{
				"tenacity":       "Tenacity",
				"backoff":        "backoff",
				"retry":          "retry",
				"retrying":       "retrying",
				"urllib3":        "urllib3 (Retry)",
				"aiohttp":        "aiohttp",
				"httpx":          "httpx",
				"pybreaker":      "pybreaker",
				"circuitbreaker": "circuitbreaker",
			}
			for dep, name := range pythonRetryLibs {
				if strings.Contains(strings.ToLower(string(content)), dep) {
					if !containsString(found, name) {
						found = append(found, name)
					}
				}
			}
			break
		}
	}

	// Check Java dependencies
	javaFiles := []string{"pom.xml", "build.gradle", "build.gradle.kts"}
	for _, file := range javaFiles {
		if content, err := safepath.ReadFile(path, file); err == nil {
			javaRetryLibs := map[string]string{
				"spring-retry":    "Spring Retry",
				"resilience4j":    "Resilience4j",
				"failsafe":        "Failsafe",
				"netflix.hystrix": "Hystrix",
				"guava":           "Guava (Retryer)",
				"bucket4j":        "Bucket4j",
			}
			for dep, name := range javaRetryLibs {
				if strings.Contains(strings.ToLower(string(content)), dep) {
					if !containsString(found, name) {
						found = append(found, name)
					}
				}
			}
			break
		}
	}

	// Check Rust dependencies
	if safepath.Exists(path, "Cargo.toml") {
		if content, err := safepath.ReadFile(path, "Cargo.toml"); err == nil {
			rustRetryLibs := map[string]string{
				"backoff":       "backoff",
				"retry":         "retry",
				"tokio-retry":   "tokio-retry",
				"again":         "again",
				"backon":        "backon",
				"reqwest-retry": "reqwest-retry",
				"failsafe":      "failsafe",
			}
			for dep, name := range rustRetryLibs {
				if strings.Contains(string(content), `"`+dep+`"`) || strings.Contains(string(content), dep+" =") {
					if !containsString(found, name) {
						found = append(found, name)
					}
				}
			}
		}
	}

	// If no library found, scan source code for custom retry implementations
	if len(found) == 0 {
		if c.hasCustomRetryCode(path) {
			found = append(found, "custom implementation")
		}
	}

	// Build result
	if len(found) > 0 {
		return rb.Pass("Retry/resilience: " + strings.Join(found, ", ")), nil
	}
	return rb.Warn("No retry logic found (consider adding for external calls)"), nil
}

// retryCodePatterns are code-level indicators of custom retry logic.
var retryCodePatterns = []string{
	"withRetry",
	"retryWithBackoff",
	"retryRequest",
	"retryable",
	"maxRetries",
	"max_retries",
	"retryCount",
	"retry_count",
	"retryDelay",
	"retry_delay",
	"backoffMs",
	"backoff_ms",
	"exponentialBackoff",
	"exponential_backoff",
	"@Retryable",
	"@retry",
}

func (c *RetryCheck) hasCustomRetryCode(path string) bool {
	// Scan source files across all ecosystems
	filePatterns := []string{
		// Go
		"*.go", "cmd/*.go", "cmd/*/*.go", "internal/*.go", "internal/*/*.go", "pkg/*.go", "pkg/*/*.go",
		// Node/TS
		"*.js", "*.ts", "src/*.js", "src/*.ts", "src/*/*.js", "src/*/*.ts", "src/*/*/*.js", "src/*/*/*.ts",
		"lib/*.js", "lib/*.ts", "lib/*/*.js", "lib/*/*.ts",
		// Python
		"*.py", "src/*.py", "src/*/*.py", "app/*.py", "app/*/*.py",
		// Java
		"src/main/java/*.java", "src/main/java/*/*.java", "src/main/java/*/*/*.java",
		"src/main/java/*/*/*/*.java", "src/main/java/*/*/*/*/*.java",
		// Rust
		"src/*.rs", "src/*/*.rs",
	}

	for _, pattern := range filePatterns {
		files, err := safepath.Glob(path, pattern)
		if err != nil {
			continue
		}
		for _, file := range files {
			baseName := filepath.Base(file)

			// Skip test files, vendored code, node_modules
			if strings.Contains(baseName, "test") || strings.Contains(baseName, "spec") ||
				strings.Contains(file, "node_modules") || strings.Contains(file, "vendor") ||
				strings.Contains(file, "venv") {
				continue
			}

			// Strong signal: file named retry or backoff
			nameNoExt := strings.TrimSuffix(baseName, filepath.Ext(baseName))
			if nameNoExt == "retry" || nameNoExt == "backoff" ||
				nameNoExt == "retrier" || nameNoExt == "retry_utils" || nameNoExt == "retryUtils" {
				return true
			}

			// Scan file contents for retry patterns
			content, err := safepath.ReadFileAbs(file)
			if err != nil {
				continue
			}
			text := string(content)
			for _, p := range retryCodePatterns {
				if strings.Contains(text, p) {
					return true
				}
			}
		}
	}

	return false
}
