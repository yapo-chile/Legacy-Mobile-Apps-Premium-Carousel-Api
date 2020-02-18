package loggers

import "github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"

type getUserProductsLogger struct {
	logger Logger
}

func (l *getUserProductsLogger) LogErrorGettingUserProducts(email string, err error) {
	l.logger.Error("Error getting user products data: email %s, error: %+v", email, err)
}

// MakeGetUserProductsLogger sets up a GetUserProductsLogger instrumented
// via the provided logger
func MakeGetUserProductsLogger(logger Logger) usecases.GetUserProductsLogger {
	return &getUserProductsLogger{
		logger: logger,
	}
}
