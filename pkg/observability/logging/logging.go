package logging

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"

	"github.com/ChargePi/ChargePi-go/internal/pkg/models/settings"
	"github.com/ChargePi/ChargePi-go/internal/pkg/util"
)

const (
	LogFileName = "chargepi.log"
	LogFileDir  = "/var/log/chargepi"
)

// SetupZap sets up zap logger with the given configuration
func SetupZap(loggingConfig settings.Logging, isDebug bool) *zap.Logger {
	// Determine log level
	logLevel := zap.InfoLevel
	if isDebug {
		logLevel = zap.DebugLevel
	}

	// Create encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// Create cores
	var cores []zapcore.Core

	// Console core
	consoleEncoder := zapcore.NewJSONEncoder(encoderConfig)
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), logLevel)
	cores = append(cores, consoleCore)

	// File core
	fileCore := createFileCore(encoderConfig, logLevel)
	if fileCore != nil {
		cores = append(cores, fileCore)
	}

	// Remote logging cores
	for _, logType := range loggingConfig.LogTypes {
		switch LogType(logType.Type) {
		case RemoteLogging:
			if !util.IsNilInterfaceOrPointer(logType.Address) && !util.IsNilInterfaceOrPointer(logType.Format) {
				remoteCore := createRemoteCore(encoderConfig, logLevel, *logType.Address, LogFormat(*logType.Format))
				if remoteCore != nil {
					cores = append(cores, remoteCore)
				}
			}
		case ConsoleLogging:
			// Console logging is already handled above
		}
	}

	// Create logger
	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger
}

func createFileCore(encoderConfig zapcore.EncoderConfig, level zapcore.Level) zapcore.Core {
	// Ensure log directory exists
	if err := os.MkdirAll(LogFileDir, 0755); err != nil {
		return nil
	}

	// Create lumberjack writer for file rotation
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(LogFileDir, LogFileName),
		MaxSize:    200, // megabytes
		MaxBackups: 20,
		MaxAge:     1, // days
		Compress:   false,
	})

	encoder := zapcore.NewJSONEncoder(encoderConfig)
	return zapcore.NewCore(encoder, writer, level)
}

func createRemoteCore(encoderConfig zapcore.EncoderConfig, level zapcore.Level, address string, format LogFormat) zapcore.Core {
	// For now, we'll implement basic syslog support
	// Graylog support would require additional dependencies
	switch format {
	case Syslog:
		// Create a network writer for syslog
		writer, _, err := zap.Open(fmt.Sprintf("tcp://%s", address))
		if err != nil {
			return nil
		}
		encoder := zapcore.NewJSONEncoder(encoderConfig)
		return zapcore.NewCore(encoder, writer, level)
	case Gelf:
		// Graylog GELF format - would need additional implementation
		// For now, return nil to skip this core
		return nil
	default:
		return nil
	}
}

// Sync flushes any buffered log entries
func Sync(logger *zap.Logger) {
	_ = logger.Sync()
}
