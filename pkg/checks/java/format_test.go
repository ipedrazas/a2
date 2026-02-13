package javacheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type FormatCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *FormatCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "java-format-test-*")
	s.Require().NoError(err)
}

func (s *FormatCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *FormatCheckTestSuite) TestIDAndName() {
	check := &FormatCheck{}
	s.Equal("java:format", check.ID())
	s.Equal("Java Format", check.Name())
}

func (s *FormatCheckTestSuite) TestRun_NoFormatter() {
	check := &FormatCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "No formatter configuration found")
}

func (s *FormatCheckTestSuite) TestRun_ResultLanguage() {
	check := &FormatCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangJava, result.Language)
}

func (s *FormatCheckTestSuite) TestRun_GoogleJavaFormat_Maven() {
	content := `<project>
  <build>
    <plugins>
      <plugin>
        <groupId>com.coveo</groupId>
        <artifactId>fmt-maven-plugin</artifactId>
        <configuration>
          <style>google-java-format</style>
        </configuration>
      </plugin>
    </plugins>
  </build>
</project>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &FormatCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "google-java-format")
}

func (s *FormatCheckTestSuite) TestRun_GoogleJavaFormat_Gradle() {
	content := `plugins {
    id 'com.github.sherter.google-java-format' version '0.9'
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "build.gradle"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &FormatCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "google-java-format")
}

func (s *FormatCheckTestSuite) TestRun_GoogleJavaFormat_GradleKts() {
	content := `plugins {
    id("com.github.sherter.google-java-format") version "0.9"
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "build.gradle.kts"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &FormatCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "google-java-format")
}

func (s *FormatCheckTestSuite) TestRun_Spotless_Maven() {
	content := `<project>
  <build>
    <plugins>
      <plugin>
        <groupId>com.diffplug.spotless</groupId>
        <artifactId>spotless-maven-plugin</artifactId>
      </plugin>
    </plugins>
  </build>
</project>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &FormatCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Spotless")
}

func (s *FormatCheckTestSuite) TestRun_Spotless_Gradle() {
	content := `plugins {
    id 'com.diffplug.spotless' version '6.0.0'
}

spotless {
    java {
        googleJavaFormat()
    }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "build.gradle"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &FormatCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Spotless")
}

func (s *FormatCheckTestSuite) TestRun_EditorConfig() {
	content := `root = true

[*.java]
indent_style = space
indent_size = 4
`
	err := os.WriteFile(filepath.Join(s.tempDir, ".editorconfig"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &FormatCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "EditorConfig")
}

func (s *FormatCheckTestSuite) TestRun_EditorConfig_WithoutJava() {
	content := `root = true

[*.py]
indent_style = space
indent_size = 4
`
	err := os.WriteFile(filepath.Join(s.tempDir, ".editorconfig"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &FormatCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
}

func (s *FormatCheckTestSuite) TestRun_IntelliJFormatter() {
	// Create .idea/codeStyles directory with a file
	codeStyleDir := filepath.Join(s.tempDir, ".idea", "codeStyles")
	err := os.MkdirAll(codeStyleDir, 0755)
	s.Require().NoError(err)

	err = os.WriteFile(filepath.Join(codeStyleDir, "Project.xml"), []byte("<code_scheme/>"), 0644)
	s.Require().NoError(err)

	check := &FormatCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "IntelliJ")
}

func (s *FormatCheckTestSuite) TestRun_IntelliJCodeStyleSettings() {
	// Create .idea directory with codeStyleSettings.xml
	ideaDir := filepath.Join(s.tempDir, ".idea")
	err := os.MkdirAll(ideaDir, 0755)
	s.Require().NoError(err)

	err = os.WriteFile(filepath.Join(ideaDir, "codeStyleSettings.xml"), []byte("<settings/>"), 0644)
	s.Require().NoError(err)

	check := &FormatCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "IntelliJ")
}

func (s *FormatCheckTestSuite) TestRun_EclipseFormatter_Settings() {
	// Create .settings directory
	err := os.MkdirAll(filepath.Join(s.tempDir, ".settings"), 0755)
	s.Require().NoError(err)

	check := &FormatCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Eclipse")
}

func (s *FormatCheckTestSuite) TestRun_EclipseFormatter_File() {
	err := os.WriteFile(filepath.Join(s.tempDir, "eclipse-formatter.xml"), []byte("<profiles/>"), 0644)
	s.Require().NoError(err)

	check := &FormatCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Eclipse")
}

func (s *FormatCheckTestSuite) TestRun_MultipleFormatters() {
	// Spotless in pom.xml
	pomContent := `<project>
  <build>
    <plugins>
      <plugin>
        <artifactId>spotless-maven-plugin</artifactId>
      </plugin>
    </plugins>
  </build>
</project>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(pomContent), 0644)
	s.Require().NoError(err)

	// EditorConfig
	editorConfig := `[*.java]
indent_style = space
`
	err = os.WriteFile(filepath.Join(s.tempDir, ".editorconfig"), []byte(editorConfig), 0644)
	s.Require().NoError(err)

	check := &FormatCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Spotless")
	s.Contains(result.Reason, "EditorConfig")
}

func TestFormatCheckTestSuite(t *testing.T) {
	suite.Run(t, new(FormatCheckTestSuite))
}
