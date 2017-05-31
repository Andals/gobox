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
)

var bfr *bufFlushRoutine

// must be called first
func InitBufferAutoFlushRoutine(maxBufNum int, timeInterval time.Duration) {
	bfr = &bufFlushRoutine{
		buffers: make(map[string]*Buffer),

		bufAddCh: make(chan *bufAddChItem, maxBufNum),
		bufDelCh: make(chan string, maxBufNum),
	}

	go bfr.run(timeInterval)
}

/**
* @name auto flush routine
* @{ */

type bufAddChItem struct {
	key string
	buf *Buffer
}

type bufFlushRoutine struct {
	buffers map[string]*Buffer

	bufAddCh chan *bufAddChItem
	bufDelCh chan string
}

func (this *bufFlushRoutine) addBuffer(key string, buf *Buffer) {
	this.bufAddCh <- &bufAddChItem{key, buf}
}

func (this *bufFlushRoutine) delBuffer(key string) {
	this.bufDelCh <- key
}

func (this *bufFlushRoutine) run(timeInterval time.Duration) {
	ticker := time.NewTicker(timeInterval)

	for {
		select {
		case item, _ := <-this.bufAddCh:
			this.buffers[item.key] = item.buf
		case key, _ := <-this.bufDelCh:
			delete(this.buffers, key)
		case <-ticker.C:
			for key, buf := range this.buffers {
				if buf == nil {
					delete(this.buffers, key)
				} else {
					buf.Flush()
				}
			}
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
	key    string
}

func NewBuffer(w IWriter, bufsize int) *Buffer {
	this := &Buffer{
		w:   w,
		buf: bufio.NewWriterSize(w, bufsize),

		lockCh: make(chan int, 1),
	}

	this.key = fmt.Sprintf("%p", this)
	this.lockCh <- 1
	bfr.addBuffer(this.key, this)

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

	bfr.delBuffer(this.key)
}

/**  @} */
