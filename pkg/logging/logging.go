package logging

import (
	"log/syslog"

	graylog "github.com/gemnasium/logrus-graylog-hook/v3"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	"github.com/lorenzodonini/ocpp-go/ws"
	log "github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
)

// Setup setup logs
func Setup(logger *log.Logger, loggingConfig settings.Logging, isDebug bool) {
	// Default logging settings
	logLevel := log.InfoLevel
	formatter := &log.JSONFormatter{}
	logger.SetFormatter(formatter)

	if isDebug {
		logLevel = log.DebugLevel
		ocppj.SetLogger(logger)
		ws.SetLogger(logger)
	}

	logger.SetLevel(logLevel)

	for _, logType := range loggingConfig.LogTypes {
		switch LogType(logType.Type) {
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
