package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type DockerfileCheckTestSuite struct {
	suite.Suite
	tempDir string
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

func (s *DockerfileCheckTestSuite) TestDockerfileExists() {
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
