package typescriptcheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type FormatTestSuite struct {
	suite.Suite
	tempDir string
	check   *FormatCheck
}

func (s *FormatTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "ts-format-test-*")
	s.NoError(err)
	s.tempDir = dir
	s.check = &FormatCheck{}
}

func (s *FormatTestSuite) TearDownTest() {
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

func (s *FormatTestSuite) writeFile(name, content string) {
	path := filepath.Join(s.tempDir, name)
	dir := filepath.Dir(path)
	if dir != s.tempDir {
		err := os.MkdirAll(dir, 0755)
		s.NoError(err)
	}
	err := os.WriteFile(path, []byte(content), 0644)
	s.NoError(err)
}

func (s *FormatTestSuite) TestID() {
	s.Equal("typescript:format", s.check.ID())
}

func (s *FormatTestSuite) TestName() {
	s.Equal("TypeScript Format", s.check.Name())
}

func (s *FormatTestSuite) TestRun_NoTsconfig() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Message, "No tsconfig.json found")
}

func (s *FormatTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangTypeScript, result.Language)
}

func (s *FormatTestSuite) TestDetectFormatter_Prettier() {
	s.writeFile(".prettierrc", "{}")
	formatter := s.check.detectFormatter(s.tempDir)
	s.Equal("prettier", formatter)
}

func (s *FormatTestSuite) TestDetectFormatter_Biome() {
	s.writeFile("biome.json", "{}")
	formatter := s.check.detectFormatter(s.tempDir)
	s.Equal("biome", formatter)
}

func (s *FormatTestSuite) TestDetectFormatter_Dprint() {
	s.writeFile("dprint.json", "{}")
	formatter := s.check.detectFormatter(s.tempDir)
	s.Equal("dprint", formatter)
}

func (s *FormatTestSuite) TestDetectFormatter_None() {
	formatter := s.check.detectFormatter(s.tempDir)
	s.Equal("", formatter)
}

func (s *FormatTestSuite) TestCountUnformattedFiles() {
	output := `
src/index.ts
src/utils.tsx
lib/helpers.js
`
	count := countUnformattedFiles(output)
	s.Equal(3, count)
}

func (s *FormatTestSuite) TestCountUnformattedFiles_NoFiles() {
	output := "All files are formatted"
	count := countUnformattedFiles(output)
	s.Equal(0, count)
}

func TestFormatTestSuite(t *testing.T) {
	suite.Run(t, new(FormatTestSuite))
}
