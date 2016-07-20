package log

import (
	//     "fmt"
	"testing"
	"time"

	logWriter "andals/gobox/log/writer"
)

func TestSimpleLogger(t *testing.T) {
	w, _ := logWriter.NewFileWriter("/tmp/test_simple_logger.log")
	logger, _ := NewSyncSimpleFileLogger(w, LEVEL_INFO)

	msg := []byte("test simple logger")

	testLogger(logger, msg)
	logger.Free()
}

func TestSimpleBufferLogger(t *testing.T) {
	EnableBufferAutoFlush(time.Second * 1)

	w, _ := logWriter.NewFileWriter("/tmp/test_simple_buffer_logger.log")
	logger, _ := NewSyncSimpleBufferFileLogger(w, 1024, LEVEL_INFO)

	msg := []byte("test simple buffer logger")

	testLogger(logger, msg)
	logger.Free()

	DisableBufferAutoFlush()
}
