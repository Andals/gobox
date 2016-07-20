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

/**
* @name buffer auto flush
* @{ */

var bufferAutoFlushControl struct {
	enabled bool

	lockCh    chan int
	disableCh chan int

	bufferList []*Buffer
}

func init() {
	bufferAutoFlushControl.enabled = false
	bufferAutoFlushControl.lockCh = make(chan int, 1)
	bufferAutoFlushControl.lockCh <- 1
	bufferAutoFlushControl.disableCh = make(chan int)
}

func EnableBufferAutoFlush(timeInterval time.Duration) {
	if bufferAutoFlushControl.enabled || timeInterval <= 0 {
		return
	}

	<-bufferAutoFlushControl.lockCh
	if !bufferAutoFlushControl.enabled {
		go bufferAutoFlushRoutine(timeInterval)

		bufferAutoFlushControl.enabled = true
	}

	bufferAutoFlushControl.lockCh <- 1
}

func bufferAutoFlushRoutine(timeInterval time.Duration) {
	for {
		select {
		case <-time.After(timeInterval):
			for _, buf := range bufferAutoFlushControl.bufferList {
				buf.Flush()
			}
		case <-bufferAutoFlushControl.disableCh:
			bufferAutoFlushControl.disableCh <- 1
			return
		}
	}
}

func DisableBufferAutoFlush() {
	if !bufferAutoFlushControl.enabled {
		return
	}

	<-bufferAutoFlushControl.lockCh
	if bufferAutoFlushControl.enabled {
		bufferAutoFlushControl.disableCh <- 1
		<-bufferAutoFlushControl.disableCh

		bufferAutoFlushControl.enabled = false
	}

	bufferAutoFlushControl.lockCh <- 1
}

/**  @} */

type Buffer struct {
	w   IWriter
	buf *bufio.Writer

	lockCh chan int
}

func NewBufferWriter(writer IWriter, bufsize int) *Buffer {
	this := &Buffer{
		w:   writer,
		buf: bufio.NewWriterSize(writer, bufsize),

		lockCh: make(chan int, 1),
	}

	this.lockCh <- 1
	bufferAutoFlushControl.bufferList = append(bufferAutoFlushControl.bufferList, this)

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
	this.Flush()
	this.w.Free()
}
