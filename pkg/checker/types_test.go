package checker

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// CheckerTypesTestSuite is the test suite for the checker types package.
type CheckerTypesTestSuite struct {
	suite.Suite
}

// SetupTest is called before each test method.
func (suite *CheckerTypesTestSuite) SetupTest() {
	// Setup code if needed
}

// TestStatus_String_Pass tests that Pass status returns "PASS".
func (suite *CheckerTypesTestSuite) TestStatus_String_Pass() {
	status := Pass
	suite.Equal("PASS", status.String())
}

// TestStatus_String_Warn tests that Warn status returns "WARN".
func (suite *CheckerTypesTestSuite) TestStatus_String_Warn() {
	status := Warn
	suite.Equal("WARN", status.String())
}

// TestStatus_String_Fail tests that Fail status returns "FAIL".
func (suite *CheckerTypesTestSuite) TestStatus_String_Fail() {
	status := Fail
	suite.Equal("FAIL", status.String())
}

// TestStatus_String_Info tests that Info status returns "INFO".
func (suite *CheckerTypesTestSuite) TestStatus_String_Info() {
	status := Info
	suite.Equal("INFO", status.String())
}

// TestStatus_String_Unknown tests that unknown status returns "UNKNOWN".
func (suite *CheckerTypesTestSuite) TestStatus_String_Unknown() {
	// Test with a status value that doesn't exist
	status := Status(99)
	suite.Equal("UNKNOWN", status.String())
}

// TestStatus_String_AllValues tests all defined status values.
func (suite *CheckerTypesTestSuite) TestStatus_String_AllValues() {
	suite.Equal("PASS", Pass.String())
	suite.Equal("WARN", Warn.String())
	suite.Equal("FAIL", Fail.String())
	suite.Equal("INFO", Info.String())
}

// TestCheckerTypesTestSuite runs all the tests in the suite.
func TestCheckerTypesTestSuite(t *testing.T) {
	suite.Run(t, new(CheckerTypesTestSuite))
}
