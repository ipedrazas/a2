package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type IntegrationCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *IntegrationCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "integration-test-*")
	s.Require().NoError(err)
}

func (s *IntegrationCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *IntegrationCheckTestSuite) TestIDAndName() {
	check := &IntegrationCheck{}
	s.Equal("common:integration", check.ID())
	s.Equal("Integration Tests", check.Name())
}

func (s *IntegrationCheckTestSuite) TestIntegrationDirectory() {
	integrationDir := filepath.Join(s.tempDir, "tests", "integration")
	err := os.MkdirAll(integrationDir, 0755)
	s.Require().NoError(err)

	// Create a test file
	err = os.WriteFile(filepath.Join(integrationDir, "test_api.py"), []byte("def test_api(): pass"), 0644)
	s.Require().NoError(err)

	check := &IntegrationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "tests/integration/")
}

func (s *IntegrationCheckTestSuite) TestE2EDirectory() {
	e2eDir := filepath.Join(s.tempDir, "e2e")
	err := os.MkdirAll(e2eDir, 0755)
	s.Require().NoError(err)

	err = os.WriteFile(filepath.Join(e2eDir, "test.spec.ts"), []byte("test()"), 0644)
	s.Require().NoError(err)

	check := &IntegrationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "e2e/")
}

func (s *IntegrationCheckTestSuite) TestIntegrationTestsDirectory() {
	integrationDir := filepath.Join(s.tempDir, "integration_tests")
	err := os.MkdirAll(integrationDir, 0755)
	s.Require().NoError(err)

	err = os.WriteFile(filepath.Join(integrationDir, "test.go"), []byte("package integration"), 0644)
	s.Require().NoError(err)

	check := &IntegrationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "integration_tests/")
}

func (s *IntegrationCheckTestSuite) TestGoIntegrationTestFile() {
	content := `package myapp

import "testing"

func TestIntegration_UserAPI(t *testing.T) {
	// Integration test
}
`
	err := os.WriteFile(filepath.Join(s.tempDir, "user_integration_test.go"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &IntegrationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Go integration tests")
}

func (s *IntegrationCheckTestSuite) TestDockerComposeTest() {
	content := `version: '3.8'
services:
  db:
    image: postgres:15
  app:
    build: .
    depends_on:
      - db
`
	err := os.WriteFile(filepath.Join(s.tempDir, "docker-compose.test.yml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &IntegrationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "docker-compose.test.yml")
}

func (s *IntegrationCheckTestSuite) TestDockerComposeE2E() {
	content := `version: '3.8'
services:
  selenium:
    image: selenium/standalone-chrome
`
	err := os.WriteFile(filepath.Join(s.tempDir, "docker-compose.e2e.yaml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &IntegrationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "docker-compose.e2e.yaml")
}

func (s *IntegrationCheckTestSuite) TestTestcontainersGo() {
	content := `module myapp

go 1.21

require github.com/testcontainers/testcontainers-go v0.25.0
`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &IntegrationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "testcontainers")
}

func (s *IntegrationCheckTestSuite) TestTestcontainersPython() {
	content := `[project]
name = "myapp"
dependencies = [
    "testcontainers>=3.7.0",
    "pytest>=7.0.0",
]
`
	err := os.WriteFile(filepath.Join(s.tempDir, "pyproject.toml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &IntegrationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "testcontainers")
}

func (s *IntegrationCheckTestSuite) TestTestcontainersNode() {
	content := `{
  "name": "myapp",
  "devDependencies": {
    "testcontainers": "^10.0.0",
    "jest": "^29.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &IntegrationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "testcontainers")
}

func (s *IntegrationCheckTestSuite) TestCypressConfig() {
	content := `const { defineConfig } = require('cypress')

module.exports = defineConfig({
  e2e: {
    baseUrl: 'http://localhost:3000',
  },
})
`
	err := os.WriteFile(filepath.Join(s.tempDir, "cypress.config.js"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &IntegrationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Cypress")
}

func (s *IntegrationCheckTestSuite) TestPlaywrightConfig() {
	content := `import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: './tests',
  use: {
    baseURL: 'http://localhost:3000',
  },
});
`
	err := os.WriteFile(filepath.Join(s.tempDir, "playwright.config.ts"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &IntegrationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Playwright")
}

func (s *IntegrationCheckTestSuite) TestWebdriverIOConfig() {
	content := `exports.config = {
  runner: 'local',
  specs: ['./test/specs/**/*.js'],
  capabilities: [{
    browserName: 'chrome',
  }],
}
`
	err := os.WriteFile(filepath.Join(s.tempDir, "wdio.conf.js"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &IntegrationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "WebdriverIO")
}

func (s *IntegrationCheckTestSuite) TestSeleniumInPackageJson() {
	content := `{
  "name": "myapp",
  "devDependencies": {
    "selenium-webdriver": "^4.0.0",
    "mocha": "^10.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &IntegrationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Selenium")
}

func (s *IntegrationCheckTestSuite) TestMultipleIntegrationIndicators() {
	// Integration directory
	integrationDir := filepath.Join(s.tempDir, "tests", "integration")
	err := os.MkdirAll(integrationDir, 0755)
	s.Require().NoError(err)
	err = os.WriteFile(filepath.Join(integrationDir, "test.py"), []byte(""), 0644)
	s.Require().NoError(err)

	// Cypress config
	content := `module.exports = { e2e: {} }`
	err = os.WriteFile(filepath.Join(s.tempDir, "cypress.config.js"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &IntegrationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "tests/integration/")
	s.Contains(result.Reason, "Cypress")
}

func (s *IntegrationCheckTestSuite) TestNoIntegrationTests() {
	// Create some unit test file but no integration tests
	err := os.WriteFile(filepath.Join(s.tempDir, "user_test.go"), []byte("package myapp\n\nfunc TestUser(t *testing.T) {}"), 0644)
	s.Require().NoError(err)

	check := &IntegrationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "No integration tests found")
}

func (s *IntegrationCheckTestSuite) TestEmptyDirectory() {
	check := &IntegrationCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "No integration tests found")
}

func TestIntegrationCheckTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationCheckTestSuite))
}
