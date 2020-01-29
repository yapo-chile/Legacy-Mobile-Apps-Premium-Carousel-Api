package loggers

import (
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/interfaces/handlers"
)

type getUserDataHandlerPrometheusDefaultLogger struct {
	logger Logger
}

func (l *getUserDataHandlerPrometheusDefaultLogger) LogBadRequest(input interface{}) {
	l.logger.Error("Bad request with input: %+v", input)
}

func (l *getUserDataHandlerPrometheusDefaultLogger) LogErrorGettingInternalData(err error) {
	l.logger.Error("Error getting internal data: %+v ", err)
}

// MakeGetUserDataHandlerLogger sets up a InternalUserDataHandlerLogger instrumented
// via the provided logger
func MakeGetUserDataHandlerLogger(logger Logger) handlers.GetUserDataHandlerPrometheusDefaultLogger {
	return &getUserDataHandlerPrometheusDefaultLogger{
		logger: logger,
	}
}
