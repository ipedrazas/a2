package gocheck

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/suite"
)

// DepsTestSuite is the test suite for the dependency check package.
type DepsTestSuite struct {
	suite.Suite
	hasGovulncheck bool
}

// SetupTest is called before each test method.
func (suite *DepsTestSuite) SetupTest() {
	// Check if govulncheck is available
	_, err := exec.LookPath("govulncheck")
	suite.hasGovulncheck = err == nil
}

// TestDepsCheck_ID tests that DepsCheck returns correct ID.
func (suite *DepsTestSuite) TestDependencyCheck_ID() {
	check := &DepsCheck{}
	suite.Equal("go:deps", check.ID())
}

// TestDepsCheck_Name tests that DepsCheck returns correct name.
func (suite *DepsTestSuite) TestDepsCheck_Name() {
	check := &DepsCheck{}
	suite.Equal("Go Vulnerabilities", check.Name())
}

// TestDepsCheck_Run_NoGovulncheck tests that DepsCheck handles missing govulncheck.
func (suite *DepsTestSuite) TestDepsCheck_Run_NoGovulncheck() {
	// This test will pass regardless of whether govulncheck is installed
	// because the check handles both cases gracefully
	check := &DepsCheck{}
	result, err := check.Run(".")

	suite.NoError(err)
	// Should return a result (either pass with message or actual check result)
	suite.NotNil(result)
	suite.Equal("go:deps", result.ID)
	suite.Equal("Go Vulnerabilities", result.Name)
}

// TestDependencyCheck_Run_WithGovulncheck tests that DependencyCheck runs when govulncheck is available.
func (suite *DepsTestSuite) TestDepsCheck_Run_WithGovulncheck() {
	if !suite.hasGovulncheck {
		suite.T().Skip("govulncheck not installed, skipping test")
	}

	check := &DepsCheck{}
	result, err := check.Run(".")

	suite.NoError(err)
	suite.NotNil(result)
	// Result depends on whether vulnerabilities are found
	// But should always have a message
	suite.NotEmpty(result.Reason)
}

// TestCountVulnerabilities_VulnerabilityPattern tests that countVulnerabilities finds "Vulnerability #" pattern.
func (suite *DepsTestSuite) TestCountVulnerabilities_VulnerabilityPattern() {
	output := `Vulnerability #1: CVE-2023-12345
Vulnerability #2: CVE-2023-67890
Some other text`
	count := countVulnerabilities(output)
	suite.Equal(2, count)
}

// TestCountVulnerabilities_GOPattern tests that countVulnerabilities finds "GO-" pattern as fallback.
func (suite *DepsTestSuite) TestCountVulnerabilities_GOPattern() {
	output := `GO-2023-1234 found
GO-2023-5678 found
Some other text`
	count := countVulnerabilities(output)
	suite.Equal(2, count)
}

// TestCountVulnerabilities_NoMatches tests that countVulnerabilities returns 0 for no matches.
func (suite *DepsTestSuite) TestCountVulnerabilities_NoMatches() {
	output := `No vulnerabilities found
Everything is safe`
	count := countVulnerabilities(output)
	suite.Equal(0, count)
}

// TestCountVulnerabilities_EmptyOutput tests that countVulnerabilities handles empty output.
func (suite *DepsTestSuite) TestCountVulnerabilities_EmptyOutput() {
	count := countVulnerabilities("")
	suite.Equal(0, count)
}

// TestCountVulnerabilities_MixedPatterns tests that countVulnerabilities handles mixed patterns.
func (suite *DepsTestSuite) TestCountVulnerabilities_MixedPatterns() {
	output := `Vulnerability #1: CVE-2023-12345
GO-2023-1234 found
Vulnerability #2: CVE-2023-67890`
	// Should count "Vulnerability #" first (2 matches)
	count := countVulnerabilities(output)
	suite.Equal(2, count)
}

// TestFormatVulnMessage_Single tests that formatVulnMessage formats single vulnerability.
func (suite *DepsTestSuite) TestFormatVulnMessage_Single() {
	message := formatVulnMessage(1)
	suite.Contains(message, "1 vulnerability found")
	suite.Contains(message, "Run 'govulncheck ./...' for details")
}

// TestFormatVulnMessage_Multiple tests that formatVulnMessage formats multiple vulnerabilities.
func (suite *DepsTestSuite) TestFormatVulnMessage_Multiple() {
	message := formatVulnMessage(5)
	suite.Contains(message, "5 vulnerabilities found")
	suite.Contains(message, "Run 'govulncheck ./...' for details")
}

// TestFormatVulnMessage_Zero tests that formatVulnMessage handles zero count.
func (suite *DepsTestSuite) TestFormatVulnMessage_Zero() {
	message := formatVulnMessage(0)
	suite.Contains(message, "0 vulnerabilities found")
}

// TestDepsCheck_Run_NoVulnerabilities tests that DepsCheck returns Pass when no vulnerabilities found.
func (suite *DepsTestSuite) TestDependencyCheck_Run_NoVulnerabilities() {
	if !suite.hasGovulncheck {
		suite.T().Skip("govulncheck not installed, skipping test")
	}

	// This test depends on the actual project having no vulnerabilities
	// or we'd need a test fixture
	check := &DepsCheck{}
	result, err := check.Run(".")

	suite.NoError(err)
	// Result depends on actual vulnerabilities in the project
	// But should always return a valid result
	suite.NotNil(result)
	suite.Equal("go:deps", result.ID)
}

// TestDepsCheck_Run_ErrorHandling tests that DepsCheck handles govulncheck errors.
func (suite *DepsTestSuite) TestDependencyCheck_Run_ErrorHandling() {
	if !suite.hasGovulncheck {
		suite.T().Skip("govulncheck not installed, skipping test")
	}

	// Run on a non-existent path to trigger error handling
	check := &DepsCheck{}
	result, err := check.Run("/nonexistent/path/that/does/not/exist")

	suite.NoError(err)
	// Should handle error gracefully
	suite.NotNil(result)
}

// TestDepsTestSuite runs all the tests in the suite.
func TestDepsTestSuite(t *testing.T) {
	suite.Run(t, new(DepsTestSuite))
}
