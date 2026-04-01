package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

type NamingCheckTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *NamingCheckTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "naming-test-*")
	s.Require().NoError(err)
}

func (s *NamingCheckTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *NamingCheckTestSuite) createFile(name string) {
	filePath := filepath.Join(s.tempDir, name)
	dir := filepath.Dir(filePath)
	if dir != s.tempDir {
		err := os.MkdirAll(dir, 0755)
		s.Require().NoError(err)
	}
	err := os.WriteFile(filePath, []byte("// placeholder"), 0644)
	s.Require().NoError(err)
}

func (s *NamingCheckTestSuite) TestIDAndName() {
	check := &NamingCheck{}
	s.Equal("common:naming", check.ID())
	s.Equal("Naming Consistency", check.Name())
}

func (s *NamingCheckTestSuite) TestConsistentSnakeCase() {
	s.createFile("user_service.go")
	s.createFile("user_handler.go")
	s.createFile("auth_middleware.go")
	s.createFile("db_connection.go")

	check := &NamingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "snake_case")
	s.Contains(result.Reason, "Consistent")
}

func (s *NamingCheckTestSuite) TestConsistentCamelCase() {
	s.createFile("userService.js")
	s.createFile("authHandler.js")
	s.createFile("dbConnection.js")
	s.createFile("apiClient.js")

	check := &NamingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "camelCase")
}

func (s *NamingCheckTestSuite) TestConsistentPascalCase() {
	s.createFile("UserService.java")
	s.createFile("AuthHandler.java")
	s.createFile("DbConnection.java")
	s.createFile("ApiClient.java")

	check := &NamingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "PascalCase")
}

func (s *NamingCheckTestSuite) TestConsistentKebabCase() {
	s.createFile("user-service.ts")
	s.createFile("auth-handler.ts")
	s.createFile("db-connection.ts")
	s.createFile("api-client.ts")

	check := &NamingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "kebab-case")
}

func (s *NamingCheckTestSuite) TestMixedConventions() {
	// Create a mix of naming styles with no dominant one
	s.createFile("user_service.go")
	s.createFile("authHandler.go")
	s.createFile("DbConnection.go")
	s.createFile("api-client.go")

	check := &NamingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "Mixed")
}

func (s *NamingCheckTestSuite) TestMostlyConsistent() {
	// 3 snake_case, 1 camelCase = 75% dominant
	s.createFile("user_service.go")
	s.createFile("auth_handler.go")
	s.createFile("db_connection.go")
	s.createFile("apiClient.go")

	check := &NamingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.False(result.Passed)
	s.Equal(checker.Warn, result.Status)
	s.Contains(result.Reason, "Mostly")
	s.Contains(result.Reason, "snake_case")
}

func (s *NamingCheckTestSuite) TestTooFewFiles() {
	s.createFile("main.go")
	s.createFile("util.go")

	check := &NamingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Too few")
}

func (s *NamingCheckTestSuite) TestEmptyDirectory() {
	check := &NamingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Too few")
}

func (s *NamingCheckTestSuite) TestSkipsNodeModules() {
	// Files in node_modules should be ignored
	s.createFile("user_service.go")
	s.createFile("auth_handler.go")
	s.createFile("db_connection.go")
	s.createFile("node_modules/somePackage.js")
	s.createFile("node_modules/anotherPackage.js")

	check := &NamingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "snake_case")
}

func (s *NamingCheckTestSuite) TestSkipsVendor() {
	s.createFile("user_service.go")
	s.createFile("auth_handler.go")
	s.createFile("db_connection.go")
	s.createFile("vendor/SomePackage.go")

	check := &NamingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "snake_case")
}

func (s *NamingCheckTestSuite) TestIgnoresNonCodeFiles() {
	s.createFile("user_service.go")
	s.createFile("auth_handler.go")
	s.createFile("db_connection.go")
	s.createFile("MyDocument.md")
	s.createFile("SomeConfig.yaml")

	check := &NamingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "snake_case")
}

