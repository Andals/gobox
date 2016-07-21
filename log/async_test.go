package log

import (
	//     "fmt"
	"sync"
	"testing"
	"time"

	"andals/gobox/log/buffer"
	"andals/gobox/log/writer"
)

func TestAsyncLogger(t *testing.T) {
	defer FreeAsyncLogRoutines()

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
	bw := buffer.NewBuffer(fw, 1024)
	sl, _ := NewSimpleLogger(bw, LEVEL_INFO, new(SimpleFormater))
	logger := NewAsyncLogger(sl, NewAsyncLogRoutine(10))

	msg := []byte("test async simple logger")

	testLogger(logger, msg)
	time.Sleep(time.Second * 3)

	logger.Free()
}

func asyncWebLogger(wg *sync.WaitGroup) {
	defer wg.Done()

	fw, _ := writer.NewFileWriter("/tmp/test_async_web_logger.log")
	bw := buffer.NewBuffer(fw, 1024)
	sl, _ := NewSimpleLogger(bw, LEVEL_INFO, NewWebFormater([]byte("async_web")))
	logger := NewAsyncLogger(sl, NewAsyncLogRoutine(10))

	msg := []byte("test async web logger")

	testLogger(logger, msg)

	logger.Free()
}
