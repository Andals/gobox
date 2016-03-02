/**
* @file buffer.go
* @brief writer with buffer
* @author ligang
* @date 2016-02-04
 */

package writer

import (
	"bufio"
	"time"
)

type Buffer struct {
	lockCh            chan int
	freeCh            chan int
	autoFlushInterval time.Duration
	w                 IWriter
	buf               *bufio.Writer
}

func NewBufferWriter(writer IWriter, bufsize int, autoFlushInterval time.Duration) *Buffer {
	this := &Buffer{
		lockCh:            make(chan int, 1),
		freeCh:            make(chan int),
		autoFlushInterval: autoFlushInterval,
		w:                 writer,
		buf:               bufio.NewWriterSize(writer, bufsize),
	}

	this.lockCh <- 1
	go this.flushRoutine()

	return this
}

func (this *Buffer) Write(p []byte) (n int, err error) {
	defer func() {
		this.lockCh <- 1
	}()

	<-this.lockCh
	return this.buf.Write(p)
}

func (this *Buffer) Flush() error {
	defer func() {
		this.lockCh <- 1
	}()

	<-this.lockCh
	return this.buf.Flush()
}

func (this *Buffer) Free() {
	this.freeCh <- 1
	<-this.freeCh

	this.Flush()
	this.w.Free()
	this.buf = nil
}

func (this *Buffer) flushRoutine() {
	defer func() {
		this.freeCh <- 1
	}()

	for {
		select {
		case <-time.After(this.autoFlushInterval):
			this.Flush()
		case <-this.freeCh:
			return
		}
	}
}
