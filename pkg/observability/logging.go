package observability

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	LogFileName = "chargepi.log"
	LogFileDir  = "/var/log/chargepi"
)

// SetupLogger sets up zap logger with the given configuration
func SetupLogger(isDebug bool) *zap.Logger {
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

// Sync flushes any buffered log entries
func Sync(logger *zap.Logger) {
	_ = logger.Sync()
}
