package file

import (
	"fmt"
	"os"
	"sync"
)

var accessLogFile *os.File
var onceAccessLog sync.Once

// GetAccessLogFile returns the *os.File for access log.
// If the file is not already open, it opens and returns the file.
func GetAccessLogFile() *os.File {
	onceAccessLog.Do(func() {
		var err error
		accessLogFile, err = os.OpenFile("log/access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(fmt.Sprintf("failed to open access.log: %v", err))
		}
	})
	return accessLogFile
}
