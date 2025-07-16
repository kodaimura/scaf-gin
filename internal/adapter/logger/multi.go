package logger

import (
	"io"
	"log"
	"os"

	"scaf-gin/internal/core"
)

// MultiLogger writes log messages to both standard output (stdout) and a specified file.
type MultiLogger struct {
	level logLevel
}

func NewMultiLogger(file *os.File) core.LoggerI {
	log.SetFlags(0)                                // Disable default timestamps and flags in the log output
	log.SetOutput(io.MultiWriter(os.Stdout, file)) // Write log output to both stdout and the provided file
	return &MultiLogger{
		level: getLogLevel(),
	}
}

// Debug logs a debug-level message to both stdout and the file.
func (l *MultiLogger) Debug(format string, v ...any) {
	l.logf(DEBUG, "DEBUG", format, v...)
}

// Info logs an info-level message to both stdout and the file.
func (l *MultiLogger) Info(format string, v ...any) {
	l.logf(INFO, "INFO", format, v...)
}

// Warn logs a warning-level message to both stdout and the file.
func (l *MultiLogger) Warn(format string, v ...any) {
	l.logf(WARN, "WARN", format, v...)
}

// Error logs an error-level message to both stdout and the file.
func (l *MultiLogger) Error(format string, v ...any) {
	l.logf(ERROR, "ERROR", format, v...)
}

func (l *MultiLogger) logf(level logLevel, tag, format string, v ...any) {
	if l.level <= level {
		log.Printf("["+tag+"] "+format, v...)
	}
}
