package core

type LoggerI interface {
	Debug(format string, v ...any)
	Info(format string, v ...any)
	Warn(format string, v ...any)
	Error(format string, v ...any)
}

var Logger LoggerI = &noopLogger{}

func SetLogger(l LoggerI) {
	Logger = l
}

type noopLogger struct{}

func (n *noopLogger) Debug(format string, v ...any) {}
func (n *noopLogger) Info(format string, v ...any)  {}
func (n *noopLogger) Warn(format string, v ...any)  {}
func (n *noopLogger) Error(format string, v ...any) {}
