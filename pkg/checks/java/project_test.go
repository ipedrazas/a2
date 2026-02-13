package javacheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type ProjectCheckTestSuite struct {
	suite.Suite
	check   *ProjectCheck
	tempDir string
}

func (s *ProjectCheckTestSuite) SetupTest() {
	s.check = &ProjectCheck{}
	tempDir, err := os.MkdirTemp("", "java-project-test-*")
	s.Require().NoError(err)
	s.tempDir = tempDir
}

func (s *ProjectCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *ProjectCheckTestSuite) TestID() {
	s.Equal("java:project", s.check.ID())
}

func (s *ProjectCheckTestSuite) TestName() {
	s.Equal("Java Project", s.check.Name())
}

func (s *ProjectCheckTestSuite) TestRun_MavenProject() {
	pomXml := filepath.Join(s.tempDir, "pom.xml")
	err := os.WriteFile(pomXml, []byte("<project></project>"), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Maven")
}

func (s *ProjectCheckTestSuite) TestRun_MavenProjectWithWrapper() {
	pomXml := filepath.Join(s.tempDir, "pom.xml")
	err := os.WriteFile(pomXml, []byte("<project></project>"), 0644)
	s.Require().NoError(err)

	mvnw := filepath.Join(s.tempDir, "mvnw")
	err = os.WriteFile(mvnw, []byte("#!/bin/bash"), 0755)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Maven")
	s.Contains(result.Reason, "wrapper")
}

func (s *ProjectCheckTestSuite) TestRun_GradleProject() {
	buildGradle := filepath.Join(s.tempDir, "build.gradle")
	err := os.WriteFile(buildGradle, []byte("plugins {}"), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Gradle")
	s.Contains(result.Reason, "Groovy")
}

func (s *ProjectCheckTestSuite) TestRun_GradleKotlinDSL() {
	buildGradle := filepath.Join(s.tempDir, "build.gradle.kts")
	err := os.WriteFile(buildGradle, []byte("plugins {}"), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Gradle")
	s.Contains(result.Reason, "Kotlin")
}

func (s *ProjectCheckTestSuite) TestRun_GradleProjectWithWrapper() {
	buildGradle := filepath.Join(s.tempDir, "build.gradle")
	err := os.WriteFile(buildGradle, []byte("plugins {}"), 0644)
	s.Require().NoError(err)

	gradlew := filepath.Join(s.tempDir, "gradlew")
	err = os.WriteFile(gradlew, []byte("#!/bin/bash"), 0755)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Gradle")
	s.Contains(result.Reason, "wrapper")
}

func (s *ProjectCheckTestSuite) TestRun_NoProject() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Reason, "No Java project")
}

func (s *ProjectCheckTestSuite) TestRun_GradlePreferredOverMaven() {
	// When both exist, Gradle is preferred (more modern)
	pomXml := filepath.Join(s.tempDir, "pom.xml")
	err := os.WriteFile(pomXml, []byte("<project></project>"), 0644)
	s.Require().NoError(err)

	buildGradle := filepath.Join(s.tempDir, "build.gradle")
	err = os.WriteFile(buildGradle, []byte("plugins {}"), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Gradle")
}

func (s *ProjectCheckTestSuite) TestRun_ResultLanguage() {
	pomXml := filepath.Join(s.tempDir, "pom.xml")
	err := os.WriteFile(pomXml, []byte("<project></project>"), 0644)
	s.Require().NoError(err)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangJava, result.Language)
}

func (s *ProjectCheckTestSuite) TestDetectBuildTool_Maven() {
	pomXml := filepath.Join(s.tempDir, "pom.xml")
	err := os.WriteFile(pomXml, []byte("<project></project>"), 0644)
	s.Require().NoError(err)

	tool := detectBuildTool(s.tempDir)
	s.Equal("maven", tool)
}

func (s *ProjectCheckTestSuite) TestDetectBuildTool_Gradle() {
	buildGradle := filepath.Join(s.tempDir, "build.gradle")
	err := os.WriteFile(buildGradle, []byte("plugins {}"), 0644)
	s.Require().NoError(err)

	tool := detectBuildTool(s.tempDir)
	s.Equal("gradle", tool)
}

func (s *ProjectCheckTestSuite) TestDetectBuildTool_GradleWrapper() {
	gradlew := filepath.Join(s.tempDir, "gradlew")
	err := os.WriteFile(gradlew, []byte("#!/bin/bash"), 0755)
	s.Require().NoError(err)

	tool := detectBuildTool(s.tempDir)
	s.Equal("gradle", tool)
}

func (s *ProjectCheckTestSuite) TestDetectBuildTool_MavenWrapper() {
	mvnw := filepath.Join(s.tempDir, "mvnw")
	err := os.WriteFile(mvnw, []byte("#!/bin/bash"), 0755)
	s.Require().NoError(err)

	tool := detectBuildTool(s.tempDir)
	s.Equal("maven", tool)
}

func (s *ProjectCheckTestSuite) TestDetectBuildTool_None() {
	tool := detectBuildTool(s.tempDir)
	s.Equal("", tool)
}

func TestProjectCheckTestSuite(t *testing.T) {
	suite.Run(t, new(ProjectCheckTestSuite))
}
