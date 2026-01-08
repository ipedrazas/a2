package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type ErrorsCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *ErrorsCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "errors-test-*")
	s.Require().NoError(err)
}

func (s *ErrorsCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *ErrorsCheckTestSuite) TestIDAndName() {
	check := &ErrorsCheck{}
	s.Equal("common:errors", check.ID())
	s.Equal("Error Tracking", check.Name())
}

func (s *ErrorsCheckTestSuite) TestGoSentry() {
	content := `module myapp

go 1.21

require github.com/getsentry/sentry-go v0.25.0
`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ErrorsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Sentry")
}

func (s *ErrorsCheckTestSuite) TestGoRollbar() {
	content := `module myapp

go 1.21

require github.com/rollbar/rollbar-go v1.4.5
`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ErrorsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Rollbar")
}

func (s *ErrorsCheckTestSuite) TestGoBugsnag() {
	content := `module myapp

go 1.21

require github.com/bugsnag/bugsnag-go v2.2.0
`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ErrorsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Bugsnag")
}

func (s *ErrorsCheckTestSuite) TestPythonSentrySDK() {
	content := `[project]
name = "myapp"
dependencies = [
    "sentry-sdk>=1.30.0",
]
`
	err := os.WriteFile(filepath.Join(s.tempDir, "pyproject.toml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ErrorsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Sentry")
}

func (s *ErrorsCheckTestSuite) TestPythonRollbar() {
	content := `sentry-sdk
rollbar>=0.16.0
`
	err := os.WriteFile(filepath.Join(s.tempDir, "requirements.txt"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ErrorsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Rollbar")
}

func (s *ErrorsCheckTestSuite) TestNodeSentryNode() {
	content := `{
  "name": "myapp",
  "dependencies": {
    "express": "^4.18.0",
    "@sentry/node": "^7.80.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ErrorsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Sentry")
}

func (s *ErrorsCheckTestSuite) TestNodeBugsnag() {
	content := `{
  "name": "myapp",
  "dependencies": {
    "@bugsnag/js": "^7.22.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ErrorsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Bugsnag")
}

func (s *ErrorsCheckTestSuite) TestNodeHoneybadger() {
	content := `{
  "name": "myapp",
  "dependencies": {
    "@honeybadger-io/js": "^6.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ErrorsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Honeybadger")
}

func (s *ErrorsCheckTestSuite) TestJavaSentry() {
	content := `<project>
  <dependencies>
    <dependency>
      <groupId>io.sentry</groupId>
      <artifactId>sentry</artifactId>
      <version>6.30.0</version>
    </dependency>
  </dependencies>
</project>
`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ErrorsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Sentry")
}

func (s *ErrorsCheckTestSuite) TestJavaRollbar() {
	content := `plugins {
    id 'java'
}

dependencies {
    implementation 'com.rollbar:rollbar-java:1.10.0'
}
`
	err := os.WriteFile(filepath.Join(s.tempDir, "build.gradle"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ErrorsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Rollbar")
}

func (s *ErrorsCheckTestSuite) TestSentryCLIConfig() {
	content := `[defaults]
org=my-org
project=my-project
`
	err := os.WriteFile(filepath.Join(s.tempDir, ".sentryclirc"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ErrorsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Sentry CLI")
}

func (s *ErrorsCheckTestSuite) TestSentryProperties() {
	content := `defaults.url=https://sentry.io/
defaults.org=my-org
defaults.project=my-project
`
	err := os.WriteFile(filepath.Join(s.tempDir, "sentry.properties"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ErrorsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Sentry")
}

func (s *ErrorsCheckTestSuite) TestEnvExampleWithSentryDSN() {
	content := `# Application settings
DEBUG=false
PORT=8080

# Error tracking
SENTRY_DSN=
`
	err := os.WriteFile(filepath.Join(s.tempDir, ".env.example"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ErrorsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "error tracking env vars")
}

func (s *ErrorsCheckTestSuite) TestEnvExampleWithRollbarToken() {
	content := `DEBUG=false
ROLLBAR_TOKEN=
`
	err := os.WriteFile(filepath.Join(s.tempDir, ".env.sample"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ErrorsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "error tracking env vars")
}

func (s *ErrorsCheckTestSuite) TestMultipleErrorTrackers() {
	// Go module with both Sentry and Rollbar
	content := `module myapp

go 1.21

require (
    github.com/getsentry/sentry-go v0.25.0
    github.com/rollbar/rollbar-go v1.4.5
)
`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ErrorsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Sentry")
	s.Contains(result.Message, "Rollbar")
}

func (s *ErrorsCheckTestSuite) TestNoErrorTracking() {
	// Create a project without error tracking
	content := `module myapp

go 1.21

require github.com/gin-gonic/gin v1.9.0
`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &ErrorsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "No error tracking found")
}

func (s *ErrorsCheckTestSuite) TestEmptyDirectory() {
	check := &ErrorsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "No error tracking found")
}

func TestErrorsCheckTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorsCheckTestSuite))
}
