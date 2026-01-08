package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type TracingCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *TracingCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "tracing-test-*")
	s.Require().NoError(err)
}

func (s *TracingCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *TracingCheckTestSuite) TestIDAndName() {
	check := &TracingCheck{}
	s.Equal("common:tracing", check.ID())
	s.Equal("Distributed Tracing", check.Name())
}

func (s *TracingCheckTestSuite) TestOtelCollectorConfig() {
	err := os.WriteFile(filepath.Join(s.tempDir, "otel-collector-config.yaml"), []byte(`
receivers:
  otlp:
    protocols:
      grpc:
      http:
`), 0644)
	s.Require().NoError(err)

	check := &TracingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "OpenTelemetry Collector")
}

func (s *TracingCheckTestSuite) TestJaegerConfig() {
	err := os.WriteFile(filepath.Join(s.tempDir, "jaeger-config.yaml"), []byte(`
sampling:
  default: true
`), 0644)
	s.Require().NoError(err)

	check := &TracingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Jaeger")
}

func (s *TracingCheckTestSuite) TestGoOpenTelemetry() {
	content := `module myapp

go 1.21

require (
	go.opentelemetry.io/otel v1.21.0
	go.opentelemetry.io/otel/sdk v1.21.0
)`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &TracingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "OpenTelemetry")
}

func (s *TracingCheckTestSuite) TestGoJaeger() {
	content := `module myapp

go 1.21

require (
	github.com/jaegertracing/jaeger-client-go v2.30.0
)`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &TracingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Jaeger")
}

func (s *TracingCheckTestSuite) TestGoDatadog() {
	content := `module myapp

go 1.21

require (
	github.com/DataDog/dd-trace-go v1.0.0
)`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &TracingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Datadog")
}

func (s *TracingCheckTestSuite) TestGoZipkin() {
	content := `module myapp

go 1.21

require (
	github.com/openzipkin/zipkin-go v0.4.0
)`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &TracingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Zipkin")
}

func (s *TracingCheckTestSuite) TestNodeOpenTelemetry() {
	content := `{
  "name": "my-app",
  "dependencies": {
    "@opentelemetry/sdk-node": "^0.45.0",
    "@opentelemetry/auto-instrumentations-node": "^0.40.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &TracingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "OpenTelemetry")
}

func (s *TracingCheckTestSuite) TestNodeDatadog() {
	content := `{
  "name": "my-app",
  "dependencies": {
    "dd-trace": "^4.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &TracingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Datadog")
}

func (s *TracingCheckTestSuite) TestNodeSentry() {
	content := `{
  "name": "my-app",
  "dependencies": {
    "@sentry/tracing": "^7.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &TracingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Sentry")
}

func (s *TracingCheckTestSuite) TestNodeElasticAPM() {
	content := `{
  "name": "my-app",
  "dependencies": {
    "elastic-apm-node": "^4.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &TracingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Elastic APM")
}

func (s *TracingCheckTestSuite) TestPythonOpenTelemetry() {
	content := `opentelemetry-sdk==1.21.0
opentelemetry-instrumentation==0.42b0`
	err := os.WriteFile(filepath.Join(s.tempDir, "requirements.txt"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &TracingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "OpenTelemetry")
}

func (s *TracingCheckTestSuite) TestPythonDatadog() {
	content := `ddtrace==2.0.0`
	err := os.WriteFile(filepath.Join(s.tempDir, "requirements.txt"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &TracingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Datadog")
}

func (s *TracingCheckTestSuite) TestPythonSentry() {
	content := `sentry-sdk==1.35.0`
	err := os.WriteFile(filepath.Join(s.tempDir, "requirements.txt"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &TracingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Sentry")
}

func (s *TracingCheckTestSuite) TestJavaOpenTelemetry() {
	content := `<dependency>
    <groupId>io.opentelemetry</groupId>
    <artifactId>opentelemetry-api</artifactId>
</dependency>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &TracingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "OpenTelemetry")
}

func (s *TracingCheckTestSuite) TestJavaSpringSleuth() {
	content := `<dependency>
    <groupId>org.springframework.cloud</groupId>
    <artifactId>spring-cloud-sleuth</artifactId>
</dependency>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &TracingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Spring Sleuth")
}

func (s *TracingCheckTestSuite) TestJavaMicrometer() {
	content := `implementation 'io.micrometer:micrometer-tracing:1.2.0'`
	err := os.WriteFile(filepath.Join(s.tempDir, "build.gradle"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &TracingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Micrometer")
}

func (s *TracingCheckTestSuite) TestCITracing() {
	githubDir := filepath.Join(s.tempDir, ".github", "workflows")
	err := os.MkdirAll(githubDir, 0755)
	s.Require().NoError(err)

	content := `name: CI
on: [push]
jobs:
  test:
    runs-on: ubuntu-latest
    env:
      OTEL_EXPORTER_OTLP_ENDPOINT: http://localhost:4318`
	err = os.WriteFile(filepath.Join(githubDir, "ci.yml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &TracingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "CI tracing")
}

func (s *TracingCheckTestSuite) TestNoTracingFound() {
	content := `module myapp

go 1.21

require (
	github.com/gin-gonic/gin v1.9.0
)`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &TracingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "No distributed tracing found")
}

func (s *TracingCheckTestSuite) TestEmptyDirectory() {
	check := &TracingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "No distributed tracing found")
}

func TestTracingCheckTestSuite(t *testing.T) {
	suite.Run(t, new(TracingCheckTestSuite))
}
