package cache

// Logger allows for logging errors in the asynchronous calls
type Logger interface {
	// Build returns the data for the supplied key by populating dest
	Log(message string, args ...interface{})
}

// LoggerFunc implements Logger as a function
type LoggerFunc func(message string, args ...interface{})

// Log implements Logger
func (l LoggerFunc) Log(message string, args ...interface{}) {
	l(message, args...)
}

// No op implementation of Logger
var noopLogger = LoggerFunc(func(message string, args ...interface{}) {
	// intentionally do nothing
})
