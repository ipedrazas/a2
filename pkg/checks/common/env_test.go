package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type EnvCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *EnvCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "env-test-*")
	s.Require().NoError(err)
}

func (s *EnvCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *EnvCheckTestSuite) TestIDAndName() {
	check := &EnvCheck{}
	s.Equal("common:env", check.ID())
	s.Equal("Environment Config", check.Name())
}

func (s *EnvCheckTestSuite) TestEnvExample() {
	content := `# Database
DATABASE_URL=
# API Keys
API_KEY=
`
	err := os.WriteFile(filepath.Join(s.tempDir, ".env.example"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &EnvCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, ".env.example")
}

func (s *EnvCheckTestSuite) TestEnvSample() {
	content := `DEBUG=false
PORT=8080
`
	err := os.WriteFile(filepath.Join(s.tempDir, ".env.sample"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &EnvCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, ".env.sample")
}

func (s *EnvCheckTestSuite) TestEnvTemplate() {
	err := os.WriteFile(filepath.Join(s.tempDir, ".env.template"), []byte("KEY=value"), 0644)
	s.Require().NoError(err)

	check := &EnvCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, ".env.template")
}

func (s *EnvCheckTestSuite) TestExampleEnv() {
	err := os.WriteFile(filepath.Join(s.tempDir, "example.env"), []byte("KEY=value"), 0644)
	s.Require().NoError(err)

	check := &EnvCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "example.env")
}

func (s *EnvCheckTestSuite) TestEnvInGitignore() {
	gitignore := `.env
.env.local
node_modules/
`
	err := os.WriteFile(filepath.Join(s.tempDir, ".gitignore"), []byte(gitignore), 0644)
	s.Require().NoError(err)

	check := &EnvCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, ".env in .gitignore")
}

func (s *EnvCheckTestSuite) TestEnvInGitignore_Wildcard() {
	gitignore := `.env*
node_modules/
`
	err := os.WriteFile(filepath.Join(s.tempDir, ".gitignore"), []byte(gitignore), 0644)
	s.Require().NoError(err)

	check := &EnvCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, ".env in .gitignore")
}

func (s *EnvCheckTestSuite) TestEnvExistsNotIgnored() {
	// Create .env file without gitignore
	err := os.WriteFile(filepath.Join(s.tempDir, ".env"), []byte("SECRET=value"), 0644)
	s.Require().NoError(err)

	check := &EnvCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, ".env file exists but not in .gitignore")
}

func (s *EnvCheckTestSuite) TestEnvExistsButIgnored() {
	// Create .env file
	err := os.WriteFile(filepath.Join(s.tempDir, ".env"), []byte("SECRET=value"), 0644)
	s.Require().NoError(err)

	// Create .gitignore with .env
	err = os.WriteFile(filepath.Join(s.tempDir, ".gitignore"), []byte(".env\n"), 0644)
	s.Require().NoError(err)

	check := &EnvCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, ".env in .gitignore")
}

func (s *EnvCheckTestSuite) TestGoDotenv() {
	content := `module myapp

go 1.21

require github.com/joho/godotenv v1.5.1
`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &EnvCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Go dotenv")
}

func (s *EnvCheckTestSuite) TestGoEnvconfig() {
	content := `module myapp

go 1.21

require github.com/kelseyhightower/envconfig v1.4.0
`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &EnvCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Go dotenv")
}

func (s *EnvCheckTestSuite) TestGoViper() {
	content := `module myapp

go 1.21

require github.com/spf13/viper v1.18.0
`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &EnvCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Go dotenv")
}

