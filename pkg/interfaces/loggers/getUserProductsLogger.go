package loggers

import "gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/usecases"

type getUserProductsLogger struct {
	logger Logger
}

func (l *getUserProductsLogger) LogErrorGettingUserProducts(err error) {
	l.logger.Error("error getting user products data - error: %+v", err)
}

func (l *getUserProductsLogger) LogErrorGettingUserProductsByEmail(email string, err error) {
	l.logger.Error("error getting user products data: email %s - error: %+v", email, err)
}

// MakeGetUserProductsLogger sets up a GetUserProductsLogger instrumented
// via the provided logger
func MakeGetUserProductsLogger(logger Logger) usecases.GetUserProductsLogger {
	return &getUserProductsLogger{
		logger: logger,
	}
}
