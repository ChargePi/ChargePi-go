package logging

import (
	"fmt"
	"log/syslog"
	"time"

	graylog "github.com/gemnasium/logrus-graylog-hook/v3"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
)

const logFilePath = "/var/log/chargepi/"

// Setup setup logs
func Setup(logger *log.Logger, loggingConfig settings.Logging, isDebug bool) {
	var (
		// Default logging settings
		logLevel                = log.WarnLevel
		formatter log.Formatter = &log.JSONFormatter{}
	)

	if isDebug {
		logLevel = log.DebugLevel
	}

	logger.SetFormatter(formatter)
	logger.SetLevel(logLevel)

	for _, logType := range loggingConfig.LogTypes {
		switch LogType(logType.Type) {
		case FileLogging:
			fileLogging(logger, isDebug, logFilePath)
		case RemoteLogging:
			if util.IsNilInterfaceOrPointer(logType.Address) && util.IsNilInterfaceOrPointer(logType.Format) {
				remoteLogging(logger, *logType.Address, LogFormat(*logType.Format))
			}
		case ConsoleLogging:
		}
	}
}

// remoteLogging sends logs remotely to Graylog or any Syslog receiver.
func remoteLogging(logger *log.Logger, address string, format LogFormat) {
	var (
		hook log.Hook
		err  error
	)

	switch format {
	case Gelf:
		graylogHook := graylog.NewAsyncGraylogHook(address, map[string]interface{}{})
		hook = graylogHook
	case Syslog:
		hook, err = lSyslog.NewSyslogHook(
			"tcp",
			address,
			syslog.LOG_WARNING,
			"chargePi",
		)
	default:
		return
	}

	if err == nil {
		logger.AddHook(hook)
	}
}

// fileLogging sets up the logging to file.
func fileLogging(logger *log.Logger, isDebug bool, path string) {
	fileName := fmt.Sprintf("%s.%s.chargepi.log", path, ".%Y%m%d%H%M")
	writer, err := rotatelogs.New(
		fileName,
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(86400)*time.Second),
	)
	if err != nil {
		return
	}

	writerMap := make(lfshook.WriterMap)
	writerMap[log.InfoLevel] = writer
	writerMap[log.ErrorLevel] = writer

	if isDebug {
		writerMap[log.DebugLevel] = writer
	}

	hook := lfshook.NewHook(
		writerMap,
		&log.JSONFormatter{},
	)

	logger.AddHook(hook)
}
