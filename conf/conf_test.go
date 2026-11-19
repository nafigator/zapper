package conf

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	ss "github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

const (
	OSOpen patch = iota + 1
	OSOpenErr
	OSOpenNil
	IOReadAll
	IOReadAllErr
	yamlUnmarshalErr
	yamlUnmarshalErrOnce
	cfgBuildErr

	unmarshalErrMsg = "unmarshal error stub"
)

type patch int

type suite struct {
	ss.Suite

	tmpDir string
}

type testCase struct {
	expected      *zap.Config
	name          string
	expectedErr   string
	yml           string
	path          string
	patches       []patch
	expectedPanic bool
	verbose       bool
}

func TestRun(t *testing.T) {
	t.Parallel()

	s := suite{
		tmpDir: t.TempDir(),
	}

	ss.Run(t, &s)
}

func (s *suite) TestNew() {
	for _, c := range newDataProvider() {
		s.Run(c.name, func() {
			if len(c.patches) > 0 {
				patches := applyPatches(c.patches, s.tmpDir)
				defer resetAll(patches)
			}

			actual, err := New(c.path)

			if c.expectedErr != "" {
				s.Require().EqualError(err, c.expectedErr, "conf returns unexpected error")
			}

			if c.expected != nil {
				s.True(
					cmp.Equal(c.expected, actual, opts()...),
					"conf.New() invalid response",
					cmp.Diff(actual, c.expected, opts()...),
				)
			}
		})
	}
}

func (s *suite) TestYML() {
	for _, c := range YMLDataProvider() {
		s.Run(c.name, func() {
			if len(c.patches) > 0 {
				patches := applyPatches(c.patches, s.tmpDir)
				defer resetAll(patches)
			}

			if c.expectedPanic {
				s.Panics(func() {
					MustYML(c.yml)
				}, "expected panic didn't catch")

				return
			}

			actual := MustYML(c.yml)

			s.True(
				cmp.Equal(c.expected, actual, opts()...),
				"conf.New() invalid response",
				cmp.Diff(actual, c.expected, opts()...),
			)
		})
	}
}

func (s *suite) TestDefault() {
	for _, c := range defaultDataProvider() {
		s.Run(c.name, func() {
			if len(c.patches) > 0 {
				patches := applyPatches(c.patches, s.tmpDir)
				defer resetAll(patches)
			}

			if c.expectedPanic {
				s.Panics(func() {
					Default()
				}, "expected panic didn't catch")

				return
			}

			actual := Default()

			if c.expected != nil {
				s.True(
					cmp.Equal(c.expected, actual, opts()...),
					"conf.Default() invalid response",
					cmp.Diff(actual, c.expected, opts()...),
				)
			}
		})
	}
}

func (s *suite) TestMust() {
	for _, c := range mustDataProvider() {
		s.Run(c.name, func() {
			if len(c.patches) > 0 {
				patches := applyPatches(c.patches, s.tmpDir)
				defer resetAll(patches)
			}

			if c.expectedPanic {
				s.Panics(func() {
					_ = Must()
				}, "expected panic didn't catch")

				return
			}

			actual := Must()

			if c.expected != nil {
				s.True(
					cmp.Equal(c.expected, actual, opts()...),
					"conf.Must() invalid response",
					cmp.Diff(actual, c.expected, opts()...),
				)
			}
		})
	}
}

func newDataProvider() []*testCase {
	return []*testCase{
		{
			name:     "successful",
			expected: Default(),
			patches:  []patch{OSOpen, IOReadAll},
		},
		{
			name:        "open file error",
			expectedErr: os.ErrNotExist.Error(),
			patches:     []patch{OSOpenErr},
		},
		{
			name:        "open file nil",
			expectedErr: ErrNotOpen.Error(),
			patches:     []patch{OSOpenNil},
		},
		{
			name:        "read file error",
			expectedErr: os.ErrPermission.Error(),
			patches:     []patch{OSOpen, IOReadAllErr},
		},
		{
			name:        "yml unmarshal error",
			expectedErr: unmarshalErrMsg,
			patches:     []patch{OSOpen, IOReadAll, yamlUnmarshalErr},
		},
	}
}

func YMLDataProvider() []*testCase {
	return []*testCase{
		{
			name:     "successful",
			expected: Default(),
			yml:      DefaultConf,
		},
		{
			name:          "yml unmarshal panic",
			expectedPanic: true,
			yml:           DefaultConf,
			patches:       []patch{yamlUnmarshalErr},
		},
	}
}

