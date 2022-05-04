package loggers

import "gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/usecases"

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
