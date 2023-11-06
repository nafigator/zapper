package zapper_test

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"bou.ke/monkey"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	. "github.com/nafigator/zapper"
	ss "github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

const (
	OSOpen patch = iota + 1
	IOReadAll
	IOReadAllErr
	yamlUnmarshalErr
	cfgBuildErr
)

type patch int

type fallbackLogger interface {
	Printf(string, ...any)
}

type suite struct {
	tmpDir string
	ss.Suite
}

type testCase struct {
	name          string
	expectedErr   string
	expectedPanic bool
	expected      *zap.SugaredLogger
	path          *string
	fb            fallbackLogger
	patches       []patch
}

func TestRun(t *testing.T) {
	s := suite{
		tmpDir: t.TempDir(),
	}

	ss.Run(t, &s)
}

func (s *suite) TestNew() {
	for _, c := range newDataProvider() {
		s.Run(c.name, func() {
			if len(c.patches) > 0 {
				applyPatches(c.patches, s.tmpDir)
				defer func() {
					monkey.UnpatchAll()
				}()
			}

			actual, err := New(c.path, c.fb)

			if c.expectedErr != "" {
				s.EqualError(err, c.expectedErr, "writer returns unexpected error")
			}

			if c.expected != nil {
				s.True(cmp.Equal(c.expected, actual, opts()...), "invalid response writer result")
			}
		})
	}
}

func (s *suite) TestMust() {
	for _, c := range mustDataProvider() {
		s.Run(c.name, func() {
			if len(c.patches) > 0 {
				applyPatches(c.patches, s.tmpDir)
				defer func() {
					monkey.UnpatchAll()
				}()
			}

			if c.expectedPanic {
				s.Panics(func() {
					Must(c.path, c.fb)
				}, "expected panic didn't catch")

				return
			}

			actual := Must(c.path, c.fb)

			s.True(cmp.Equal(c.expected, actual, opts()...), "invalid response writer result")
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
			name:     "provide invalid path and fallback logger",
			expected: defaultZapper(),
			path:     toPtr("/invalid/path/config.yml"),
			fb:       log.Default(),
		},
		{
			name:     "provide valid path and fallback logger with read ok",
			expected: defaultZapper(),
			path:     toPtr("/valid/path/config.yml"),
			fb:       log.Default(),
			patches:  []patch{OSOpen, IOReadAll},
		},
		{
			name:     "provide valid path and fallback logger with read err",
			expected: defaultZapper(),
			path:     toPtr("/valid/path/config.yml"),
			fb:       log.Default(),
			patches:  []patch{OSOpen, IOReadAllErr},
		},
		{
			name:        "provide valid path and fallback logger with read ok and unmarshall error",
			path:        toPtr("/valid/path/config.yml"),
			fb:          log.Default(),
			patches:     []patch{OSOpen, IOReadAll, yamlUnmarshalErr},
			expectedErr: "unmarshal error stub",
		},
		{
			name:        "provide valid path and fallback logger with read ok and conf build error",
			path:        toPtr("/valid/path/config.yml"),
			fb:          log.Default(),
			patches:     []patch{OSOpen, IOReadAll, cfgBuildErr},
			expectedErr: "config build error stub",
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
			name:          "provide valid path and fallback logger with read ok and unmarshall error",
			path:          toPtr("/valid/path/config.yml"),
			fb:            log.Default(),
			patches:       []patch{OSOpen, IOReadAll, yamlUnmarshalErr},
			expectedPanic: true,
		},
		{
			name:          "provide valid path and fallback logger with read ok and conf build error",
			path:          toPtr("/valid/path/config.yml"),
			fb:            log.Default(),
			patches:       []patch{OSOpen, IOReadAll, cfgBuildErr},
			expectedPanic: true,
		},
	}
}

func defaultZapper() *zap.SugaredLogger {
	var cfg zap.Config
	var logger *zap.Logger
	var err error

	if err = yaml.Unmarshal([]byte(DefaultConf), &cfg); err != nil {
		panic(err)
	}

	if logger, err = cfg.Build(); err != nil {
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

func toPtr(s string) *string {
	return &s
}

func applyPatches(p []patch, d string) {
	for _, v := range p {
		switch v {
		case OSOpen:
			monkey.Patch(os.Open, func(name string) (*os.File, error) {
				f, err := os.CreateTemp(d, "*")
				if err != nil {
					panic(err)
				}

				return f, nil
			})
		case IOReadAll:
			monkey.Patch(io.ReadAll, func(io.Reader) ([]byte, error) {
				return []byte(DefaultConf), nil
			})
		case IOReadAllErr:
			monkey.Patch(io.ReadAll, func(io.Reader) ([]byte, error) {
				return nil, os.ErrPermission
			})
		case yamlUnmarshalErr:
			monkey.Patch(yaml.Unmarshal, func([]byte, interface{}) error {
				return fmt.Errorf("unmarshal error stub")
			})
		case cfgBuildErr:
			monkey.Patch(zap.Config.Build, func(zap.Config, ...zap.Option) (*zap.Logger, error) {
				return nil, fmt.Errorf("config build error stub")
			})
		default:
			continue
		}
	}
}
