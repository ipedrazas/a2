package typescriptcheck

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type LoggingTestSuite struct {
	suite.Suite
	tempDir string
	check   *LoggingCheck
}

func (s *LoggingTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "ts-logging-test-*")
	s.NoError(err)
	s.tempDir = dir
	s.check = &LoggingCheck{}
}

func (s *LoggingTestSuite) TearDownTest() {
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

func (s *LoggingTestSuite) writeFile(name, content string) {
	path := filepath.Join(s.tempDir, name)
	dir := filepath.Dir(path)
	if dir != s.tempDir {
		err := os.MkdirAll(dir, 0755)
		s.NoError(err)
	}
	err := os.WriteFile(path, []byte(content), 0644)
	s.NoError(err)
}

func (s *LoggingTestSuite) TestID() {
	s.Equal("typescript:logging", s.check.ID())
}

func (s *LoggingTestSuite) TestName() {
	s.Equal("TypeScript Logging", s.check.Name())
}

func (s *LoggingTestSuite) TestRun_NoTsconfig() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Fail, result.Status)
	s.Contains(result.Reason, "No tsconfig.json found")
}

func (s *LoggingTestSuite) TestRun_ResultLanguage() {
	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.Equal(checker.LangTypeScript, result.Language)
}

func (s *LoggingTestSuite) TestRun_WithLogLayer() {
	s.writeFile("tsconfig.json", `{}`)
	s.writeFile("package.json", `{
		"name": "test-app",
		"dependencies": {
			"loglayer": "^6.0.0"
		}
	}`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "LogLayer")
}

func (s *LoggingTestSuite) TestRun_WithWinston() {
	s.writeFile("tsconfig.json", `{}`)
	s.writeFile("package.json", `{
		"name": "test-app",
		"dependencies": {
			"winston": "^3.0.0"
		}
	}`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Winston")
}

func (s *LoggingTestSuite) TestRun_LogLayerInDevDependencies() {
	s.writeFile("tsconfig.json", `{}`)
	s.writeFile("package.json", `{
		"name": "test-app",
		"devDependencies": {
			"loglayer": "^6.0.0"
		}
	}`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "LogLayer")
}

func (s *LoggingTestSuite) TestRun_WithLoggingLibAndConsole() {
	s.writeFile("tsconfig.json", `{}`)
	s.writeFile("package.json", `{
		"name": "test-app",
		"dependencies": {
			"loglayer": "^6.0.0"
		}
	}`)
	s.writeFile("src/index.ts", `console.log("debug")`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "LogLayer")
	s.Contains(result.Reason, "console.log")
}

func (s *LoggingTestSuite) TestRun_NoLoggingLib() {
	s.writeFile("tsconfig.json", `{}`)
	s.writeFile("package.json", `{"name": "test-app"}`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "No logging library detected")
}

func (s *LoggingTestSuite) TestRun_ConsoleLogOnly() {
	s.writeFile("tsconfig.json", `{}`)
	s.writeFile("package.json", `{"name": "test-app"}`)
	s.writeFile("src/index.ts", `
		function main() {
			console.log("Hello world");
			console.error("error");
		}
	`)

	result, err := s.check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "console.log")
}

func TestLoggingTestSuite(t *testing.T) {
	suite.Run(t, new(LoggingTestSuite))
}
