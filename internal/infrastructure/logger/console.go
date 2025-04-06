package logger

import (
	"log"

	"goscaf/config"
	"goscaf/internal/core"
)

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
)

func LogLevel() int {
	switch config.LogLevel {
	case "DEBUG", "debug":
		return DEBUG
	case "INFO", "info":
		return INFO
	case "WARN", "warn":
		return WARN
	case "ERROR", "error":
		return ERROR
	default:
		return INFO
	}
}

type ConsoleLogger struct {
	logLevel int
}

func NewConsoleLogger() core.LoggerI {
	log.SetFlags(0)
	return &ConsoleLogger{
		logLevel: LogLevel(),
	}
}

func (l *ConsoleLogger) Debug(format string, v ...interface{}) {
	if l.logLevel <= DEBUG {
		log.Printf("[DEBUG] "+format, v...)
	}
}

func (l *ConsoleLogger) Info(format string, v ...interface{}) {
	if l.logLevel <= INFO {
		log.Printf("[INFO] "+format, v...)
	}
}

func (l *ConsoleLogger) Warn(format string, v ...interface{}) {
	if l.logLevel <= WARN {
		log.Printf("[WARN] "+format, v...)
	}
}

func (l *ConsoleLogger) Error(format string, v ...interface{}) {
	if l.logLevel <= ERROR {
		log.Printf("[ERROR] "+format , v...)
	}
}