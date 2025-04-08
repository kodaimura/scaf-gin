package logger

import (
	"strings"

	"goscaf/config"
)

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