package output

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ProgressTestSuite struct {
	suite.Suite
}

func (s *ProgressTestSuite) TestProgressReporter_UpdateAndDone() {
	pr := NewProgressReporter()

	// Should not panic when called multiple times
	pr.Update(1, 10)
	pr.Update(5, 10)
	pr.Update(10, 10)
	pr.Done()
}

func (s *ProgressTestSuite) TestProgressReporter_DoneWithoutUpdate() {
	pr := NewProgressReporter()

	// Should not panic when Done is called without any Update
	pr.Done()
}

func (s *ProgressTestSuite) TestClearLine() {
	result := clearLine(5)
	s.Equal("     ", result)
	s.Len(result, 5)

	result = clearLine(0)
	s.Equal("", result)
	s.Len(result, 0)
}

func TestProgressTestSuite(t *testing.T) {
	suite.Run(t, new(ProgressTestSuite))
}
