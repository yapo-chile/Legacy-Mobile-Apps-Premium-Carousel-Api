package loggers

import "gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/usecases"

type getAdLogger struct {
	logger Logger
}

func (l *getAdLogger) LogWarnGettingCache(listID string, err error) {
	l.logger.Warn("not able to get ad cache listID: %s - %+v", listID, err)
}

func (l *getAdLogger) LogWarnSettingCache(listID string, err error) {
	l.logger.Warn("not able to set ad cache for listID: %s - %+v", listID, err)
}

func (l *getAdLogger) LogErrorGettingAd(listID string, err error) {
	l.logger.Error("Error getting ad data listID: %s error: %+v", listID, err)
}

// MakeGetAdLogger sets up a GetAdLogger instrumented
// via the provided logger
func MakeGetAdLogger(logger Logger) usecases.GetAdLogger {
	return &getAdLogger{
		logger: logger,
	}
}
