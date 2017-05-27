/**
* @file buffer.go
* @brief writer with buffer
* @author ligang
* @date 2016-02-04
 */

package buffer

import (
	"bufio"
	"fmt"
	"sync"
	"time"

	"andals/gobox/log/writer"
)

var fr *flushRoutine

// must be called first
func Init(maxBufNum int, timeInterval time.Duration) {
	fr = &flushRoutine{
		buffers: make(map[string]*Buffer),

		bufAddCh: make(chan *Buffer, maxBufNum),
		bufDelCh: make(chan *Buffer, maxBufNum),
	}

	go fr.run(timeInterval)
}

/**
* @name auto flush routine
* @{ */

type flushRoutine struct {
	buffers map[string]*Buffer
	brwMux  sync.RWMutex

	bufAddCh chan *Buffer
	bufDelCh chan *Buffer
}

func (this *flushRoutine) addBuffer(buf *Buffer) {
	this.bufAddCh <- buf
}

func (this *flushRoutine) delBuffer(buf *Buffer) {
	this.bufDelCh <- buf
}

func (this *flushRoutine) run(timeInterval time.Duration) {
	for {
		select {
		case buf, _ := <-this.bufAddCh:
			key := bfrKey(buf)

			this.brwMux.Lock()
			this.buffers[key] = buf
			this.brwMux.Unlock()
		case buf, _ := <-this.bufDelCh:
			key := bfrKey(buf)

			this.brwMux.Lock()
			delete(this.buffers, key)
			this.brwMux.Unlock()
		case <-time.After(timeInterval):
			this.brwMux.RLock()
			for _, buf := range this.buffers {
				fr.brwMux.RUnlock()
				if buf == nil {
					this.delBuffer(buf)
				} else {
					buf.Flush()
				}
				this.brwMux.RLock()
			}
			this.brwMux.RUnlock()
		}
	}
}

func bfrKey(buf *Buffer) string {
	return fmt.Sprintf("%p", buf)
}

/**  @} */

/**
* @name buffer
* @{ */

type Buffer struct {
	w   writer.IWriter
	buf *bufio.Writer

	lockCh chan int
}

func NewBuffer(w writer.IWriter, bufsize int) *Buffer {
	this := &Buffer{
		w:   w,
		buf: bufio.NewWriterSize(w, bufsize),

		lockCh: make(chan int, 1),
	}

	this.lockCh <- 1
	fr.addBuffer(this)

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

	fr.delBuffer(this)
}

/**  @} */
