package writer

import "go.uber.org/zap"

type ErrWriter struct {
	logger *zap.SugaredLogger
}

func (w *ErrWriter) Write(p []byte) (int, error) {
	w.logger.Error(string(p))

	return len(p), nil
}

// New creates writer for usage in standard http package.
func New(log *zap.SugaredLogger) *ErrWriter {
	return &ErrWriter{logger: log}
}
