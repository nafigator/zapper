package writer_test

import (
	"testing"

	. "github.com/nafigator/zapper/writer"
	ss "github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type suite struct {
	ss.Suite
}

type testCase struct {
	logger      ErrLogger
	msg         string
	expectedErr error
	expected    int
}

// TestRun run tests suite.
func TestRun(t *testing.T) {
	s := suite{}

	ss.Run(t, &s)
}

// TestWrite tests writer.Write() func.
func (s *suite) TestWrite() {
	for _, c := range writeDataProvider() {
		actual, err := New(c.logger).Write([]byte(c.msg))

		s.Equal(c.expectedErr, err, "writer returns unexpected error")
		s.Equal(c.expected, actual, "invalid response writer result")
	}
}

func writeDataProvider() []*testCase {
	return []*testCase{
		{
			logger:   zap.NewNop().Sugar(),
			msg:      "test msg",
			expected: 8,
		},
	}
}
