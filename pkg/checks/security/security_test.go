package security

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/stretchr/testify/suite"
)

// SecurityTestSuite is test suite for security checks.
type SecurityTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupTest is called before each test method.
func (suite *SecurityTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "a2-security-test-*")
	suite.NoError(err)
	suite.tempDir = dir
}

// TearDownTest is called after each test method.
func (suite *SecurityTestSuite) TearDownTest() {
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

// createTempFile creates a temporary file with given content.
func (suite *SecurityTestSuite) createTempFile(name, content string) string {
	filePath := filepath.Join(suite.tempDir, name)
	err := os.WriteFile(filePath, []byte(content), 0644)
	suite.NoError(err)
	return filePath
}

// TestShellInjectionCheck_ID tests that ShellInjectionCheck returns correct ID.
func (suite *SecurityTestSuite) TestShellInjectionCheck_ID() {
	check := &ShellInjectionCheck{}
	suite.Equal("security:shell_injection", check.ID())
}

// TestShellInjectionCheck_Name tests that ShellInjectionCheck returns correct name.
func (suite *SecurityTestSuite) TestShellInjectionCheck_Name() {
	check := &ShellInjectionCheck{}
	suite.Equal("Shell Injection Detection", check.Name())
}

// TestShellInjectionCheck_Run_SafeCode tests that safe code passes the check.
func (suite *SecurityTestSuite) TestShellInjectionCheck_Run_SafeCode() {
	suite.createTempFile("main.go", `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`)

	check := &ShellInjectionCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
	suite.Contains(result.Message, "No dangerous")
}

// TestShellInjectionCheck_Run_ShellInjection tests detection of shell injection.
func (suite *SecurityTestSuite) TestShellInjectionCheck_Run_ShellInjection() {
	suite.createTempFile("main.go", `package main

import (
	"os/exec"
	"fmt"
)

func main() {
	userInput := "some command"
	cmd := exec.Command("sh", "-c", userInput)
	fmt.Println(cmd)
}
`)

	check := &ShellInjectionCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	suite.Contains(result.Message, "pattern")
}

// TestFileSystemCheck_ID tests that FileSystemCheck returns correct ID.
func (suite *SecurityTestSuite) TestFileSystemCheck_ID() {
	check := &FileSystemCheck{}
	suite.Equal("security:filesystem", check.ID())
}

// TestFileSystemCheck_Name tests that FileSystemCheck returns correct name.
func (suite *SecurityTestSuite) TestFileSystemCheck_Name() {
	check := &FileSystemCheck{}
	suite.Equal("File System Safety", check.Name())
}

// TestFileSystemCheck_Run_SafeCode tests that safe code passes the check.
func (suite *SecurityTestSuite) TestFileSystemCheck_Run_SafeCode() {
	suite.createTempFile("main.go", `package main

import (
	"io/ioutil"
	"fmt"
)

func main() {
	data, _ := ioutil.ReadFile("config.json")
	fmt.Println(data)
}
`)

	check := &FileSystemCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
	suite.Contains(result.Message, "No path traversal")
}

// TestFileSystemCheck_Run_PathTraversal tests detection of path traversal.
func (suite *SecurityTestSuite) TestFileSystemCheck_Run_PathTraversal() {
	suite.createTempFile("main.go", `package main

import (
	"io/ioutil"
	"fmt"
)

func main() {
	data, _ := ioutil.ReadFile("../../../etc/passwd")
	fmt.Println(data)
}
`)

	check := &FileSystemCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	suite.Contains(result.Message, "Path traversal")
}

// TestNetworkCheck_ID tests that NetworkCheck returns correct ID.
func (suite *SecurityTestSuite) TestNetworkCheck_ID() {
	check := &NetworkCheck{}
	suite.Equal("security:network", check.ID())
}

// TestNetworkCheck_Name tests that NetworkCheck returns correct name.
func (suite *SecurityTestSuite) TestNetworkCheck_Name() {
	check := &NetworkCheck{}
	suite.Equal("Network Exfiltration Detection", check.Name())
}

// TestNetworkCheck_Run_SafeCode tests that safe code passes the check.
func (suite *SecurityTestSuite) TestNetworkCheck_Run_SafeCode() {
	suite.createTempFile("main.go", `package main

import "fmt"

func main() {
	fmt.Println("No network operations here")
}
`)

	check := &NetworkCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
	suite.Contains(result.Message, "No suspicious")
}

// TestNetworkCheck_Run_SuspiciousNetwork tests detection of suspicious network operations.
func (suite *SecurityTestSuite) TestNetworkCheck_Run_SuspiciousNetwork() {
	suite.createTempFile("main.go", `package main

import (
	"net/http"
	"fmt"
)

func main() {
	userURL := "http://suspicious-domain.com/data"
	resp, _ := http.Get(userURL)
	fmt.Println(resp)
}
`)

	check := &NetworkCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Warn, result.Status) // Network check is Warn by default
	suite.Contains(result.Message, "network operation")
}

