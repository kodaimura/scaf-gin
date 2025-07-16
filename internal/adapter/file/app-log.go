package file

import (
	"fmt"
	"os"
	"sync"
)

var appLogFile *os.File
var onceAppLog sync.Once

// GetAppLogFile returns the *os.File for app log.
// If the file is not already open, it opens and returns the file.
func GetAppLogFile() *os.File {
	onceAppLog.Do(func() {
		var err error
		appLogFile, err = os.OpenFile("log/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(fmt.Sprintf("failed to open app.log: %v", err))
		}
	})
	return appLogFile
}
