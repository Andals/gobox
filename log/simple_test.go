package log

import (
	//     "fmt"
	"testing"
	"time"

	"andals/gobox/log/buffer"
	"andals/gobox/log/writer"
)

func TestSimpleLogger(t *testing.T) {
	fw, _ := writer.NewFileWriter("/tmp/test_simple_logger.log")
	logger, _ := NewSimpleLogger(fw, LEVEL_INFO, new(SimpleFormater))

	msg := []byte("test simple logger")

	testLogger(logger, msg)

	logger.Free()
}

func TestSimpleBufferLogger(t *testing.T) {
	buffer.Init(1024, time.Second*7)

	fw, _ := writer.NewFileWriter("/tmp/test_simple_buffer_logger.log")
	bw := buffer.NewBuffer(fw, 1024)
	logger, _ := NewSimpleLogger(bw, LEVEL_INFO, new(SimpleFormater))

	msg := []byte("test simple buffer logger")

	testLogger(logger, msg)

	logger.Free()
}
