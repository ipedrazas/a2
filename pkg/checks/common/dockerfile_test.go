package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/stretchr/testify/suite"
)

type DockerfileCheckTestSuite struct {
	suite.Suite
	tempDir        string
	trivyInstalled bool
}

func (s *DockerfileCheckTestSuite) SetupSuite() {
	s.trivyInstalled = checkutil.ToolAvailable("trivy")
}

func (s *DockerfileCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "dockerfile-test-*")
	s.Require().NoError(err)
}

func (s *DockerfileCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *DockerfileCheckTestSuite) TestIDAndName() {
	check := &DockerfileCheck{}
	s.Equal("common:dockerfile", check.ID())
	s.Equal("Container Ready", check.Name())
}

func (s *DockerfileCheckTestSuite) TestNoDockerfile() {
	check := &DockerfileCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "No Dockerfile or Containerfile found")
}

// Tests that run when trivy IS installed
func (s *DockerfileCheckTestSuite) TestTrivyInstalled_ScansDockerfile() {
	if !s.trivyInstalled {
		s.T().Skip("trivy not installed")
	}

	// Create a simple Dockerfile
	dockerfile := `FROM alpine:latest
RUN apk add --no-cache curl
`
	err := os.WriteFile(filepath.Join(s.tempDir, "Dockerfile"), []byte(dockerfile), 0644)
	s.Require().NoError(err)

	check := &DockerfileCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	// Result should mention trivy
	s.Contains(result.Message, "trivy:")
	s.Contains(result.Message, "Dockerfile")
}

func (s *DockerfileCheckTestSuite) TestTrivyInstalled_WithDockerignore() {
	if !s.trivyInstalled {
		s.T().Skip("trivy not installed")
	}

	dockerfile := `FROM alpine:latest
`
	err := os.WriteFile(filepath.Join(s.tempDir, "Dockerfile"), []byte(dockerfile), 0644)
	s.Require().NoError(err)

	err = os.WriteFile(filepath.Join(s.tempDir, ".dockerignore"), []byte("node_modules\n"), 0644)
	s.Require().NoError(err)

	check := &DockerfileCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.Contains(result.Message, "trivy:")
	s.Contains(result.Message, ".dockerignore")
}

func (s *DockerfileCheckTestSuite) TestTrivyInstalled_FindsIssues() {
	if !s.trivyInstalled {
		s.T().Skip("trivy not installed")
	}

	// Create a Dockerfile with common issues (running as root, no healthcheck, etc.)
	dockerfile := `FROM ubuntu:latest
RUN apt-get update && apt-get install -y curl
`
	err := os.WriteFile(filepath.Join(s.tempDir, "Dockerfile"), []byte(dockerfile), 0644)
	s.Require().NoError(err)

	check := &DockerfileCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.Contains(result.Message, "trivy:")
	// This Dockerfile should have some issues (e.g., using latest tag, running as root)
	// But we don't assert on specific issues as trivy rules may change
}

// Tests for basic check when trivy is NOT installed
func (s *DockerfileCheckTestSuite) TestDockerfileExists() {
	if s.trivyInstalled {
		s.T().Skip("trivy installed - this test checks basic fallback")
	}

	err := os.WriteFile(filepath.Join(s.tempDir, "Dockerfile"), []byte("FROM alpine\n"), 0644)
	s.Require().NoError(err)

	check := &DockerfileCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Dockerfile found")
}

func (s *DockerfileCheckTestSuite) TestDockerfileLowercase() {
	if s.trivyInstalled {
		s.T().Skip("trivy installed - this test checks basic fallback")
	}

	err := os.WriteFile(filepath.Join(s.tempDir, "dockerfile"), []byte("FROM alpine\n"), 0644)
	s.Require().NoError(err)

	check := &DockerfileCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	// On case-insensitive filesystems, may match "Dockerfile" first
	s.Contains(result.Message, "found")
}

func (s *DockerfileCheckTestSuite) TestContainerfile() {
	if s.trivyInstalled {
		s.T().Skip("trivy installed - this test checks basic fallback")
	}

	err := os.WriteFile(filepath.Join(s.tempDir, "Containerfile"), []byte("FROM alpine\n"), 0644)
	s.Require().NoError(err)

	check := &DockerfileCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Containerfile found")
}

func (s *DockerfileCheckTestSuite) TestDockerfileWithIgnore() {
	if s.trivyInstalled {
		s.T().Skip("trivy installed - this test checks basic fallback")
	}

	err := os.WriteFile(filepath.Join(s.tempDir, "Dockerfile"), []byte("FROM alpine\n"), 0644)
	s.Require().NoError(err)

	err = os.WriteFile(filepath.Join(s.tempDir, ".dockerignore"), []byte("node_modules\n"), 0644)
	s.Require().NoError(err)

	check := &DockerfileCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, ".dockerignore")
}

func (s *DockerfileCheckTestSuite) TestDockerfileWithoutIgnore() {
	if s.trivyInstalled {
		s.T().Skip("trivy installed - this test checks basic fallback")
	}

	err := os.WriteFile(filepath.Join(s.tempDir, "Dockerfile"), []byte("FROM alpine\n"), 0644)
	s.Require().NoError(err)

	check := &DockerfileCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "consider adding .dockerignore")
}

func TestDockerfileCheckTestSuite(t *testing.T) {
	suite.Run(t, new(DockerfileCheckTestSuite))
}
