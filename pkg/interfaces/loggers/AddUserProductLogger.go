package loggers

import "github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"

type addUserProductLogger struct {
	logger Logger
}

func (l *addUserProductLogger) LogWarnSettingCache(userID string, err error) {
	l.logger.Warn("Error setting product cache userID: %s error: %+v", userID, err)
}

func (l *addUserProductLogger) LogErrorAddingProduct(userID string, err error) {
	l.logger.Error("Error adding product to userID: %s error: %+v", userID, err)
}

// MakeAddUserProductLogger sets up a AddUserProductLogger instrumented
// via the provided logger
func MakeAddUserProductLogger(logger Logger) usecases.AddUserProductLogger {
	return &addUserProductLogger{
		logger: logger,
	}
}
