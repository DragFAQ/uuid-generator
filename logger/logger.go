package log

// A fallback Logger with minimalistic setup
var log Logger

const (
	// Debug has verbose message
	Debug = "debug"
	// Info is default log level
	Info = "info"
	// Warn is for log messages about possible issues
	Warn = "warn"
	// Error is for log errors
	Error = "error"
	// Fatal is for log fatal messages. The system shuts down after log the message.
	Fatal = "fatal"
)

// Logger is our contract for the logger
type Logger interface {
	Debugf(format string, args ...interface{})

	Infof(format string, args ...interface{})

	Warnf(format string, args ...interface{})

	Errorf(format string, args ...interface{})

	Fatalf(format string, args ...interface{})

	Panicf(format string, args ...interface{})
}

// Configuration stores the config for the logger
// For some loggers there can only be one level across writers, for such the level of Console is picked by default
type Configuration struct {
	EnableConsole     bool
	ConsoleJSONFormat bool
	ConsoleLevel      string
}
