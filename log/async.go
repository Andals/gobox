/**
* @file async.go
* @brief async log use goroutine, one logger one routine
* @author ligang
* @date 2016-02-06
 */

package log

import (
	"errors"
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
	key   string
	msgCh chan *asyncMsg
	stCh  chan int
	wg    *sync.WaitGroup

	*simpleLogger
}

// prevent asyncLoggerContainer race
var acch chan int
var asyncLoggerContainer map[string]*asyncLogger

func init() {
	asyncLoggerContainer = make(map[string]*asyncLogger)

	acch = make(chan int, 1)
	acch <- 1
}

func NewAsyncLogger(key string, logger *simpleLogger, queueLen int) (*asyncLogger, error) {
	defer func() {
		acch <- 1
	}()

	<-acch

	_, ok := asyncLoggerContainer[key]
	if ok {
		return nil, errors.New("key exists")
	}

	this := &asyncLogger{
		key:   key,
		msgCh: make(chan *asyncMsg, queueLen),
		stCh:  make(chan int),
		wg:    new(sync.WaitGroup),

		simpleLogger: logger,
	}

	this.wg.Add(1)
	go this.logRoutine()

	asyncLoggerContainer[key] = this

	return this, nil
}

func FreeAllAsyncLogger() {
	for _, logger := range asyncLoggerContainer {
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
			this.simpleLogger.Log(am.level, am.msg)
		case st, _ := <-this.stCh:
			switch st {
			case ST_FLUSH:
				this.simpleLogger.Flush()
			case ST_FREE:
				for 0 != len(this.msgCh) {
					am, _ := <-this.msgCh
					this.simpleLogger.Log(am.level, am.msg)
				}
				this.simpleLogger.Free()
				return
			}
		}
	}
}
