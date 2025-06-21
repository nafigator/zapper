// Package conf provides Zap Logger different config build approaches.
package conf

import (
	"io"
	"log" //nolint:depguard // Allowed in tests
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

const (
	verboseEnv    = "ZAPPER_CONF_VERBOSE"
	openConfigMsg = "Zap config %s open failure: %s. Fallback to default"
	readConfigMsg = "Zap config read failure: %s. Fallback to default"
	findConfigMsg = "Zap configs not found in paths: %+v. Fallback to default"
	yamlConfigMsg = "Zap config %s parsing error: %s. Fallback to default"
	localPath     = "./zapper.yml"
	systemPath    = "/etc/zap/zapper.yml"
	zeroString    = "0"
	falseString   = "false"
	DefaultConf   = `
level: info
encoding: console
outputPaths:
  - stdout
errorOutputPaths:
  - stderr
encoderConfig:
  messageKey: message
  levelKey:   level
  timeKey:    time
  callerKey:  line
  levelEncoder: capitalColor
  timeEncoder:
    layout: 2006-01-02 15:04:05.000
  durationEncoder: string
  callerEncoder: default
`
	PlainText = `
level: info
encoding: console
outputPaths:
  - stdout
errorOutputPaths:
  - stderr
encoderConfig:
  messageKey: message
  levelKey:   level
  timeKey:    time
  levelEncoder: capital
  timeEncoder:
    layout: 2006-01-02 15:04:05.000
  durationEncoder: string
`
	Debug = `
level: debug
encoding: console
outputPaths:
  - stdout
errorOutputPaths:
  - stderr
encoderConfig:
  messageKey: message
  levelKey:   level
  timeKey:    time
  callerKey:  line
  levelEncoder: capitalColor
  timeEncoder:
    layout: 2006-01-02 15:04:05.000
  durationEncoder: string
  callerEncoder: default
`
)

// New creates Zap Logger config from custom file path.
func New(p string) (*zap.Config, error) {
	var b []byte
	var c zap.Config
	var err error
	var f *os.File

	if f, err = os.Open(p); err != nil {
		return nil, err
	}

	if b, err = read(f); err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

// MustYML creates Zap Logger config from custom yml-config.
// On any error call panic.
func MustYML(p string) *zap.Config {
	var c zap.Config
	var err error

	if err = yaml.Unmarshal([]byte(p), &c); err != nil {
		panic(err)
	}

	return &c
}

// Default initializes zap.Config from default constant yml config.
func Default() *zap.Config {
	var c zap.Config
	var err error

	if err = yaml.Unmarshal([]byte(DefaultConf), &c); err != nil {
		panic(err)
	}

	return &c
}

// Must returns Zap Logger config from custom file path.
// On any error return default config.
// All errors suppressed but this can be changed by ZAPPER_CONF_VERBOSE=1
// env variable.
func Must() *zap.Config {
	var b []byte
	var c zap.Config
	var err error
	var f *os.File

	if silent() {
		log.SetOutput(io.Discard)

		defer log.SetOutput(os.Stderr)
	}

	p := []string{
		localPath,
		systemPath,
	}

	if f, err = findFile(p); err != nil {
		log.Printf(findConfigMsg, p)

		return Default()
	}

	if b, err = read(f); err != nil {
		log.Printf(readConfigMsg, err)

		return Default()
	}

	if err = yaml.Unmarshal(b, &c); err != nil {
		log.Printf(yamlConfigMsg, p, err)

		return Default()
	}

	return &c
}

func read(f *os.File) ([]byte, error) {
	var b []byte
	var err error

	if f == nil {
		return nil, ErrNotOpen
	}

	if b, err = io.ReadAll(f); err != nil {
		return nil, err
	}

	return b, nil
}

func findFile(p []string) (*os.File, error) {
	var f *os.File
	var err error

	for _, pp := range p {
		if f, err = os.Open(pp); err == nil {
			break
		}

		log.Printf(openConfigMsg, pp, err)

		continue
	}

	if f == nil {
		return nil, ErrNotFound
	}

	return f, nil
}

func silent() bool {
	v, ok := os.LookupEnv(verboseEnv)

	return !ok || v == falseString || v == zeroString
}
