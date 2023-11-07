package zapper_test

import (
	"fmt"
	"testing"

	"bou.ke/monkey"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	. "github.com/nafigator/zapper"
	"github.com/nafigator/zapper/conf"
	ss "github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

const (
	cfgBuildErr patch = iota + 1

	cfgBuildErrMsg = "config build error stub"
)

type patch int

type suite struct {
	ss.Suite
}

type testCase struct {
	name          string
	expectedErr   string
	expectedPanic bool
	expected      *zap.SugaredLogger
	patches       []patch
}

func TestRun(t *testing.T) {
	ss.Run(t, &suite{})
}

func (s *suite) TestNew() {
	for _, c := range newDataProvider() {
		s.Run(c.name, func() {
			if len(c.patches) > 0 {
				applyPatches(c.patches)
				defer func() {
					monkey.UnpatchAll()
				}()
			}

			actual, err := New(conf.Default())

			if c.expectedErr != "" {
				s.EqualError(err, c.expectedErr, "New() returns unexpected error")
			}

			if c.expected != nil {
				s.True(cmp.Equal(c.expected, actual, opts()...), "New() invalid response")
			}
		})
	}
}

func (s *suite) TestMust() {
	for _, c := range mustDataProvider() {
		s.Run(c.name, func() {
			if len(c.patches) > 0 {
				applyPatches(c.patches)
				defer func() {
					monkey.UnpatchAll()
				}()
			}

			if c.expectedPanic {
				s.Panics(func() {
					Must(conf.Default())
				}, "expected panic didn't catch")

				return
			}

			actual := Must(conf.Default())

			s.True(cmp.Equal(c.expected, actual, opts()...), "Must() invalid response")
		})
	}
}

func newDataProvider() []*testCase {
	return []*testCase{
		{
			name:     "default config",
			expected: defaultZapper(),
		},
		{
			name:        "config build error",
			expectedErr: cfgBuildErrMsg,
			patches:     []patch{cfgBuildErr},
		},
	}
}

func mustDataProvider() []*testCase {
	return []*testCase{
		{
			name:     "default config",
			expected: defaultZapper(),
		},
		{
			name:          "config build error",
			patches:       []patch{cfgBuildErr},
			expectedPanic: true,
		},
	}
}

func defaultZapper() *zap.SugaredLogger {
	var logger *zap.Logger
	var err error

	if logger, err = conf.Default().Build(); err != nil {
		panic(err)
	}

	return logger.Sugar()
}

func opts() cmp.Options {
	return []cmp.Option{
		cmp.AllowUnexported(zap.SugaredLogger{}),
		cmpopts.IgnoreUnexported(zap.Logger{}),
	}
}

func applyPatches(p []patch) {
	for _, v := range p {
		switch v {
		case cfgBuildErr:
			monkey.Patch(zap.Config.Build, func(zap.Config, ...zap.Option) (*zap.Logger, error) {
				return nil, fmt.Errorf(cfgBuildErrMsg)
			})
		default:
			continue
		}
	}
}
