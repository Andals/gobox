/**
* @file logger.go
* @author ligang
* @date 2016-02-04
 */

package log

import (
	"errors"

	logWriter "andals/gobox/log/writer"
)

type simpleLogger struct {
	globalLevel  int
	w            logWriter.IWriter
	levelWriters map[int]logWriter.IWriter
	formater     IFormater

	//replace mutex when logging
	lockCh chan int
}

func NewSimpleLogger(writer logWriter.IWriter, globalLevel int, formater IFormater) (*simpleLogger, error) {
	_, ok := logLevels[globalLevel]
	if !ok {
		return nil, errors.New("Global level not exists")
	}

	this := &simpleLogger{
		globalLevel:  globalLevel,
		w:            writer,
		levelWriters: make(map[int]logWriter.IWriter),

		lockCh: make(chan int, 1),
	}

	noopWriter := new(logWriter.Noop)
	for level, _ := range logLevels {
		if level < globalLevel {
			this.levelWriters[level] = noopWriter
		} else {
			this.levelWriters[level] = this.w
		}
	}

	if formater == nil {
		formater = new(NoopFormater)
	}
	this.formater = formater

	this.lockCh <- 1

	return this, nil
}

func (this *simpleLogger) Debug(msg []byte) {
	this.Log(LEVEL_DEBUG, msg)
}

func (this *simpleLogger) Info(msg []byte) {
	this.Log(LEVEL_INFO, msg)
}

func (this *simpleLogger) Notice(msg []byte) {
	this.Log(LEVEL_NOTICE, msg)
}

func (this *simpleLogger) Warning(msg []byte) {
	this.Log(LEVEL_WARNING, msg)
}

func (this *simpleLogger) Error(msg []byte) {
	this.Log(LEVEL_ERROR, msg)
}

func (this *simpleLogger) Critical(msg []byte) {
	this.Log(LEVEL_CRITICAL, msg)
}

func (this *simpleLogger) Alert(msg []byte) {
	this.Log(LEVEL_ALERT, msg)
}

func (this *simpleLogger) Emergency(msg []byte) {
	this.Log(LEVEL_EMERGENCY, msg)
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
	close(this.lockCh)
	this.levelWriters = nil
}
