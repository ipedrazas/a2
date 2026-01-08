package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type RetryCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *RetryCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "retry-test-*")
	s.Require().NoError(err)
}

func (s *RetryCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *RetryCheckTestSuite) TestIDAndName() {
	check := &RetryCheck{}
	s.Equal("common:retry", check.ID())
	s.Equal("Retry Logic", check.Name())
}

func (s *RetryCheckTestSuite) TestGoBackoff() {
	content := `module myapp

go 1.21

require (
	github.com/cenkalti/backoff/v4 v4.2.0
)`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &RetryCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "backoff")
}

func (s *RetryCheckTestSuite) TestGoRetryGo() {
	content := `module myapp

go 1.21

require (
	github.com/avast/retry-go/v4 v4.5.0
)`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &RetryCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "retry-go")
}

func (s *RetryCheckTestSuite) TestGoRetryableHttp() {
	content := `module myapp

go 1.21

require (
	github.com/hashicorp/go-retryablehttp v0.7.0
)`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &RetryCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "go-retryablehttp")
}

func (s *RetryCheckTestSuite) TestGoCircuitBreaker() {
	content := `module myapp

go 1.21

require (
	github.com/sony/gobreaker v0.5.0
)`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &RetryCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "gobreaker")
}

func (s *RetryCheckTestSuite) TestGoHystrix() {
	content := `module myapp

go 1.21

require (
	github.com/afex/hystrix-go v0.0.0
)`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &RetryCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "hystrix-go")
}

func (s *RetryCheckTestSuite) TestGoFailsafe() {
	content := `module myapp

go 1.21

require (
	github.com/failsafe-go/failsafe-go v0.6.0
)`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &RetryCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "failsafe-go")
}

func (s *RetryCheckTestSuite) TestNodeAsyncRetry() {
	content := `{
  "name": "my-app",
  "dependencies": {
    "async-retry": "^1.3.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &RetryCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "async-retry")
}

func (s *RetryCheckTestSuite) TestNodeAxiosRetry() {
	content := `{
  "name": "my-app",
  "dependencies": {
    "axios-retry": "^4.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &RetryCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "axios-retry")
}

func (s *RetryCheckTestSuite) TestNodePRetry() {
	content := `{
  "name": "my-app",
  "dependencies": {
    "p-retry": "^6.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &RetryCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "p-retry")
}

func (s *RetryCheckTestSuite) TestNodeGot() {
	content := `{
  "name": "my-app",
  "dependencies": {
    "got": "^14.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &RetryCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "got")
}

func (s *RetryCheckTestSuite) TestNodeCockatiel() {
	content := `{
  "name": "my-app",
  "dependencies": {
    "cockatiel": "^3.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &RetryCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Cockatiel")
}

func (s *RetryCheckTestSuite) TestNodeOpossum() {
	content := `{
  "name": "my-app",
  "dependencies": {
    "opossum": "^8.0.0"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, "package.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &RetryCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Opossum")
}

func (s *RetryCheckTestSuite) TestPythonTenacity() {
	content := `tenacity==8.2.0`
	err := os.WriteFile(filepath.Join(s.tempDir, "requirements.txt"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &RetryCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Tenacity")
}

func (s *RetryCheckTestSuite) TestPythonBackoff() {
	content := `backoff==2.2.0`
	err := os.WriteFile(filepath.Join(s.tempDir, "requirements.txt"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &RetryCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "backoff")
}

func (s *RetryCheckTestSuite) TestPythonUrllib3() {
	content := `urllib3==2.1.0`
	err := os.WriteFile(filepath.Join(s.tempDir, "requirements.txt"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &RetryCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "urllib3")
}

func (s *RetryCheckTestSuite) TestPythonCircuitBreaker() {
	content := `circuitbreaker==2.0.0`
	err := os.WriteFile(filepath.Join(s.tempDir, "requirements.txt"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &RetryCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "circuitbreaker")
}

func (s *RetryCheckTestSuite) TestJavaSpringRetry() {
	content := `<dependency>
    <groupId>org.springframework.retry</groupId>
    <artifactId>spring-retry</artifactId>
</dependency>`
	err := os.WriteFile(filepath.Join(s.tempDir, "pom.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &RetryCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Spring Retry")
}

func (s *RetryCheckTestSuite) TestJavaResilience4j() {
	content := `implementation 'io.github.resilience4j:resilience4j-retry:2.0.0'`
	err := os.WriteFile(filepath.Join(s.tempDir, "build.gradle"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &RetryCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Resilience4j")
}

func (s *RetryCheckTestSuite) TestJavaFailsafe() {
	content := `implementation 'dev.failsafe:failsafe:3.3.0'`
	err := os.WriteFile(filepath.Join(s.tempDir, "build.gradle"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &RetryCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Failsafe")
}

func (s *RetryCheckTestSuite) TestRustBackoff() {
	content := `[package]
name = "myapp"

[dependencies]
backoff = "0.4"`
	err := os.WriteFile(filepath.Join(s.tempDir, "Cargo.toml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &RetryCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "backoff")
}

func (s *RetryCheckTestSuite) TestRustTokioRetry() {
	content := `[package]
name = "myapp"

[dependencies]
tokio-retry = "0.3"`
	err := os.WriteFile(filepath.Join(s.tempDir, "Cargo.toml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &RetryCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "tokio-retry")
}

func (s *RetryCheckTestSuite) TestRustAgain() {
	content := `[package]
name = "myapp"

[dependencies]
again = "0.1"`
	err := os.WriteFile(filepath.Join(s.tempDir, "Cargo.toml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &RetryCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "again")
}

func (s *RetryCheckTestSuite) TestMultipleLibraries() {
	content := `module myapp

go 1.21

require (
	github.com/cenkalti/backoff/v4 v4.2.0
	github.com/sony/gobreaker v0.5.0
)`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &RetryCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "backoff")
	s.Contains(result.Message, "gobreaker")
}

func (s *RetryCheckTestSuite) TestNoRetryLogicFound() {
	content := `module myapp

go 1.21

require (
	github.com/gin-gonic/gin v1.9.0
)`
	err := os.WriteFile(filepath.Join(s.tempDir, "go.mod"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &RetryCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "No retry logic found")
}

func (s *RetryCheckTestSuite) TestEmptyDirectory() {
	check := &RetryCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "No retry logic found")
}

func TestRetryCheckTestSuite(t *testing.T) {
	suite.Run(t, new(RetryCheckTestSuite))
}
