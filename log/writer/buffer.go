/**
* @file buffer.go
* @brief writer with buffer
* @author ligang
* @date 2016-02-04
 */

package writer

import (
	"bufio"
	"errors"
	"time"
)

/**
* @name buffer writer
* @{ */

type Buffer struct {
	w IWriter

	*bufio.Writer
}

func NewBufferWriter(writer IWriter, bufsize int) *Buffer {
	this := &Buffer{
		w:      writer,
		Writer: bufio.NewWriterSize(writer, bufsize),
	}

	return this
}

func (this *Buffer) Free() {
	this.Writer.Flush()
	this.w.Free()
	this.Writer = nil
}

/**  @} */

/**
* @name buffer writer with time flush
* @{ */

type BufferWithTimeFlush struct {
	lockCh       chan int
	freeCh       chan int
	timeInterval time.Duration
	buf          *Buffer
}

func NewBufferWriterWithTimeFlush(writer IWriter, bufsize int, timeInterval time.Duration) (*BufferWithTimeFlush, error) {
	if timeInterval == 0 {
		return nil, errors.New("time interval equal 0")
	}

	this := &BufferWithTimeFlush{
		lockCh:       make(chan int, 1),
		freeCh:       make(chan int),
		timeInterval: timeInterval,
		buf:          NewBufferWriter(writer, bufsize),
	}

	this.lockCh <- 1
	go this.flushRoutine()

	return this, nil
}

func (this *BufferWithTimeFlush) Write(p []byte) (n int, err error) {
	defer func() {
		this.lockCh <- 1
	}()

	<-this.lockCh
	return this.buf.Write(p)
}

func (this *BufferWithTimeFlush) Flush() error {
	defer func() {
		this.lockCh <- 1
	}()

	<-this.lockCh
	return this.buf.Flush()
}

func (this *BufferWithTimeFlush) Free() {
	this.freeCh <- 1
	<-this.freeCh

	this.Flush()
	this.buf.Free()
}

func (this *BufferWithTimeFlush) flushRoutine() {
	defer func() {
		this.freeCh <- 1
	}()

	for {
		select {
		case <-time.After(this.timeInterval):
			this.Flush()
		case <-this.freeCh:
			return
		}
	}
}

/**  @} */
