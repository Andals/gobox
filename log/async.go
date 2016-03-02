/**
* @file async.go
* @brief async log use goroutine, one logger one routine
* @author ligang
* @date 2016-02-06
 */

package log

import (
	"sync"
)

const (
	ST_FLUSH = 1
	ST_FREE  = 2
)

type asyncMsg struct {
	level int
	msg   []byte
}

type asyncLogger struct {
	msgCh chan *asyncMsg
	stCh  chan int
	wg    *sync.WaitGroup

	ILogger
}

type asyncLoggerList struct {
	// prevent asyncLoggerList race
	lch chan int

	loggers []*asyncLogger
}

var allist asyncLoggerList

func init() {
	allist.lch = make(chan int, 1)
	allist.lch <- 1
}

func NewAsyncLogger(logger ILogger, queueLen int) (*asyncLogger, error) {
	defer func() {
		allist.lch <- 1
	}()

	<-allist.lch

	this := &asyncLogger{
		msgCh: make(chan *asyncMsg, queueLen),
		stCh:  make(chan int),
		wg:    new(sync.WaitGroup),

		ILogger: logger,
	}

	this.wg.Add(1)
	go this.logRoutine()

	allist.loggers = append(allist.loggers, this)

	return this, nil
}

func FreeAllAsyncLogger() {
	for _, logger := range allist.loggers {
		logger.Free()
	}
}

func (this *asyncLogger) Log(level int, msg []byte) error {
	am := &asyncMsg{
		level: level,
		msg:   msg,
	}

	this.msgCh <- am

	return nil
}

func (this *asyncLogger) Flush() error {
	this.stCh <- ST_FLUSH

	return nil
}

func (this *asyncLogger) Free() {
	this.stCh <- ST_FREE

	this.wg.Wait()
}

func (this *asyncLogger) logRoutine() {
	defer this.wg.Done()

	for {
		select {
		case am, _ := <-this.msgCh:
			this.ILogger.Log(am.level, am.msg)
		case st, _ := <-this.stCh:
			switch st {
			case ST_FLUSH:
				this.ILogger.Flush()
			case ST_FREE:
				for 0 != len(this.msgCh) {
					am, _ := <-this.msgCh
					this.ILogger.Log(am.level, am.msg)
				}
				this.ILogger.Free()
				return
			}
		}
	}
}
