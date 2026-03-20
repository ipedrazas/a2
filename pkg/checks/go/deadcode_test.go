package gocheck

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
	tempDir           string
	deadcodeInstalled bool
}

func (s *DeadcodeCheckTestSuite) SetupSuite() {
	s.deadcodeInstalled = checkutil.ToolAvailable("deadcode")
}

func (s *DeadcodeCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "deadcode-go-test-*")
	s.Require().NoError(err)
}

func (s *DeadcodeCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *DeadcodeCheckTestSuite) TestIDAndName() {
	check := &DeadcodeCheck{}
	s.Equal("go:deadcode", check.ID())
	s.Equal("Go Dead Code", check.Name())
}

func (s *DeadcodeCheckTestSuite) TestToolNotInstalled() {
	if s.deadcodeInstalled {
		s.T().Skip("deadcode is installed")
	}

	check := &DeadcodeCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Reason, "deadcode")
}

func (s *DeadcodeCheckTestSuite) TestCleanCode() {
	if !s.deadcodeInstalled {
		s.T().Skip("deadcode not installed")
	}

	// Create a simple Go project with no dead code
	goMod := "module test\n\ngo 1.21\n"
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(goMod), 0644)
	s.Require().NoError(err)

	mainGo := `package main

func main() {
	println("hello")
}
`
	err = os.WriteFile(filepath.Join(s.tempDir, "main.go"), []byte(mainGo), 0644)
	s.Require().NoError(err)

	check := &DeadcodeCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.NotEqual(checker.Fail, result.Status)
}

func (s *DeadcodeCheckTestSuite) TestCountNonEmptyLines() {
	s.Equal(0, countNonEmptyLines(""))
	s.Equal(0, countNonEmptyLines("  \n  \n"))
	s.Equal(2, countNonEmptyLines("foo\nbar\n"))
	s.Equal(1, countNonEmptyLines("  foo  \n\n"))
}

func TestDeadcodeCheckTestSuite(t *testing.T) {
	suite.Run(t, new(DeadcodeCheckTestSuite))
}
