package file

import (
	"fmt"
	"os"
	"sync"
)

var errorLogFile *os.File
var onceErrorLog sync.Once

// GetErrorLogFile returns the *os.File for error log.
// If the file is not already open, it opens and returns the file.
func GetErrorLogFile() *os.File {
	onceErrorLog.Do(func() {
		var err error
		errorLogFile, err = os.OpenFile("log/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(fmt.Sprintf("failed to open error.log: %v", err))
		}
	})
	return errorLogFile
}
