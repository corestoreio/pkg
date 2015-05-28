package utils

var (
	logger Logger

	nullLog        = &NullLogger{}
	_       Logger = (*NullLogger)(nil)
)

func init() {
	logger = nullLog
}

// SetLogger sets your preferred Logger
func SetLogger(l Logger) {
	if logger != nullLog {
		panic("logger already initialized")
	}
	logger = l
}

func GetLogger() Logger {
	return logger
}

// Default Logger is always the NullLogger
// @see https://github.com/mgutz/logxi

// Logger is the interface for logging.
type Logger interface {
	Trace(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{}) error
	Fatal(msg string, args ...interface{})
	Log(level int, msg string, args []interface{})

	SetLevel(int)
	IsTrace() bool
	IsDebug() bool
	IsInfo() bool
	IsWarn() bool
	// Error, Fatal not needed, those SHOULD always be logged
}

// NullLogger is the default logger for this package.
type NullLogger struct{}

// Trace logs a debug entry.
func (l *NullLogger) Trace(msg string, args ...interface{}) {}

// Debug logs a debug entry.
func (l *NullLogger) Debug(msg string, args ...interface{}) {}

// Info logs an info entry.
func (l *NullLogger) Info(msg string, args ...interface{}) {}

// Warn logs a warn entry.
func (l *NullLogger) Warn(msg string, args ...interface{}) {}

// Error logs an error entry.
func (l *NullLogger) Error(msg string, args ...interface{}) error { return nil }

// Fatal logs a fatal entry then panics.
func (l *NullLogger) Fatal(msg string, args ...interface{}) { panic("exit due to fatal error") }

// Log logs a leveled entry.
func (l *NullLogger) Log(level int, msg string, args []interface{}) {}

// IsTrace determines if this logger logs a trace statement.
func (l *NullLogger) IsTrace() bool { return false }

// IsDebug determines if this logger logs a debug statement.
func (l *NullLogger) IsDebug() bool { return false }

// IsInfo determines if this logger logs an info statement.
func (l *NullLogger) IsInfo() bool { return false }

// IsWarn determines if this logger logs a warning statement.
func (l *NullLogger) IsWarn() bool { return false }

// SetLevel sets the level of this logger.
func (l *NullLogger) SetLevel(level int) {}
