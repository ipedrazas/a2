package rustcheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type CoverageTestSuite struct {
	suite.Suite
	tempDir string
	check   *CoverageCheck
}

func (s *CoverageTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "rust-coverage-test-*")
	s.NoError(err)
	s.tempDir = dir
	s.check = &CoverageCheck{Threshold: 80.0}
}

func (s *CoverageTestSuite) TearDownTest() {
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

func (s *CoverageTestSuite) writeFile(name, content string) {
	path := filepath.Join(s.tempDir, name)
	dir := filepath.Dir(path)
	if dir != s.tempDir {
		err := os.MkdirAll(dir, 0755)
		s.NoError(err)
	}
	err := os.WriteFile(path, []byte(content), 0644)
	s.NoError(err)
}

func (s *CoverageTestSuite) TestID() {
	s.Equal("rust:coverage", s.check.ID())
}

func (s *CoverageTestSuite) TestName() {
	s.Equal("Rust Coverage", s.check.Name())
}

func (s *CoverageTestSuite) TestRun_NoCargoToml() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Message, "No Cargo.toml found")
}

func (s *CoverageTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangRust, result.Language)
}

func (s *CoverageTestSuite) TestContains() {
	list := []string{"apple", "banana", "cherry"}
	s.True(contains(list, "banana"))
	s.False(contains(list, "grape"))
}

func (s *CoverageTestSuite) TestFindCoverageReports_NoCoverage() {
	cov := s.check.findCoverageReports(s.tempDir)
	s.Equal(-1.0, cov)
}

func (s *CoverageTestSuite) TestFindCoverageReports_WithReports() {
	// Create a cobertura.xml with valid coverage format
	s.writeFile("target/tarpaulin/cobertura.xml", `<?xml version="1.0"?>
<coverage line-rate="0.85" branch-rate="0.75">
</coverage>`)

	cov := s.check.findCoverageReports(s.tempDir)
	s.Equal(85.0, cov)
}

func TestCoverageTestSuite(t *testing.T) {
	suite.Run(t, new(CoverageTestSuite))
}
