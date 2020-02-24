package loggers

import "github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"

type setConfigLogger struct {
	logger Logger
}

func (l *setConfigLogger) LogWarnSettingCache(userID string, err error) {
	l.logger.Warn("Error setting product cache userID: %s error: %+v", userID, err)
}

func (l *setConfigLogger) LogErrorSettingConfig(userProductID int, err error) {
	l.logger.Error("Error setting config to userProductID: %d error: %+v", userProductID, err)
}

// MakeSetConfigLogger sets up a SetConfigLogger instrumented
// via the provided logger
func MakeSetConfigLogger(logger Logger) usecases.SetConfigLogger {
	return &setConfigLogger{
		logger: logger,
	}
}
