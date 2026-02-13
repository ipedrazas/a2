package javacheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type LintCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *LintCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "java-lint-test-*")
	s.Require().NoError(err)
}

func (s *LintCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *LintCheckTestSuite) TestIDAndName() {
	check := &LintCheck{}
	s.Equal("java:lint", check.ID())
	s.Equal("Java Lint", check.Name())
}

func (s *LintCheckTestSuite) TestRun_NoLinters() {
	check := &LintCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "No static analysis tools configured")
}

func (s *LintCheckTestSuite) TestRun_ResultLanguage() {
	check := &LintCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangJava, result.Language)
}

func (s *LintCheckTestSuite) TestRun_Checkstyle_ConfigFile() {
	err := os.WriteFile(filepath.Join(s.tempDir, "checkstyle.xml"), []byte("<module/>"), 0644)
	s.Require().NoError(err)

	check := &LintCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Checkstyle")
}

func (s *LintCheckTestSuite) TestRun_Checkstyle_ConfigSubdirectory() {
	configDir := filepath.Join(s.tempDir, "config", "checkstyle")
	err := os.MkdirAll(configDir, 0755)
	s.Require().NoError(err)

	err = os.WriteFile(filepath.Join(configDir, "checkstyle.xml"), []byte("<module/>"), 0644)
	s.Require().NoError(err)

	check := &LintCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Checkstyle")
}

func (s *LintCheckTestSuite) TestRun_Checkstyle_Maven() {
	content := `<project>
  <build>
    <plugins>
      <plugin>
        <groupId>org.apache.maven.plugins</groupId>
        <artifactId>maven-checkstyle-plugin</artifactId>
      </plugin>
    </plugins>
  </build>
</project>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LintCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Checkstyle")
}

func (s *LintCheckTestSuite) TestRun_Checkstyle_Gradle() {
	content := `plugins {
    id 'java'
    id 'checkstyle'
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "build.gradle"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LintCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Checkstyle")
}

func (s *LintCheckTestSuite) TestRun_SpotBugs_ConfigFile() {
	err := os.WriteFile(filepath.Join(s.tempDir, "spotbugs.xml"), []byte("<FindBugsFilter/>"), 0644)
	s.Require().NoError(err)

	check := &LintCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "SpotBugs")
}

func (s *LintCheckTestSuite) TestRun_SpotBugs_ExcludeFile() {
	err := os.WriteFile(filepath.Join(s.tempDir, "spotbugs-exclude.xml"), []byte("<FindBugsFilter/>"), 0644)
	s.Require().NoError(err)

	check := &LintCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "SpotBugs")
}

func (s *LintCheckTestSuite) TestRun_SpotBugs_Maven() {
	content := `<project>
  <build>
    <plugins>
      <plugin>
        <groupId>com.github.spotbugs</groupId>
        <artifactId>spotbugs-maven-plugin</artifactId>
      </plugin>
    </plugins>
  </build>
</project>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LintCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "SpotBugs")
}

func (s *LintCheckTestSuite) TestRun_SpotBugs_Gradle() {
	content := `plugins {
    id 'java'
    id 'com.github.spotbugs' version '5.0.0'
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "build.gradle"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LintCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "SpotBugs")
}

func (s *LintCheckTestSuite) TestRun_FindBugs_Legacy() {
	content := `<project>
  <build>
    <plugins>
      <plugin>
        <groupId>org.codehaus.mojo</groupId>
        <artifactId>findbugs-maven-plugin</artifactId>
      </plugin>
    </plugins>
  </build>
</project>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LintCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "SpotBugs")
}

func (s *LintCheckTestSuite) TestRun_PMD_ConfigFile() {
	err := os.WriteFile(filepath.Join(s.tempDir, "pmd.xml"), []byte("<ruleset/>"), 0644)
	s.Require().NoError(err)

	check := &LintCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "PMD")
}

func (s *LintCheckTestSuite) TestRun_PMD_RulesetFile() {
	err := os.WriteFile(filepath.Join(s.tempDir, "ruleset.xml"), []byte("<ruleset/>"), 0644)
	s.Require().NoError(err)

	check := &LintCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "PMD")
}

func (s *LintCheckTestSuite) TestRun_PMD_Maven() {
	content := `<project>
  <build>
    <plugins>
      <plugin>
        <groupId>org.apache.maven.plugins</groupId>
        <artifactId>maven-pmd-plugin</artifactId>
      </plugin>
    </plugins>
  </build>
</project>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LintCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "PMD")
}

func (s *LintCheckTestSuite) TestRun_PMD_Gradle() {
	content := `plugins {
    id 'java'
    id 'pmd'
}

pmd {
    toolVersion = "6.55.0"
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "build.gradle"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LintCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "PMD")
}

func (s *LintCheckTestSuite) TestRun_ErrorProne_Maven() {
	content := `<project>
  <build>
    <plugins>
      <plugin>
        <groupId>org.apache.maven.plugins</groupId>
        <artifactId>maven-compiler-plugin</artifactId>
        <configuration>
          <annotationProcessorPaths>
            <path>
              <groupId>com.google.errorprone</groupId>
              <artifactId>error_prone_core</artifactId>
            </path>
          </annotationProcessorPaths>
        </configuration>
      </plugin>
    </plugins>
  </build>
</project>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LintCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Error Prone")
}

func (s *LintCheckTestSuite) TestRun_ErrorProne_Gradle() {
	content := `plugins {
    id 'java'
    id 'net.ltgt.errorprone' version '3.1.0'
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "build.gradle"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LintCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Error Prone")
}

func (s *LintCheckTestSuite) TestRun_SonarQube_PropertiesFile() {
	err := os.WriteFile(filepath.Join(s.tempDir, "sonar-project.properties"), []byte("sonar.projectKey=myapp"), 0644)
	s.Require().NoError(err)

	check := &LintCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "SonarQube")
}

func (s *LintCheckTestSuite) TestRun_SonarQube_Maven() {
	content := `<project>
  <build>
    <plugins>
      <plugin>
        <groupId>org.sonarsource.scanner.maven</groupId>
        <artifactId>sonar-maven-plugin</artifactId>
      </plugin>
    </plugins>
  </build>
</project>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LintCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "SonarQube")
}

func (s *LintCheckTestSuite) TestRun_SonarQube_Gradle() {
	content := `plugins {
    id 'java'
    id 'org.sonarqube' version '4.0.0'
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "build.gradle"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LintCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "SonarQube")
}

func (s *LintCheckTestSuite) TestRun_MultipleLinters() {
	// Create checkstyle.xml
	err := os.WriteFile(filepath.Join(s.tempDir, "checkstyle.xml"), []byte("<module/>"), 0644)
	s.Require().NoError(err)

	// PMD in pom.xml
	pomContent := `<project>
  <build>
    <plugins>
      <plugin>
        <artifactId>maven-pmd-plugin</artifactId>
      </plugin>
    </plugins>
  </build>
</project>`
	err = os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(pomContent), 0644)
	s.Require().NoError(err)

	check := &LintCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Checkstyle")
	s.Contains(result.Reason, "PMD")
}

func (s *LintCheckTestSuite) TestRun_GradleKts() {
	content := `plugins {
    java
    checkstyle
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "build.gradle.kts"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &LintCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Reason, "Checkstyle")
}

func TestLintCheckTestSuite(t *testing.T) {
	suite.Run(t, new(LintCheckTestSuite))
}
