package logger

import (
	"log"

	"goscaf/internal/core"
)

// ConsoleLogger logs messages to the standard output.
type ConsoleLogger struct {
	level logLevel
}

func NewConsoleLogger() core.LoggerI {
	log.SetFlags(0) // Disable default timestamps and flags in the log output
	return &ConsoleLogger{
		level: getLogLevel(),
	}
}

// Debug logs a debug-level message.
func (l *ConsoleLogger) Debug(format string, v ...interface{}) {
	l.logf(DEBUG, "DEBUG", format, v...)
}

// Info logs an info-level message.
func (l *ConsoleLogger) Info(format string, v ...interface{}) {
	l.logf(INFO, "INFO", format, v...)
}

// Warn logs a warning-level message.
func (l *ConsoleLogger) Warn(format string, v ...interface{}) {
	l.logf(WARN, "WARN", format, v...)
}

// Error logs an error-level message.
func (l *ConsoleLogger) Error(format string, v ...interface{}) {
	l.logf(ERROR, "ERROR", format, v...)
}

func (l *ConsoleLogger) logf(level logLevel, tag, format string, v ...interface{}) {
	if l.level <= level {
		log.Printf("["+tag+"] "+format, v...)
	}
}
