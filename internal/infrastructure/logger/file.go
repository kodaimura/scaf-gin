package logger

import (
	"os"
	"log"

	"goscaf/internal/core"
)

type FileLogger struct {
	logLevel int
}

func NewFileLogger(file *os.File) core.LoggerI {
	log.SetFlags(0)
	log.SetOutput(file)
	return &FileLogger{
		logLevel: LogLevel(),
	}
}

func (l *FileLogger) Debug(format string, v ...interface{}) {
	if l.logLevel <= DEBUG {
		log.Printf("[DEBUG] "+format, v...)
	}
}

func (l *FileLogger) Info(format string, v ...interface{}) {
	if l.logLevel <= INFO {
		log.Printf("[INFO] "+format, v...)
	}
}

func (l *FileLogger) Warn(format string, v ...interface{}) {
	if l.logLevel <= WARN {
		log.Printf("[WARN] "+format, v...)
	}
}

func (l *FileLogger) Error(format string, v ...interface{}) {
	if l.logLevel <= ERROR {
		log.Printf("[ERROR] "+format , v...)
	}
}