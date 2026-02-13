package swiftcheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type ProjectTestSuite struct {
	suite.Suite
	tempDir string
	check   *ProjectCheck
}

func (s *ProjectTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "swift-project-test-*")
	s.NoError(err)
	s.tempDir = dir
	s.check = &ProjectCheck{}
}

func (s *ProjectTestSuite) TearDownTest() {
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

func (s *ProjectTestSuite) writeFile(name, content string) {
	path := filepath.Join(s.tempDir, name)
	dir := filepath.Dir(path)
	if dir != s.tempDir {
		err := os.MkdirAll(dir, 0755)
		s.NoError(err)
	}
	err := os.WriteFile(path, []byte(content), 0644)
	s.NoError(err)
}

func (s *ProjectTestSuite) TestID() {
	s.Equal("swift:project", s.check.ID())
}

func (s *ProjectTestSuite) TestName() {
	s.Equal("Swift Project", s.check.Name())
}

func (s *ProjectTestSuite) TestRun_NoPackageSwift() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Reason, "No Package.swift found")
}

func (s *ProjectTestSuite) TestRun_BasicPackageSwift() {
	s.writeFile("Package.swift", `// swift-tools-version:5.7
import PackageDescription

let package = Package(
    name: "MyApp",
    targets: [
        .executableTarget(name: "MyApp"),
    ]
)
`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "MyApp")
}

func (s *ProjectTestSuite) TestRun_WithResolved() {
	s.writeFile("Package.swift", `// swift-tools-version:5.7
import PackageDescription

let package = Package(
    name: "TestPkg"
)
`)
	s.writeFile("Package.resolved", `{
  "version": 2,
  "pins": []
}`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "dependencies resolved")
}

func (s *ProjectTestSuite) TestRun_MinimalPackageSwift() {
	s.writeFile("Package.swift", `// swift-tools-version:5.7
import PackageDescription

let package = Package(
    name: "test"
)
`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
}

func (s *ProjectTestSuite) TestRun_ResultLanguage() {
	s.writeFile("Package.swift", `// swift-tools-version:5.7
import PackageDescription

let package = Package(
    name: "test"
)
`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangSwift, result.Language)
}

func (s *ProjectTestSuite) TestExtractPackageName() {
	tests := []struct {
		content  string
		expected string
	}{
		{`name: "MyApp"`, "MyApp"},
		{`    name: "SpacedApp",`, "SpacedApp"},
		{`name:"NoSpace"`, "NoSpace"},
		{`other: "value"`, ""},
	}

	for _, tc := range tests {
		result := extractPackageName(tc.content)
		s.Equal(tc.expected, result, "for content: %s", tc.content)
	}
}

func TestProjectTestSuite(t *testing.T) {
	suite.Run(t, new(ProjectTestSuite))
}
