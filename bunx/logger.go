package bunx

import "go.uber.org/zap"

type Logger interface {
	Info(args ...interface{})
}

func DefaultLogger() Logger {
	l, _ := zap.NewProduction()
	return l.Sugar()
}

type NoopLogger struct{}

func NewNoopLogger() Logger {
	return new(NoopLogger)
}

func (n *NoopLogger) Info(_ ...interface{}) {}
