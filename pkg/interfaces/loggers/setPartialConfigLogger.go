package loggers

import "github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"

type setPartialConfigLogger struct {
	logger Logger
}

func (l *setPartialConfigLogger) LogWarnSettingCache(userID int, err error) {
	l.logger.Warn("unable to set product cache userID: %d - %+v", userID, err)
}

func (l *setPartialConfigLogger) LogErrorSettingPartialConfig(userProductID int, err error) {
	l.logger.Error("error setting partial config for userProductID: %s - %+v", userProductID, err)
}

// MakeSetPartialConfigLogger sets up a SetPartialConfigLogger instrumented
// via the provided logger
func MakeSetPartialConfigLogger(logger Logger) usecases.SetPartialConfigLogger {
	return &setPartialConfigLogger{
		logger: logger,
	}
}
