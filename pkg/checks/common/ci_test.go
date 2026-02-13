package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type CICheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *CICheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "ci-test-*")
	s.Require().NoError(err)
}

func (s *CICheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *CICheckTestSuite) TestIDAndName() {
	check := &CICheck{}
	s.Equal("common:ci", check.ID())
	s.Equal("CI Pipeline", check.Name())
}

func (s *CICheckTestSuite) TestNoCIConfig() {
	check := &CICheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "No CI/CD configuration found")
}

func (s *CICheckTestSuite) TestGitHubActions() {
	// Create GitHub Actions workflow
	workflowDir := filepath.Join(s.tempDir, ".github", "workflows")
	err := os.MkdirAll(workflowDir, 0755)
	s.Require().NoError(err)

	err = os.WriteFile(filepath.Join(workflowDir, "ci.yml"), []byte("name: CI\n"), 0644)
	s.Require().NoError(err)

	check := &CICheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "GitHub Actions")
}

func (s *CICheckTestSuite) TestGitLabCI() {
	err := os.WriteFile(filepath.Join(s.tempDir, ".gitlab-ci.yml"), []byte("stages:\n"), 0644)
	s.Require().NoError(err)

	check := &CICheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "GitLab CI")
}

func (s *CICheckTestSuite) TestJenkinsfile() {
	err := os.WriteFile(filepath.Join(s.tempDir, "Jenkinsfile"), []byte("pipeline {\n}\n"), 0644)
	s.Require().NoError(err)

	check := &CICheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Jenkins")
}

func (s *CICheckTestSuite) TestCircleCI() {
	circleDir := filepath.Join(s.tempDir, ".circleci")
	err := os.MkdirAll(circleDir, 0755)
	s.Require().NoError(err)

	err = os.WriteFile(filepath.Join(circleDir, "config.yml"), []byte("version: 2.1\n"), 0644)
	s.Require().NoError(err)

	check := &CICheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "CircleCI")
}

func (s *CICheckTestSuite) TestTravisCI() {
	err := os.WriteFile(filepath.Join(s.tempDir, ".travis.yml"), []byte("language: go\n"), 0644)
	s.Require().NoError(err)

	check := &CICheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Travis CI")
}

func (s *CICheckTestSuite) TestAzurePipelines() {
	err := os.WriteFile(filepath.Join(s.tempDir, "azure-pipelines.yml"), []byte("trigger:\n"), 0644)
	s.Require().NoError(err)

	check := &CICheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Azure Pipelines")
}

func (s *CICheckTestSuite) TestMultipleCIConfigs() {
	// Create both GitLab CI and Jenkinsfile
	err := os.WriteFile(filepath.Join(s.tempDir, ".gitlab-ci.yml"), []byte("stages:\n"), 0644)
	s.Require().NoError(err)

	err = os.WriteFile(filepath.Join(s.tempDir, "Jenkinsfile"), []byte("pipeline {}\n"), 0644)
	s.Require().NoError(err)

	check := &CICheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	// Should list both
	s.Contains(result.Reason, "GitLab CI")
	s.Contains(result.Reason, "Jenkins")
}

func (s *CICheckTestSuite) TestTaskfile() {
	err := os.WriteFile(filepath.Join(s.tempDir, "Taskfile.yml"), []byte("version: 3\n"), 0644)
	s.Require().NoError(err)

	check := &CICheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Taskfile")
}

func TestCICheckTestSuite(t *testing.T) {
	suite.Run(t, new(CICheckTestSuite))
}
