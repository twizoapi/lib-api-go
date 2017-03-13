package twizo

import (
	"io/ioutil"
	"log"
	"os"
)

type LogLevel int

const (
	Debug   LogLevel = 0
	Info    LogLevel = 1
	Warning LogLevel = 2
	Error   LogLevel = 3
)

type APILoggerInterface interface {
	Debug()
	Info()
	Warning()
	Error()
}

type APILogger struct {
	APILoggerInterface
	Loggers map[LogLevel]*log.Logger
}

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

func (l *APILogger) Set(level LogLevel, logger *log.Logger) {
	l.Loggers[level] = logger
}

func (l *APILogger) Get(level LogLevel) *log.Logger {
	return l.Loggers[level]
}

func (l *APILogger) Info() *log.Logger {
	return l.Get(Info)
}

func (l *APILogger) Debug() *log.Logger {
	return l.Get(Debug)
}

func (l *APILogger) Warning() *log.Logger {
	return l.Get(Warning)
}

func (l *APILogger) Error() *log.Logger {
	return l.Get(Error)
}
