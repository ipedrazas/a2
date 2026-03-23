package pythoncheck

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type DepsFreshnessCheckTestSuite struct {
	suite.Suite
}

func (s *DepsFreshnessCheckTestSuite) TestIDAndName() {
	check := &DepsFreshnessCheck{}
	s.Equal("python:deps_freshness", check.ID())
	s.Equal("Python Dependency Freshness", check.Name())
}

func (s *DepsFreshnessCheckTestSuite) TestCountOutdatedPipPackages_None() {
	output := `Package    Version  Latest   Type
--------   -------  ------   ----`
	s.Equal(0, countOutdatedPipPackages(output))
}

func (s *DepsFreshnessCheckTestSuite) TestCountOutdatedPipPackages_Some() {
	output := `Package    Version  Latest   Type
--------   -------  ------   ----
requests   2.28.0   2.31.0   wheel
flask      2.2.0    3.0.0    wheel`
	s.Equal(2, countOutdatedPipPackages(output))
}

func (s *DepsFreshnessCheckTestSuite) TestCountOutdatedPipPackages_Empty() {
	s.Equal(0, countOutdatedPipPackages(""))
}

func TestDepsFreshnessCheckTestSuite(t *testing.T) {
	suite.Run(t, new(DepsFreshnessCheckTestSuite))
}
