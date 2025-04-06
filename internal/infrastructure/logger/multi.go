package logger

import (
	"io"
	"os"
	"log"

	"goscaf/internal/core"
)

type MultiLogger struct {
	logLevel int
}

func NewMultiLogger(file *os.File) core.LoggerI {
	log.SetFlags(0)
	log.SetOutput(io.MultiWriter(os.Stdout, file))
	return &MultiLogger{
		logLevel: LogLevel(),
	}
}

func (l *MultiLogger) Debug(format string, v ...interface{}) {
	if l.logLevel <= DEBUG {
		log.Printf("[DEBUG] "+format, v...)
	}
}

func (l *MultiLogger) Info(format string, v ...interface{}) {
	if l.logLevel <= INFO {
		log.Printf("[INFO] "+format, v...)
	}
}

func (l *MultiLogger) Warn(format string, v ...interface{}) {
	if l.logLevel <= WARN {
		log.Printf("[WARN] "+format, v...)
	}
}

func (l *MultiLogger) Error(format string, v ...interface{}) {
	if l.logLevel <= ERROR {
		log.Printf("[ERROR] "+format , v...)
	}
}