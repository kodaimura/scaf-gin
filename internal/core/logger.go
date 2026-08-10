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
}

// JSONLogger logs JSON lines to stdout.
type JSONLogger struct {
	name  string
	level logLevel
	mu    sync.Mutex
}

type logEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Logger    string `json:"logger"`
	Message   string `json:"message"`
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

func (l *JSONLogger) logf(level logLevel, tag, format string, v ...any) {
	if l.level <= level {
		entry := logEntry{
			Timestamp: time.Now().UTC().Format("2006-01-02T15:04:05.000-07:00"),
			Level:     tag,
			Logger:    l.name,
			Message:   fmt.Sprintf(format, v...),
		}
		body, err := json.Marshal(entry)
		if err != nil {
			body = []byte(fmt.Sprintf(
				`{"timestamp":"%s","level":"ERROR","logger":"%s","message":"failed to marshal log entry: %s"}`,
				time.Now().UTC().Format("2006-01-02T15:04:05.000-07:00"),
				l.name,
				err.Error(),
			))
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
