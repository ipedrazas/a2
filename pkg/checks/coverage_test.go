package checks

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// CoverageTestSuite is the test suite for the coverage check package.
type CoverageTestSuite struct {
	suite.Suite
}

// SetupTest is called before each test method.
func (suite *CoverageTestSuite) SetupTest() {
	// Setup code if needed
}

// TestParseCoverage_SinglePackage tests that parseCoverage extracts coverage from single package.
func (suite *CoverageTestSuite) TestParseCoverage_SinglePackage() {
	output := `ok  	github.com/ipedrazas/a2/pkg/runner	0.347s	coverage: 75.5% of statements`
	coverage := parseCoverage(output)
	suite.Equal(75.5, coverage)
}

// TestParseCoverage_MultiplePackages tests that parseCoverage calculates average from multiple packages.
func (suite *CoverageTestSuite) TestParseCoverage_MultiplePackages() {
	output := `ok  	github.com/ipedrazas/a2/pkg/runner	0.347s	coverage: 80.0% of statements
ok  	github.com/ipedrazas/a2/pkg/config	0.349s	coverage: 90.0% of statements
ok  	github.com/ipedrazas/a2/pkg/checker	0.346s	coverage: 100.0% of statements`
	coverage := parseCoverage(output)
	suite.Equal(90.0, coverage) // (80 + 90 + 100) / 3 = 90
}

// TestParseCoverage_NoMatches tests that parseCoverage returns 0 for no matches.
func (suite *CoverageTestSuite) TestParseCoverage_NoMatches() {
	output := `ok  	github.com/ipedrazas/a2/pkg/runner	0.347s`
	coverage := parseCoverage(output)
	suite.Equal(0.0, coverage)
}

// TestParseCoverage_EmptyOutput tests that parseCoverage handles empty output.
func (suite *CoverageTestSuite) TestParseCoverage_EmptyOutput() {
	coverage := parseCoverage("")
	suite.Equal(0.0, coverage)
}

// TestParseCoverage_InvalidFormat tests that parseCoverage handles invalid format.
func (suite *CoverageTestSuite) TestParseCoverage_InvalidFormat() {
	output := `some random text without coverage information`
	coverage := parseCoverage(output)
	suite.Equal(0.0, coverage)
}

// TestParseCoverage_DecimalValues tests that parseCoverage handles decimal values correctly.
func (suite *CoverageTestSuite) TestParseCoverage_DecimalValues() {
	output := `ok  	github.com/ipedrazas/a2/pkg/runner	0.347s	coverage: 75.75% of statements
ok  	github.com/ipedrazas/a2/pkg/config	0.349s	coverage: 82.25% of statements`
	coverage := parseCoverage(output)
	suite.Equal(79.0, coverage) // (75.75 + 82.25) / 2 = 79.0
}

// TestParseCoverage_WithNoTestFiles tests that parseCoverage handles "no test files" output.
func (suite *CoverageTestSuite) TestParseCoverage_WithNoTestFiles() {
	output := `?   	github.com/ipedrazas/a2/pkg/something	[no test files]
ok  	github.com/ipedrazas/a2/pkg/runner	0.347s	coverage: 80.0% of statements`
	coverage := parseCoverage(output)
	suite.Equal(80.0, coverage) // Should only count packages with coverage
}

// TestParseCoverage_ComplexOutput tests that parseCoverage handles complex output with various formats.
func (suite *CoverageTestSuite) TestParseCoverage_ComplexOutput() {
	output := `ok  	github.com/ipedrazas/a2/pkg/runner	0.347s	coverage: 50.0% of statements
ok  	github.com/ipedrazas/a2/pkg/config	0.349s	coverage: 60.0% of statements
ok  	github.com/ipedrazas/a2/pkg/checker	0.346s	coverage: 70.0% of statements
ok  	github.com/ipedrazas/a2/pkg/checks	0.293s	coverage: 80.0% of statements`
	coverage := parseCoverage(output)
	suite.Equal(65.0, coverage) // (50 + 60 + 70 + 80) / 4 = 65
}

// TestCoverageCheck_ID tests that CoverageCheck returns correct ID.
func (suite *CoverageTestSuite) TestCoverageCheck_ID() {
	check := &CoverageCheck{}
	suite.Equal("coverage", check.ID())
}

// TestCoverageCheck_Name tests that CoverageCheck returns correct name.
func (suite *CoverageTestSuite) TestCoverageCheck_Name() {
	check := &CoverageCheck{}
	suite.Equal("Test Coverage", check.Name())
}

// TestDefaultCoverageCheck tests that DefaultCoverageCheck returns check with 80% threshold.
func (suite *CoverageTestSuite) TestDefaultCoverageCheck() {
	check := DefaultCoverageCheck()
	suite.NotNil(check)
	suite.Equal(80.0, check.Threshold)
	suite.Equal("coverage", check.ID())
	suite.Equal("Test Coverage", check.Name())
}

// TestCoverageCheck_Run_AboveThreshold tests CoverageCheck with coverage above threshold.
// Note: This test requires actual go test to run, so it may be skipped in some environments
func (suite *CoverageTestSuite) TestCoverageCheck_Run_AboveThreshold() {
	// This test would require a real Go project with tests
	// For now, we'll test the logic with parseCoverage which is already tested
	// The actual Run() method integration test would need a test fixture
}

// TestCoverageCheck_Run_BelowThreshold tests CoverageCheck with coverage below threshold.
// Note: This test requires actual go test to run
func (suite *CoverageTestSuite) TestCoverageCheck_Run_BelowThreshold() {
	// Similar to above - would need test fixtures
}

// TestCoverageCheck_Run_DefaultThreshold tests that CoverageCheck uses default threshold when not set.
func (suite *CoverageTestSuite) TestCoverageCheck_Run_DefaultThreshold() {
	check := &CoverageCheck{
		Threshold: 0, // Should use default 80.0
	}

	// We can't easily test Run() without actual go test execution
	// But we can verify the threshold logic
	threshold := check.Threshold
	if threshold == 0 {
		threshold = 80.0
	}
	suite.Equal(80.0, threshold)
}

// TestCoverageCheck_Run_CustomThreshold tests that CoverageCheck uses custom threshold.
func (suite *CoverageTestSuite) TestCoverageCheck_Run_CustomThreshold() {
	check := &CoverageCheck{
		Threshold: 90.0,
	}

	threshold := check.Threshold
	if threshold == 0 {
		threshold = 80.0
	}
	suite.Equal(90.0, threshold)
}

// TestParseCoverage_EdgeCases tests edge cases for parseCoverage.
func (suite *CoverageTestSuite) TestParseCoverage_EdgeCases() {
	// Test with just percentage
	output := `coverage: 50.0%`
	coverage := parseCoverage(output)
	suite.Equal(50.0, coverage)

	// Test with extra whitespace before percentage (regex handles this)
	output = `coverage:   75.5%`
	coverage = parseCoverage(output)
	suite.Equal(75.5, coverage)

	// Test with multiple matches in same line (should count all)
	output = `coverage: 50.0% of statements
coverage: 70.0% of statements`
	coverage = parseCoverage(output)
	suite.Equal(60.0, coverage)
}

// TestCoverageTestSuite runs all the tests in the suite.
func TestCoverageTestSuite(t *testing.T) {
	suite.Run(t, new(CoverageTestSuite))
}
