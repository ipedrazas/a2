package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type HealthCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *HealthCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "health-test-*")
	s.Require().NoError(err)
}

func (s *HealthCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *HealthCheckTestSuite) TestIDAndName() {
	check := &HealthCheck{}
	s.Equal("common:health", check.ID())
	s.Equal("Health Endpoint", check.Name())
}

func (s *HealthCheckTestSuite) TestNoHealthEndpoint() {
	// Create a file without health patterns
	code := `package main

func main() {
	fmt.Println("Hello")
}
`
	err := os.WriteFile(filepath.Join(s.tempDir, "main.go"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &HealthCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "No health endpoint pattern detected")
}

func (s *HealthCheckTestSuite) TestHealthEndpointGo() {
	code := `package main

import "net/http"

func main() {
	http.HandleFunc("/health", healthHandler)
}
`
	err := os.WriteFile(filepath.Join(s.tempDir, "main.go"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &HealthCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "/health")
}

func (s *HealthCheckTestSuite) TestHealthzEndpoint() {
	// Use /readiness which doesn't overlap with other patterns
	code := `package main

func setupRoutes() {
	r.GET("/readiness", readinessHandler)
}
`
	err := os.WriteFile(filepath.Join(s.tempDir, "routes.go"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &HealthCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "/readiness")
}

func (s *HealthCheckTestSuite) TestReadinessEndpoint() {
	code := `from flask import Flask

app = Flask(__name__)

@app.route('/ready')
def ready():
    return 'OK'
`
	err := os.WriteFile(filepath.Join(s.tempDir, "app.py"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &HealthCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "/ready")
}

func (s *HealthCheckTestSuite) TestHealthCheckFunction() {
	code := `package main

func HealthCheck() error {
	return nil
}
`
	err := os.WriteFile(filepath.Join(s.tempDir, "health.go"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &HealthCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "HealthCheck")
}

func (s *HealthCheckTestSuite) TestHealthCheckInPython() {
	code := `def health_check():
    return {"status": "ok"}
`
	err := os.WriteFile(filepath.Join(s.tempDir, "health.py"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &HealthCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "health_check")
}

func (s *HealthCheckTestSuite) TestHealthInTypeScript() {
	code := `import express from 'express';

const app = express();

app.get('/health', (req, res) => {
  res.send('OK');
});
`
	err := os.WriteFile(filepath.Join(s.tempDir, "server.ts"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &HealthCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "/health")
}

func (s *HealthCheckTestSuite) TestSkipsNodeModules() {
	// Create node_modules with health endpoint (should be skipped)
	nodeModules := filepath.Join(s.tempDir, "node_modules", "some-lib")
	err := os.MkdirAll(nodeModules, 0755)
	s.Require().NoError(err)

	code := `app.get('/health', handler);`
	err = os.WriteFile(filepath.Join(nodeModules, "index.js"), []byte(code), 0644)
	s.Require().NoError(err)

	// Create a regular file without health
	mainCode := `console.log('hello');`
	err = os.WriteFile(filepath.Join(s.tempDir, "main.js"), []byte(mainCode), 0644)
	s.Require().NoError(err)

	check := &HealthCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
}

func (s *HealthCheckTestSuite) TestLivenessProbeInYAML() {
	code := `apiVersion: v1
kind: Pod
spec:
  containers:
  - name: app
    livenessProbe:
      httpGet:
        path: /healthz
`
	err := os.WriteFile(filepath.Join(s.tempDir, "pod.yaml"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &HealthCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
}

func (s *HealthCheckTestSuite) TestEmptyDirectory() {
	check := &HealthCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
}

func TestHealthCheckTestSuite(t *testing.T) {
	suite.Run(t, new(HealthCheckTestSuite))
}
