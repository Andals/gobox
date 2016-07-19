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
	logger, _ := NewSyncSimpleFileLogger(w, LEVEL_INFO)

	msg := []byte("test simple logger")

	testLogger(logger, msg)
	logger.Free()
}

func TestSimpleBufferLogger(t *testing.T) {
	w, _ := logWriter.NewFileWriter("/tmp/test_simple_buffer_logger.log")
	logger, _ := NewSyncSimpleBufferFileLogger(w, 1024, LEVEL_INFO, time.Second*1)

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
	logger, _ := NewAsyncSimpleBufferFileLogger(w, 1024, LEVEL_INFO, 10, time.Second*2)

	msg := []byte("test async simple logger")

	testLogger(logger, msg)
	time.Sleep(time.Second * 3)
	testLogger(logger, msg)
}

func asyncWebLogger(wg *sync.WaitGroup) {
	defer wg.Done()

	w, _ := logWriter.NewFileWriter("/tmp/test_async_web_logger.log")
	logger, _ := NewAsyncSimpleWebBufferFileLogger(w, []byte("async_web"), 1024, LEVEL_INFO, 10, time.Second*2)

	msg := []byte("test async web logger")

	testLogger(logger, msg)
}
