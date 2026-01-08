package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type MetricsCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *MetricsCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "metrics-test-*")
	s.Require().NoError(err)
}

func (s *MetricsCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *MetricsCheckTestSuite) TestIDAndName() {
	check := &MetricsCheck{}
	s.Equal("common:metrics", check.ID())
	s.Equal("Metrics Instrumentation", check.Name())
}

func (s *MetricsCheckTestSuite) TestGoPrometheus() {
	content := `module myapp

go 1.21

require github.com/prometheus/client_golang v1.17.0
`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MetricsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Prometheus")
}

func (s *MetricsCheckTestSuite) TestGoOpenTelemetry() {
	content := `module myapp

go 1.21

require go.opentelemetry.io/otel v1.19.0
require go.opentelemetry.io/otel/sdk v1.19.0
`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MetricsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "OpenTelemetry")
}

func (s *MetricsCheckTestSuite) TestGoDatadog() {
	content := `module myapp

go 1.21

require gopkg.in/DataDog/dd-trace-go.v1 v1.56.0
`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MetricsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Datadog")
}

func (s *MetricsCheckTestSuite) TestPythonPrometheus() {
	content := `[project]
name = "myapp"
dependencies = [
    "prometheus_client>=0.17.0",
]
`
	err := os.WriteFile(filepath.Join(s.tempDir, "pyproject.toml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MetricsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Prometheus")
}

func (s *MetricsCheckTestSuite) TestPythonStatsD() {
	content := `prometheus-client
statsd>=4.0.0
`
	err := os.WriteFile(filepath.Join(s.tempDir, "requirements.txt"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MetricsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "StatsD")
}

func (s *MetricsCheckTestSuite) TestNodePromClient() {
	content := `{
  "name": "myapp",
  "dependencies": {
    "express": "^4.18.0",
    "prom-client": "^15.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MetricsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Prometheus")
}

func (s *MetricsCheckTestSuite) TestNodeOpenTelemetry() {
	content := `{
  "name": "myapp",
  "dependencies": {
    "@opentelemetry/sdk-metrics": "^1.17.0",
    "@opentelemetry/api": "^1.6.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MetricsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "OpenTelemetry")
}

func (s *MetricsCheckTestSuite) TestNodeDatadog() {
	content := `{
  "name": "myapp",
  "dependencies": {
    "dd-trace": "^4.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MetricsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Datadog")
}

func (s *MetricsCheckTestSuite) TestJavaMicrometer() {
	content := `<project>
  <dependencies>
    <dependency>
      <groupId>io.micrometer</groupId>
      <artifactId>micrometer-core</artifactId>
      <version>1.12.0</version>
    </dependency>
  </dependencies>
</project>
`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MetricsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Micrometer")
}

func (s *MetricsCheckTestSuite) TestJavaDropwizardMetrics() {
	content := `plugins {
    id 'java'
}

dependencies {
    implementation 'io.dropwizard.metrics:metrics-core:4.2.0'
}
`
	err := os.WriteFile(filepath.Join(s.tempDir, "build.gradle"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MetricsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Dropwizard Metrics")
}

func (s *MetricsCheckTestSuite) TestPrometheusConfigFile() {
	content := `global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'myapp'
    static_configs:
      - targets: ['localhost:8080']
`
	err := os.WriteFile(filepath.Join(s.tempDir, "prometheus.yml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MetricsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "prometheus.yml")
}

func (s *MetricsCheckTestSuite) TestOtelCollectorConfig() {
	content := `receivers:
  otlp:
    protocols:
      grpc:
      http:

exporters:
  prometheus:
    endpoint: "0.0.0.0:8889"

service:
  pipelines:
    metrics:
      receivers: [otlp]
      exporters: [prometheus]
`
	err := os.WriteFile(filepath.Join(s.tempDir, "otel-collector-config.yaml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MetricsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "otel-collector-config.yaml")
}

func (s *MetricsCheckTestSuite) TestGrafanaDashboards() {
	grafanaDir := filepath.Join(s.tempDir, "grafana")
	err := os.MkdirAll(grafanaDir, 0755)
	s.Require().NoError(err)

	content := `{
  "title": "My Dashboard",
  "panels": []
}`
	err = os.WriteFile(filepath.Join(grafanaDir, "dashboard.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MetricsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Grafana dashboards")
}

func (s *MetricsCheckTestSuite) TestMultipleMetricsLibraries() {
	// Go module with both Prometheus and OpenTelemetry
	content := `module myapp

go 1.21

require (
    github.com/prometheus/client_golang v1.17.0
    go.opentelemetry.io/otel v1.19.0
)
`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MetricsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Prometheus")
	s.Contains(result.Message, "OpenTelemetry")
}

func (s *MetricsCheckTestSuite) TestNoMetrics() {
	// Create a project without metrics
	content := `module myapp

go 1.21

require github.com/gin-gonic/gin v1.9.0
`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &MetricsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "No metrics instrumentation found")
}

func (s *MetricsCheckTestSuite) TestEmptyDirectory() {
	check := &MetricsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "No metrics instrumentation found")
}

func TestMetricsCheckTestSuite(t *testing.T) {
	suite.Run(t, new(MetricsCheckTestSuite))
}
