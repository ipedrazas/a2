package common

import (
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/safepath"
)

// MetricsCheck verifies that metrics instrumentation exists.
type MetricsCheck struct{}

func (c *MetricsCheck) ID() string   { return "common:metrics" }
func (c *MetricsCheck) Name() string { return "Metrics Instrumentation" }

// Run checks for metrics libraries and configuration.
func (c *MetricsCheck) Run(path string) (checker.Result, error) {
	result := checker.Result{
		Name:     c.Name(),
		ID:       c.ID(),
		Language: checker.LangCommon,
	}

	var found []string

	// Check for Go metrics libraries
	if safepath.Exists(path, "go.mod") {
		if content, err := safepath.ReadFile(path, "go.mod"); err == nil {
			goMetricsLibs := []struct {
				pattern string
				name    string
			}{
				{"prometheus/client_golang", "Prometheus"},
				{"go-metrics", "go-metrics"},
				{"go.opentelemetry.io/otel", "OpenTelemetry"},
				{"DataDog/dd-trace-go", "Datadog"},
				{"newrelic/go-agent", "New Relic"},
				{"go.opencensus.io", "OpenCensus"},
			}
			for _, lib := range goMetricsLibs {
				if strings.Contains(string(content), lib.pattern) {
					found = append(found, lib.name)
				}
			}
		}
	}

	// Check for Python metrics libraries
	pythonConfigs := []string{"pyproject.toml", "requirements.txt", "setup.py"}
	for _, cfg := range pythonConfigs {
		if safepath.Exists(path, cfg) {
			if content, err := safepath.ReadFile(path, cfg); err == nil {
				contentLower := strings.ToLower(string(content))
				pythonMetricsLibs := []struct {
					pattern string
					name    string
				}{
					{"prometheus_client", "Prometheus"},
					{"prometheus-client", "Prometheus"},
					{"statsd", "StatsD"},
					{"opentelemetry", "OpenTelemetry"},
					{"ddtrace", "Datadog"},
					{"newrelic", "New Relic"},
				}
				for _, lib := range pythonMetricsLibs {
					if strings.Contains(contentLower, lib.pattern) {
						found = append(found, lib.name)
					}
				}
			}
			break
		}
	}

	// Check for Node.js metrics libraries
	if safepath.Exists(path, "package.json") {
		if content, err := safepath.ReadFile(path, "package.json"); err == nil {
			nodeMetricsLibs := []struct {
				pattern string
				name    string
			}{
				{"prom-client", "Prometheus"},
				{"hot-shots", "StatsD"},
				{"@opentelemetry/sdk-metrics", "OpenTelemetry"},
				{"dd-trace", "Datadog"},
				{"newrelic", "New Relic"},
				{"appmetrics", "appmetrics"},
			}
			for _, lib := range nodeMetricsLibs {
				if strings.Contains(string(content), lib.pattern) {
					found = append(found, lib.name)
				}
			}
		}
	}

	// Check for Java metrics libraries
	javaConfigs := []string{"pom.xml", "build.gradle", "build.gradle.kts"}
	for _, cfg := range javaConfigs {
		if safepath.Exists(path, cfg) {
			if content, err := safepath.ReadFile(path, cfg); err == nil {
				contentLower := strings.ToLower(string(content))
				javaMetricsLibs := []struct {
					pattern string
					name    string
				}{
					{"micrometer", "Micrometer"},
					{"prometheus", "Prometheus"},
					{"dropwizard.metrics", "Dropwizard Metrics"},
					{"metrics-core", "Dropwizard Metrics"},
					{"opentelemetry", "OpenTelemetry"},
					{"dd-trace", "Datadog"},
				}
				for _, lib := range javaMetricsLibs {
					if strings.Contains(contentLower, lib.pattern) {
						found = append(found, lib.name)
					}
				}
			}
			break
		}
	}

	// Check for metrics configuration files
	metricsConfigs := []string{
		"prometheus.yml",
		"prometheus.yaml",
		"metrics.yml",
		"metrics.yaml",
		"otel-collector-config.yaml",
		"otel-collector-config.yml",
	}
	for _, cfg := range metricsConfigs {
		if safepath.Exists(path, cfg) {
			found = append(found, cfg)
		}
	}

	// Check for Grafana dashboards
	if safepath.IsDir(path, "grafana") {
		if files, err := safepath.Glob(path+"/grafana", "*.json"); err == nil && len(files) > 0 {
			found = append(found, "Grafana dashboards")
		}
	}
	if safepath.IsDir(path, "dashboards") {
		if files, err := safepath.Glob(path+"/dashboards", "*.json"); err == nil && len(files) > 0 {
			found = append(found, "Grafana dashboards")
		}
	}

	// Check for /metrics endpoint in code (common pattern)
	if c.hasMetricsEndpoint(path) {
		found = append(found, "/metrics endpoint")
	}

	// Build result
	found = unique(found)
	if len(found) > 0 {
		result.Passed = true
		result.Status = checker.Pass
		result.Message = "Metrics instrumentation found: " + strings.Join(found, ", ")
	} else {
		result.Passed = false
		result.Status = checker.Warn
		result.Message = "No metrics instrumentation found (consider adding Prometheus or OpenTelemetry)"
	}

	return result, nil
}

func (c *MetricsCheck) hasMetricsEndpoint(path string) bool {
	// Common patterns for metrics endpoints
	patterns := []string{
		`"/metrics"`,
		`'/metrics'`,
		`/metrics`,
		`metricsHandler`,
		`promhttp.Handler`,
		`metrics_endpoint`,
	}

	// Check Go files
	if files, err := safepath.Glob(path, "**/*.go"); err == nil {
		for _, file := range files {
			if content, err := safepath.ReadFileAbs(file); err == nil {
				for _, pattern := range patterns {
					if strings.Contains(string(content), pattern) {
						return true
					}
				}
			}
		}
	}

	// Check Python files
	if files, err := safepath.Glob(path, "**/*.py"); err == nil {
		for _, file := range files {
			if content, err := safepath.ReadFileAbs(file); err == nil {
				for _, pattern := range patterns {
					if strings.Contains(string(content), pattern) {
						return true
					}
				}
			}
		}
	}

	// Check Node.js/TypeScript files
	nodeExts := []string{"*.js", "*.ts"}
	for _, ext := range nodeExts {
		if files, err := safepath.Glob(path, "**/"+ext); err == nil {
			for _, file := range files {
				// Skip node_modules
				if strings.Contains(file, "node_modules") {
					continue
				}
				if content, err := safepath.ReadFileAbs(file); err == nil {
					for _, pattern := range patterns {
						if strings.Contains(string(content), pattern) {
							return true
						}
					}
				}
			}
		}
	}

	return false
}
