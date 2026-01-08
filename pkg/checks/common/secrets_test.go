package common

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type SecretsCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *SecretsCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "secrets-test-*")
	s.Require().NoError(err)
}

func (s *SecretsCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *SecretsCheckTestSuite) TestIDAndName() {
	check := &SecretsCheck{}
	s.Equal("common:secrets", check.ID())
	s.Equal("Secrets Detection", check.Name())
}

func (s *SecretsCheckTestSuite) TestGitleaksConfigured() {
	// Create gitleaks config
	err := os.WriteFile(filepath.Join(s.tempDir, ".gitleaks.toml"), []byte(`
title = "gitleaks config"

[[rules]]
id = "aws-access-key"
description = "AWS Access Key"
regex = "AKIA[0-9A-Z]{16}"
`), 0644)
	s.Require().NoError(err)

	check := &SecretsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "Gitleaks")
}

func (s *SecretsCheckTestSuite) TestTruffleHogConfigured() {
	err := os.WriteFile(filepath.Join(s.tempDir, ".trufflehog.yml"), []byte(`
detectors:
  - AWS
`), 0644)
	s.Require().NoError(err)

	check := &SecretsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "TruffleHog")
}

func (s *SecretsCheckTestSuite) TestDetectSecretsConfigured() {
	err := os.WriteFile(filepath.Join(s.tempDir, ".secrets.baseline"), []byte(`{
    "version": "1.0.0",
    "results": {}
}`), 0644)
	s.Require().NoError(err)

	check := &SecretsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "detect-secrets")
}

func (s *SecretsCheckTestSuite) TestPreCommitWithGitleaks() {
	err := os.WriteFile(filepath.Join(s.tempDir, ".pre-commit-config.yaml"), []byte(`
repos:
  - repo: https://github.com/gitleaks/gitleaks
    rev: v8.18.0
    hooks:
      - id: gitleaks
`), 0644)
	s.Require().NoError(err)

	check := &SecretsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Message, "pre-commit")
}

func (s *SecretsCheckTestSuite) TestNoScannerNoSecrets() {
	// Create a file without secrets
	code := `package main

func main() {
	fmt.Println("Hello, World!")
}
`
	err := os.WriteFile(filepath.Join(s.tempDir, "main.go"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &SecretsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "No secret scanning configured")
}

func (s *SecretsCheckTestSuite) TestDetectsAWSAccessKey() {
	code := `package main

const awsKey = "AKIAIOSFODNN7EXAMPLE"
`
	err := os.WriteFile(filepath.Join(s.tempDir, "config.go"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &SecretsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "AWS Access Key")
}

func (s *SecretsCheckTestSuite) TestDetectsPrivateKey() {
	code := `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA0Z...
-----END RSA PRIVATE KEY-----
`
	err := os.WriteFile(filepath.Join(s.tempDir, "key.pem"), []byte(code), 0644)
	s.Require().NoError(err)

	// Need a code file extension
	err = os.WriteFile(filepath.Join(s.tempDir, "config.go"), []byte(`package main

var key = `+"`"+`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA...
-----END RSA PRIVATE KEY-----`+"`"+`
`), 0644)
	s.Require().NoError(err)

	check := &SecretsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "Private Key")
}

func (s *SecretsCheckTestSuite) TestDetectsAPIKey() {
	// Using a test pattern that matches our API key regex
	suffix := "abcdefghij1234567890abcd"
	prefix := "sk_"
	apiKey := fmt.Sprintf("%slive_%s", prefix, suffix)
	code := fmt.Sprintf(`
api_key = "%s"
`, apiKey)

	err := os.WriteFile(filepath.Join(s.tempDir, "settings.py"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &SecretsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "Stripe Key")
}

func (s *SecretsCheckTestSuite) TestDetectsGitHubToken() {
	code := `const token = "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx";
`
	err := os.WriteFile(filepath.Join(s.tempDir, "auth.js"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &SecretsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "GitHub Token")
}

func (s *SecretsCheckTestSuite) TestDetectsDatabaseURL() {
	code := `DATABASE_URL=postgres://user:secretpassword@localhost:5432/mydb
`
	err := os.WriteFile(filepath.Join(s.tempDir, "config.yaml"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &SecretsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "Database URL")
}

func (s *SecretsCheckTestSuite) TestSkipsEnvExample() {
	// Secrets in .env.example should be ignored (they're templates)
	code := `API_KEY="your_api_key_here_abcdefghij12345678"
SECRET="your_secret_key_here_abcdefghij1234"
`
	err := os.WriteFile(filepath.Join(s.tempDir, ".env.example"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &SecretsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	// Should warn about no scanner, not about the secrets
	s.Contains(result.Message, "No secret scanning configured")
}

func (s *SecretsCheckTestSuite) TestSkipsNodeModules() {
	// Create node_modules with secrets (should be skipped)
	nodeModules := filepath.Join(s.tempDir, "node_modules", "some-lib")
	err := os.MkdirAll(nodeModules, 0755)
	s.Require().NoError(err)

	code := `const key = "AKIAIOSFODNN7EXAMPLE";`
	err = os.WriteFile(filepath.Join(nodeModules, "index.js"), []byte(code), 0644)
	s.Require().NoError(err)

	// Create a regular file without secrets
	mainCode := `console.log('hello');`
	err = os.WriteFile(filepath.Join(s.tempDir, "main.js"), []byte(mainCode), 0644)
	s.Require().NoError(err)

	check := &SecretsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	// Should only warn about no scanner, not about secrets in node_modules
	s.Contains(result.Message, "No secret scanning configured")
}

func (s *SecretsCheckTestSuite) TestSkipsVendorDirectory() {
	// Create vendor with secrets (should be skipped)
	vendor := filepath.Join(s.tempDir, "vendor", "github.com", "lib")
	err := os.MkdirAll(vendor, 0755)
	s.Require().NoError(err)

	code := `const key = "AKIAIOSFODNN7EXAMPLE";`
	err = os.WriteFile(filepath.Join(vendor, "auth.go"), []byte(code), 0644)
	s.Require().NoError(err)

	// Create a regular file without secrets
	mainCode := `package main`
	err = os.WriteFile(filepath.Join(s.tempDir, "main.go"), []byte(mainCode), 0644)
	s.Require().NoError(err)

	check := &SecretsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.Contains(result.Message, "No secret scanning configured")
}

func (s *SecretsCheckTestSuite) TestMultipleSecrets() {
	code := `package main

const (
	awsKey = "AKIAIOSFODNN7EXAMPLE"
	token = "ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
)
`
	err := os.WriteFile(filepath.Join(s.tempDir, "config.go"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &SecretsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	// Should mention multiple secrets found
	s.Contains(result.Message, "secrets found")
}

func (s *SecretsCheckTestSuite) TestEmptyDirectory() {
	check := &SecretsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "No secret scanning configured")
}

func (s *SecretsCheckTestSuite) TestJWTToken() {
	code := `const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c";
`
	err := os.WriteFile(filepath.Join(s.tempDir, "auth.js"), []byte(code), 0644)
	s.Require().NoError(err)

	check := &SecretsCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Message, "JWT Token")
}

func TestSecretsCheckTestSuite(t *testing.T) {
	suite.Run(t, new(SecretsCheckTestSuite))
}
