package log

import (
	//     "fmt"
	"sync"
	"testing"
	"time"

	"andals/gobox/log/writer"
)

func TestAsyncLogger(t *testing.T) {
	InitAsyncLogRoutine(4096)
	defer FreeAsyncLogRoutine()
	writer.InitBufferAutoFlushRoutine(1024, time.Second*7)
	writer.FreeBuffers()

	wg := new(sync.WaitGroup)

	wg.Add(2)

	go asyncSimpleLogger(wg)
	go asyncWebLogger(wg)

	wg.Wait()

	time.Sleep(time.Second * 8)
}

func asyncSimpleLogger(wg *sync.WaitGroup) {
	defer wg.Done()

	fw, _ := writer.NewFileWriter("/tmp/test_async_simple_logger.log")
	bw := writer.NewBuffer(fw, 1024)
	sl, _ := NewSimpleLogger(bw, LEVEL_INFO, new(SimpleFormater))
	logger := NewAsyncLogger(sl)

	msg := []byte("test async simple logger")

	testLogger(logger, msg)
	time.Sleep(time.Second * 3)

	logger.Free()
}

func asyncWebLogger(wg *sync.WaitGroup) {
	defer wg.Done()

	fw, _ := writer.NewFileWriter("/tmp/test_async_web_logger.log")
	bw := writer.NewBuffer(fw, 1024)
	sl, _ := NewSimpleLogger(bw, LEVEL_INFO, NewWebFormater([]byte("async_web"), []byte("127.0.0.1")))
	logger := NewAsyncLogger(sl)

	msg := []byte("test async web logger")

	testLogger(logger, msg)

	logger.Free()
}
