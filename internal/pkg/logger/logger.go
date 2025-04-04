package logger

import (
	"log"
)

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
)

var logLevel int

func init() {
	log.SetFlags(0)

	level := "INFO" //:= config.GetConfig().LogLevel
	switch level {
	case "DEBUG", "debug":
		logLevel = DEBUG
	case "INFO", "info":
		logLevel = INFO
	case "WARN", "warn":
		logLevel = WARN
	case "ERROR", "error":
		return LOG_LEVEL_ERROR
	default:
		return LOG_LEVEL_INFO
	}
}

func Debug(format string, v ...interface{}) {
	if logLevel <= DEBUG {
		log.Printf("[DEBUG] "+format, v...)
	}
}

func Info(format string, v ...interface{}) {
	if logLevel <= INFO {
		log.Printf("[INFO] "+format, v...)
	}
}

func Warn(format string, v ...interface{}) {
	if logLevel <= WARN {
		log.Printf("[WARN] "+format, v...)
	}
}

func Error(format string, v ...interface{}) {
	if logLevel <= ERROR {
		log.Printf("[ERROR] "+format , v...)
	}
}