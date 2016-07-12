package log

import (
	"andals/gobox/log/formater"
	logWriter "andals/gobox/log/writer"
)

func NewSyncSimpleBufferFileLogger(path string, bufsize, level int) (ILogger, error) {
	w, err := logWriter.NewFileWriter(path)
	if err != nil {
		return nil, err
	}

	writer := logWriter.NewBufferWriter(w, bufsize)

	return NewSimpleLogger(writer, level, new(formater.Simple))
}

func NewAsyncSimpleBufferFileLogger(path string, bufsize, level, queueLen int) (ILogger, error) {
	w, err := logWriter.NewFileWriter(path)
	if err != nil {
		return nil, err
	}

	writer := logWriter.NewBufferWriter(w, bufsize)

	l, err := NewSimpleLogger(writer, level, new(formater.Simple))
	if err != nil {
		return nil, err
	}

	return NewAsyncLogger(l, queueLen)
}

func NewSyncSimpleWebBufferFileLogger(path, logId string, bufsize, level int) (ILogger, error) {
	w, err := logWriter.NewFileWriter(path)
	if err != nil {
		return nil, err
	}

	writer := logWriter.NewBufferWriter(w, bufsize)

	return NewSimpleLogger(writer, level, formater.NewWeb(logId))
}

func NewAsyncSimpleWebBufferFileLogger(path, logId string, bufsize, level, queueLen int) (ILogger, error) {
	w, err := logWriter.NewFileWriter(path)
	if err != nil {
		return nil, err
	}

	writer := logWriter.NewBufferWriter(w, bufsize)

	l, err := NewSimpleLogger(writer, level, formater.NewWeb(logId))
	if err != nil {
		return nil, err
	}

	return NewAsyncLogger(l, queueLen)
}
