package log

import (
	logWriter "andals/gobox/log/writer"

	"time"
)

type fileWriters struct {
	lockCh chan int
	ws     map[string]*logWriter.File
}

var fwriters fileWriters

func init() {
	fwriters.lockCh = make(chan int, 1)
	fwriters.lockCh <- 1
	fwriters.ws = make(map[string]*logWriter.File)
}

func NewSyncSimpleFileLogger(path string, level int) (ILogger, error) {
	w, err := getFileWriter(path)
	if err != nil {
		return nil, err
	}

	return NewSimpleLogger(w, level, new(SimpleFormater))
}

func NewSyncSimpleBufferFileLogger(path string, bufsize, level int, flushTimeInterval time.Duration) (ILogger, error) {
	w, err := getFileWriter(path)
	if err != nil {
		return nil, err
	}

	writer := logWriter.NewBufferWriter(w, bufsize, flushTimeInterval)
	if err != nil {
		return nil, err
	}

	return NewSimpleLogger(writer, level, new(SimpleFormater))
}

func NewAsyncSimpleBufferFileLogger(path string, bufsize, level, queueLen int, flushTimeInterval time.Duration) (ILogger, error) {
	w, err := getFileWriter(path)
	if err != nil {
		return nil, err
	}

	writer := logWriter.NewBufferWriter(w, bufsize, flushTimeInterval)
	if err != nil {
		return nil, err
	}

	l, err := NewSimpleLogger(writer, level, new(SimpleFormater))
	if err != nil {
		return nil, err
	}

	return NewAsyncLogger(l, queueLen)
}

func NewAsyncSimpleWebBufferFileLogger(path string, logId []byte, bufsize, level, queueLen int, flushTimeInterval time.Duration) (ILogger, error) {
	w, err := getFileWriter(path)
	if err != nil {
		return nil, err
	}

	writer := logWriter.NewBufferWriter(w, bufsize, flushTimeInterval)
	if err != nil {
		return nil, err
	}

	l, err := NewSimpleLogger(writer, level, NewWebFormater(logId))
	if err != nil {
		return nil, err
	}

	return NewAsyncLogger(l, queueLen)
}

func getFileWriter(path string) (*logWriter.File, error) {
	<-fwriters.lockCh

	w, ok := fwriters.ws[path]
	if !ok {
		var err error

		w, err = logWriter.NewFileWriter(path)
		if err != nil {
			fwriters.lockCh <- 1
			return nil, err
		}

		fwriters.ws[path] = w
	}

	fwriters.lockCh <- 1
	return w, nil
}
