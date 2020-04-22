package loggers

import "github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"

type getReportLogger struct {
	logger Logger
}

func (l *getReportLogger) LogErrorGettingReport(err error) {
	l.logger.Error("error getting report data - error: %+v", err)
}

// MakeGetReportLogger sets up a GetReportLogger instrumented
// via the provided logger
func MakeGetReportLogger(logger Logger) usecases.GetReportLogger {
	return &getReportLogger{
		logger: logger,
	}
}
