package logger

import (
	"github.com/sirupsen/logrus"
	"io"
)

type Opts func(*logrus.Logger)

// New builds a new Logger for the application.
func New(opts ...Opts) *logrus.Logger {
	logger := logrus.New()
	for _, opt := range opts {
		opt(logger)
	}
	return logger
}

func WithLogLevel(level logrus.Level) Opts {
	return func(logger *logrus.Logger) {
		logger.SetLevel(level)
	}
}

func WithNoOutput() Opts {
	return func(logger *logrus.Logger) {
		logger.SetOutput(io.Discard)
	}
}
