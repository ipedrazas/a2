package pythoncheck

import (
	"os"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type DepsTestSuite struct {
	suite.Suite
	tempDir string
	check   *DepsCheck
}

func (s *DepsTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "python-deps-test-*")
	s.NoError(err)
	s.tempDir = dir
	s.check = &DepsCheck{}
}

func (s *DepsTestSuite) TearDownTest() {
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

func (s *DepsTestSuite) TestID() {
	s.Equal("python:deps", s.check.ID())
}

func (s *DepsTestSuite) TestName() {
	s.Equal("Python Vulnerabilities", s.check.Name())
}

func (s *DepsTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangPython, result.Language)
}

func (s *DepsTestSuite) TestCountPythonVulnerabilities_PipAudit() {
	output := `
Name         Version  ID              Fix Versions
----------  -------- --------------  -------------
requests    2.25.0   PYSEC-2021-123  2.26.0
urllib3     1.24.0   CVE-2019-11324  1.25.0
`
	count := countPythonVulnerabilities(output, "pip-audit")
	s.Equal(2, count)
}

func (s *DepsTestSuite) TestCountPythonVulnerabilities_Safety() {
	output := `
requests -> requests-2.25.0 - vulnerability found
urllib3 -> urllib3-1.24.0 - vulnerability found
`
	count := countPythonVulnerabilities(output, "safety")
	s.Equal(2, count)
}

func (s *DepsTestSuite) TestCountPythonVulnerabilities_NoVulns() {
	output := "No vulnerabilities found"
	count := countPythonVulnerabilities(output, "pip-audit")
	s.Equal(0, count)
}

func TestDepsTestSuite(t *testing.T) {
	suite.Run(t, new(DepsTestSuite))
}
