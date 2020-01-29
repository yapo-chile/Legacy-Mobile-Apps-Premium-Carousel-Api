package loggers

import (
	"github.mpi-internal.com/Yapo/goms/pkg/domain"
	"github.mpi-internal.com/Yapo/goms/pkg/usecases"
)

type fibonacciPrometheusDefaultLogger struct {
	logger Logger
}

func (l *fibonacciPrometheusDefaultLogger) LogBadInput(n int) {
	l.logger.Error("GetNth doesn't like N < 1. Input: %d", n)
}

func (l *fibonacciPrometheusDefaultLogger) LogRepositoryError(i int, x domain.Fibonacci, err error) {
	l.logger.Error("Repository refused to save (%d, %d): %s", i, x, err)
}

// MakeFibonacciLogger sets up a FibonacciLogger instrumented via the provided logger
func MakeFibonacciLogger(logger Logger) usecases.FibonacciPrometheusLogger {
	return &fibonacciPrometheusDefaultLogger{
		logger: logger,
	}
}
