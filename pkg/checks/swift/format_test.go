package swiftcheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/config"
	"github.com/stretchr/testify/suite"
)

type FormatTestSuite struct {
	suite.Suite
	tempDir string
	check   *FormatCheck
}

func (s *FormatTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "swift-format-test-*")
	s.NoError(err)
	s.tempDir = dir
	s.check = &FormatCheck{Config: &config.SwiftLanguageConfig{}}
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
	s.Equal("swift:format", s.check.ID())
}

func (s *FormatTestSuite) TestName() {
	s.Equal("Swift Format", s.check.Name())
}

func (s *FormatTestSuite) TestRun_NoPackageSwift() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Message, "No Package.swift found")
}

func (s *FormatTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangSwift, result.Language)
}

func (s *FormatTestSuite) TestHasFormatterConfig_SwiftFormat() {
	s.writeFile(".swift-format", "{}")
	s.True(s.check.hasFormatterConfig(s.tempDir, "swift-format"))
}

func (s *FormatTestSuite) TestHasFormatterConfig_Swiftformat() {
	s.writeFile(".swiftformat", "")
	s.True(s.check.hasFormatterConfig(s.tempDir, "swiftformat"))
}

func (s *FormatTestSuite) TestHasFormatterConfig_NoConfig() {
	s.False(s.check.hasFormatterConfig(s.tempDir, "swift-format"))
	s.False(s.check.hasFormatterConfig(s.tempDir, "swiftformat"))
	s.False(s.check.hasFormatterConfig(s.tempDir, "unknown"))
}

func (s *FormatTestSuite) TestFormatCount() {
	s.Equal("1", formatCount(1))
	s.Equal("5", formatCount(5))
	s.Equal("9", formatCount(9))
	s.Equal("multiple", formatCount(10))
	s.Equal("multiple", formatCount(100))
}

func TestFormatTestSuite(t *testing.T) {
	suite.Run(t, new(FormatTestSuite))
}
