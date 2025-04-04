package logger

import (
	"os"
	"log"
)

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
)

var LogLevel = DEBUG

func init() {
	log.SetFlags(0)
}

func SetOutput(file *os.File) {
	log.SetOutput(file)
}

func SetLevel(level string) {
	switch level {
	case "DEBUG", "debug":
		LogLevel = DEBUG
	case "INFO", "info":
		LogLevel = INFO
	case "WARN", "warn":
		LogLevel = WARN
	case "ERROR", "error":
		LogLevel = ERROR
	default:
		LogLevel = INFO
	}
}

func Debug(format string, v ...interface{}) {
	if LogLevel <= DEBUG {
		log.Printf("[DEBUG] "+format, v...)
	}
}

func Info(format string, v ...interface{}) {
	if LogLevel <= INFO {
		log.Printf("[INFO] "+format, v...)
	}
}

func Warn(format string, v ...interface{}) {
	if LogLevel <= WARN {
		log.Printf("[WARN] "+format, v...)
	}
}

func Error(format string, v ...interface{}) {
	if LogLevel <= ERROR {
		log.Printf("[ERROR] "+format , v...)
	}
}