func (s *NamingCheckTestSuite) TestStripsTestSuffixes() {
	// Go test files should have _test stripped before classification
	s.createFile("user_service.go")
	s.createFile("user_service_test.go")
	s.createFile("auth_handler.go")
	s.createFile("auth_handler_test.go")

	check := &NamingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "snake_case")
}

func (s *NamingCheckTestSuite) TestESLintNamingConvention() {
	content := `{
  "rules": {
    "@typescript-eslint/naming-convention": "error"
  }
}`
	err := os.WriteFile(filepath.Join(s.tempDir, ".eslintrc.json"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &NamingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "ESLint naming-convention")
}

func (s *NamingCheckTestSuite) TestPylintNaming() {
	content := `[FORMAT]
variable-naming-style=snake_case
function-naming-style=snake_case
`
	err := os.WriteFile(filepath.Join(s.tempDir, ".pylintrc"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &NamingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Pylint naming")
}

func (s *NamingCheckTestSuite) TestRubocopNaming() {
	content := `Naming/MethodName:
  EnforcedStyle: snake_case
`
	err := os.WriteFile(filepath.Join(s.tempDir, ".rubocop.yml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &NamingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "RuboCop Naming")
}

func (s *NamingCheckTestSuite) TestCheckstyleNaming() {
	content := `<?xml version="1.0"?>
<!DOCTYPE module PUBLIC "-//Checkstyle//DTD Checkstyle Configuration 1.3//EN" "https://checkstyle.org/dtds/configuration_1_3.dtd">
<module name="Checker">
  <module name="TreeWalker">
    <module name="NamingConvention"/>
  </module>
</module>`
	err := os.WriteFile(filepath.Join(s.tempDir, "checkstyle.xml"), []byte(content), 0644)
	s.Require().NoError(err)

	check := &NamingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Checkstyle naming")
}

func (s *NamingCheckTestSuite) TestClippyConfig() {
	err := os.WriteFile(filepath.Join(s.tempDir, "clippy.toml"), []byte("too-many-arguments-threshold = 10"), 0644)
	s.Require().NoError(err)

	check := &NamingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "Clippy")
}

func (s *NamingCheckTestSuite) TestConsistentFlatNames() {
	s.createFile("main.go")
	s.createFile("server.go")
	s.createFile("handler.go")
	s.createFile("config.go")

	check := &NamingCheck{}
	result, err := check.Run(s.tempDir)

	s.NoError(err)
	s.True(result.Passed)
	s.Equal(checker.Pass, result.Status)
	s.Contains(result.Reason, "lowercase")
}

func TestNamingCheckTestSuite(t *testing.T) {
	suite.Run(t, new(NamingCheckTestSuite))
}

// Unit tests for classifyName
func TestClassifyName(t *testing.T) {
	tests := []struct {
		name     string
		expected namingConvention
	}{
		{"user_service", conventionSnakeCase},
		{"auth_handler", conventionSnakeCase},
		{"my_file_name", conventionSnakeCase},
		{"userService", conventionCamelCase},
		{"authHandler", conventionCamelCase},
		{"UserService", conventionPascalCase},
		{"AuthHandler", conventionPascalCase},
		{"user-service", conventionKebabCase},
		{"auth-handler", conventionKebabCase},
		{"main", conventionFlat},
		{"server", conventionFlat},
		{"config", conventionFlat},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := classifyName(tt.name)
			if result != tt.expected {
				t.Errorf("classifyName(%q) = %d (%s), want %d (%s)",
					tt.name, result, conventionName(result), tt.expected, conventionName(tt.expected))
			}
		})
	}
}

func TestHasUpperAfterFirst(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"userService", true},
		{"UserService", true},
		{"main", false},
		{"server", false},
		{"A", false},
		{"Ab", false},
		{"aB", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasUpperAfterFirst(tt.name)
			if result != tt.expected {
				t.Errorf("hasUpperAfterFirst(%q) = %v, want %v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestStripTestAffixes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user_service_test", "user_service"},
		{"test_user_service", "user_service"},
		{"userService.test", "userService"},
		{"userService.spec", "userService"},
		{"main", "main"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := stripTestAffixes(tt.input)
			if result != tt.expected {
				t.Errorf("stripTestAffixes(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
