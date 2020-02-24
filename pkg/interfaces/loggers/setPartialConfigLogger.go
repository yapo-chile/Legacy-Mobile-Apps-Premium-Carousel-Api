package loggers

import "github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"

type setPartialConfigLogger struct {
	logger Logger
}

func (l *setPartialConfigLogger) LogWarnSettingCache(userID string, err error) {
	l.logger.Warn("Error setting product cache userID: %s error: %+v", userID, err)
}

func (l *setPartialConfigLogger) LogErrorSettingPartialConfig(userProductID int, err error) {
	l.logger.Error("Error setting partial config userProductID: %s error: %+v", userProductID, err)
}

// MakeSetPartialConfigLogger sets up a SetPartialConfigLogger instrumented
// via the provided logger
func MakeSetPartialConfigLogger(logger Logger) usecases.SetPartialConfigLogger {
	return &setPartialConfigLogger{
		logger: logger,
	}
}
