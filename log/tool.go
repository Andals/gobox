package log

import (
	logWriter "andals/gobox/log/writer"

	"time"
)

func NewSyncSimpleFileLogger(path string, level int) (ILogger, error) {
	w, err := logWriter.NewFileWriter(path)
	if err != nil {
		return nil, err
	}

	return NewSimpleLogger(w, level, new(SimpleFormater))
}

func NewSyncSimpleBufferFileLogger(path string, bufsize, level int, timeInterval time.Duration) (ILogger, error) {
	w, err := logWriter.NewFileWriter(path)
	if err != nil {
		return nil, err
	}

	writer, err := logWriter.NewBufferWriterWithTimeFlush(w, bufsize, timeInterval)
	if err != nil {
		return nil, err
	}

	return NewSimpleLogger(writer, level, new(SimpleFormater))
}

func NewAsyncSimpleBufferFileLogger(path string, bufsize, level, queueLen int, timeInterval time.Duration) (ILogger, error) {
	w, err := logWriter.NewFileWriter(path)
	if err != nil {
		return nil, err
	}

	writer, err := logWriter.NewBufferWriterWithTimeFlush(w, bufsize, timeInterval)
	if err != nil {
		return nil, err
	}

	l, err := NewSimpleLogger(writer, level, new(SimpleFormater))
	if err != nil {
		return nil, err
	}

	return NewAsyncLogger(l, queueLen)
}

func NewAsyncSimpleWebBufferFileLogger(path, logId string, bufsize, level, queueLen int, timeInterval time.Duration) (ILogger, error) {
	w, err := logWriter.NewFileWriter(path)
	if err != nil {
		return nil, err
	}

	writer, err := logWriter.NewBufferWriterWithTimeFlush(w, bufsize, timeInterval)
	if err != nil {
		return nil, err
	}

	l, err := NewSimpleLogger(writer, level, NewWebFormater(logId))
	if err != nil {
		return nil, err
	}

	return NewAsyncLogger(l, queueLen)
}
