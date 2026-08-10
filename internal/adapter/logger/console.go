package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"scaf-gin/config"
	"scaf-gin/internal/core"
)

// ConsoleLogger logs JSON lines to stdout.
type ConsoleLogger struct {
	level logLevel
	mu    sync.Mutex
}

func NewConsoleLogger() core.LoggerI {
	return &ConsoleLogger{
		level: getLogLevel(),
	}
}

func NewJSONLogger() core.LoggerI {
	return NewConsoleLogger()
}

// Debug logs a debug-level message.
func (l *ConsoleLogger) Debug(format string, v ...any) {
	l.logf(DEBUG, "DEBUG", format, v...)
}

// Info logs an info-level message.
func (l *ConsoleLogger) Info(format string, v ...any) {
	l.logf(INFO, "INFO", format, v...)
}

// Warn logs a warning-level message.
func (l *ConsoleLogger) Warn(format string, v ...any) {
	l.logf(WARN, "WARN", format, v...)
}

// Error logs an error-level message.
func (l *ConsoleLogger) Error(format string, v ...any) {
	l.logf(ERROR, "ERROR", format, v...)
}

func (l *ConsoleLogger) logf(level logLevel, tag, format string, v ...any) {
	if l.level <= level {
		entry := map[string]any{
			"time":    time.Now().UTC().Format(time.RFC3339Nano),
			"level":   strings.ToLower(tag),
			"message": fmt.Sprintf(format, v...),
		}
		body, err := json.Marshal(entry)
		if err != nil {
			body = []byte(fmt.Sprintf(`{"level":"error","message":"failed to marshal log entry: %s"}`, err.Error()))
		}

		l.mu.Lock()
		defer l.mu.Unlock()
		fmt.Fprintln(os.Stdout, string(body))
	}
}

// ===============================
// Common for the logger package.
// ===============================
type logLevel int

const (
	DEBUG logLevel = iota
	INFO
	WARN
	ERROR
)

func getLogLevel() logLevel {
	switch strings.ToLower(config.LogLevel) {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn":
		return WARN
	case "error":
		return ERROR
	default:
		return INFO
	}
}
