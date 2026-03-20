package nodecheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type DeadcodeCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *DeadcodeCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "deadcode-node-test-*")
	s.Require().NoError(err)
}

func (s *DeadcodeCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *DeadcodeCheckTestSuite) TestIDAndName() {
	check := &DeadcodeCheck{}
	s.Equal("node:deadcode", check.ID())
	s.Equal("Node Dead Code", check.Name())
}

func (s *DeadcodeCheckTestSuite) TestNoPackageJSON() {
	check := &DeadcodeCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Reason, "package.json")
}

func (s *DeadcodeCheckTestSuite) TestCountKnipIssues() {
	s.Equal(0, countKnipIssues(""))
	// Section headers should not be counted
	s.Equal(0, countKnipIssues("Unused files (2)\n---"))
	// File paths should be counted
	s.Equal(2, countKnipIssues("Unused files (2)\nsrc/foo.ts\nsrc/bar.ts\n"))
}

func (s *DeadcodeCheckTestSuite) TestHasKnipDep() {
	// Create package.json with knip as devDependency
	pkg := `{"name":"test","devDependencies":{"knip":"^5.0.0"}}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(pkg), 0644)
	s.Require().NoError(err)

	check := &DeadcodeCheck{}
	s.True(check.hasKnipDep(s.tempDir))
}

func (s *DeadcodeCheckTestSuite) TestHasKnipDep_NotPresent() {
	pkg := `{"name":"test","devDependencies":{"eslint":"^8.0.0"}}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(pkg), 0644)
	s.Require().NoError(err)

	check := &DeadcodeCheck{}
	s.False(check.hasKnipDep(s.tempDir))
}

func TestDeadcodeCheckTestSuite(t *testing.T) {
	suite.Run(t, new(DeadcodeCheckTestSuite))
}
