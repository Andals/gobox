/**
* @file buffer.go
* @brief writer with buffer
* @author ligang
* @date 2016-02-04
 */

package writer

import (
	"bufio"
	"fmt"
	"time"

	"andals/gobox/shardmap"
)

/**
* @name buffer auto flush
* @{ */

const (
	BUFFER_MAP_SHARD_CNT = 32
)

var bufferAutoFlushControl struct {
	enabled bool

	lockCh    chan int
	disableCh chan int

	buffers *shardmap.ShardMap
}

func init() {
	bufferAutoFlushControl.enabled = false
	bufferAutoFlushControl.lockCh = make(chan int, 1)
	bufferAutoFlushControl.lockCh <- 1
	bufferAutoFlushControl.disableCh = make(chan int)
	bufferAutoFlushControl.buffers = shardmap.New(BUFFER_MAP_SHARD_CNT)
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

func bufferAutoFlushRoutine(timeInterval time.Duration) {
	for {
		select {
		case <-time.After(timeInterval):
			bufferAutoFlushControl.buffers.Walk(func(k string, v interface{}) {
				buf, ok := v.(*Buffer)
				if buf == nil || !ok {
					bufferAutoFlushControl.buffers.Del(k)
				} else {
					buf.Flush()
				}
			})
		case <-bufferAutoFlushControl.disableCh:
			bufferAutoFlushControl.disableCh <- 1
			return
		}
	}
}

func addAutoFlushBuffer(buf *Buffer) {
	bufferAutoFlushControl.buffers.Set(bufferAutoFlushKey(buf), buf)
}

func delAutoFlushBuffer(buf *Buffer) {
	bufferAutoFlushControl.buffers.Del(bufferAutoFlushKey(buf))
}

func bufferAutoFlushKey(buf *Buffer) string {
	return fmt.Sprintf("%x", buf)
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
	if bufferAutoFlushControl.enabled {
		addAutoFlushBuffer(this)
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
	this.Flush()
	this.w.Free()

	if bufferAutoFlushControl.enabled {
		delAutoFlushBuffer(this)
	}
}
