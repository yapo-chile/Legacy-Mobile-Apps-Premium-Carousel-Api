package loggers

import "github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"

type expireProductsLogger struct {
	logger Logger
}

func (l *expireProductsLogger) LogExpireProductsError(err error) {
	l.logger.Error("error expiring products: %+v", err)
}

// MakeExpireProductsLogger sets up a ExpireProductsLogger instrumented
// via the provided logger
func MakeExpireProductsLogger(logger Logger) usecases.ExpireProductsLogger {
	return &expireProductsLogger{
		logger: logger,
	}
}
