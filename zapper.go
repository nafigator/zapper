package zapper

import (
	"io"
	"os"

	_ "github.com/kontera-technologies/zap-net-sink"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

const (
	openConfigWarning = "Zap config open failure: %s. Fallback to default"
	readConfigWarning = "Zap config read failure: %s. Fallback to default"
	localPath         = "./config.yml"
	systemPath        = "/etc/zap/config.yml"
	defaultConf       = `
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
)

type fallbackLogger interface {
	Printf(string, ...any)
}

// New creates logger instance.
func New(path *string, fl fallbackLogger) (*zap.SugaredLogger, error) {
	var cfg zap.Config
	var logger *zap.Logger
	var err error

	if err = yaml.Unmarshal(getYaml(path, fl), &cfg); err != nil {
		return nil, err
	}

	if logger, err = cfg.Build(); err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

// Must function creates logger instance without error returning.
func Must(path *string, fl fallbackLogger) *zap.SugaredLogger {
	var cfg zap.Config
	var logger *zap.Logger
	var err error

	if err = yaml.Unmarshal(getYaml(path, fl), &cfg); err != nil {
		if fl != nil {
			fl.Printf("Zapper unmarshal failure: %s", err)
		}

		l, _ := zap.NewProductionConfig().Build()

		return l.Sugar()
	}

	if logger, err = cfg.Build(); err != nil {
		l, _ := zap.NewProductionConfig().Build()

		return l.Sugar()
	}

	return logger.Sugar()
}

func getYaml(path *string, fl fallbackLogger) []byte {
	var bytes []byte
	var err error
	var yamlFile *os.File

	p := getPathList(path)

	if yamlFile = openFile(p, fl); yamlFile == nil {
		return []byte(defaultConf)
	}

	if bytes, err = io.ReadAll(yamlFile); err != nil {
		if fl != nil {
			fl.Printf(readConfigWarning, err)
		}

		return []byte(defaultConf)
	}

	return bytes
}

func openFile(p []string, fl fallbackLogger) *os.File {
	var yamlFile *os.File
	var err error

	for _, v := range p {
		if yamlFile, err = os.Open(v); err != nil && fl != nil {
			fl.Printf(openConfigWarning, err)
		}
	}

	return yamlFile
}

func getPathList(path *string) []string {
	if path != nil {
		return []string{*path}
	}

	return []string{
		localPath,
		systemPath,
	}
}