func (s *EnvCheckTestSuite) TestPythonDotenv() {
	content := `[project]
name = "myapp"
dependencies = [
    "python-dotenv>=1.0.0",
]
`
	err := os.WriteFile(filepath.Join(s.tempDir, "pyproject.toml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &EnvCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Python dotenv")
}

func (s *EnvCheckTestSuite) TestPythonDotenv_Requirements() {
	content := `flask>=2.0.0
python-dotenv>=1.0.0
`
	err := os.WriteFile(filepath.Join(s.tempDir, "requirements.txt"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &EnvCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Python dotenv")
}

func (s *EnvCheckTestSuite) TestPythonPydanticSettings() {
	content := `[project]
name = "myapp"
dependencies = [
    "pydantic-settings>=2.0.0",
]
`
	err := os.WriteFile(filepath.Join(s.tempDir, "pyproject.toml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &EnvCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Python dotenv")
}

func (s *EnvCheckTestSuite) TestPythonDjangoEnviron() {
	content := `django-environ>=0.11.0
`
	err := os.WriteFile(filepath.Join(s.tempDir, "requirements.txt"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &EnvCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Python dotenv")
}

func (s *EnvCheckTestSuite) TestNodeDotenv() {
	content := `{
  "name": "myapp",
  "dependencies": {
    "express": "^4.18.0",
    "dotenv": "^16.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &EnvCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Node.js dotenv")
}

func (s *EnvCheckTestSuite) TestNodeDotenvSafe() {
	content := `{
  "name": "myapp",
  "dependencies": {
    "dotenv-safe": "^8.2.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &EnvCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Node.js dotenv")
}

func (s *EnvCheckTestSuite) TestNodeCrossEnv() {
	content := `{
  "name": "myapp",
  "devDependencies": {
    "cross-env": "^7.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &EnvCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Node.js dotenv")
}

func (s *EnvCheckTestSuite) TestSpringBootConfig() {
	// Create Spring Boot application.properties
	resourcesDir := filepath.Join(s.tempDir, "src", "main", "resources")
	err := os.MkdirAll(resourcesDir, 0755)
	s.Require().NoError(err)

	content := `spring.datasource.url=${DATABASE_URL}
server.port=${PORT:8080}
`
	err = os.WriteFile(filepath.Join(resourcesDir, "application.properties"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &EnvCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Spring Boot config")
}

func (s *EnvCheckTestSuite) TestSpringBootYaml() {
	// Create Spring Boot application.yml
	resourcesDir := filepath.Join(s.tempDir, "src", "main", "resources")
	err := os.MkdirAll(resourcesDir, 0755)
	s.Require().NoError(err)

	content := `spring:
  datasource:
    url: ${DATABASE_URL}
server:
  port: ${PORT:8080}
`
	err = os.WriteFile(filepath.Join(resourcesDir, "application.yml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &EnvCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Spring Boot config")
}

func (s *EnvCheckTestSuite) TestJavaDotenv() {
	content := `plugins {
    id 'java'
}

dependencies {
    implementation 'io.github.cdimascio:dotenv-java:3.0.0'
}
`
	err := os.WriteFile(filepath.Join(s.tempDir, "build.gradle"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &EnvCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, "Java config")
}

func (s *EnvCheckTestSuite) TestMultipleFindings() {
	// .env.example
	err := os.WriteFile(filepath.Join(s.tempDir, ".env.example"), []byte("KEY=value"), 0644)
	s.Require().NoError(err)

	// .gitignore with .env
	err = os.WriteFile(filepath.Join(s.tempDir, ".gitignore"), []byte(".env\n"), 0644)
	s.Require().NoError(err)

	// Go dotenv
	goMod := `module myapp

go 1.21

require github.com/joho/godotenv v1.5.1
`
	err = os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(goMod), 0644)
	s.Require().NoError(err)

	check := &EnvCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, ".env.example")
	s.Contains(result.Message, "Go dotenv")
	s.Contains(result.Message, ".env in .gitignore")
}

func (s *EnvCheckTestSuite) TestEnvExampleWithIssue() {
	// .env.example exists
	err := os.WriteFile(filepath.Join(s.tempDir, ".env.example"), []byte("KEY=value"), 0644)
	s.Require().NoError(err)

	// .env exists but not ignored
	err = os.WriteFile(filepath.Join(s.tempDir, ".env"), []byte("KEY=secret"), 0644)
	s.Require().NoError(err)

	check := &EnvCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, ".env.example")
	s.Contains(result.Message, ".env file exists but not in .gitignore")
}

func (s *EnvCheckTestSuite) TestNoEnvConfig() {
	check := &EnvCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "No environment configuration found")
}

func (s *EnvCheckTestSuite) TestResultLanguage() {
	check := &EnvCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangCommon, result.Language)
}

func (s *EnvCheckTestSuite) TestGitignoreWithComments() {
	gitignore := `# Environment files
.env
# Don't ignore example
!.env.example
`
	err := os.WriteFile(filepath.Join(s.tempDir, ".gitignore"), []byte(gitignore), 0644)
	s.Require().NoError(err)

	check := &EnvCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Contains(result.Message, ".env in .gitignore")
}

func TestEnvCheckTestSuite(t *testing.T) {
	suite.Run(t, new(EnvCheckTestSuite))
}
