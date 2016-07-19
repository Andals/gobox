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
	w   IWriter
	buf *bufio.Writer

	lockCh            chan int
	freeCh            chan int
	flushTimeInterval time.Duration
}

func NewBufferWriter(writer IWriter, bufsize int, flushTimeInterval time.Duration) *Buffer {
	this := &Buffer{
		w:   writer,
		buf: bufio.NewWriterSize(writer, bufsize),

		lockCh: make(chan int, 1),
	}

	this.lockCh <- 1
	if flushTimeInterval > 0 {
		this.freeCh = make(chan int)
		go this.flushRoutine()
	}

	return this
}

func (this *Buffer) Write(p []byte) (int, error) {
	<-this.lockCh
	n, err := this.buf.Write(p)
	this.lockCh <- 1

	return n, err
}

func (this *Buffer) Flush() error {
	<-this.lockCh
	err := this.buf.Flush()
	this.lockCh <- 1

	return err
}

func (this *Buffer) Free() {
	if this.freeCh != nil {
		this.freeCh <- 1
		<-this.freeCh
	}

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
		case <-time.After(this.flushTimeInterval):
			this.Flush()
		case <-this.freeCh:
			return
		}
	}
}
