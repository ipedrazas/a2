package common

import (
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// TracingCheck verifies distributed tracing is implemented.
type TracingCheck struct{}

func (c *TracingCheck) ID() string   { return "common:tracing" }
func (c *TracingCheck) Name() string { return "Distributed Tracing" }

// Run checks for tracing configuration and libraries.
func (c *TracingCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangCommon)

	var found []string

	// Check for OpenTelemetry config files
	otelConfigs := []string{
		"otel-collector-config.yaml",
		"otel-collector-config.yml",
		"opentelemetry-config.yaml",
		"collector-config.yaml",
	}
	for _, cfg := range otelConfigs {
		if safepath.Exists(path, cfg) {
			found = append(found, "OpenTelemetry Collector")
			break
		}
	}

	// Check for Jaeger config
	if safepath.Exists(path, "jaeger-config.yaml") || safepath.Exists(path, "jaeger.yaml") {
		found = append(found, "Jaeger")
	}

	// Check Go dependencies
	if safepath.Exists(path, "go.mod") {
		if content, err := safepath.ReadFile(path, "go.mod"); err == nil {
			goTracing := map[string]string{
				"go.opentelemetry.io/otel":        "OpenTelemetry",
				"github.com/jaegertracing/":       "Jaeger",
				"github.com/openzipkin/zipkin-go": "Zipkin",
				"gopkg.in/DataDog/dd-trace-go":    "Datadog",
				"github.com/DataDog/dd-trace-go":  "Datadog",
				"go.elastic.co/apm":               "Elastic APM",
				"github.com/newrelic/go-agent":    "New Relic",
				"github.com/lightstep/":           "Lightstep",
			}
			for dep, name := range goTracing {
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
			nodeTracing := map[string]string{
				"@opentelemetry/":                 "OpenTelemetry",
				"dd-trace":                        "Datadog",
				"jaeger-client":                   "Jaeger",
				"zipkin":                          "Zipkin",
				"@elastic/apm-rum":                "Elastic APM",
				"elastic-apm-node":                "Elastic APM",
				"newrelic":                        "New Relic",
				"@sentry/tracing":                 "Sentry",
				"lightstep-tracer":                "Lightstep",
				"@honeycombio/opentelemetry-node": "Honeycomb",
			}
			for dep, name := range nodeTracing {
				if strings.Contains(string(content), `"`+dep) {
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
			pythonTracing := map[string]string{
				"opentelemetry":     "OpenTelemetry",
				"ddtrace":           "Datadog",
				"jaeger-client":     "Jaeger",
				"py-zipkin":         "Zipkin",
				"elastic-apm":       "Elastic APM",
				"newrelic":          "New Relic",
				"sentry-sdk":        "Sentry",
				"honeycomb-beeline": "Honeycomb",
			}
			for dep, name := range pythonTracing {
				if strings.Contains(strings.ToLower(string(content)), dep) {
					if !containsString(found, name) {
						found = append(found, name)
					}
				}
			}
			break
		}
	}

	// Check Java dependencies (pom.xml, build.gradle)
	javaFiles := []string{"pom.xml", "build.gradle", "build.gradle.kts"}
	for _, file := range javaFiles {
		if content, err := safepath.ReadFile(path, file); err == nil {
			javaTracing := map[string]string{
				"opentelemetry":       "OpenTelemetry",
				"io.jaegertracing":    "Jaeger",
				"io.zipkin":           "Zipkin",
				"dd-trace":            "Datadog",
				"elastic-apm-agent":   "Elastic APM",
				"newrelic":            "New Relic",
				"spring-cloud-sleuth": "Spring Sleuth",
				"micrometer-tracing":  "Micrometer",
			}
			for dep, name := range javaTracing {
				if strings.Contains(string(content), dep) {
					if !containsString(found, name) {
						found = append(found, name)
					}
				}
			}
			break
		}
	}

	// Check for tracing in CI/CD configs
	ciFiles := []string{
		".github/workflows/ci.yml",
		".github/workflows/ci.yaml",
		".gitlab-ci.yml",
	}
	for _, ciFile := range ciFiles {
		if content, err := safepath.ReadFile(path, ciFile); err == nil {
			if strings.Contains(strings.ToLower(string(content)), "otel") ||
				strings.Contains(strings.ToLower(string(content)), "jaeger") ||
				strings.Contains(strings.ToLower(string(content)), "tracing") {
				if !containsString(found, "CI tracing") {
					found = append(found, "CI tracing")
				}
			}
		}
	}

	// Build result
	if len(found) > 0 {
		return rb.Pass("Tracing configured: " + strings.Join(found, ", ")), nil
	}
	return rb.Warn("No distributed tracing found (consider OpenTelemetry)"), nil
}
