package pythoncheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/stretchr/testify/suite"
)

type DeadcodeCheckTestSuite struct {
	suite.Suite
	tempDir          string
	vultureInstalled bool
}

func (s *DeadcodeCheckTestSuite) SetupSuite() {
	s.vultureInstalled = checkutil.ToolAvailable("vulture")
}

func (s *DeadcodeCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "deadcode-python-test-*")
	s.Require().NoError(err)
}

func (s *DeadcodeCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *DeadcodeCheckTestSuite) TestIDAndName() {
	check := &DeadcodeCheck{}
	s.Equal("python:deadcode", check.ID())
	s.Equal("Python Dead Code", check.Name())
}

func (s *DeadcodeCheckTestSuite) TestNoPythonProject() {
	check := &DeadcodeCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
}

func (s *DeadcodeCheckTestSuite) TestToolNotInstalled() {
	if s.vultureInstalled {
		s.T().Skip("vulture is installed")
	}

	err := os.WriteFile(filepath.Join(s.tempDir, "pyproject.toml"), []byte("[project]\nname = \"test\"\n"), 0644)
	s.Require().NoError(err)

	check := &DeadcodeCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Reason, "vulture")
}

func (s *DeadcodeCheckTestSuite) TestCountFindings() {
	s.Equal(0, countFindings(""))
	s.Equal(0, countFindings("some random text\n"))
	s.Equal(1, countFindings("main.py:5: unused function 'foo' (60% confidence)\n"))
	s.Equal(2, countFindings("a.py:1: unused import 'os' (90% confidence)\nb.py:3: unused variable 'x' (60% confidence)\n"))
}

func TestDeadcodeCheckTestSuite(t *testing.T) {
	suite.Run(t, new(DeadcodeCheckTestSuite))
}
