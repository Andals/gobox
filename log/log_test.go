package log

import (
	//     "fmt"
	logWriter "andals/gobox/log/writer"
	"sync"
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
	defer FreeAllAsyncLogger()

	wg := new(sync.WaitGroup)

	wg.Add(2)

	go asyncLogger1(wg)
	go asyncLogger2(wg)

	wg.Wait()
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

func asyncLogger1(wg *sync.WaitGroup) {
	defer wg.Done()

	path := "/tmp/test_a1.log"

	w, _ := logWriter.NewFileWriter(path)
	writer := logWriter.NewBufferWriter(w, 1024)

	l, _ := NewSimpleLogger(writer, LEVEL_INFO)
	logger, _ := NewAsyncLogger("test_a1", l, 10)
	msg := []byte("test async1 logger\n")

	testLogger(logger, msg)
}

func asyncLogger2(wg *sync.WaitGroup) {
	defer wg.Done()

	path := "/tmp/test_a2.log"

	w, _ := logWriter.NewFileWriter(path)
	writer := logWriter.NewBufferWriter(w, 1024)

	l, _ := NewSimpleLogger(writer, LEVEL_INFO)
	logger, _ := NewAsyncLogger("test_a2", l, 10)
	msg := []byte("test async2 logger\n")

	testLogger(logger, msg)
}
