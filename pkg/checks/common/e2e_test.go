package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type E2ECheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *E2ECheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "e2e-test-*")
	s.Require().NoError(err)
}

func (s *E2ECheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *E2ECheckTestSuite) TestIDAndName() {
	check := &E2ECheck{}
	s.Equal("common:e2e", check.ID())
	s.Equal("E2E Tests", check.Name())
}

func (s *E2ECheckTestSuite) TestCypressConfigJs() {
	err := os.WriteFile(filepath.Join(s.tempDir, "cypress.config.js"), []byte(`
module.exports = {
  e2e: {
    baseUrl: 'http://localhost:3000'
  }
};
`), 0644)
	s.Require().NoError(err)

	check := &E2ECheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Cypress")
}

func (s *E2ECheckTestSuite) TestCypressConfigTs() {
	err := os.WriteFile(filepath.Join(s.tempDir, "cypress.config.ts"), []byte(`
export default defineConfig({});
`), 0644)
	s.Require().NoError(err)

	check := &E2ECheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Cypress")
}

func (s *E2ECheckTestSuite) TestPlaywrightConfig() {
	err := os.WriteFile(filepath.Join(s.tempDir, "playwright.config.ts"), []byte(`
import { defineConfig } from '@playwright/test';
export default defineConfig({});
`), 0644)
	s.Require().NoError(err)

	check := &E2ECheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Playwright")
}

func (s *E2ECheckTestSuite) TestWebdriverIO() {
	err := os.WriteFile(filepath.Join(s.tempDir, "wdio.conf.js"), []byte(`
exports.config = {};
`), 0644)
	s.Require().NoError(err)

	check := &E2ECheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "WebdriverIO")
}

func (s *E2ECheckTestSuite) TestNightwatch() {
	err := os.WriteFile(filepath.Join(s.tempDir, "nightwatch.conf.js"), []byte(`
module.exports = {};
`), 0644)
	s.Require().NoError(err)

	check := &E2ECheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Nightwatch")
}

func (s *E2ECheckTestSuite) TestTestCafe() {
	err := os.WriteFile(filepath.Join(s.tempDir, ".testcaferc.json"), []byte(`
{
  "browsers": ["chrome"]
}
`), 0644)
	s.Require().NoError(err)

	check := &E2ECheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "TestCafe")
}

func (s *E2ECheckTestSuite) TestCodeceptJS() {
	err := os.WriteFile(filepath.Join(s.tempDir, "codecept.conf.js"), []byte(`
exports.config = {};
`), 0644)
	s.Require().NoError(err)

	check := &E2ECheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "CodeceptJS")
}

func (s *E2ECheckTestSuite) TestDetox() {
	err := os.WriteFile(filepath.Join(s.tempDir, ".detoxrc.json"), []byte(`
{
  "testRunner": "jest"
}
`), 0644)
	s.Require().NoError(err)

	check := &E2ECheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Detox")
}

func (s *E2ECheckTestSuite) TestE2EDirectory() {
	e2eDir := filepath.Join(s.tempDir, "e2e")
	err := os.MkdirAll(e2eDir, 0755)
	s.Require().NoError(err)

	err = os.WriteFile(filepath.Join(e2eDir, "test.spec.ts"), []byte(`
describe('E2E', () => {});
`), 0644)
	s.Require().NoError(err)

	check := &E2ECheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "e2e directory")
}

func (s *E2ECheckTestSuite) TestCypressDirectory() {
	cypressDir := filepath.Join(s.tempDir, "cypress")
	err := os.MkdirAll(cypressDir, 0755)
	s.Require().NoError(err)

	check := &E2ECheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "cypress directory")
}

func (s *E2ECheckTestSuite) TestPlaywrightDirectory() {
	playwrightDir := filepath.Join(s.tempDir, "playwright")
	err := os.MkdirAll(playwrightDir, 0755)
	s.Require().NoError(err)

	check := &E2ECheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "playwright directory")
}

func (s *E2ECheckTestSuite) TestE2EDependencyInPackageJson() {
	content := `{
  "name": "my-app",
  "devDependencies": {
    "cypress": "^12.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &E2ECheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Cypress")
}

func (s *E2ECheckTestSuite) TestPlaywrightTestInPackageJson() {
	content := `{
  "name": "my-app",
  "devDependencies": {
    "@playwright/test": "^1.40.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &E2ECheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Playwright")
}

func (s *E2ECheckTestSuite) TestPuppeteerInPackageJson() {
	content := `{
  "name": "my-app",
  "devDependencies": {
    "puppeteer": "^21.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &E2ECheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Puppeteer")
}

func (s *E2ECheckTestSuite) TestSeleniumInPackageJson() {
	content := `{
  "name": "my-app",
  "devDependencies": {
    "selenium-webdriver": "^4.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &E2ECheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Selenium")
}

func (s *E2ECheckTestSuite) TestPythonSelenium() {
	content := `selenium==4.15.0
pytest==7.4.0`
	err := os.WriteFile(filepath.Join(s.tempDir, "requirements.txt"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &E2ECheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Selenium")
}

func (s *E2ECheckTestSuite) TestPythonPlaywright() {
	content := `playwright==1.40.0`
	err := os.WriteFile(filepath.Join(s.tempDir, "requirements.txt"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &E2ECheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Playwright")
}

func (s *E2ECheckTestSuite) TestPythonRobot() {
	content := `robotframework==6.0`
	err := os.WriteFile(filepath.Join(s.tempDir, "requirements.txt"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &E2ECheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Robot")
}

func (s *E2ECheckTestSuite) TestBehave() {
	err := os.WriteFile(filepath.Join(s.tempDir, "behave.ini"), []byte(`
[behave]
format=pretty
`), 0644)
	s.Require().NoError(err)

	check := &E2ECheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Behave")
}

func (s *E2ECheckTestSuite) TestCucumber() {
	supportDir := filepath.Join(s.tempDir, "features", "support")
	err := os.MkdirAll(supportDir, 0755)
	s.Require().NoError(err)

	err = os.WriteFile(filepath.Join(supportDir, "env.rb"), []byte(`
require 'cucumber'
`), 0644)
	s.Require().NoError(err)

	check := &E2ECheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Cucumber")
}

func (s *E2ECheckTestSuite) TestNoE2EFound() {
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(`{
  "name": "my-app",
  "dependencies": {}
}`), 0644)
	s.Require().NoError(err)

	check := &E2ECheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "No E2E tests found")
}

func (s *E2ECheckTestSuite) TestEmptyDirectory() {
	check := &E2ECheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "No E2E tests found")
}

func TestE2ECheckTestSuite(t *testing.T) {
	suite.Run(t, new(E2ECheckTestSuite))
}
