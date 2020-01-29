package loggers

import "github.mpi-internal.com/Yapo/goms/pkg/usecases"

type gomsRepositoryLogger struct {
	logger Logger
}

func (l *gomsRepositoryLogger) LogURI(s string) {
	l.logger.Debug("%s", s)
}

func (l *gomsRepositoryLogger) LogRequestErr(e error) {
	l.logger.Error("Error obtaining healthcheck information %+v", e)
}

func (l *gomsRepositoryLogger) LogHealthcheckOK(s string) {
	l.logger.Info("%s", s)
}

// MakeGomsRepoLogger sets up a SyncLogger instrumented via the provided logger
func MakeGomsRepoLogger(logger Logger) usecases.HealthcheckPrometheusLogger {
	return &gomsRepositoryLogger{
		logger: logger,
	}
}
