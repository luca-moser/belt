package log

import (
	"fmt"
	"io"
	"os"

	colorable "github.com/mattn/go-colorable"
	log "gopkg.in/inconshreveable/log15.v2"
)

// InitLogger modifies the root log15-logger to log to the console and a logfile with the current PID (optional).
// It returns the log file's io.Writer.
func InitLogger(appName string, usePID bool) (io.Writer, error) {

	// auto create the log file directory
	os.Mkdir("./logs", 0777)

	// open a new logfile
	var logFileName string
	if usePID {
		logFileName = fmt.Sprintf("./logs/%s_%d.log", appName, os.Getpid())
	} else {
		logFileName = fmt.Sprintf("./logs/%s.log", appName)
	}
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	handler := log.MultiHandler(
		log.StreamHandler(logFile, log.LogfmtFormat()),
		log.StreamHandler(colorable.NewColorableStdout(), log.TerminalFormat()),
	)

	log.Root().SetHandler(handler)
	return logFile, nil
}
