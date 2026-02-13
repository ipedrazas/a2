package typescriptcheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type CoverageTestSuite struct {
	suite.Suite
	tempDir string
	check   *CoverageCheck
}

func (s *CoverageTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "ts-coverage-test-*")
	s.NoError(err)
	s.tempDir = dir
	s.check = &CoverageCheck{Threshold: 80.0}
}

func (s *CoverageTestSuite) TearDownTest() {
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

func (s *CoverageTestSuite) writeFile(name, content string) {
	path := filepath.Join(s.tempDir, name)
	dir := filepath.Dir(path)
	if dir != s.tempDir {
		err := os.MkdirAll(dir, 0755)
		s.NoError(err)
	}
	err := os.WriteFile(path, []byte(content), 0644)
	s.NoError(err)
}

func (s *CoverageTestSuite) TestID() {
	s.Equal("typescript:coverage", s.check.ID())
}

func (s *CoverageTestSuite) TestName() {
	s.Equal("TypeScript Coverage", s.check.Name())
}

func (s *CoverageTestSuite) TestRun_NoTsconfig() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Reason, "No tsconfig.json found")
}

func (s *CoverageTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangTypeScript, result.Language)
}

func (s *CoverageTestSuite) TestHasJestCoverage_WithConfig() {
	s.writeFile("jest.config.js", `module.exports = { collectCoverage: true }`)
	s.True(s.check.hasJestCoverage(s.tempDir))
}

func (s *CoverageTestSuite) TestHasJestCoverage_NoConfig() {
	s.False(s.check.hasJestCoverage(s.tempDir))
}

func (s *CoverageTestSuite) TestHasVitestCoverage_WithConfig() {
	s.writeFile("vitest.config.ts", `export default { coverage: {} }`)
	s.True(s.check.hasVitestCoverage(s.tempDir))
}

func (s *CoverageTestSuite) TestHasVitestCoverage_ViteConfig() {
	// Vitest often configured in vite.config.ts via vitest/config
	s.writeFile("vite.config.ts", `import { defineConfig } from "vitest/config";
export default defineConfig({ test: { environment: "jsdom" } });`)
	s.True(s.check.hasVitestCoverage(s.tempDir))
}

func (s *CoverageTestSuite) TestHasC8_WithDep() {
	s.writeFile("package.json", `{"devDependencies": {"c8": "^7.0.0"}}`)
	s.True(s.check.hasC8(s.tempDir))
}

func (s *CoverageTestSuite) TestHasC8_WithConfig() {
	s.writeFile(".c8rc.json", `{}`)
	s.True(s.check.hasC8(s.tempDir))
}

func (s *CoverageTestSuite) TestHasNYC_WithDep() {
	s.writeFile("package.json", `{"devDependencies": {"nyc": "^15.0.0"}}`)
	s.True(s.check.hasNYC(s.tempDir))
}

func (s *CoverageTestSuite) TestCheckCICoverage_Codecov() {
	s.writeFile(".github/workflows/ci.yml", `
jobs:
  test:
    steps:
      - uses: codecov/codecov-action@v3
`)
	result := s.check.checkCICoverage(s.tempDir)
	s.Equal("Codecov (CI)", result)
}

func (s *CoverageTestSuite) TestCheckCICoverage_None() {
	result := s.check.checkCICoverage(s.tempDir)
	s.Equal("", result)
}

func (s *CoverageTestSuite) TestFindCoverageReports_NoCoverage() {
	cov := s.check.findCoverageReports(s.tempDir)
	s.Equal(-1.0, cov)
}

func (s *CoverageTestSuite) TestExtractCoverage_Cobertura() {
	content := `<?xml version="1.0"?>
<coverage line-rate="0.85" branch-rate="0.75">
</coverage>`
	cov := extractCoverage(content)
	s.Equal(85.0, cov)
}

func (s *CoverageTestSuite) TestExtractCoverage_LCOV() {
	content := `SF:/path/to/file.ts
LF:100
LH:85
end_of_record`
	cov := extractCoverage(content)
	s.Equal(85.0, cov)
}

func (s *CoverageTestSuite) TestExtractCoverage_JSON() {
	content := `{
  "total": {
    "lines": {"total": 100, "covered": 85, "pct": 85.0}
  }
}`
	cov := extractCoverage(content)
	s.Equal(85.0, cov)
}

func (s *CoverageTestSuite) TestExtractCoverage_Invalid() {
	cov := extractCoverage("invalid content")
	s.Equal(-1.0, cov)
}

func TestCoverageTestSuite(t *testing.T) {
	suite.Run(t, new(CoverageTestSuite))
}
