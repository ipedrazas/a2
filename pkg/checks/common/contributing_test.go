package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type ContributingCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *ContributingCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "contributing-test-*")
	s.Require().NoError(err)
}

func (s *ContributingCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *ContributingCheckTestSuite) TestIDAndName() {
	check := &ContributingCheck{}
	s.Equal("common:contributing", check.ID())
	s.Equal("Contributing Guidelines", check.Name())
}

func (s *ContributingCheckTestSuite) TestContributingMd() {
	err := os.WriteFile(filepath.Join(s.tempDir, "CONTRIBUTING.md"), []byte(`
# Contributing

Thank you for contributing!
`), 0644)
	s.Require().NoError(err)

	check := &ContributingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "CONTRIBUTING.md")
}

func (s *ContributingCheckTestSuite) TestContributingTxt() {
	err := os.WriteFile(filepath.Join(s.tempDir, "CONTRIBUTING.txt"), []byte(`
Contributing guidelines here.
`), 0644)
	s.Require().NoError(err)

	check := &ContributingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "CONTRIBUTING.txt")
}

func (s *ContributingCheckTestSuite) TestContributingInGithub() {
	githubDir := filepath.Join(s.tempDir, ".github")
	err := os.MkdirAll(githubDir, 0755)
	s.Require().NoError(err)

	err = os.WriteFile(filepath.Join(githubDir, "CONTRIBUTING.md"), []byte(`
# How to Contribute
`), 0644)
	s.Require().NoError(err)

	check := &ContributingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, ".github/CONTRIBUTING.md")
}

func (s *ContributingCheckTestSuite) TestPRTemplate() {
	githubDir := filepath.Join(s.tempDir, ".github")
	err := os.MkdirAll(githubDir, 0755)
	s.Require().NoError(err)

	err = os.WriteFile(filepath.Join(githubDir, "PULL_REQUEST_TEMPLATE.md"), []byte(`
## Description
`), 0644)
	s.Require().NoError(err)

	check := &ContributingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "PR template")
}

func (s *ContributingCheckTestSuite) TestPRTemplateLowercase() {
	githubDir := filepath.Join(s.tempDir, ".github")
	err := os.MkdirAll(githubDir, 0755)
	s.Require().NoError(err)

	err = os.WriteFile(filepath.Join(githubDir, "pull_request_template.md"), []byte(`
## Description
`), 0644)
	s.Require().NoError(err)

	check := &ContributingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "PR template")
}

func (s *ContributingCheckTestSuite) TestIssueTemplateFile() {
	githubDir := filepath.Join(s.tempDir, ".github")
	err := os.MkdirAll(githubDir, 0755)
	s.Require().NoError(err)

	err = os.WriteFile(filepath.Join(githubDir, "ISSUE_TEMPLATE.md"), []byte(`
## Bug Report
`), 0644)
	s.Require().NoError(err)

	check := &ContributingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "issue templates")
}

func (s *ContributingCheckTestSuite) TestIssueTemplateDirectory() {
	issueDir := filepath.Join(s.tempDir, ".github", "ISSUE_TEMPLATE")
	err := os.MkdirAll(issueDir, 0755)
	s.Require().NoError(err)

	err = os.WriteFile(filepath.Join(issueDir, "bug_report.md"), []byte(`
name: Bug Report
`), 0644)
	s.Require().NoError(err)

	check := &ContributingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "issue templates")
}

func (s *ContributingCheckTestSuite) TestCodeowners() {
	err := os.WriteFile(filepath.Join(s.tempDir, "CODEOWNERS"), []byte(`
* @team-lead
`), 0644)
	s.Require().NoError(err)

	check := &ContributingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "CODEOWNERS")
}

func (s *ContributingCheckTestSuite) TestCodeownersInGithub() {
	githubDir := filepath.Join(s.tempDir, ".github")
	err := os.MkdirAll(githubDir, 0755)
	s.Require().NoError(err)

	err = os.WriteFile(filepath.Join(githubDir, "CODEOWNERS"), []byte(`
*.go @go-team
`), 0644)
	s.Require().NoError(err)

	check := &ContributingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "CODEOWNERS")
}

func (s *ContributingCheckTestSuite) TestMultipleFound() {
	// Create CONTRIBUTING.md
	err := os.WriteFile(filepath.Join(s.tempDir, "CONTRIBUTING.md"), []byte(`# Contributing`), 0644)
	s.Require().NoError(err)

	// Create CODEOWNERS
	err = os.WriteFile(filepath.Join(s.tempDir, "CODEOWNERS"), []byte(`* @owner`), 0644)
	s.Require().NoError(err)

	// Create .github directory with templates
	githubDir := filepath.Join(s.tempDir, ".github")
	err = os.MkdirAll(githubDir, 0755)
	s.Require().NoError(err)

	err = os.WriteFile(filepath.Join(githubDir, "PULL_REQUEST_TEMPLATE.md"), []byte(`## PR`), 0644)
	s.Require().NoError(err)

	check := &ContributingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "CONTRIBUTING.md")
	s.Contains(result.Reason, "CODEOWNERS")
	s.Contains(result.Reason, "PR template")
}

func (s *ContributingCheckTestSuite) TestNoContributingFound() {
	// Create some other files
	err := os.WriteFile(filepath.Join(s.tempDir, "README.md"), []byte(`# Project`), 0644)
	s.Require().NoError(err)

	check := &ContributingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "No contribution guidelines found")
}

func (s *ContributingCheckTestSuite) TestEmptyDirectory() {
	check := &ContributingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "No contribution guidelines found")
}

func TestContributingCheckTestSuite(t *testing.T) {
	suite.Run(t, new(ContributingCheckTestSuite))
}
