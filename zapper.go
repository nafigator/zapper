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
func New(fl fallbackLogger, path *string) (*zap.SugaredLogger, error) {
	var cfg zap.Config
	var logger *zap.Logger
	var err error

	if err = yaml.Unmarshal(getYaml(fl, path), &cfg); err != nil {
		return nil, err
	}

	if logger, err = cfg.Build(); err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

// Must function creates logger instance without error returning.
func Must(fl fallbackLogger, path *string) *zap.SugaredLogger {
	var cfg zap.Config
	var logger *zap.Logger
	var err error

	if err = yaml.Unmarshal(getYaml(fl, path), &cfg); err != nil {
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

func defaultPathList() []string {
	return []string{
		localPath,
		systemPath,
	}
}

func getYaml(fl fallbackLogger, path *string) []byte {
	var bytes []byte
	var err error
	var yamlFile *os.File
	var p []string

	if path != nil {
		p = []string{*path}
	} else {
		p = defaultPathList()
	}

	for _, v := range p {
		if yamlFile, err = os.Open(v); err != nil {
			if fl != nil {
				fl.Printf(openConfigWarning, err)
			}

			continue
		}
	}

	if yamlFile == nil {
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
