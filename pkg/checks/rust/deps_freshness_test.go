package rustcheck

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type DepsFreshnessCheckTestSuite struct {
	suite.Suite
}

func (s *DepsFreshnessCheckTestSuite) TestIDAndName() {
	check := &DepsFreshnessCheck{}
	s.Equal("rust:deps_freshness", check.ID())
	s.Equal("Rust Dependency Freshness", check.Name())
}

func (s *DepsFreshnessCheckTestSuite) TestCountOutdatedCrates_None() {
	output := `Name             Project  Compat   Latest   Kind
----             -------  ------   ------   ----`
	s.Equal(0, countOutdatedCrates(output))
}

func (s *DepsFreshnessCheckTestSuite) TestCountOutdatedCrates_Some() {
	output := `Name             Project  Compat   Latest   Kind
----             -------  ------   ------   ----
serde            1.0.100  1.0.200  1.0.200  Normal
tokio            1.20.0   1.35.0   1.35.0   Normal`
	s.Equal(2, countOutdatedCrates(output))
}

func (s *DepsFreshnessCheckTestSuite) TestCountOutdatedCrates_Empty() {
	s.Equal(0, countOutdatedCrates(""))
}

func TestDepsFreshnessCheckTestSuite(t *testing.T) {
	suite.Run(t, new(DepsFreshnessCheckTestSuite))
}
