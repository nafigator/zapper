// Package zapper provides Zap Logger easy configuration by predefined config.
package zapper

import (
	_ "github.com/nafigator/zap-net-sink" // Init UDP and TCP zap logger sinks
	"go.uber.org/zap"
)

// New creates logger instance.
func New(c *zap.Config) (*zap.SugaredLogger, error) {
	logger, err := c.Build()
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

// Must function creates logger instance with panic on failure.
func Must(c *zap.Config) *zap.SugaredLogger {
	logger, err := c.Build()
	if err != nil {
		panic(err)
	}

	return logger.Sugar()
}
