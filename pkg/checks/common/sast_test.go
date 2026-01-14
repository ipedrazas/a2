package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
	"github.com/stretchr/testify/suite"
)

type SASTCheckTestSuite struct {
	suite.Suite
	tempDir          string
	semgrepInstalled bool
}

func (s *SASTCheckTestSuite) SetupSuite() {
	s.semgrepInstalled = checkutil.ToolAvailable("semgrep")
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

// Tests that run when semgrep IS installed
func (s *SASTCheckTestSuite) TestSemgrepInstalled_RunsWithDefaults() {
	if !s.semgrepInstalled {
		s.T().Skip("semgrep not installed")
	}

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "semgrep:")
	s.Contains(result.Message, "auto rules")
}

func (s *SASTCheckTestSuite) TestSemgrepInstalled_UsesConfig() {
	if !s.semgrepInstalled {
		s.T().Skip("semgrep not installed")
	}

	err := os.WriteFile(filepath.Join(s.tempDir, ".semgrep.yml"), []byte("rules: []"), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "semgrep:")
	s.Contains(result.Message, ".semgrep.yml")
}

// Tests for config detection when semgrep is NOT installed
func (s *SASTCheckTestSuite) TestSemgrepYml() {
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

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
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

	err := os.WriteFile(filepath.Join(s.tempDir, "semgrep.yaml"), []byte("rules: []"), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Semgrep")
}

func (s *SASTCheckTestSuite) TestSemgrepDir() {
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

	err := os.MkdirAll(filepath.Join(s.tempDir, ".semgrep"), 0755)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Semgrep")
}

// The following tests check config detection and only run when semgrep is NOT installed
// When semgrep IS installed, it runs semgrep instead of detecting other tools

func (s *SASTCheckTestSuite) TestSonarProperties() {
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

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
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

	err := os.WriteFile(filepath.Join(s.tempDir, ".sonarcloud.properties"), []byte("sonar.projectKey=test"), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "SonarQube")
}

func (s *SASTCheckTestSuite) TestSonarInGradle() {
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

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
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

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
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

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
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

	err := os.WriteFile(filepath.Join(s.tempDir, "snyk.json"), []byte("{}"), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Snyk")
}

func (s *SASTCheckTestSuite) TestCodeQLConfig() {
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

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
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

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
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

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
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

	err := os.WriteFile(filepath.Join(s.tempDir, "checkmarx.config"), []byte("{}"), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Checkmarx")
}

func (s *SASTCheckTestSuite) TestVeracode() {
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

	err := os.WriteFile(filepath.Join(s.tempDir, "veracode.json"), []byte("{}"), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Veracode")
}

func (s *SASTCheckTestSuite) TestFortify() {
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

	err := os.WriteFile(filepath.Join(s.tempDir, "fortify-sca.properties"), []byte(""), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Fortify")
}

func (s *SASTCheckTestSuite) TestBearer() {
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

	err := os.WriteFile(filepath.Join(s.tempDir, "bearer.yml"), []byte(""), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Bearer")
}

func (s *SASTCheckTestSuite) TestHorusec() {
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

	err := os.WriteFile(filepath.Join(s.tempDir, "horusec-config.json"), []byte("{}"), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Horusec")
}

func (s *SASTCheckTestSuite) TestGosecInMakefile() {
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

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
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

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
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

	err := os.WriteFile(filepath.Join(s.tempDir, ".bandit"), []byte("[bandit]\nexclude_dirs = tests"), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Bandit")
}

func (s *SASTCheckTestSuite) TestBanditInPyproject() {
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

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
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

	err := os.WriteFile(filepath.Join(s.tempDir, ".safety-policy.yml"), []byte(""), 0644)
	s.Require().NoError(err)

	check := &SASTCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Safety")
}

func (s *SASTCheckTestSuite) TestEslintPluginSecurity() {
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

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
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

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
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

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
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

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
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

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
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

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
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

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
	if s.semgrepInstalled {
		s.T().Skip("semgrep installed - this test checks config detection fallback")
	}

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
