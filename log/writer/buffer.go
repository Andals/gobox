/**
* @file buffer.go
* @brief writer with buffer
* @author ligang
* @date 2016-02-04
 */

package writer

import (
	"bufio"
	"sync"
	"time"
)

var bfr *bufFlushRoutine
var bfrInitMutex sync.Mutex

// must be called first
func InitBufferAutoFlushRoutine(maxBufNum int, timeInterval time.Duration) {
	bfrInitMutex.Lock()

	if bfr == nil {
		bfr = &bufFlushRoutine{
			buffers: make(map[uint64]*Buffer),

			bufAddCh: make(chan *Buffer, maxBufNum),
			bufDelCh: make(chan *Buffer, maxBufNum),
			freeCh:   make(chan int),
		}

		go bfr.run(timeInterval)
	}

	bfrInitMutex.Unlock()
}

func FreeBuffers() {
	bfr.freeCh <- 1
	<-bfr.freeCh
}

/**
* @name auto flush routine
* @{ */

type bufFlushRoutine struct {
	curIndex uint64
	buffers  map[uint64]*Buffer

	bufAddCh chan *Buffer
	bufDelCh chan *Buffer
	freeCh   chan int
}

func (this *bufFlushRoutine) addBuffer(buf *Buffer) {
	this.bufAddCh <- buf
}

func (this *bufFlushRoutine) delBuffer(buf *Buffer) {
	this.bufDelCh <- buf
}

func (this *bufFlushRoutine) flushAll() {
	for index, buf := range this.buffers {
		if buf == nil || buf.buf == nil {
			delete(this.buffers, index)
		} else {
			buf.Flush()
		}
	}
}

func (this *bufFlushRoutine) run(timeInterval time.Duration) {
	ticker := time.NewTicker(timeInterval)

	for {
		select {
		case buf, _ := <-this.bufAddCh:
			buf.index = this.curIndex
			this.buffers[this.curIndex] = buf
			this.curIndex++
		case buf, _ := <-this.bufDelCh:
			delete(this.buffers, buf.index)
			buf.buf = nil
		case <-ticker.C:
			this.flushAll()
		case <-this.freeCh:
			this.flushAll()
			this.freeCh <- 1
			return
		}
	}
}

/**  @} */

/**
* @name buffer
* @{ */

type Buffer struct {
	w   IWriter
	buf *bufio.Writer

	lockCh chan int
	index  uint64
}

func NewBuffer(w IWriter, bufsize int) *Buffer {
	this := &Buffer{
		w:   w,
		buf: bufio.NewWriterSize(w, bufsize),

		lockCh: make(chan int, 1),
	}

	this.lockCh <- 1
	bfr.addBuffer(this)

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

	bfr.delBuffer(this)
}

/**  @} */
