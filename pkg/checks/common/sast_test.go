package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type SASTCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *SASTCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "sast-test-*")
	s.Require().NoError(err)
}

func (s *SASTCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *SASTCheckTestSuite) TestIDAndName() {
	check := &SASTCheck{}
	s.Equal("common:sast", check.ID())
	s.Equal("SAST Security Scanning", check.Name())
}

func (s *SASTCheckTestSuite) TestSemgrepYml() {
	err := os.WriteFile(filepath.Join(s.tempDir, ".semgrep.yml"), []byte("rules: []"), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Semgrep")
}

func (s *SASTCheckTestSuite) TestSemgrepYaml() {
	err := os.WriteFile(filepath.Join(s.tempDir, "semgrep.yaml"), []byte("rules: []"), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Semgrep")
}

func (s *SASTCheckTestSuite) TestSemgrepDir() {
	err := os.MkdirAll(filepath.Join(s.tempDir, ".semgrep"), 0755)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Semgrep")
}

func (s *SASTCheckTestSuite) TestSonarProperties() {
	content := `sonar.projectKey=myproject
sonar.organization=myorg
`
	err := os.WriteFile(filepath.Join(s.tempDir, "sonar-project.properties"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "SonarQube")
}

func (s *SASTCheckTestSuite) TestSonarCloudProperties() {
	err := os.WriteFile(filepath.Join(s.tempDir, ".sonarcloud.properties"), []byte("sonar.projectKey=test"), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "SonarQube")
}

func (s *SASTCheckTestSuite) TestSonarInGradle() {
	content := `plugins {
    id 'java'
    id 'org.sonarqube' version '4.0.0'
}
`
	err := os.WriteFile(filepath.Join(s.tempDir, "build.gradle"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "SonarQube")
}

func (s *SASTCheckTestSuite) TestSonarInMaven() {
	content := `<?xml version="1.0" encoding="UTF-8"?>
<project>
    <properties>
        <sonar.projectKey>myproject</sonar.projectKey>
    </properties>
</project>
`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "SonarQube")
}

func (s *SASTCheckTestSuite) TestSnykConfig() {
	err := os.WriteFile(filepath.Join(s.tempDir, ".snyk"), []byte("version: v1.0.0"), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Snyk")
}

func (s *SASTCheckTestSuite) TestSnykJson() {
	err := os.WriteFile(filepath.Join(s.tempDir, "snyk.json"), []byte("{}"), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Snyk")
}

func (s *SASTCheckTestSuite) TestCodeQLConfig() {
	codeqlDir := filepath.Join(s.tempDir, ".github", "codeql")
	err := os.MkdirAll(codeqlDir, 0755)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "CodeQL")
}

func (s *SASTCheckTestSuite) TestCodeQLWorkflow() {
	workflowsDir := filepath.Join(s.tempDir, ".github", "workflows")
	err := os.MkdirAll(workflowsDir, 0755)
	s.Require().NoError(err)

	content := `name: CodeQL
on: [push]
jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: github/codeql-action/init@v2
      - uses: github/codeql-action/analyze@v2
`
	err = os.WriteFile(filepath.Join(workflowsDir, "codeql.yml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "CodeQL")
}

func (s *SASTCheckTestSuite) TestCodeQLWorkflowByContent() {
	workflowsDir := filepath.Join(s.tempDir, ".github", "workflows")
	err := os.MkdirAll(workflowsDir, 0755)
	s.Require().NoError(err)

	content := `name: Security
on: [push]
jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: github/codeql-action/init@v3
`
	err = os.WriteFile(filepath.Join(workflowsDir, "security.yml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "CodeQL")
}

func (s *SASTCheckTestSuite) TestCheckmarx() {
	err := os.WriteFile(filepath.Join(s.tempDir, "checkmarx.config"), []byte("{}"), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Checkmarx")
}

func (s *SASTCheckTestSuite) TestVeracode() {
	err := os.WriteFile(filepath.Join(s.tempDir, "veracode.json"), []byte("{}"), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Veracode")
}

func (s *SASTCheckTestSuite) TestFortify() {
	err := os.WriteFile(filepath.Join(s.tempDir, "fortify-sca.properties"), []byte(""), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Fortify")
}

func (s *SASTCheckTestSuite) TestBearer() {
	err := os.WriteFile(filepath.Join(s.tempDir, "bearer.yml"), []byte(""), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Bearer")
}

func (s *SASTCheckTestSuite) TestHorusec() {
	err := os.WriteFile(filepath.Join(s.tempDir, "horusec-config.json"), []byte("{}"), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Horusec")
}

func (s *SASTCheckTestSuite) TestGosecInMakefile() {
	// Create go.mod to indicate Go project
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte("module test\n\ngo 1.21"), 0644)
	s.Require().NoError(err)

	makefile := `lint:
	gosec ./...
`
	err = os.WriteFile(filepath.Join(s.tempDir, "Makefile"), []byte(makefile), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "gosec")
}

func (s *SASTCheckTestSuite) TestGosecInTaskfile() {
	// Create go.mod to indicate Go project
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte("module test\n\ngo 1.21"), 0644)
	s.Require().NoError(err)

	taskfile := `version: '3'
tasks:
  security:
    cmds:
      - gosec ./...
`
	err = os.WriteFile(filepath.Join(s.tempDir, "Taskfile.yml"), []byte(taskfile), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "gosec")
}

func (s *SASTCheckTestSuite) TestBanditConfig() {
	err := os.WriteFile(filepath.Join(s.tempDir, ".bandit"), []byte("[bandit]\nexclude_dirs = tests"), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Bandit")
}

func (s *SASTCheckTestSuite) TestBanditInPyproject() {
	content := `[tool.bandit]
exclude_dirs = ["tests"]
`
	err := os.WriteFile(filepath.Join(s.tempDir, "pyproject.toml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Bandit")
}

func (s *SASTCheckTestSuite) TestSafetyPolicy() {
	err := os.WriteFile(filepath.Join(s.tempDir, ".safety-policy.yml"), []byte(""), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Safety")
}

func (s *SASTCheckTestSuite) TestEslintPluginSecurity() {
	content := `{
  "name": "myapp",
  "devDependencies": {
    "eslint-plugin-security": "^1.7.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "eslint-plugin-security")
}

func (s *SASTCheckTestSuite) TestAuditCI() {
	content := `{
  "name": "myapp",
  "devDependencies": {
    "audit-ci": "^6.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "audit-ci")
}

func (s *SASTCheckTestSuite) TestFindSecBugs() {
	content := `<?xml version="1.0" encoding="UTF-8"?>
<project>
    <build>
        <plugins>
            <plugin>
                <groupId>com.github.spotbugs</groupId>
                <artifactId>spotbugs-maven-plugin</artifactId>
                <configuration>
                    <plugins>
                        <plugin>
                            <groupId>com.h3xstream.findsecbugs</groupId>
                            <artifactId>findsecbugs-plugin</artifactId>
                        </plugin>
                    </plugins>
                </configuration>
            </plugin>
        </plugins>
    </build>
</project>
`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "FindSecBugs")
}

func (s *SASTCheckTestSuite) TestCISemgrep() {
	workflowsDir := filepath.Join(s.tempDir, ".github", "workflows")
	err := os.MkdirAll(workflowsDir, 0755)
	s.Require().NoError(err)

	content := `name: Security
on: [push]
jobs:
  semgrep:
    runs-on: ubuntu-latest
    steps:
      - uses: returntocorp/semgrep-action@v1
`
	err = os.WriteFile(filepath.Join(workflowsDir, "security.yml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Semgrep")
}

func (s *SASTCheckTestSuite) TestCITrivy() {
	workflowsDir := filepath.Join(s.tempDir, ".github", "workflows")
	err := os.MkdirAll(workflowsDir, 0755)
	s.Require().NoError(err)

	content := `name: Security
on: [push]
jobs:
  trivy:
    runs-on: ubuntu-latest
    steps:
      - uses: aquasecurity/trivy-action@master
`
	err = os.WriteFile(filepath.Join(workflowsDir, "security.yml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Trivy")
}

func (s *SASTCheckTestSuite) TestGitLabSAST() {
	content := `include:
  - template: Security/SAST.gitlab-ci.yml

sast:
  stage: test
`
	err := os.WriteFile(filepath.Join(s.tempDir, ".gitlab-ci.yml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "GitLab SAST")
}

func (s *SASTCheckTestSuite) TestMultipleFindings() {
	// Create Semgrep config
	err := os.WriteFile(filepath.Join(s.tempDir, ".semgrep.yml"), []byte("rules: []"), 0644)
	s.Require().NoError(err)

	// Create Snyk config
	err = os.WriteFile(filepath.Join(s.tempDir, ".snyk"), []byte(""), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Semgrep")
	s.Contains(result.Message, "Snyk")
}

func (s *SASTCheckTestSuite) TestNoSASTTooling() {
	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "No SAST tooling found")
}

func (s *SASTCheckTestSuite) TestResultLanguage() {
	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangCommon, result.Language)
}

func TestSASTCheckTestSuite(t *testing.T) {
	suite.Run(t, new(SASTCheckTestSuite))
}
