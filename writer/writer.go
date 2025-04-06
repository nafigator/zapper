// Package writer provides err logger wrapper for using Zap Logger as net/http server logger.
package writer

type ErrLogger interface {
	Error(...interface{})
}

type ErrWriter struct {
	logger ErrLogger
}

// Write writes byes to inner error logger.
func (w *ErrWriter) Write(p []byte) (int, error) {
	w.logger.Error(string(p))

	return len(p), nil
}

// New creates writer for usage in standard http package.
func New(log ErrLogger) *ErrWriter {
	return &ErrWriter{logger: log}
}
