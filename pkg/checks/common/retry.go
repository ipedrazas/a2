package common

import (
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

	// Build result
	if len(found) > 0 {
		return rb.Pass("Retry/resilience: " + strings.Join(found, ", ")), nil
	}
	return rb.Warn("No retry logic found (consider adding for external calls)"), nil
}
