package log

import (
	//     "fmt"
	"testing"
	"time"

	logWriter "andals/gobox/log/writer"
)

func TestSimpleLogger(t *testing.T) {
	fw, _ := logWriter.NewFileWriter("/tmp/test_simple_logger.log")
	logger, _ := NewSimpleLogger(fw, LEVEL_INFO, new(SimpleFormater))

	msg := []byte("test simple logger")

	testLogger(logger, msg)

	logger.Free()
	logger = nil

}

func TestSimpleBufferLogger(t *testing.T) {
	logWriter.EnableBufferAutoFlush(time.Second * 1)

	fw, _ := logWriter.NewFileWriter("/tmp/test_simple_buffer_logger.log")
	bw := logWriter.NewBufferWriter(fw, 1024)
	logger, _ := NewSimpleLogger(bw, LEVEL_INFO, new(SimpleFormater))

	msg := []byte("test simple buffer logger")

	testLogger(logger, msg)

	logger.Free()
	logger = nil

	logWriter.DisableBufferAutoFlush()
}
