/**
* @file logger.go
* @author ligang
* @date 2016-02-04
 */

package log

import (
	"errors"

	logFormater "andals/gobox/log/formater"
	logLevel "andals/gobox/log/level"
	logWriter "andals/gobox/log/writer"
)

type simpleLogger struct {
	globalLevel  int
	w            logWriter.IWriter
	levelWriters map[int]logWriter.IWriter
	formater     logFormater.IFormater

	//replace mutex when logging
	lockCh chan int
}

func NewSimpleLogger(writer logWriter.IWriter, globalLevel int, formater logFormater.IFormater) (*simpleLogger, error) {
	_, ok := logLevel.LogLevels[globalLevel]
	if !ok {
		errors.New("Global level not exists")
	}

	this := &simpleLogger{
		globalLevel:  globalLevel,
		w:            writer,
		levelWriters: make(map[int]logWriter.IWriter),

		lockCh: make(chan int, 1),
	}

	noopWriter := new(logWriter.Noop)
	for level, _ := range logLevel.LogLevels {
		if level < globalLevel {
			this.levelWriters[level] = noopWriter
		} else {
			this.levelWriters[level] = this.w
		}
	}

	if formater == nil {
		formater = new(logFormater.Noop)
	}
	this.formater = formater

	this.lockCh <- 1

	return this, nil
}

func (this *simpleLogger) Debug(msg []byte) {
	this.Log(logLevel.LEVEL_DEBUG, msg)
}

func (this *simpleLogger) Info(msg []byte) {
	this.Log(logLevel.LEVEL_INFO, msg)
}

func (this *simpleLogger) Notice(msg []byte) {
	this.Log(logLevel.LEVEL_NOTICE, msg)
}

func (this *simpleLogger) Warning(msg []byte) {
	this.Log(logLevel.LEVEL_WARNING, msg)
}

func (this *simpleLogger) Error(msg []byte) {
	this.Log(logLevel.LEVEL_ERROR, msg)
}

func (this *simpleLogger) Critical(msg []byte) {
	this.Log(logLevel.LEVEL_CRITICAL, msg)
}

func (this *simpleLogger) Alert(msg []byte) {
	this.Log(logLevel.LEVEL_ALERT, msg)
}

func (this *simpleLogger) Emergency(msg []byte) {
	this.Log(logLevel.LEVEL_EMERGENCY, msg)
}

func (this *simpleLogger) Log(level int, msg []byte) error {
	writer, ok := this.levelWriters[level]
	if !ok {
		errors.New("Level not exists")
	}

	msg = this.formater.Format(level, msg)

	<-this.lockCh
	writer.Write(msg)
	this.lockCh <- 1

	return nil
}

func (this *simpleLogger) Flush() error {
	return this.w.Flush()
}

func (this *simpleLogger) Free() {
	this.w.Free()
}
