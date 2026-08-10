package core

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type Logger interface {
	Debug(format string, v ...any)
	Info(format string, v ...any)
	Warn(format string, v ...any)
	Error(format string, v ...any)
	DebugFields(message string, fields map[string]any)
	InfoFields(message string, fields map[string]any)
	WarnFields(message string, fields map[string]any)
	ErrorFields(message string, fields map[string]any)
}

// JSONLogger logs JSON lines to stdout.
type JSONLogger struct {
	name  string
	level logLevel
	mu    sync.Mutex
}

func NewConsoleLogger(logLevel string) Logger {
	return NewJSONLogger(logLevel)
}

func NewJSONLogger(logLevel string, names ...string) Logger {
	name := "app"
	if len(names) > 0 && names[0] != "" {
		name = names[0]
	}
	return &JSONLogger{
		name:  name,
		level: getLogLevel(logLevel),
	}
}

func (l *JSONLogger) Debug(format string, v ...any) {
	l.logf(DEBUG, "DEBUG", format, v...)
}

func (l *JSONLogger) Info(format string, v ...any) {
	l.logf(INFO, "INFO", format, v...)
}

func (l *JSONLogger) Warn(format string, v ...any) {
	l.logf(WARN, "WARN", format, v...)
}

func (l *JSONLogger) Error(format string, v ...any) {
	l.logf(ERROR, "ERROR", format, v...)
}

func (l *JSONLogger) DebugFields(message string, fields map[string]any) {
	l.log(DEBUG, "DEBUG", message, fields)
}

func (l *JSONLogger) InfoFields(message string, fields map[string]any) {
	l.log(INFO, "INFO", message, fields)
}

func (l *JSONLogger) WarnFields(message string, fields map[string]any) {
	l.log(WARN, "WARN", message, fields)
}

func (l *JSONLogger) ErrorFields(message string, fields map[string]any) {
	l.log(ERROR, "ERROR", message, fields)
}

func (l *JSONLogger) logf(level logLevel, tag, format string, v ...any) {
	l.log(level, tag, fmt.Sprintf(format, v...), nil)
}

func (l *JSONLogger) log(level logLevel, tag, message string, fields map[string]any) {
	if l.level <= level {
		entry := map[string]any{
			"timestamp": time.Now().UTC().Format("2006-01-02T15:04:05.000-07:00"),
			"level":     tag,
			"logger":    l.name,
			"message":   message,
		}
		for key, value := range fields {
			entry[key] = value
		}
		body, err := json.Marshal(entry)
		if err != nil {
			fallback := map[string]any{
				"timestamp": time.Now().UTC().Format("2006-01-02T15:04:05.000-07:00"),
				"level":     "ERROR",
				"logger":    l.name,
				"message":   fmt.Sprintf("failed to marshal log entry: %s", err.Error()),
			}
			body, _ = json.Marshal(fallback)
		}

		l.mu.Lock()
		defer l.mu.Unlock()
		fmt.Fprintln(os.Stdout, string(body))
	}
}

type logLevel int

const (
	DEBUG logLevel = iota
	INFO
	WARN
	ERROR
)

func getLogLevel(value string) logLevel {
	switch strings.ToLower(value) {
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
