package loggers

import "github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"

type addUserProductLogger struct {
	logger Logger
}

func (l *addUserProductLogger) LogWarnSettingCache(userID int, err error) {
	l.logger.Warn("not able to set product cache userID: %d - %+v", userID, err)
}

func (l *addUserProductLogger) LogErrorAddingProduct(userID int, err error) {
	l.logger.Error("Error adding product to userID: %d error: %+v", userID, err)
}

func (l *addUserProductLogger) LogWarnPushingEvent(productID int, err error) {
	l.logger.Warn("not able to push event to queue productID: %d - %+v", productID, err)
}

// MakeAddUserProductLogger sets up a AddUserProductLogger instrumented
// via the provided logger
func MakeAddUserProductLogger(logger Logger) usecases.AddUserProductLogger {
	return &addUserProductLogger{
		logger: logger,
	}
}
