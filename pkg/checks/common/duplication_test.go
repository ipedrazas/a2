package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/stretchr/testify/suite"
)

type DuplicationCheckTestSuite struct {
	suite.Suite
	tempDir        string
	jscpdInstalled bool
}

func (s *DuplicationCheckTestSuite) SetupSuite() {
	s.jscpdInstalled = checkutil.ToolAvailable("jscpd")
}

func (s *DuplicationCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "duplication-test-*")
	s.Require().NoError(err)
}

func (s *DuplicationCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *DuplicationCheckTestSuite) TestIDAndName() {
	check := &DuplicationCheck{}
	s.Equal("common:duplication", check.ID())
	s.Equal("Code Duplication", check.Name())
}

func (s *DuplicationCheckTestSuite) TestNoToolNoConfig() {
	if s.jscpdInstalled {
		s.T().Skip("jscpd is installed - this test checks the fallback path")
	}

	check := &DuplicationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Reason, "jscpd")
}

func (s *DuplicationCheckTestSuite) TestConfigFileDetected_JscpdJSON() {
	if s.jscpdInstalled {
		s.T().Skip("jscpd is installed - this test checks config detection fallback")
	}

	err := os.WriteFile(filepath.Join(s.tempDir, ".jscpd.json"), []byte(`{"threshold": 5}`), 0644)
	s.Require().NoError(err)

	check := &DuplicationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "jscpd")
}

func (s *DuplicationCheckTestSuite) TestConfigFileDetected_SonarQube() {
	if s.jscpdInstalled {
		s.T().Skip("jscpd is installed - this test checks config detection fallback")
	}

	err := os.WriteFile(filepath.Join(s.tempDir, "sonar-project.properties"), []byte("sonar.projectKey=test\n"), 0644)
	s.Require().NoError(err)

	check := &DuplicationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "SonarQube")
}

func (s *DuplicationCheckTestSuite) TestJscpdInstalled_NoDuplication() {
	if !s.jscpdInstalled {
		s.T().Skip("jscpd not installed")
	}

	// Create a simple file with no duplication
	err := os.WriteFile(filepath.Join(s.tempDir, "main.go"), []byte(`package main

func main() {
	println("hello")
}
`), 0644)
	s.Require().NoError(err)

	check := &DuplicationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	// Should pass or at least not fail
	s.NotEqual(checker.Fail, result.Status)
}

func (s *DuplicationCheckTestSuite) TestParseOutput_Valid() {
	check := &DuplicationCheck{}

	jsonOutput := `{"statistics":{"clones":3,"duplicatedLines":45,"sources":10,"percentage":12.5}}`
	stats, err := check.parseOutput(jsonOutput)

	s.NoError(err)
	s.Equal(3, stats.Clones)
	s.Equal(45, stats.Duplicates)
	s.Equal(10, stats.Sources)
	s.InDelta(12.5, stats.Percentage, 0.01)
}

func (s *DuplicationCheckTestSuite) TestParseOutput_Empty() {
	check := &DuplicationCheck{}

	_, err := check.parseOutput("")
	s.Error(err)
}

func (s *DuplicationCheckTestSuite) TestParseOutput_InvalidJSON() {
	check := &DuplicationCheck{}

	_, err := check.parseOutput("not json")
	s.Error(err)
}

func (s *DuplicationCheckTestSuite) TestParseOutput_NoDuplication() {
	check := &DuplicationCheck{}

	jsonOutput := `{"statistics":{"clones":0,"duplicatedLines":0,"sources":5,"percentage":0}}`
	stats, err := check.parseOutput(jsonOutput)

	s.NoError(err)
	s.Equal(0, stats.Clones)
	s.InDelta(0.0, stats.Percentage, 0.01)
}

func TestDuplicationCheckTestSuite(t *testing.T) {
	suite.Run(t, new(DuplicationCheckTestSuite))
}
