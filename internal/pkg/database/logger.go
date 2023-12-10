package database

import log "github.com/sirupsen/logrus"

type Logger struct {
	logger log.FieldLogger
}

func newLogger() *Logger {
	logger := log.StandardLogger().WithField("component", "database")
	return &Logger{
		logger: logger,
	}
}

func (l *Logger) Errorf(s string, i ...interface{}) {
	l.logger.Errorf(s, i)
}

func (l *Logger) Warningf(s string, i ...interface{}) {
	l.logger.Warningf(s, i)
}

func (l *Logger) Infof(s string, i ...interface{}) {
	l.logger.Infof(s, i)
}

func (l *Logger) Debugf(s string, i ...interface{}) {
	l.logger.Debugf(s, i)
}
