package logger

import (
	"log"
	"os"

	"scaf-gin/internal/core"
)

// FileLogger writes log messages to a file with different log levels.
type FileLogger struct {
	level logLevel
}

func NewFileLogger(file *os.File) core.LoggerI {
	log.SetFlags(0)     // Disable default timestamps and flags in the log output
	log.SetOutput(file) // Set the output destination to the provided file
	return &FileLogger{
		level: getLogLevel(),
	}
}

// Debug logs a debug-level message to the file.
func (l *FileLogger) Debug(format string, v ...interface{}) {
	l.logf(DEBUG, "DEBUG", format, v...)
}

// Info logs an info-level message to the file.
func (l *FileLogger) Info(format string, v ...interface{}) {
	l.logf(INFO, "INFO", format, v...)
}

// Warn logs a warning-level message to the file.
func (l *FileLogger) Warn(format string, v ...interface{}) {
	l.logf(WARN, "WARN", format, v...)
}

// Error logs an error-level message to the file.
func (l *FileLogger) Error(format string, v ...interface{}) {
	l.logf(ERROR, "ERROR", format, v...)
}

func (l *FileLogger) logf(level logLevel, tag, format string, v ...interface{}) {
	if l.level <= level {
		log.Printf("["+tag+"] "+format, v...)
	}
}
