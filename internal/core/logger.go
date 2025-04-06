package core

type LoggerI interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
}

var Logger LoggerI = &noopLogger{}

func SetLogger(l LoggerI) {
	Logger = l
}

type noopLogger struct {}

func (n *noopLogger) Debug(format string, v ...interface{}) {}
func (n *noopLogger) Info(format string, v ...interface{}) {}
func (n *noopLogger) Warn(format string, v ...interface{}) {}
func (n *noopLogger) Error(format string, v ...interface{}) {}