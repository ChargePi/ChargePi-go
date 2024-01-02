package logging

type (
	LogType   string
	LogFormat string
)

const (
	RemoteLogging  = LogType("remote")
	ConsoleLogging = LogType("console")

	Syslog = LogFormat("syslog")
	Gelf   = LogFormat("gelf")
)
