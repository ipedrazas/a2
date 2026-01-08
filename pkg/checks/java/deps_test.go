package javacheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type DepsCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *DepsCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "java-deps-test-*")
	s.Require().NoError(err)
}

func (s *DepsCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *DepsCheckTestSuite) TestIDAndName() {
	check := &DepsCheck{}
	s.Equal("java:deps", check.ID())
	s.Equal("Java Dependencies", check.Name())
}

func (s *DepsCheckTestSuite) TestRun_NoDependencyScanning() {
	check := &DepsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "No dependency scanning configured")
}

func (s *DepsCheckTestSuite) TestRun_ResultLanguage() {
	check := &DepsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangJava, result.Language)
}

func (s *DepsCheckTestSuite) TestRun_OWASPDependencyCheck_Maven() {
	content := `<project>
  <build>
    <plugins>
      <plugin>
        <groupId>org.owasp</groupId>
        <artifactId>dependency-check-maven</artifactId>
      </plugin>
    </plugins>
  </build>
</project>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &DepsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "OWASP Dependency-Check")
}

func (s *DepsCheckTestSuite) TestRun_OWASPDependencyCheck_Gradle() {
	content := `plugins {
    id 'java'
    id 'org.owasp.dependency-check' version '8.0.0'
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "build.gradle"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &DepsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "OWASP Dependency-Check")
}

func (s *DepsCheckTestSuite) TestRun_OWASPDependencyCheck_GradleKts() {
	content := `plugins {
    java
    id("org.owasp.dependency-check") version "8.0.0"
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "build.gradle.kts"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &DepsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "OWASP Dependency-Check")
}

func (s *DepsCheckTestSuite) TestRun_Snyk_ConfigFile() {
	err := os.WriteFile(filepath.Join(s.tempDir, ".snyk"), []byte("version: v1.0.0"), 0644)
	s.Require().NoError(err)

	check := &DepsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Snyk")
}

func (s *DepsCheckTestSuite) TestRun_Snyk_Maven() {
	content := `<project>
  <build>
    <plugins>
      <plugin>
        <groupId>io.snyk</groupId>
        <artifactId>snyk-maven-plugin</artifactId>
      </plugin>
    </plugins>
  </build>
</project>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &DepsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Snyk")
}

func (s *DepsCheckTestSuite) TestRun_Dependabot_Yml() {
	githubDir := filepath.Join(s.tempDir, ".github")
	err := os.MkdirAll(githubDir, 0755)
	s.Require().NoError(err)

	content := `version: 2
updates:
  - package-ecosystem: "maven"
    directory: "/"
    schedule:
      interval: "weekly"`
	err = os.WriteFile(filepath.Join(githubDir, "dependabot.yml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &DepsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Dependabot")
}

func (s *DepsCheckTestSuite) TestRun_Dependabot_Yaml() {
	githubDir := filepath.Join(s.tempDir, ".github")
	err := os.MkdirAll(githubDir, 0755)
	s.Require().NoError(err)

	content := `version: 2
updates:
  - package-ecosystem: "gradle"
    directory: "/"
    schedule:
      interval: "daily"`
	err = os.WriteFile(filepath.Join(githubDir, "dependabot.yaml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &DepsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Dependabot")
}

func (s *DepsCheckTestSuite) TestRun_Renovate_Json() {
	content := `{
  "extends": ["config:base"],
  "packageRules": [
    {
      "matchPackagePatterns": ["*"]
    }
  ]
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "renovate.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &DepsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Renovate")
}

func (s *DepsCheckTestSuite) TestRun_Renovate_Json5() {
	content := `{
  // Renovate config
  "extends": ["config:base"]
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "renovate.json5"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &DepsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Renovate")
}

func (s *DepsCheckTestSuite) TestRun_Renovate_Rc() {
	err := os.WriteFile(filepath.Join(s.tempDir, ".renovaterc"), []byte(`{"extends": ["config:base"]}`), 0644)
	s.Require().NoError(err)

	check := &DepsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Renovate")
}

func (s *DepsCheckTestSuite) TestRun_Renovate_RcJson() {
	err := os.WriteFile(filepath.Join(s.tempDir, ".renovaterc.json"), []byte(`{}`), 0644)
	s.Require().NoError(err)

	check := &DepsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Renovate")
}

func (s *DepsCheckTestSuite) TestRun_MavenDependencyPlugin() {
	content := `<project>
  <build>
    <plugins>
      <plugin>
        <groupId>org.apache.maven.plugins</groupId>
        <artifactId>maven-dependency-plugin</artifactId>
        <executions>
          <execution>
            <goals>
              <goal>analyze</goal>
            </goals>
          </execution>
        </executions>
      </plugin>
    </plugins>
  </build>
</project>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &DepsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Maven Dependency Plugin")
}

func (s *DepsCheckTestSuite) TestRun_GradleDependencyVerification() {
	gradleDir := filepath.Join(s.tempDir, "gradle")
	err := os.MkdirAll(gradleDir, 0755)
	s.Require().NoError(err)

	content := `<?xml version="1.0" encoding="UTF-8"?>
<verification-metadata>
  <components/>
</verification-metadata>`
	err = os.WriteFile(filepath.Join(gradleDir, "verification-metadata.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &DepsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Gradle Dependency Verification")
}

func (s *DepsCheckTestSuite) TestRun_GradleDependencyVerification_LegacyFile() {
	gradleDir := filepath.Join(s.tempDir, "gradle")
	err := os.MkdirAll(gradleDir, 0755)
	s.Require().NoError(err)

	err = os.WriteFile(filepath.Join(gradleDir, "dependency-verification.xml"), []byte("<?xml/>"), 0644)
	s.Require().NoError(err)

	check := &DepsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Gradle Dependency Verification")
}

func (s *DepsCheckTestSuite) TestRun_MultipleTools() {
	// OWASP in pom.xml
	pomContent := `<project>
  <build>
    <plugins>
      <plugin>
        <artifactId>dependency-check-maven</artifactId>
      </plugin>
    </plugins>
  </build>
</project>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(pomContent), 0644)
	s.Require().NoError(err)

	// Dependabot config
	githubDir := filepath.Join(s.tempDir, ".github")
	err = os.MkdirAll(githubDir, 0755)
	s.Require().NoError(err)
	err = os.WriteFile(filepath.Join(githubDir, "dependabot.yml"), []byte("version: 2"), 0644)
	s.Require().NoError(err)

	check := &DepsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "OWASP Dependency-Check")
	s.Contains(result.Message, "Dependabot")
}

func TestDepsCheckTestSuite(t *testing.T) {
	suite.Run(t, new(DepsCheckTestSuite))
}