func defaultDataProvider() []*testCase {
	return []*testCase{
		{
			name:     "successful",
			expected: Default(),
		},
		{
			name:          "yml unmarshal panic",
			expectedPanic: true,
			patches:       []patch{yamlUnmarshalErr},
		},
	}
}

func mustDataProvider() []*testCase {
	return []*testCase{
		{
			name:     "successful",
			expected: Default(),
			patches:  []patch{OSOpen, IOReadAll},
		},
		{
			name:     "successful verbose",
			expected: Default(),
			patches:  []patch{OSOpen, IOReadAll},
			verbose:  true,
		},
		{
			name:     "open file error",
			expected: Default(),
			patches:  []patch{OSOpenErr},
		},
		{
			name:     "open file error verbose",
			expected: Default(),
			patches:  []patch{OSOpenErr},
			verbose:  true,
		},
		{
			name:     "open file nil",
			expected: Default(),
			patches:  []patch{OSOpenNil},
		},
		{
			name:        "read file error",
			expectedErr: os.ErrPermission.Error(),
			patches:     []patch{OSOpen, IOReadAllErr},
		},
		{
			name:        "read file error verbose",
			expectedErr: os.ErrPermission.Error(),
			patches:     []patch{OSOpen, IOReadAllErr},
			verbose:     true,
		},
		{
			name:     "yml unmarshal error",
			expected: Default(),
			patches:  []patch{OSOpen, IOReadAll, yamlUnmarshalErrOnce},
		},
		{
			name:     "yml unmarshal error verbose",
			expected: Default(),
			patches:  []patch{OSOpen, IOReadAll, yamlUnmarshalErrOnce},
			verbose:  true,
		},
		{
			name:          "yml unmarshal panic",
			expectedPanic: true,
			patches:       []patch{OSOpen, IOReadAll, yamlUnmarshalErr},
		},
		{
			name:          "yml unmarshal panic verbose",
			expectedPanic: true,
			patches:       []patch{OSOpen, IOReadAll, yamlUnmarshalErr},
			verbose:       true,
		},
	}
}

func applyPatches(p []patch, d string) []*gomonkey.Patches {
	result := make([]*gomonkey.Patches, 0, len(p))

	for _, v := range p {
		var pt *gomonkey.Patches

		switch v {
		case OSOpen:
			pt = gomonkey.ApplyFunc(os.Open, func(_ string) (*os.File, error) {
				f, err := os.CreateTemp(d, "*")
				if err != nil {
					panic(err)
				}

				return f, nil
			})
		case OSOpenErr:
			pt = gomonkey.ApplyFunc(os.Open, func(_ string) (*os.File, error) {
				return nil, os.ErrNotExist
			})
		case OSOpenNil:
			pt = gomonkey.ApplyFunc(os.Open, func(_ string) (*os.File, error) {
				return nil, error(nil)
			})
		case IOReadAll:
			pt = gomonkey.ApplyFunc(io.ReadAll, func(io.Reader) ([]byte, error) {
				return []byte(DefaultConf), nil
			})
		case IOReadAllErr:
			pt = gomonkey.ApplyFunc(io.ReadAll, func(io.Reader) ([]byte, error) {
				return nil, os.ErrPermission
			})
		case yamlUnmarshalErr:
			pt = gomonkey.ApplyFunc(yaml.Unmarshal, func([]byte, any) error {
				return errors.New(unmarshalErrMsg)
			})
		case yamlUnmarshalErrOnce:
			pt = gomonkey.ApplyFunc(yaml.Unmarshal, func([]byte, any) error {
				pt.Reset()

				return errors.New(unmarshalErrMsg)
			})
		case cfgBuildErr:
			pt = gomonkey.ApplyFunc(zap.Config.Build, func(zap.Config, ...zap.Option) (*zap.Logger, error) {
				return nil, errors.New("config build error stub")
			})
		default:
			continue
		}

		if pt != nil {
			result = append(result, pt)
		}
	}

	return result
}

func resetAll(patches []*gomonkey.Patches) {
	for _, p := range patches {
		if p != nil {
			p.Reset()
		}
	}
}

func opts() cmp.Options {
	return []cmp.Option{
		cmpopts.IgnoreUnexported(zap.AtomicLevel{}),
		cmpopts.IgnoreFields(zapcore.EncoderConfig{}, "EncodeLevel", "EncodeTime", "EncodeDuration", "EncodeCaller"),
	}
}
