package twizo

import (
	"io/ioutil"
	"log"
	"os"
)

// LogLevel the log levels defined
type LogLevel int

// List of Loglevels
const (
	Debug   LogLevel = 0
	Info    LogLevel = 1
	Warning LogLevel = 2
	Error   LogLevel = 3
)

// APILoggerInterface api logger interface
type APILoggerInterface interface {
	Debug()
	Info()
	Warning()
	Error()
}

// APILogger actual logger
type APILogger struct {
	APILoggerInterface
	Loggers map[LogLevel]*log.Logger
}

// InitLoggers initializes the logger
func InitLoggers() *APILogger {
	l := &APILogger{
		Loggers: make(map[LogLevel]*log.Logger),
	}
	l.Set(Debug, log.New(ioutil.Discard, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile))
	l.Set(Info, log.New(ioutil.Discard, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile))
	l.Set(Warning, log.New(ioutil.Discard, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile))
	l.Set(Error, log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile))
	return l
}

// Set binds a loglevel to a logger
func (l *APILogger) Set(level LogLevel, logger *log.Logger) {
	l.Loggers[level] = logger
}

// Get returns the logger for a loglevel
func (l *APILogger) Get(level LogLevel) *log.Logger {
	return l.Loggers[level]
}

// Debug gets the debug logger
func (l *APILogger) Debug() *log.Logger {
	return l.Get(Debug)
}

// Info gets the info logger
func (l *APILogger) Info() *log.Logger {
	return l.Get(Info)
}

// Warning gets the warning logger
func (l *APILogger) Warning() *log.Logger {
	return l.Get(Warning)
}

// Error gets the error logger
func (l *APILogger) Error() *log.Logger {
	return l.Get(Error)
}
