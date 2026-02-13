package swiftcheck

import (
	"os"
	"path/filepath"
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
	dir, err := os.MkdirTemp("", "swift-deps-test-*")
	s.NoError(err)
	s.tempDir = dir
	s.check = &DepsCheck{}
}

func (s *DepsTestSuite) TearDownTest() {
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

func (s *DepsTestSuite) writeFile(name, content string) {
	path := filepath.Join(s.tempDir, name)
	dir := filepath.Dir(path)
	if dir != s.tempDir {
		err := os.MkdirAll(dir, 0755)
		s.NoError(err)
	}
	err := os.WriteFile(path, []byte(content), 0644)
	s.NoError(err)
}

func (s *DepsTestSuite) TestID() {
	s.Equal("swift:deps", s.check.ID())
}

func (s *DepsTestSuite) TestName() {
	s.Equal("Swift Vulnerabilities", s.check.Name())
}

func (s *DepsTestSuite) TestRun_NoPackageSwift() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Message, "No Package.swift found")
}

func (s *DepsTestSuite) TestRun_NoPackageResolved() {
	s.writeFile("Package.swift", `// swift-tools-version:5.7
import PackageDescription
let package = Package(name: "Test")
`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Info, result.Status)
	s.Contains(result.Message, "No Package.resolved")
}

func (s *DepsTestSuite) TestRun_EmptyDeps() {
	s.writeFile("Package.swift", `// swift-tools-version:5.7
import PackageDescription
let package = Package(name: "Test")
`)
	s.writeFile("Package.resolved", `{
  "version": 2,
  "pins": []
}`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "No dependencies")
}

func (s *DepsTestSuite) TestRun_WithDepsV2() {
	s.writeFile("Package.swift", `// swift-tools-version:5.7
import PackageDescription
let package = Package(name: "Test")
`)
	s.writeFile("Package.resolved", `{
  "version": 2,
  "pins": [
    {"identity": "swift-argument-parser", "state": {"version": "1.0.0"}},
    {"identity": "swift-log", "state": {"version": "1.4.0"}}
  ]
}`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "2 dependencies")
}

func (s *DepsTestSuite) TestRun_WithDepsV1() {
	s.writeFile("Package.swift", `// swift-tools-version:5.7
import PackageDescription
let package = Package(name: "Test")
`)
	s.writeFile("Package.resolved", `{
  "version": 1,
  "object": {
    "pins": [
      {"package": "Alamofire", "state": {"version": "5.0.0"}}
    ]
  }
}`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "1 dependency")
}

func (s *DepsTestSuite) TestRun_InvalidJSON() {
	s.writeFile("Package.swift", `// swift-tools-version:5.7
import PackageDescription
let package = Package(name: "Test")
`)
	s.writeFile("Package.resolved", `{invalid json}`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "Cannot parse")
}

func (s *DepsTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangSwift, result.Language)
}

func (s *DepsTestSuite) TestParsePackageResolved_V2() {
	data := []byte(`{
  "version": 2,
  "pins": [
    {"identity": "pkg1", "state": {}},
    {"identity": "pkg2", "state": {}}
  ]
}`)

	deps, err := parsePackageResolved(data)
	s.NoError(err)
	s.Len(deps, 2)
	s.Contains(deps, "pkg1")
	s.Contains(deps, "pkg2")
}

func (s *DepsTestSuite) TestParsePackageResolved_V1() {
	data := []byte(`{
  "version": 1,
  "object": {
    "pins": [
      {"package": "Alamofire", "state": {}},
      {"package": "SwiftyJSON", "state": {}}
    ]
  }
}`)

	deps, err := parsePackageResolved(data)
	s.NoError(err)
	s.Len(deps, 2)
	s.Contains(deps, "Alamofire")
	s.Contains(deps, "SwiftyJSON")
}

func (s *DepsTestSuite) TestFormatDepsMessage() {
	s.Contains(formatDepsMessage(0), "0 dependencies")
	s.Contains(formatDepsMessage(1), "1 dependency")
	s.Contains(formatDepsMessage(5), "5 dependencies")
	s.Contains(formatDepsMessage(15), "15 dependencies")
}

func TestDepsTestSuite(t *testing.T) {
	suite.Run(t, new(DepsTestSuite))
}
