package log

import (
	//     "fmt"
	logWriter "andals/gobox/log/writer"
	"testing"
)

func TestSimpleLogger(t *testing.T) {
	path := "/tmp/test.log"

	w, _ := logWriter.NewFileWriter(path)
	writer := logWriter.NewBufferWriter(w, 1024)

	logger, _ := NewSimpleLogger(writer, LEVEL_INFO)
	msg := []byte("test simple logger\n")

	testLogger(logger, msg)
	logger.Free()
}

func TestAsyncLogger(t *testing.T) {
	path := "/tmp/test.log"

	w, _ := logWriter.NewFileWriter(path)
	writer := logWriter.NewBufferWriter(w, 1024)

	l, _ := NewSimpleLogger(writer, LEVEL_INFO)
	logger, _ := NewAsyncLogger("test", l, 10)
	msg := []byte("test async logger\n")

	testLogger(logger, msg)
}

func testLogger(logger ILogger, msg []byte) {
	defer FreeAllAsyncLogger()

	logger.Debug(msg)
	logger.Info(msg)
	logger.Notice(msg)
	logger.Warning(msg)
	logger.Error(msg)
	logger.Critical(msg)
	logger.Alert(msg)
	logger.Emergency(msg)
}
