package javacheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/stretchr/testify/suite"
)

type TestsCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *TestsCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "java-tests-test-*")
	s.Require().NoError(err)
}

func (s *TestsCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *TestsCheckTestSuite) TestIDAndName() {
	check := &TestsCheck{}
	s.Equal("java:tests", check.ID())
	s.Equal("Java Tests", check.Name())
}

func (s *TestsCheckTestSuite) TestRun_NoBuildTool() {
	check := &TestsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Message, "No build tool detected")
}

func (s *TestsCheckTestSuite) TestRun_ResultLanguage() {
	// Create pom.xml so we have a build tool
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte("<project></project>"), 0644)
	s.Require().NoError(err)

	check := &TestsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangJava, result.Language)
}

func (s *TestsCheckTestSuite) TestDetectBuildTool_ConfigOverride() {
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte("<project></project>"), 0644)
	s.Require().NoError(err)

	check := &TestsCheck{
		Config: &config.JavaLanguageConfig{
			BuildTool: "gradle",
		},
	}
	tool := check.detectBuildTool(s.tempDir)
	s.Equal("gradle", tool)
}

func (s *TestsCheckTestSuite) TestGetMavenTestCommand_WithWrapper() {
	err := os.WriteFile(filepath.Join(s.tempDir, "mvnw"), []byte("#!/bin/bash"), 0755)
	s.Require().NoError(err)

	check := &TestsCheck{}
	cmd := check.getMavenTestCommand(s.tempDir)

	s.Equal("./mvnw", cmd.Args[0])
	s.Contains(cmd.Args, "test")
}

func (s *TestsCheckTestSuite) TestGetMavenTestCommand_WithoutWrapper() {
	check := &TestsCheck{}
	cmd := check.getMavenTestCommand(s.tempDir)

	s.Equal("mvn", cmd.Args[0])
	s.Contains(cmd.Args, "test")
}

func (s *TestsCheckTestSuite) TestGetGradleTestCommand_WithWrapper() {
	err := os.WriteFile(filepath.Join(s.tempDir, "gradlew"), []byte("#!/bin/bash"), 0755)
	s.Require().NoError(err)

	check := &TestsCheck{}
	cmd := check.getGradleTestCommand(s.tempDir)

	s.Equal("./gradlew", cmd.Args[0])
	s.Contains(cmd.Args, "test")
}

func (s *TestsCheckTestSuite) TestGetGradleTestCommand_WithoutWrapper() {
	check := &TestsCheck{}
	cmd := check.getGradleTestCommand(s.tempDir)

	s.Equal("gradle", cmd.Args[0])
	s.Contains(cmd.Args, "test")
}

func (s *TestsCheckTestSuite) TestExtractTestSummary_MavenPass() {
	output := `
[INFO] Running com.example.AppTest
[INFO] Tests run: 5, Failures: 0, Errors: 0, Skipped: 0
[INFO] BUILD SUCCESS
`
	summary := extractTestSummary(output, "maven")
	s.Equal("5 tests passed", summary)
}

func (s *TestsCheckTestSuite) TestExtractTestSummary_MavenFail() {
	output := `
[INFO] Running com.example.AppTest
[ERROR] Tests run: 10, Failures: 2, Errors: 1, Skipped: 1
[INFO] BUILD FAILURE
`
	summary := extractTestSummary(output, "maven")
	s.Contains(summary, "10 tests")
	s.Contains(summary, "2 failures")
	s.Contains(summary, "1 errors")
}

func (s *TestsCheckTestSuite) TestExtractTestSummary_MavenMultipleModules() {
	output := `
[INFO] Running com.example.module1.Test1
[INFO] Tests run: 5, Failures: 0, Errors: 0, Skipped: 0
[INFO] Running com.example.module2.Test2
[INFO] Tests run: 3, Failures: 1, Errors: 0, Skipped: 1
`
	summary := extractTestSummary(output, "maven")
	s.Contains(summary, "8 tests")
	s.Contains(summary, "1 failures")
}

func (s *TestsCheckTestSuite) TestExtractTestSummary_Gradle() {
	output := `
> Task :test
5 tests completed
BUILD SUCCESSFUL
`
	summary := extractTestSummary(output, "gradle")
	s.Equal("5 tests completed", summary)
}

func (s *TestsCheckTestSuite) TestExtractTestSummary_NoMatch() {
	output := "Some random output without test results"
	summary := extractTestSummary(output, "maven")
	s.Equal("", summary)
}

func (s *TestsCheckTestSuite) TestParseInt() {
	s.Equal(0, parseInt("0"))
	s.Equal(5, parseInt("5"))
	s.Equal(123, parseInt("123"))
	s.Equal(0, parseInt(""))
}

func (s *TestsCheckTestSuite) TestIntToStr() {
	s.Equal("0", intToStr(0))
	s.Equal("5", intToStr(5))
	s.Equal("123", intToStr(123))
	s.Equal("1000", intToStr(1000))
}

func TestTestsCheckTestSuite(t *testing.T) {
	suite.Run(t, new(TestsCheckTestSuite))
}