// TestObfuscationCheck_ID tests that ObfuscationCheck returns correct ID.
func (suite *SecurityTestSuite) TestObfuscationCheck_ID() {
	check := &ObfuscationCheck{}
	suite.Equal("security:obfuscation", check.ID())
}

// TestObfuscationCheck_Name tests that ObfuscationCheck returns correct name.
func (suite *SecurityTestSuite) TestObfuscationCheck_Name() {
	check := &ObfuscationCheck{}
	suite.Equal("Code Obfuscation Detection", check.Name())
}

// TestObfuscationCheck_Run_SafeCode tests that safe code passes the check.
func (suite *SecurityTestSuite) TestObfuscationCheck_Run_SafeCode() {
	suite.createTempFile("main.go", `package main

import "fmt"

func main() {
	message := "Hello, World!"
	fmt.Println(message)
}
`)

	check := &ObfuscationCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.True(result.Passed)
	suite.Equal(checker.Pass, result.Status)
	suite.Contains(result.Message, "No obfuscated")
}

// TestObfuscationCheck_Run_EncodedString tests detection of encoded strings.
func (suite *SecurityTestSuite) TestObfuscationCheck_Run_EncodedString() {
	suite.createTempFile("main.go", `package main

import "encoding/base64"

func main() {
	// Use a longer base64 string to match the common pattern (40+ chars)
	encoded := "SGVsbG8gV29ybGQgV2l0aCBIaWdoIEVudHJvcHkSGVsbG8gV29ybGQgV2l0aCBIaWdoIEVudHJvcHk="
	data, _ := base64.StdEncoding.DecodeString(encoded)
	println(string(data))
}
`)

	check := &ObfuscationCheck{}
	result, err := check.Run(suite.tempDir)

	suite.NoError(err)
	suite.False(result.Passed)
	suite.Equal(checker.Fail, result.Status)
	suite.Contains(result.Message, "obfuscation")
}

// TestRegister tests that Register returns all expected checks.
func (suite *SecurityTestSuite) TestRegister() {
	regs := Register(nil)

	suite.Len(regs, 4, "Should register exactly 4 security checks")

	// Check that all expected check IDs are present
	checkIDs := make(map[string]bool)
	for _, reg := range regs {
		checkIDs[reg.Meta.ID] = true
	}

	suite.True(checkIDs["security:obfuscation"], "Should have obfuscation check")
	suite.True(checkIDs["security:shell_injection"], "Should have shell injection check")
	suite.True(checkIDs["security:filesystem"], "Should have filesystem check")
	suite.True(checkIDs["security:network"], "Should have network check")
}

// TestRegister_Metadata tests that registered checks have correct metadata.
func (suite *SecurityTestSuite) TestRegister_Metadata() {
	regs := Register(nil)

	for _, reg := range regs {
		suite.NotEmpty(reg.Meta.ID, "Check ID should not be empty")
		suite.NotEmpty(reg.Meta.Name, "Check name should not be empty")
		suite.NotEmpty(reg.Meta.Description, "Check description should not be empty")
		suite.NotEmpty(reg.Meta.Languages, "Check should have at least one language")
		suite.NotEmpty(reg.Meta.Order, "Check should have an order value")
		suite.NotEmpty(reg.Meta.Suggestion, "Check should have a suggestion")

		// All security checks should apply to LangCommon
		hasCommon := false
		for _, lang := range reg.Meta.Languages {
			if lang == checker.LangCommon {
				hasCommon = true
				break
			}
		}
		suite.True(hasCommon, "Security check should apply to LangCommon")
	}
}

// TestSecurityTestSuite runs all tests in suite.
func TestSecurityTestSuite(t *testing.T) {
	suite.Run(t, new(SecurityTestSuite))
}
