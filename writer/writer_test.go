package writer_test

import (
	"testing"

	. "github.com/nafigator/zapper/writer" //nolint: revive // In tests it's ok
	ss "github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type suite struct {
	ss.Suite
}

type testCase struct {
	logger      ErrLogger
	expectedErr error
	msg         string
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
			expected: 8,
			logger:   zap.NewNop().Sugar(),
			msg:      "test msg",
		},
	}
}
