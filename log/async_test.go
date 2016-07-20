package log

import (
	//     "fmt"
	"sync"
	"testing"
	"time"

	logWriter "andals/gobox/log/writer"
)

func TestAsyncLogger(t *testing.T) {
	defer FreeAsyncLogRoutines()

	wg := new(sync.WaitGroup)

	wg.Add(2)

	EnableBufferAutoFlush(time.Second * 2)

	go asyncSimpleLogger(wg)
	go asyncWebLogger(wg)

	wg.Wait()

	time.Sleep(time.Second * 8)
	DisableBufferAutoFlush()
}

func asyncSimpleLogger(wg *sync.WaitGroup) {
	defer wg.Done()

	w, _ := logWriter.NewFileWriter("/tmp/test_async_simple_logger.log")
	logger, _ := NewAsyncSimpleBufferFileLogger(w, 1024, LEVEL_INFO, NewAsyncLogRoutine(10))

	msg := []byte("test async simple logger")

	testLogger(logger, msg)
	time.Sleep(time.Second * 3)
	testLogger(logger, msg)

}

func asyncWebLogger(wg *sync.WaitGroup) {
	defer wg.Done()

	w, _ := logWriter.NewFileWriter("/tmp/test_async_web_logger.log")
	logger, _ := NewAsyncSimpleWebBufferFileLogger(w, []byte("async_web"), 1024, LEVEL_INFO, NewAsyncLogRoutine(10))

	msg := []byte("test async web logger")

	testLogger(logger, msg)
}
