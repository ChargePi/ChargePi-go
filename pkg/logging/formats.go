package logging

type (
	LogType   string
	LogFormat string
)

const (
	RemoteLogging  = LogType("remote")
	FileLogging    = LogType("file")
	ConsoleLogging = LogType("console")

	Syslog = LogFormat("syslog")
	Gelf   = LogFormat("gelf")
)
