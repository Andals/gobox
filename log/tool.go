package log

import (
	logWriter "andals/gobox/log/writer"

	"time"
)

func NewSyncSimpleFileLogger(w *logWriter.File, level int) (ILogger, error) {
	return NewSimpleLogger(w, level, new(SimpleFormater))
}

func NewSyncSimpleBufferFileLogger(w *logWriter.File, bufsize, level int, timeInterval time.Duration) (ILogger, error) {
	writer := logWriter.NewBufferWriter(w, bufsize, timeInterval)

	return NewSimpleLogger(writer, level, new(SimpleFormater))
}

func NewAsyncSimpleBufferFileLogger(w *logWriter.File, bufsize, level int, timeInterval time.Duration, ach AsyncLogRoutineCh) (ILogger, error) {
	writer := logWriter.NewBufferWriter(w, bufsize, timeInterval)

	l, err := NewSimpleLogger(writer, level, new(SimpleFormater))
	if err != nil {
		return nil, err
	}

	return NewAsyncLogger(l, ach), nil
}

func NewAsyncSimpleWebBufferFileLogger(w *logWriter.File, logId []byte, bufsize, level int, timeInterval time.Duration, ach AsyncLogRoutineCh) (ILogger, error) {
	writer := logWriter.NewBufferWriter(w, bufsize, timeInterval)

	l, err := NewSimpleLogger(writer, level, NewWebFormater(logId))
	if err != nil {
		return nil, err
	}

	return NewAsyncLogger(l, ach), nil
}
