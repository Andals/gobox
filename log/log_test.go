package log

import (
	//     "fmt"
	"sync"
	"testing"
	"time"

	logWriter "andals/gobox/log/writer"
)

func TestSimpleLogger(t *testing.T) {
	w, _ := logWriter.NewFileWriter("/tmp/test_simple_logger.log")
	logger, _ := NewSimpleLogger(w, LEVEL_INFO, new(SimpleFormater))

	msg := []byte("test simple logger")

	testLogger(logger, msg)
	logger.Free()
}

func TestSimpleBufferLogger(t *testing.T) {
	//     logger, _ := NewSyncSimpleBufferFileLogger(, 1024, LEVEL_INFO, time.Second*1)

	w, _ := logWriter.NewFileWriter("/tmp/test_simple_buffer_logger.log")
	writer := logWriter.NewBufferWriter(w, 1024, time.Second*1)
	logger, _ := NewSimpleLogger(writer, LEVEL_INFO, new(SimpleFormater))

	msg := []byte("test simple buffer logger")

	testLogger(logger, msg)
	logger.Free()
}

func TestAsyncLogger(t *testing.T) {
	defer FreeAllAsyncLogger()

	wg := new(sync.WaitGroup)

	wg.Add(2)

	go asyncSimpleLogger(wg)
	go asyncWebLogger(wg)

	wg.Wait()

	time.Sleep(time.Second * 8)
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

func asyncSimpleLogger(wg *sync.WaitGroup) {
	defer wg.Done()

	w, _ := logWriter.NewFileWriter("/tmp/test_async_simple_logger.log")
	writer := logWriter.NewBufferWriter(w, 1024, time.Second*2)
	l, _ := NewSimpleLogger(writer, LEVEL_INFO, new(SimpleFormater))
	logger, _ := NewAsyncLogger(l, 10)

	msg := []byte("test async simple logger")

	testLogger(logger, msg)
	time.Sleep(time.Second * 3)
	testLogger(logger, msg)
}

func asyncWebLogger(wg *sync.WaitGroup) {
	defer wg.Done()

	w, _ := logWriter.NewFileWriter("/tmp/test_async_web_logger.log")
	writer := logWriter.NewBufferWriter(w, 1024, time.Second*2)
	l, _ := NewSimpleLogger(writer, LEVEL_INFO, NewWebFormater([]byte("async_web")))
	logger, _ := NewAsyncLogger(l, 10)

	msg := []byte("test async2 logger")

	testLogger(logger, msg)
}
