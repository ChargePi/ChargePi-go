package badger

import (
	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.Logger
}

func newLogger() *Logger {
	logger := zap.L().Named("database")
	return &Logger{
		logger: logger,
	}
}

func (l *Logger) Errorf(s string, i ...interface{}) {
	l.logger.Sugar().Errorf(s, i)
}

func (l *Logger) Warningf(s string, i ...interface{}) {
	l.logger.Sugar().Warnf(s, i)
}

func (l *Logger) Infof(s string, i ...interface{}) {
	l.logger.Sugar().Infof(s, i)
}

func (l *Logger) Debugf(s string, i ...interface{}) {
	l.logger.Sugar().Debugf(s, i)
}
