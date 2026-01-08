package javacheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/stretchr/testify/suite"
)

type BuildCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *BuildCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "java-build-test-*")
	s.Require().NoError(err)
}

func (s *BuildCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *BuildCheckTestSuite) TestIDAndName() {
	check := &BuildCheck{}
	s.Equal("java:build", check.ID())
	s.Equal("Java Build", check.Name())
}

func (s *BuildCheckTestSuite) TestRun_NoBuildTool() {
	check := &BuildCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Message, "No build tool detected")
}

func (s *BuildCheckTestSuite) TestRun_ResultLanguage() {
	// Create pom.xml so we have a build tool
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte("<project></project>"), 0644)
	s.Require().NoError(err)

	check := &BuildCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangJava, result.Language)
}

func (s *BuildCheckTestSuite) TestDetectBuildTool_ConfigOverride() {
	// Without config, it should detect maven
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte("<project></project>"), 0644)
	s.Require().NoError(err)

	check := &BuildCheck{}
	tool := check.detectBuildTool(s.tempDir)
	s.Equal("maven", tool)

	// With config override to gradle
	check = &BuildCheck{
		Config: &config.JavaLanguageConfig{
			BuildTool: "gradle",
		},
	}
	tool = check.detectBuildTool(s.tempDir)
	s.Equal("gradle", tool)
}

func (s *BuildCheckTestSuite) TestDetectBuildTool_ConfigAuto() {
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte("<project></project>"), 0644)
	s.Require().NoError(err)

	check := &BuildCheck{
		Config: &config.JavaLanguageConfig{
			BuildTool: "auto",
		},
	}
	tool := check.detectBuildTool(s.tempDir)
	s.Equal("maven", tool)
}

func (s *BuildCheckTestSuite) TestGetMavenCommand_WithWrapper() {
	// Create wrapper script
	err := os.WriteFile(filepath.Join(s.tempDir, "mvnw"), []byte("#!/bin/bash"), 0755)
	s.Require().NoError(err)

	check := &BuildCheck{}
	cmd := check.getMavenCommand(s.tempDir)

	s.Equal("./mvnw", cmd.Args[0])
	s.Contains(cmd.Args, "compile")
}

func (s *BuildCheckTestSuite) TestGetMavenCommand_WithoutWrapper() {
	check := &BuildCheck{}
	cmd := check.getMavenCommand(s.tempDir)

	s.Equal("mvn", cmd.Args[0])
	s.Contains(cmd.Args, "compile")
}

func (s *BuildCheckTestSuite) TestGetGradleCommand_WithWrapper() {
	// Create wrapper script
	err := os.WriteFile(filepath.Join(s.tempDir, "gradlew"), []byte("#!/bin/bash"), 0755)
	s.Require().NoError(err)

	check := &BuildCheck{}
	cmd := check.getGradleCommand(s.tempDir)

	s.Equal("./gradlew", cmd.Args[0])
	s.Contains(cmd.Args, "compileJava")
}

func (s *BuildCheckTestSuite) TestGetGradleCommand_WithoutWrapper() {
	check := &BuildCheck{}
	cmd := check.getGradleCommand(s.tempDir)

	s.Equal("gradle", cmd.Args[0])
	s.Contains(cmd.Args, "compileJava")
}

func TestBuildCheckTestSuite(t *testing.T) {
	suite.Run(t, new(BuildCheckTestSuite))
}
