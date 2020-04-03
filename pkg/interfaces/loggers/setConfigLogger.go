package loggers

import "github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"

type setConfigLogger struct {
	logger Logger
}

func (l *setConfigLogger) LogWarnSettingCache(userID int, err error) {
	l.logger.Warn("unable to set product cache userID: %d - %+v", userID, err)
}

func (l *setConfigLogger) LogErrorSettingConfig(userProductID int, err error) {
	l.logger.Error("error setting config for userProductID: %d - %+v", userProductID, err)
}

// MakeSetConfigLogger sets up a SetConfigLogger instrumented
// via the provided logger
func MakeSetConfigLogger(logger Logger) usecases.SetConfigLogger {
	return &setConfigLogger{
		logger: logger,
	}
}
