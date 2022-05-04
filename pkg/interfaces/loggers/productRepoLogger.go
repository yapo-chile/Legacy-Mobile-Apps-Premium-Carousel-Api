package loggers

import "gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/interfaces/repository"

type productRepoLogger struct {
	logger Logger
}

func (l *productRepoLogger) LogWarnPartialConfigNotSupported(name, value string) {
	l.logger.Warn("partial config %s: %s not supported", name, value)
}

// MakeProductRepositoryLogger sets up a ProductRepositoryLogger instrumented
// via the provided logger
func MakeProductRepositoryLogger(logger Logger) repository.ProductRepositoryLogger {
	return &productRepoLogger{
		logger: logger,
	}
}
