package validation

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ValidationTestSuite struct {
	suite.Suite
}

func (suite *ValidationTestSuite) TestLevenshteinDistance_SameStrings() {
	dist := LevenshteinDistance("hello", "hello")
	suite.Equal(0, dist)
}

func (suite *ValidationTestSuite) TestLevenshteinDistance_EmptyStrings() {
	suite.Equal(0, LevenshteinDistance("", ""))
	suite.Equal(5, LevenshteinDistance("hello", ""))
	suite.Equal(5, LevenshteinDistance("", "hello"))
}

func (suite *ValidationTestSuite) TestLevenshteinDistance_SingleCharDifference() {
	// Substitution
	suite.Equal(1, LevenshteinDistance("cat", "bat"))
	suite.Equal(1, LevenshteinDistance("cat", "cut"))
	suite.Equal(1, LevenshteinDistance("cat", "car"))
}

func (suite *ValidationTestSuite) TestLevenshteinDistance_Insertion() {
	suite.Equal(1, LevenshteinDistance("cat", "cats"))
	suite.Equal(1, LevenshteinDistance("cat", "cart"))
}

func (suite *ValidationTestSuite) TestLevenshteinDistance_Deletion() {
	suite.Equal(1, LevenshteinDistance("cats", "cat"))
	suite.Equal(1, LevenshteinDistance("cart", "cat"))
}

func (suite *ValidationTestSuite) TestLevenshteinDistance_MultipleDifferences() {
	suite.Equal(3, LevenshteinDistance("kitten", "sitting"))
	suite.Equal(1, LevenshteinDistance("common:health", "common:helth")) // missing 'a'
}

func (suite *ValidationTestSuite) TestLevenshteinDistance_CheckIDs() {
	// Realistic check ID typos
	suite.Equal(2, LevenshteinDistance("common:health", "common:helath")) // swap = 2 operations
	suite.Equal(2, LevenshteinDistance("go:build", "go:biuld"))           // swap = 2 operations
	suite.Equal(1, LevenshteinDistance("go:tests", "go:test"))            // missing char
	suite.Equal(2, LevenshteinDistance("common:secrets", "common:secerts"))
}

func (suite *ValidationTestSuite) TestFindSimilar_NoSimilar() {
	validIDs := []string{"go:build", "go:tests", "common:health"}
	similar := FindSimilar("completely:different", validIDs, 3)
	suite.Empty(similar)
}

func (suite *ValidationTestSuite) TestFindSimilar_ExactMatch() {
	validIDs := []string{"go:build", "go:tests", "common:health"}
	// Exact match should not be included (distance = 0)
	similar := FindSimilar("go:build", validIDs, 3)
	suite.Empty(similar)
}

func (suite *ValidationTestSuite) TestFindSimilar_SingleTypo() {
	validIDs := []string{"go:build", "go:tests", "common:health", "common:secrets"}
	similar := FindSimilar("common:helth", validIDs, 3)
	suite.Contains(similar, "common:health")
}

func (suite *ValidationTestSuite) TestFindSimilar_MultipleSuggestions() {
	validIDs := []string{"go:build", "go:tests", "go:race", "go:vet"}
	similar := FindSimilar("go:test", validIDs, 3)
	// go:tests should be first (distance 1), others might be included
	suite.NotEmpty(similar)
	suite.Equal("go:tests", similar[0])
}

func (suite *ValidationTestSuite) TestFindSimilar_MaxThreeSuggestions() {
	validIDs := []string{
		"go:a", "go:b", "go:c", "go:d", "go:e",
	}
	similar := FindSimilar("go:x", validIDs, 3)
	// Should return at most 3 suggestions
	suite.LessOrEqual(len(similar), 3)
}

func (suite *ValidationTestSuite) TestFindSimilar_SortedByDistance() {
	validIDs := []string{"common:health", "common:helath", "common:healt"}
	similar := FindSimilar("common:helth", validIDs, 3)
	// Results should be sorted by distance (closest first)
	if len(similar) >= 2 {
		dist1 := LevenshteinDistance("common:helth", similar[0])
		dist2 := LevenshteinDistance("common:helth", similar[1])
		suite.LessOrEqual(dist1, dist2)
	}
}

func (suite *ValidationTestSuite) TestFindSimilar_MaxDistanceRespected() {
	validIDs := []string{"go:build", "common:health"}
	// With max distance 2, "go:biuld" -> "go:build" (distance 2) should match
	similar := FindSimilar("go:biuld", validIDs, 2)
	suite.Contains(similar, "go:build")
	suite.NotContains(similar, "common:health")

	// With max distance 1, nothing should match (swap is distance 2)
	similarStrict := FindSimilar("go:biuld", validIDs, 1)
	suite.Empty(similarStrict)
}

func (suite *ValidationTestSuite) TestValidationResult_DefaultState() {
	result := ValidationResult{}
	suite.False(result.Valid)
	suite.Empty(result.Errors)
	suite.Empty(result.Warnings)
}

func TestValidationTestSuite(t *testing.T) {
	suite.Run(t, new(ValidationTestSuite))
}
