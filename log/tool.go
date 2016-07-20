package log

import (
	logWriter "andals/gobox/log/writer"

	"time"
)

func EnableBufferAutoFlush(timeInterval time.Duration) {
	logWriter.EnableBufferAutoFlush(timeInterval)
}

func DisableBufferAutoFlush() {
	logWriter.DisableBufferAutoFlush()
}

func NewSyncSimpleFileLogger(w *logWriter.File, level int) (ILogger, error) {
	return NewSimpleLogger(w, level, new(SimpleFormater))
}

func NewSyncSimpleBufferFileLogger(w *logWriter.File, bufsize, level int) (ILogger, error) {
	writer := logWriter.NewBufferWriter(w, bufsize)

	return NewSimpleLogger(writer, level, new(SimpleFormater))
}

func NewAsyncSimpleBufferFileLogger(w *logWriter.File, bufsize, level int, ach *AsyncLogRoutineCh) (ILogger, error) {
	writer := logWriter.NewBufferWriter(w, bufsize)

	l, err := NewSimpleLogger(writer, level, new(SimpleFormater))
	if err != nil {
		return nil, err
	}

	return NewAsyncLogger(l, ach), nil
}

func NewAsyncSimpleWebBufferFileLogger(w *logWriter.File, logId []byte, bufsize, level int, ach *AsyncLogRoutineCh) (ILogger, error) {
	writer := logWriter.NewBufferWriter(w, bufsize)

	l, err := NewSimpleLogger(writer, level, NewWebFormater(logId))
	if err != nil {
		return nil, err
	}

	return NewAsyncLogger(l, ach), nil
}
