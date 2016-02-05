package log

import (
	//     "fmt"
	logWriter "andals/gobox/log/writer"
	"testing"
)

func TestSyncLogger(t *testing.T) {
	path := "/tmp/test.log"

	writer, _ := logWriter.NewFileWriter(path)

	logger, _ := NewSimpleLogger(writer, LEVEL_INFO)
	msg := []byte("test simple logger\n")

	testLogger(logger, msg)
}

func testLogger(logger ILogger, msg []byte) {
	logger.Debug(msg)
	logger.Info(msg)
	logger.Notice(msg)
	logger.Warning(msg)
	logger.Error(msg)
	logger.Critical(msg)
	logger.Alert(msg)
	logger.Emergency(msg)
}
